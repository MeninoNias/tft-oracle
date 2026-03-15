package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Client struct {
	rdb *redis.Client
}

// NewClient creates a Redis cache client. Returns nil, nil if redisURL is empty.
func NewClient(redisURL string) (*Client, error) {
	if redisURL == "" {
		return nil, nil
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	return c.rdb.Close()
}

// --- PUUID cache (24h TTL) ---

func (c *Client) GetPUUID(ctx context.Context, gameName, tagLine string) (string, error) {
	if c == nil {
		return "", ErrCacheMiss
	}
	val, err := c.rdb.Get(ctx, puuidKey(gameName, tagLine)).Result()
	if err != nil {
		return "", ErrCacheMiss
	}
	return val, nil
}

func (c *Client) SetPUUID(ctx context.Context, gameName, tagLine, puuid string) {
	if c == nil {
		return
	}
	_ = c.rdb.Set(ctx, puuidKey(gameName, tagLine), puuid, 24*time.Hour).Err()
}

// --- Player profile cache (10min TTL) ---

func (c *Client) GetPlayerProfile(ctx context.Context, puuid string) ([]byte, error) {
	if c == nil {
		return nil, ErrCacheMiss
	}
	val, err := c.rdb.Get(ctx, profileKey(puuid)).Bytes()
	if err != nil {
		return nil, ErrCacheMiss
	}
	return val, nil
}

func (c *Client) SetPlayerProfile(ctx context.Context, puuid string, data []byte) {
	if c == nil {
		return
	}
	_ = c.rdb.Set(ctx, profileKey(puuid), data, 10*time.Minute).Err()
}

// --- Match IDs cache (5min TTL) ---

func (c *Client) GetMatchIDs(ctx context.Context, puuid string) ([]string, error) {
	if c == nil {
		return nil, ErrCacheMiss
	}
	val, err := c.rdb.LRange(ctx, matchIDsKey(puuid), 0, -1).Result()
	if err != nil || len(val) == 0 {
		return nil, ErrCacheMiss
	}
	return val, nil
}

func (c *Client) SetMatchIDs(ctx context.Context, puuid string, ids []string) {
	if c == nil || len(ids) == 0 {
		return
	}
	key := matchIDsKey(puuid)
	pipe := c.rdb.Pipeline()
	pipe.Del(ctx, key)
	vals := make([]interface{}, len(ids))
	for i, id := range ids {
		vals[i] = id
	}
	pipe.RPush(ctx, key, vals...)
	pipe.Expire(ctx, key, 5*time.Minute)
	_, _ = pipe.Exec(ctx)
}

// --- Match detail cache (permanent — matches are immutable) ---

func (c *Client) GetMatchDetail(ctx context.Context, matchID string) ([]byte, error) {
	if c == nil {
		return nil, ErrCacheMiss
	}
	val, err := c.rdb.Get(ctx, matchKey(matchID)).Bytes()
	if err != nil {
		return nil, ErrCacheMiss
	}
	return val, nil
}

func (c *Client) SetMatchDetail(ctx context.Context, matchID string, data []byte) {
	if c == nil {
		return
	}
	_ = c.rdb.Set(ctx, matchKey(matchID), data, 0).Err() // 0 = no expiration
}

// --- Key helpers ---

func puuidKey(gameName, tagLine string) string {
	return fmt.Sprintf("puuid:%s:%s", gameName, tagLine)
}

func profileKey(puuid string) string {
	return fmt.Sprintf("profile:%s", puuid)
}

func matchIDsKey(puuid string) string {
	return fmt.Sprintf("matches:%s", puuid)
}

func matchKey(matchID string) string {
	return fmt.Sprintf("match:%s", matchID)
}
