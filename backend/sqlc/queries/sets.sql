-- name: UpsertSet :one
INSERT INTO sets (number, name, mutator, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (number) DO UPDATE SET
    name = EXCLUDED.name,
    mutator = EXCLUDED.mutator,
    updated_at = NOW()
RETURNING *;

-- name: GetLatestSet :one
SELECT * FROM sets ORDER BY number DESC LIMIT 1;

-- name: GetSetByNumber :one
SELECT * FROM sets WHERE number = $1;
