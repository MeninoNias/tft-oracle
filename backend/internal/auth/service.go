package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/internal/riot"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

var _ tftv1connect.AuthServiceHandler = (*Service)(nil)

type Service struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	riot    riot.RiotAPI
	jwt     *JWTManager
}

func NewService(db *pgxpool.Pool, riotClient riot.RiotAPI, jwtMgr *JWTManager) *Service {
	return &Service{
		db:      db,
		queries: generated.New(db),
		riot:    riotClient,
		jwt:     jwtMgr,
	}
}

func (s *Service) Register(
	ctx context.Context,
	req *connect.Request[tftv1.RegisterRequest],
) (*connect.Response[tftv1.RegisterResponse], error) {
	gameName := req.Msg.GetGameName()
	tagLine := req.Msg.GetTagLine()
	region := req.Msg.GetRegion()

	if gameName == "" || tagLine == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("game_name and tag_line are required"))
	}
	if region == "" {
		region = "br"
	}

	if !s.riot.Available() {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("Riot API not configured"))
	}

	// 1. Verify Riot ID exists
	server := riot.ResolveServer(region)
	account, err := s.riot.GetAccountByRiotID(ctx, server.Region, gameName, tagLine)
	if err != nil {
		return nil, fmt.Errorf("verify riot id: %w", err)
	}

	// 2. Check if user already exists with this PUUID
	existing, err := s.queries.GetUserByPUUID(ctx, account.PUUID)
	if err == nil {
		return nil, connect.NewError(connect.CodeAlreadyExists,
			fmt.Errorf("account %s#%s is already registered (user %s)", existing.GameName, existing.TagLine, uuidToString(existing.ID)))
	}

	// 3. Generate access key (32 bytes → base64url, ~43 chars)
	rawKey, err := generateAccessKey()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("generate key: %w", err))
	}

	// 4. Hash the key with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("hash key: %w", err))
	}

	// 5. Create user in DB
	user, err := s.queries.CreateUser(ctx, generated.CreateUserParams{
		AccessKeyHash: string(hash),
		RiotPuuid:     account.PUUID,
		GameName:      gameName,
		TagLine:       tagLine,
		Region:        region,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("create user: %w", err))
	}

	log.Printf("auth: registered user %s (%s#%s)", uuidToString(user.ID), gameName, tagLine)

	return connect.NewResponse(&tftv1.RegisterResponse{
		AccessKey: rawKey,
		User:      mapUserToProto(user),
	}), nil
}

func (s *Service) Login(
	ctx context.Context,
	req *connect.Request[tftv1.LoginRequest],
) (*connect.Response[tftv1.LoginResponse], error) {
	accessKey := req.Msg.GetAccessKey()
	if accessKey == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("access_key is required"))
	}

	// Find all users and check bcrypt (no way to query by bcrypt hash directly)
	// For a small user base this is fine. For scale, store a key prefix for lookup.
	user, err := s.findUserByAccessKey(ctx, accessKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid access key"))
	}

	// Generate JWT
	token, err := s.jwt.Generate(user.ID, user.RiotPuuid)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("generate token: %w", err))
	}

	// Update last seen
	_ = s.queries.UpdateLastSeen(ctx, user.ID)

	log.Printf("auth: login user %s (%s#%s)", uuidToString(user.ID), user.GameName, user.TagLine)

	return connect.NewResponse(&tftv1.LoginResponse{
		SessionToken: token,
		User:         mapUserToProto(user),
	}), nil
}

func (s *Service) GetCurrentUser(
	ctx context.Context,
	req *connect.Request[tftv1.GetCurrentUserRequest],
) (*connect.Response[tftv1.GetCurrentUserResponse], error) {
	claims, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("invalid user id in token"))
	}

	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	return connect.NewResponse(&tftv1.GetCurrentUserResponse{
		User: mapUserToProto(user),
	}), nil
}

// findUserByAccessKey scans users to find one matching the bcrypt hash.
// For a small-scale desktop app this is acceptable. For production scale,
// store a non-reversible key prefix (first 8 chars of SHA-256) for indexed lookup.
func (s *Service) findUserByAccessKey(ctx context.Context, accessKey string) (generated.User, error) {
	rows, err := s.db.Query(ctx, "SELECT id, access_key_hash, riot_puuid, game_name, tag_line, region, created_at, last_seen FROM users")
	if err != nil {
		return generated.User{}, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u generated.User
		if err := rows.Scan(&u.ID, &u.AccessKeyHash, &u.RiotPuuid, &u.GameName, &u.TagLine, &u.Region, &u.CreatedAt, &u.LastSeen); err != nil {
			continue
		}
		if bcrypt.CompareHashAndPassword([]byte(u.AccessKeyHash), []byte(accessKey)) == nil {
			return u, nil
		}
	}

	return generated.User{}, fmt.Errorf("no matching user")
}

func generateAccessKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func mapUserToProto(u generated.User) *tftv1.User {
	return &tftv1.User{
		Id:       uuidToString(u.ID),
		GameName: u.GameName,
		TagLine:  u.TagLine,
		Region:   u.Region,
		Puuid:    u.RiotPuuid,
	}
}
