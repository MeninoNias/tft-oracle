package ai

// CoachMatchSystemPrompt is the system message for per-match coaching analysis.
const CoachMatchSystemPrompt = `You are an expert Teamfight Tactics (TFT) coach reviewing a completed match. Your role is to analyze the player's decisions and provide constructive, specific feedback.

Rules:
- You are reviewing a COMPLETED match — this is post-game analysis, not live coaching.
- Analyze the player's final board: composition, itemization, augment choices, trait activations, and economy (gold left, level).
- Compare the player's board against the current meta tier list data provided.
- Consider the full lobby context — all 8 players' compositions. Identify if the player's comp was contested.
- Grade each skill area honestly. Be constructive but direct — no sugarcoating.
- Acknowledge what the player did well, not just mistakes.
- Be SPECIFIC: use champion names, item names, trait names. Never give generic advice like "build better items".
- The player's rank is provided — calibrate your advice accordingly. A Gold player needs different tips than a Challenger.
- If data is missing (e.g., no tier list available), note it and analyze with what you have.
- Infer economy decisions from final state: high gold_left = hoarding, level 7 at late game = under-leveled, etc.`

// CoachHistorySystemPrompt is the system message for multi-match coaching analysis.
const CoachHistorySystemPrompt = `You are an expert Teamfight Tactics (TFT) coach reviewing a player's recent match history to identify patterns and deliver a coaching plan.

Rules:
- You are analyzing N completed matches to find recurring habits, both good and bad.
- Focus on PATTERNS, not individual match details — what does this player consistently do well or poorly?
- Identify if the player is one-tricking (forcing same comp), flexing well, or forcing badly.
- Compare their most-played compositions against meta viability.
- Look for itemization habits: do they always build the same items? Are those items optimal?
- Look for augment selection patterns: any bias toward certain augments?
- Rate each skill radar dimension based on concrete evidence across matches. Cite specific numbers.
- For trends, compare the first half of the match set vs the second half.
- Give a concrete, prioritized 3-5 step improvement plan. Each step should be actionable in the next game.
- The player's rank is provided — calibrate advice to their skill level.`

// MatchCoachAnalysisSchema is the JSON Schema for per-match coaching structured output.
var MatchCoachAnalysisSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"overall_grade": map[string]any{
			"type":        "string",
			"enum":        []string{"S+", "S", "A+", "A", "B+", "B", "C+", "C", "D", "F"},
			"description": "Overall performance grade for this match",
		},
		"summary": map[string]any{
			"type":        "string",
			"description": "2-3 sentence coaching summary of the match performance",
		},
		"insights": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"category": map[string]any{
						"type": "string",
						"enum": []string{"composition", "itemization", "economy", "augments", "adaptability"},
					},
					"title":  map[string]any{"type": "string", "description": "Short insight title"},
					"detail": map[string]any{"type": "string", "description": "Detailed explanation with specific champion/item names"},
					"grade": map[string]any{
						"type": "string",
						"enum": []string{"good", "okay", "poor"},
					},
				},
				"required":             []string{"category", "title", "detail", "grade"},
				"additionalProperties": false,
			},
			"description": "One insight per category (composition, itemization, economy, augments, adaptability)",
		},
		"meta_comparison": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"closest_meta_comp": map[string]any{"type": "string", "description": "Name of the closest meta composition"},
				"meta_tier": map[string]any{
					"type": "string",
					"enum": []string{"S", "A", "B", "C", "D", "unknown"},
				},
				"missing_units":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Champions missing from the ideal comp"},
				"suboptimal_items": map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Items that differ from BIS"},
				"assessment":       map[string]any{"type": "string", "description": "How close the board was to optimal"},
			},
			"required":             []string{"closest_meta_comp", "meta_tier", "missing_units", "suboptimal_items", "assessment"},
			"additionalProperties": false,
		},
		"lobby_context": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"contested_count":           map[string]any{"type": "integer", "description": "Number of players running similar comps"},
				"lobby_strength_assessment": map[string]any{"type": "string", "description": "Assessment of lobby strength relative to this player"},
				"contested_details":         map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Details on which players contested"},
			},
			"required":             []string{"contested_count", "lobby_strength_assessment", "contested_details"},
			"additionalProperties": false,
		},
		"suggestions": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"description": map[string]any{"type": "string"},
					"priority": map[string]any{
						"type": "string",
						"enum": []string{"high", "medium", "low"},
					},
					"category": map[string]any{
						"type": "string",
						"enum": []string{"composition", "itemization", "economy", "augments", "adaptability"},
					},
				},
				"required":             []string{"description", "priority", "category"},
				"additionalProperties": false,
			},
			"description": "Prioritized improvement suggestions for next games",
		},
	},
	"required":             []string{"overall_grade", "summary", "insights", "meta_comparison", "lobby_context", "suggestions"},
	"additionalProperties": false,
}

// HistoryCoachAnalysisSchema is the JSON Schema for multi-match coaching structured output.
var HistoryCoachAnalysisSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"overall_summary": map[string]any{
			"type":        "string",
			"description": "2-3 sentence summary of the player's recent performance patterns",
		},
		"skill_radar": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"economy":      map[string]any{"type": "number", "description": "Economy management rating 0-100"},
				"itemization":  map[string]any{"type": "number", "description": "Item optimization rating 0-100"},
				"composition":  map[string]any{"type": "number", "description": "Composition quality rating 0-100"},
				"adaptability": map[string]any{"type": "number", "description": "Flexibility and pivoting rating 0-100"},
				"consistency":  map[string]any{"type": "number", "description": "Placement consistency rating 0-100"},
			},
			"required":             []string{"economy", "itemization", "composition", "adaptability", "consistency"},
			"additionalProperties": false,
		},
		"patterns": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"category": map[string]any{
						"type": "string",
						"enum": []string{"composition", "itemization", "economy", "augments", "adaptability"},
					},
					"title":  map[string]any{"type": "string"},
					"detail": map[string]any{"type": "string"},
					"sentiment": map[string]any{
						"type": "string",
						"enum": []string{"positive", "neutral", "negative"},
					},
				},
				"required":             []string{"category", "title", "detail", "sentiment"},
				"additionalProperties": false,
			},
			"description": "Recurring patterns detected across matches",
		},
		"trends": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"metric": map[string]any{"type": "string", "description": "What metric is trending"},
					"direction": map[string]any{
						"type": "string",
						"enum": []string{"improving", "stable", "declining"},
					},
					"detail": map[string]any{"type": "string"},
				},
				"required":             []string{"metric", "direction", "detail"},
				"additionalProperties": false,
			},
			"description": "Performance trends comparing first half vs second half of matches",
		},
		"improvement_plan": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"description": map[string]any{"type": "string"},
					"priority": map[string]any{
						"type": "string",
						"enum": []string{"high", "medium", "low"},
					},
					"category": map[string]any{
						"type": "string",
						"enum": []string{"composition", "itemization", "economy", "augments", "adaptability"},
					},
				},
				"required":             []string{"description", "priority", "category"},
				"additionalProperties": false,
			},
			"description": "Top 3-5 prioritized improvement actions for next games",
		},
	},
	"required":             []string{"overall_summary", "skill_radar", "patterns", "trends", "improvement_plan"},
	"additionalProperties": false,
}
