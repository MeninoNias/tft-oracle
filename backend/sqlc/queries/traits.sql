-- name: UpsertTrait :one
INSERT INTO traits (api_name, set_number, name, description, icon_url, effects, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT (api_name, set_number) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_url = EXCLUDED.icon_url,
    effects = EXCLUDED.effects,
    updated_at = NOW()
RETURNING *;

-- name: GetTraitsBySet :many
SELECT * FROM traits WHERE set_number = $1 ORDER BY name;

-- name: GetTraitByApiName :one
SELECT * FROM traits WHERE api_name = $1 AND set_number = $2;
