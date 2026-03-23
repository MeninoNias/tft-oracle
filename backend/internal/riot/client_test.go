package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/time/rate"
)

// extractConnectError unwraps a connect.Error from a potentially wrapped error.
func extractConnectError(err error) *connect.Error {
	if err == nil {
		return nil
	}
	if ce, ok := err.(*connect.Error); ok {
		return ce
	}
	// Check wrapped errors
	if unwrapped := err; unwrapped != nil {
		for {
			if ce, ok := unwrapped.(*connect.Error); ok {
				return ce
			}
			type wrapper interface{ Unwrap() error }
			w, ok := unwrapped.(wrapper)
			if !ok {
				break
			}
			unwrapped = w.Unwrap()
		}
	}
	return nil
}

// newTestClient creates a Client pointing at an httptest.Server.
func newTestClient(apiKey string, server *httptest.Server) *Client {
	return &Client{
		http:    server.Client(),
		apiKey:  apiKey,
		limiter: rate.NewLimiter(rate.Inf, 1), // no rate limiting in tests
		baseURL: server.URL,
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("")
	if c.Available() {
		t.Error("client with empty key should not be available")
	}

	c2 := NewClient("RGAPI-test")
	if !c2.Available() {
		t.Error("client with key should be available")
	}
}

func TestClient_Unavailable(t *testing.T) {
	c := NewClient("")
	_, err := c.GetAccountByRiotID(context.Background(), "americas", "test", "NA1")
	if err == nil {
		t.Fatal("expected error for unavailable client")
	}
	var connErr *connect.Error
	if ok := connect.IsNotModifiedError(err); ok {
		t.Fatal("unexpected not modified error")
	}
	connErr, ok := err.(*connect.Error)
	if !ok {
		// error is wrapped, check message
		if connErr == nil {
			t.Logf("error type: %T, message: %v", err, err)
		}
	} else if connErr.Code() != connect.CodeUnavailable {
		t.Errorf("expected CodeUnavailable, got %v", connErr.Code())
	}
}

func TestDoGet_ErrorMapping(t *testing.T) {
	tests := []struct {
		status   int
		expected connect.Code
	}{
		{http.StatusUnauthorized, connect.CodeUnauthenticated},
		{http.StatusForbidden, connect.CodePermissionDenied},
		{http.StatusNotFound, connect.CodeNotFound},
		{http.StatusInternalServerError, connect.CodeInternal},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(`{"error": "test"}`))
			}))
			defer server.Close()

			c := newTestClient("test-key", server)
			_, err := c.GetAccountByRiotID(context.Background(), "americas", "test", "NA1")
			if err == nil {
				t.Fatal("expected error")
			}

			connErr := new(connect.Error)
			connErr = extractConnectError(err)
			if connErr == nil || connErr.Code() != tt.expected {
				t.Errorf("expected code %v, got error: %v", tt.expected, err)
			}
		})
	}
}

func TestRetryOn429(t *testing.T) {
	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := count.Add(1)
		if n <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`rate limited`))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AccountDTO{PUUID: "abc", GameName: "Test", TagLine: "NA1"})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	acct, err := c.GetAccountByRiotID(context.Background(), "americas", "Test", "NA1")
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if acct.PUUID != "abc" {
		t.Errorf("expected puuid 'abc', got %q", acct.PUUID)
	}
	if count.Load() != 3 {
		t.Errorf("expected 3 requests, got %d", count.Load())
	}
}

func TestRetryOn429_Exhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`rate limited`))
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	_, err := c.GetAccountByRiotID(context.Background(), "americas", "Test", "NA1")
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	connErr := extractConnectError(err)
	if connErr == nil || connErr.Code() != connect.CodeResourceExhausted {
		t.Errorf("expected CodeResourceExhausted, got: %v", err)
	}
}

func TestRetryOn429_RespectsContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`rate limited`))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	c := newTestClient("test-key", server)
	_, err := c.GetAccountByRiotID(ctx, "americas", "Test", "NA1")
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestRetryDelay(t *testing.T) {
	tests := []struct {
		name       string
		retryAfter string
		attempt    int
		expected   time.Duration
	}{
		{"retry-after 2s", "2", 0, 2 * time.Second},
		{"no header attempt 0", "", 0, 1 * time.Second},
		{"no header attempt 1", "", 1, 2 * time.Second},
		{"no header attempt 2", "", 2, 4 * time.Second},
		{"no header attempt 5 capped", "", 5, 10 * time.Second},
		{"large retry-after capped", "999", 0, 10 * time.Second},
		{"invalid retry-after fallback", "invalid", 1, 2 * time.Second},
		{"zero retry-after fallback", "0", 0, 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := retryDelay(tt.retryAfter, tt.attempt)
			if got != tt.expected {
				t.Errorf("retryDelay(%q, %d) = %v, want %v", tt.retryAfter, tt.attempt, got, tt.expected)
			}
		})
	}
}

func TestGetAccountByRiotID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Error("missing API key header")
		}
		json.NewEncoder(w).Encode(AccountDTO{PUUID: "p1", GameName: "Player", TagLine: "NA1"})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	acct, err := c.GetAccountByRiotID(context.Background(), "americas", "Player", "NA1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acct.GameName != "Player" || acct.TagLine != "NA1" || acct.PUUID != "p1" {
		t.Errorf("unexpected account: %+v", acct)
	}
}

func TestGetSummonerByPUUID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SummonerDTO{ID: "s1", PUUID: "p1", ProfileIconID: 42, SummonerLevel: 100})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	s, err := c.GetSummonerByPUUID(context.Background(), "br1", "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ID != "s1" || s.SummonerLevel != 100 {
		t.Errorf("unexpected summoner: %+v", s)
	}
}

func TestGetLeagueByPUUID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]LeagueEntryDTO{
			{QueueType: "RANKED_TFT", Tier: "GOLD", Rank: "I", LeaguePoints: 75, Wins: 50, Losses: 30},
		})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	entries, err := c.GetLeagueByPUUID(context.Background(), "br1", "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 || entries[0].Tier != "GOLD" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestGetMatchIDsByPUUID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"BR1_match1", "BR1_match2"})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	ids, err := c.GetMatchIDsByPUUID(context.Background(), "americas", "p1", 20, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 || ids[0] != "BR1_match1" {
		t.Errorf("unexpected ids: %v", ids)
	}
}

func TestGetMatch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(MatchDTO{
			Metadata: MatchMetadataDTO{MatchID: "BR1_1", DataVersion: "5"},
			Info:     MatchInfoDTO{TFTSetNumber: 13, GameLength: 1800},
		})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	m, err := c.GetMatch(context.Background(), "americas", "BR1_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Metadata.MatchID != "BR1_1" || m.Info.TFTSetNumber != 13 {
		t.Errorf("unexpected match: %+v", m)
	}
}

func TestGetChallengerLeague_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(LeagueListDTO{
			Tier:  "CHALLENGER",
			Queue: "RANKED_TFT",
			Entries: []LeagueItemDTO{
				{SummonerID: "s1", LeaguePoints: 1200, Wins: 100, Losses: 50},
			},
		})
	}))
	defer server.Close()

	c := newTestClient("test-key", server)
	league, err := c.GetChallengerLeague(context.Background(), "br1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if league.Tier != "CHALLENGER" || len(league.Entries) != 1 {
		t.Errorf("unexpected league: %+v", league)
	}
}

func TestHealthCheck_ValidKey(t *testing.T) {
	// HealthCheck uses a hardcoded Riot URL, so we can only fully test
	// the no-key path here. The key validation logic is tested via doGet error mapping.
	t.Skip("requires live Riot API — tested indirectly via doGet error mapping tests")
}

func TestHealthCheck_NoKey(t *testing.T) {
	c := NewClient("")
	err := c.HealthCheck(context.Background())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
