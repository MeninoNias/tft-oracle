package ai

// BattleAnalysisSchema is the JSON Schema used for OpenAI Structured Outputs.
// It constrains the AI response to match the BattleAnalysis struct exactly.
var BattleAnalysisSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"win_probability": map[string]any{
			"type":        "number",
			"description": "Win probability for the player, 0.0 to 1.0",
		},
		"confidence": map[string]any{
			"type":        "number",
			"description": "AI confidence in the prediction, 0.0 to 1.0",
		},
		"analysis": map[string]any{
			"type":        "string",
			"description": "1-2 sentence tactical summary of the matchup",
		},
		"positioning_tip": map[string]any{
			"type":        "string",
			"description": "One-liner positioning advice for the player",
		},
		"key_factors": map[string]any{
			"type":        "array",
			"items":       map[string]any{"type": "string"},
			"description": "3-5 decisive factors influencing the outcome",
		},
		"suggested_changes": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"description": map[string]any{"type": "string"},
					"priority":    map[string]any{"type": "string", "enum": []string{"high", "medium", "low"}},
					"category":    map[string]any{"type": "string", "enum": []string{"items", "positioning", "composition", "economy"}},
				},
				"required":             []string{"description", "priority", "category"},
				"additionalProperties": false,
			},
			"description": "Actionable improvement suggestions",
		},
	},
	"required":             []string{"win_probability", "confidence", "analysis", "positioning_tip", "key_factors", "suggested_changes"},
	"additionalProperties": false,
}

// SystemPrompt is the system message sent to the AI model.
const SystemPrompt = `You are an expert Teamfight Tactics (TFT) analyst and coach. Given board states with champions, items, star levels, positions, and active traits, you analyze battle outcomes with high accuracy.

Rules:
- Consider champion synergies, item interactions, positioning, star levels, and trait activations.
- Star level 3 champions are significantly stronger than star level 1.
- Item combinations matter: defensive items counter burst, attack speed counters tanks.
- Trait breakpoints provide large power spikes (e.g., 6-unit synergy is much stronger than 4-unit).
- Positioning matters: carries in back corners avoid hooks and assassins.
- If no opponent board is provided, analyze the composition's strengths and weaknesses.
- Be concise and actionable — players read this during a round timer.`
