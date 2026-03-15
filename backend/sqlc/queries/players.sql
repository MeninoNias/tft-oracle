-- name: UpsertPlayer :one
INSERT INTO players (
    puuid, game_name, tag_line, region, platform, summoner_id,
    profile_icon_id, summoner_level, tier, rank, league_points,
    wins, losses, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW()
)
ON CONFLICT (puuid) DO UPDATE SET
    game_name = EXCLUDED.game_name,
    tag_line = EXCLUDED.tag_line,
    region = EXCLUDED.region,
    platform = EXCLUDED.platform,
    summoner_id = EXCLUDED.summoner_id,
    profile_icon_id = EXCLUDED.profile_icon_id,
    summoner_level = EXCLUDED.summoner_level,
    tier = EXCLUDED.tier,
    rank = EXCLUDED.rank,
    league_points = EXCLUDED.league_points,
    wins = EXCLUDED.wins,
    losses = EXCLUDED.losses,
    updated_at = NOW()
RETURNING *;

-- name: GetPlayerByPUUID :one
SELECT * FROM players WHERE puuid = $1;

-- name: GetPlayerByRiotID :one
SELECT * FROM players WHERE game_name = $1 AND tag_line = $2;
