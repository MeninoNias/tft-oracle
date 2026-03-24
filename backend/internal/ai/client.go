package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/packages/param"
)

const (
	defaultModel   = "gpt-4o-mini"
	defaultTimeout = 30 * time.Second
)

// Client wraps the OpenAI API for TFT battle analysis.
type Client struct {
	client *openai.Client
	model  string
}

// NewClient creates an OpenAI client. If apiKey is empty, all methods return an error.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		return &Client{}
	}
	c := openai.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		client: &c,
		model:  defaultModel,
	}
}

// Available returns true if the OpenAI API key is configured.
func (c *Client) Available() bool {
	return c.client != nil
}

// AnalyzeBattle sends a battle prompt to GPT and returns structured analysis.
func (c *Client) AnalyzeBattle(ctx context.Context, prompt string) (*BattleAnalysis, error) {
	if !c.Available() {
		return nil, fmt.Errorf("OPENAI_API_KEY not configured")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	completion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(c.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(SystemPrompt),
			openai.UserMessage(prompt),
		},
		Temperature: param.Opt[float64]{Value: 0.3},
		MaxCompletionTokens: param.Opt[int64]{Value: 1000},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "battle_analysis",
					Strict: param.Opt[bool]{Value: true},
					Schema: BattleAnalysisSchema,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai chat completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	// Log token usage for cost tracking
	if completion.Usage.TotalTokens > 0 {
		log.Printf("openai: %d tokens used (prompt=%d, completion=%d)",
			completion.Usage.TotalTokens,
			completion.Usage.PromptTokens,
			completion.Usage.CompletionTokens)
	}

	content := completion.Choices[0].Message.Content
	var analysis BattleAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("parse ai response: %w", err)
	}

	return &analysis, nil
}
