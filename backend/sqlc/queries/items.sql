-- name: UpsertItem :one
INSERT INTO items (api_name, set_number, name, description, composition, effects, icon_url, associated_traits, incompatible_traits, tags, is_unique, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
ON CONFLICT (api_name, set_number) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    composition = EXCLUDED.composition,
    effects = EXCLUDED.effects,
    icon_url = EXCLUDED.icon_url,
    associated_traits = EXCLUDED.associated_traits,
    incompatible_traits = EXCLUDED.incompatible_traits,
    tags = EXCLUDED.tags,
    is_unique = EXCLUDED.is_unique,
    updated_at = NOW()
RETURNING *;

-- name: GetItemsBySet :many
SELECT * FROM items WHERE set_number = $1 ORDER BY name;
