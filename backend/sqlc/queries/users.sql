-- name: CreateUser :one
INSERT INTO users (access_key_hash, riot_puuid, game_name, tag_line, region)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPUUID :one
SELECT * FROM users WHERE riot_puuid = $1;

-- name: GetUserByAccessKeyHash :one
SELECT * FROM users WHERE access_key_hash = $1;

-- name: UpdateLastSeen :exec
UPDATE users SET last_seen = NOW() WHERE id = $1;
