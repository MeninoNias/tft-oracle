-- name: UpsertChampion :one
INSERT INTO champions (api_name, set_number, name, cost, role, stats, ability, icon_url, square_icon_url, tile_icon_url, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
ON CONFLICT (api_name, set_number) DO UPDATE SET
    name = EXCLUDED.name,
    cost = EXCLUDED.cost,
    role = EXCLUDED.role,
    stats = EXCLUDED.stats,
    ability = EXCLUDED.ability,
    icon_url = EXCLUDED.icon_url,
    square_icon_url = EXCLUDED.square_icon_url,
    tile_icon_url = EXCLUDED.tile_icon_url,
    updated_at = NOW()
RETURNING *;

-- name: GetChampionsBySet :many
SELECT * FROM champions WHERE set_number = $1 ORDER BY cost, name;

-- name: GetChampionByApiName :one
SELECT * FROM champions WHERE api_name = $1 AND set_number = $2;

-- name: DeleteChampionsBySet :exec
DELETE FROM champions WHERE set_number = $1;
