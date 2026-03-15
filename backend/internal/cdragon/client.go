package cdragon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://raw.communitydragon.org/latest/cdragon/tft"

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch downloads and decodes the CommunityDragon TFT JSON for the given locale.
func (c *Client) Fetch(ctx context.Context, locale string) (*CDragonData, error) {
	if locale == "" {
		locale = "en_us"
	}

	url := fmt.Sprintf("%s/%s.json", baseURL, locale)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch cdragon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cdragon returned status %d", resp.StatusCode)
	}

	var data CDragonData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode cdragon json: %w", err)
	}

	return &data, nil
}
