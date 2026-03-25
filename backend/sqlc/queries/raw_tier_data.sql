-- name: InsertRawTierData :one
INSERT INTO raw_tier_data (
    source, patch, composition_name, tier,
    win_rate, play_rate, avg_placement,
    champion_ids, core_items, scraped_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
)
RETURNING *;

-- name: GetRawTierDataByPatch :many
SELECT * FROM raw_tier_data
WHERE patch = $1
ORDER BY source, tier, composition_name;

-- name: GetRawTierDataBySourceAndPatch :many
SELECT * FROM raw_tier_data
WHERE source = $1 AND patch = $2
ORDER BY tier, composition_name;

-- name: DeleteRawTierDataBySourceAndPatch :exec
DELETE FROM raw_tier_data
WHERE source = $1 AND patch = $2;

-- name: GetLatestScrapedAt :one
SELECT MAX(scraped_at) FROM raw_tier_data;
