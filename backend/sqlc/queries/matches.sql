-- name: UpsertMatch :one
INSERT INTO matches (
    match_id, data_version, game_datetime, game_length,
    game_version, queue_id, tft_set_number, tft_game_type,
    participant_puuids
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (match_id) DO NOTHING
RETURNING *;

-- name: GetMatchByID :one
SELECT * FROM matches WHERE match_id = $1;

-- name: GetMatchesByPUUID :many
SELECT m.* FROM matches m
WHERE m.participant_puuids @> ARRAY[$1]::TEXT[]
ORDER BY m.game_datetime DESC
LIMIT $2 OFFSET $3;

-- name: MatchExists :one
SELECT EXISTS(SELECT 1 FROM matches WHERE match_id = $1);
