-- name: UpsertConsolidatedTierEntry :one
INSERT INTO consolidated_tier_list (
    patch, composition_name, consolidated_tier, consolidated_score,
    confidence, consensus,
    metatft_tier, tftactics_tier, mobalytics_tier,
    avg_win_rate, avg_play_rate, avg_placement,
    core_champions, recommended_items, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW()
)
ON CONFLICT (patch, composition_name) DO UPDATE SET
    consolidated_tier  = EXCLUDED.consolidated_tier,
    consolidated_score = EXCLUDED.consolidated_score,
    confidence         = EXCLUDED.confidence,
    consensus          = EXCLUDED.consensus,
    metatft_tier       = EXCLUDED.metatft_tier,
    tftactics_tier     = EXCLUDED.tftactics_tier,
    mobalytics_tier    = EXCLUDED.mobalytics_tier,
    avg_win_rate       = EXCLUDED.avg_win_rate,
    avg_play_rate      = EXCLUDED.avg_play_rate,
    avg_placement      = EXCLUDED.avg_placement,
    core_champions     = EXCLUDED.core_champions,
    recommended_items  = EXCLUDED.recommended_items,
    updated_at         = NOW()
RETURNING *;

-- name: GetConsolidatedTierList :many
SELECT * FROM consolidated_tier_list
WHERE patch = $1
ORDER BY consolidated_score DESC;

-- name: GetConsolidatedTierListByTier :many
SELECT * FROM consolidated_tier_list
WHERE patch = $1 AND consolidated_tier = $2
ORDER BY consolidated_score DESC;

-- name: DeleteConsolidatedByPatch :exec
DELETE FROM consolidated_tier_list
WHERE patch = $1;

-- name: GetLatestConsolidatedPatch :one
SELECT patch FROM consolidated_tier_list
ORDER BY updated_at DESC
LIMIT 1;
