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
	defaultModel        = "gpt-4o-mini"
	defaultTimeout      = 30 * time.Second
	coachTimeout        = 60 * time.Second
	coachMatchMaxTokens = 2000
	coachHistMaxTokens  = 3000
	coachTemperature    = 0.4
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

// AnalyzeMatchCoaching sends match data to GPT and returns structured coaching analysis.
// Returns the analysis and total tokens used for cost tracking.
func (c *Client) AnalyzeMatchCoaching(ctx context.Context, prompt string) (*MatchCoachAnalysis, int64, error) {
	if !c.Available() {
		return nil, 0, fmt.Errorf("OPENAI_API_KEY not configured")
	}

	ctx, cancel := context.WithTimeout(ctx, coachTimeout)
	defer cancel()

	completion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(c.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(CoachMatchSystemPrompt),
			openai.UserMessage(prompt),
		},
		Temperature:         param.Opt[float64]{Value: coachTemperature},
		MaxCompletionTokens: param.Opt[int64]{Value: coachMatchMaxTokens},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "match_coach_analysis",
					Strict: param.Opt[bool]{Value: true},
					Schema: MatchCoachAnalysisSchema,
				},
			},
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("openai coach match analysis: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, 0, fmt.Errorf("openai returned no choices")
	}

	tokensUsed := completion.Usage.TotalTokens
	if tokensUsed > 0 {
		log.Printf("openai coach-match: %d tokens used (prompt=%d, completion=%d)",
			tokensUsed, completion.Usage.PromptTokens, completion.Usage.CompletionTokens)
	}

	content := completion.Choices[0].Message.Content
	var analysis MatchCoachAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, tokensUsed, fmt.Errorf("parse coach match response: %w", err)
	}

	return &analysis, tokensUsed, nil
}

// AnalyzeHistoryCoaching sends multi-match data to GPT and returns structured coaching analysis.
// Returns the analysis and total tokens used for cost tracking.
func (c *Client) AnalyzeHistoryCoaching(ctx context.Context, prompt string) (*HistoryCoachAnalysis, int64, error) {
	if !c.Available() {
		return nil, 0, fmt.Errorf("OPENAI_API_KEY not configured")
	}

	ctx, cancel := context.WithTimeout(ctx, coachTimeout)
	defer cancel()

	completion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(c.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(CoachHistorySystemPrompt),
			openai.UserMessage(prompt),
		},
		Temperature:         param.Opt[float64]{Value: coachTemperature},
		MaxCompletionTokens: param.Opt[int64]{Value: coachHistMaxTokens},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "history_coach_analysis",
					Strict: param.Opt[bool]{Value: true},
					Schema: HistoryCoachAnalysisSchema,
				},
			},
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("openai coach history analysis: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, 0, fmt.Errorf("openai returned no choices")
	}

	tokensUsed := completion.Usage.TotalTokens
	if tokensUsed > 0 {
		log.Printf("openai coach-history: %d tokens used (prompt=%d, completion=%d)",
			tokensUsed, completion.Usage.PromptTokens, completion.Usage.CompletionTokens)
	}

	content := completion.Choices[0].Message.Content
	var analysis HistoryCoachAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, tokensUsed, fmt.Errorf("parse coach history response: %w", err)
	}

	return &analysis, tokensUsed, nil
}
