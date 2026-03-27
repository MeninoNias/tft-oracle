package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	maxRetries     = 3
	baseRetryDelay = 1 * time.Second
	maxRetryDelay  = 10 * time.Second
	defaultTimeout = 30 * time.Second
	userAgent      = "TFTOracle/1.0 (https://github.com/MeninoNias/tft-oracle)"
)

// HTTPClient is a shared HTTP client with retry and rate limiting for web scraping.
type HTTPClient struct {
	http    *http.Client
	limiter *rate.Limiter
	baseURL string // override for tests
}

// NewHTTPClient creates a scraping HTTP client.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		http: &http.Client{
			Timeout: defaultTimeout,
		},
		// Conservative: 1 request per second per scraper
		limiter: rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

// Get fetches a URL with rate limiting, retry logic, and proper headers.
func (c *HTTPClient) Get(ctx context.Context, rawURL string) ([]byte, error) {
	if c.baseURL != "" {
		rawURL = c.baseURL + rawURL
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "text/html,application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http request: %w", err)
			if attempt < maxRetries {
				sleepWithContext(ctx, retryDelay(attempt))
				continue
			}
			return nil, lastErr
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return body, nil
		}

		lastErr = fmt.Errorf("http %d: %s", resp.StatusCode, truncate(string(body), 200))

		// Retry on server errors and rate limits
		if (resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests) && attempt < maxRetries {
			sleepWithContext(ctx, retryDelay(attempt))
			continue
		}

		return nil, lastErr
	}
	return nil, lastErr
}

func retryDelay(attempt int) time.Duration {
	d := baseRetryDelay * (1 << uint(attempt))
	if d > maxRetryDelay {
		return maxRetryDelay
	}
	return d
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// Post sends a POST request with JSON body (needed for GraphQL endpoints).
func (c *HTTPClient) Post(ctx context.Context, rawURL string, body []byte) ([]byte, error) {
	if c.baseURL != "" {
		rawURL = c.baseURL + rawURL
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http request: %w", err)
			if attempt < maxRetries {
				sleepWithContext(ctx, retryDelay(attempt))
				continue
			}
			return nil, lastErr
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return respBody, nil
		}

		lastErr = fmt.Errorf("http %d: %s", resp.StatusCode, truncate(string(respBody), 200))
		if (resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests) && attempt < maxRetries {
			sleepWithContext(ctx, retryDelay(attempt))
			continue
		}
		return nil, lastErr
	}
	return nil, lastErr
}
