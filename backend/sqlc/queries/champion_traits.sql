-- name: InsertChampionTrait :exec
INSERT INTO champion_traits (champion_api_name, trait_api_name, set_number)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: GetTraitsByChampion :many
SELECT ct.trait_api_name, t.name as trait_name
FROM champion_traits ct
JOIN traits t ON t.api_name = ct.trait_api_name AND t.set_number = ct.set_number
WHERE ct.champion_api_name = $1 AND ct.set_number = $2
ORDER BY t.name;

-- name: DeleteChampionTraitsBySet :exec
DELETE FROM champion_traits WHERE set_number = $1;
