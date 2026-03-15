-- name: UpsertMatchParticipant :exec
INSERT INTO match_participants (
    match_id, puuid, placement, level, gold_left, last_round,
    time_eliminated, total_damage_to_players, players_eliminated,
    augments, companion, traits, units
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
ON CONFLICT (match_id, puuid) DO NOTHING;

-- name: GetParticipantsByMatch :many
SELECT * FROM match_participants WHERE match_id = $1 ORDER BY placement;

-- name: GetParticipantsByPUUID :many
SELECT mp.* FROM match_participants mp
JOIN matches m ON mp.match_id = m.match_id
WHERE mp.puuid = $1
ORDER BY m.game_datetime DESC
LIMIT $2 OFFSET $3;

-- name: GetPlacementDistribution :many
SELECT placement, COUNT(*) as count
FROM match_participants
WHERE puuid = $1
GROUP BY placement
ORDER BY placement;
