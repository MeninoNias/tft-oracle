package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/time/rate"
)

const (
	maxRetries     = 3
	baseRetryDelay = 1 * time.Second
	maxRetryDelay  = 10 * time.Second
)

// RiotAPI defines the contract for Riot Games API operations.
type RiotAPI interface {
	Available() bool
	HealthCheck(ctx context.Context) error
	GetAccountByRiotID(ctx context.Context, region, gameName, tagLine string) (*AccountDTO, error)
	GetSummonerByPUUID(ctx context.Context, platform, puuid string) (*SummonerDTO, error)
	GetLeagueByPUUID(ctx context.Context, platform, puuid string) ([]LeagueEntryDTO, error)
	GetMatchIDsByPUUID(ctx context.Context, region, puuid string, count, start int32) ([]string, error)
	GetMatch(ctx context.Context, region, matchID string) (*MatchDTO, error)
	GetChallengerLeague(ctx context.Context, platform string) (*LeagueListDTO, error)
}

var _ RiotAPI = (*Client)(nil)

// Client is an HTTP client for the Riot Games API.
type Client struct {
	http    *http.Client
	apiKey  string
	limiter *rate.Limiter
	baseURL string // override for tests; empty = use real Riot URLs
}

// NewClient creates a new Riot API client.
// apiKey can be empty — methods will return CodeUnavailable.
func NewClient(apiKey string) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey: apiKey,
		// Dev key limit: 20 requests per second, 100 per 2 minutes
		limiter: rate.NewLimiter(rate.Every(50*time.Millisecond), 20),
	}
}

// Available returns true if the Riot API key is configured.
func (c *Client) Available() bool {
	return c.apiKey != ""
}

// regionURL returns the base URL for a regional endpoint.
func (c *Client) regionURL(region string) string {
	if c.baseURL != "" {
		return c.baseURL
	}
	return RegionToBaseURL(region)
}

// platformURL returns the base URL for a platform endpoint.
func (c *Client) platformURL(platform string) string {
	if c.baseURL != "" {
		return c.baseURL
	}
	return PlatformToBaseURL(platform)
}

// HealthCheck validates the API key by calling the Account V1 endpoint.
func (c *Client) HealthCheck(ctx context.Context) error {
	if !c.Available() {
		return fmt.Errorf("RIOT_API_KEY not configured")
	}

	// Use a known-good endpoint to verify the key works.
	// Account V1 on americas with a dummy lookup — 404 means key is valid, 401/403 means bad key.
	testURL := "https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/healthcheck/test"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-Riot-Token", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNotFound:
		// 200 = found (unlikely for healthcheck), 404 = not found (key is valid)
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("API key is invalid (401)")
	case http.StatusForbidden:
		return fmt.Errorf("API key is expired or forbidden (403)")
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

// --- Riot API Endpoints ---
//
// Routing:
//   Account V1:       regional (americas, europe, asia, sea)
//   TFT Summoner V1:  platform (br1, na1, euw1, kr, etc.)
//   TFT League V1:    platform (br1, na1, euw1, kr, etc.)
//   TFT Match V1:     regional (americas, europe, asia, sea)

// GetAccountByRiotID looks up an account by Riot ID (Account V1, regional).
func (c *Client) GetAccountByRiotID(ctx context.Context, region, gameName, tagLine string) (*AccountDTO, error) {
	u := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s",
		c.regionURL(region), url.PathEscape(gameName), url.PathEscape(tagLine))
	var dto AccountDTO
	if err := c.doGet(ctx, u, &dto); err != nil {
		return nil, fmt.Errorf("get account by riot id: %w", err)
	}
	return &dto, nil
}

// GetSummonerByPUUID gets summoner data (TFT Summoner V1, platform).
func (c *Client) GetSummonerByPUUID(ctx context.Context, platform, puuid string) (*SummonerDTO, error) {
	u := fmt.Sprintf("%s/tft/summoner/v1/summoners/by-puuid/%s",
		c.platformURL(platform), puuid)
	var dto SummonerDTO
	if err := c.doGet(ctx, u, &dto); err != nil {
		return nil, fmt.Errorf("get summoner by puuid: %w", err)
	}
	return &dto, nil
}

// GetLeagueByPUUID gets ranked entries (TFT League V1, platform).
func (c *Client) GetLeagueByPUUID(ctx context.Context, platform, puuid string) ([]LeagueEntryDTO, error) {
	u := fmt.Sprintf("%s/tft/league/v1/by-puuid/%s",
		c.platformURL(platform), puuid)
	var dto []LeagueEntryDTO
	if err := c.doGet(ctx, u, &dto); err != nil {
		return nil, fmt.Errorf("get league entries: %w", err)
	}
	return dto, nil
}

// GetMatchIDsByPUUID gets recent match IDs (TFT Match V1, regional).
func (c *Client) GetMatchIDsByPUUID(ctx context.Context, region, puuid string, count, start int32) ([]string, error) {
	u := fmt.Sprintf("%s/tft/match/v1/matches/by-puuid/%s/ids?count=%d&start=%d",
		c.regionURL(region), puuid, count, start)
	var ids []string
	if err := c.doGet(ctx, u, &ids); err != nil {
		return nil, fmt.Errorf("get match ids: %w", err)
	}
	return ids, nil
}

// GetMatch gets full match data (TFT Match V1, regional).
func (c *Client) GetMatch(ctx context.Context, region, matchID string) (*MatchDTO, error) {
	u := fmt.Sprintf("%s/tft/match/v1/matches/%s",
		c.regionURL(region), matchID)
	var dto MatchDTO
	if err := c.doGet(ctx, u, &dto); err != nil {
		return nil, fmt.Errorf("get match %s: %w", matchID, err)
	}
	return &dto, nil
}

// GetChallengerLeague gets the Challenger league (TFT League V1, platform).
func (c *Client) GetChallengerLeague(ctx context.Context, platform string) (*LeagueListDTO, error) {
	u := fmt.Sprintf("%s/tft/league/v1/challenger", c.platformURL(platform))
	var dto LeagueListDTO
	if err := c.doGet(ctx, u, &dto); err != nil {
		return nil, fmt.Errorf("get challenger league: %w", err)
	}
	return &dto, nil
}

func (c *Client) doGet(ctx context.Context, rawURL string, out interface{}) error {
	if !c.Available() {
		return connect.NewError(connect.CodeUnavailable, fmt.Errorf("RIOT_API_KEY not configured"))
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := c.limiter.Wait(ctx); err != nil {
			return fmt.Errorf("rate limiter: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("X-Riot-Token", c.apiKey)

		resp, err := c.http.Do(req)
		if err != nil {
			return fmt.Errorf("http request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				return fmt.Errorf("decode response: %w", err)
			}
			return nil
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusTooManyRequests || attempt == maxRetries {
			return mapHTTPError(resp.StatusCode, body)
		}

		// 429 — retry with backoff
		lastErr = mapHTTPError(resp.StatusCode, body)
		delay := retryDelay(resp.Header.Get("Retry-After"), attempt)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return lastErr
}

func retryDelay(retryAfter string, attempt int) time.Duration {
	if retryAfter != "" {
		if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 {
			d := time.Duration(seconds) * time.Second
			if d > maxRetryDelay {
				return maxRetryDelay
			}
			return d
		}
	}
	// Exponential backoff: 1s, 2s, 4s, ...
	d := baseRetryDelay * (1 << uint(attempt))
	if d > maxRetryDelay {
		return maxRetryDelay
	}
	return d
}

func mapHTTPError(status int, body []byte) error {
	msg := string(body)
	switch status {
	case http.StatusNotFound:
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("not found: %s", msg))
	case http.StatusForbidden:
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("forbidden: %s", msg))
	case http.StatusUnauthorized:
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("unauthorized: %s", msg))
	case http.StatusTooManyRequests:
		return connect.NewError(connect.CodeResourceExhausted, fmt.Errorf("rate limited: %s", msg))
	default:
		return connect.NewError(connect.CodeInternal, fmt.Errorf("riot api error %d: %s", status, msg))
	}
}
