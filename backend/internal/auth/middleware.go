package auth

import (
	"context"
	"strings"

	"connectrpc.com/connect"
)

type contextKey string

const userContextKey contextKey = "auth_user"

// UserFromContext extracts the authenticated user claims from the context.
// Returns nil if no user is authenticated.
func UserFromContext(ctx context.Context) *Claims {
	claims, _ := ctx.Value(userContextKey).(*Claims)
	return claims
}

// NewAuthInterceptor creates a Connect RPC interceptor that validates JWT tokens.
// It extracts the token from the Authorization header (Bearer scheme).
// Unauthenticated requests pass through — individual handlers decide if auth is required.
func NewAuthInterceptor(jwtMgr *JWTManager) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if jwtMgr == nil {
				return next(ctx, req)
			}

			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				return next(ctx, req)
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				// No "Bearer " prefix — skip
				return next(ctx, req)
			}

			claims, err := jwtMgr.Validate(token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = context.WithValue(ctx, userContextKey, claims)
			return next(ctx, req)
		}
	}
}

// RequireAuth is a helper that checks if the request is authenticated.
// Call this at the start of handlers that require authentication.
func RequireAuth(ctx context.Context) (*Claims, error) {
	claims := UserFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}
	return claims, nil
}
