-- name: GetCoachMatchAnalysis :one
SELECT * FROM coach_match_analyses
WHERE match_id = $1 AND puuid = $2;

-- name: UpsertCoachMatchAnalysis :one
INSERT INTO coach_match_analyses (match_id, puuid, response, model, tokens_used)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (match_id, puuid) DO UPDATE SET
    response = EXCLUDED.response,
    model = EXCLUDED.model,
    tokens_used = EXCLUDED.tokens_used,
    created_at = NOW()
RETURNING *;

-- name: GetCoachHistoryAnalysis :one
SELECT * FROM coach_history_analyses
WHERE puuid = $1 AND match_count = $2
  AND created_at > NOW() - INTERVAL '1 hour';

-- name: UpsertCoachHistoryAnalysis :one
INSERT INTO coach_history_analyses (puuid, match_count, match_ids, response, model, tokens_used)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (puuid, match_count) DO UPDATE SET
    match_ids = EXCLUDED.match_ids,
    response = EXCLUDED.response,
    model = EXCLUDED.model,
    tokens_used = EXCLUDED.tokens_used,
    created_at = NOW()
RETURNING *;
