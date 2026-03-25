package consolidation

// JaccardSimilarity computes the Jaccard coefficient between two string slices.
// Returns a value between 0.0 (no overlap) and 1.0 (identical).
func JaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}

	setA := make(map[string]struct{}, len(a))
	for _, v := range a {
		setA[v] = struct{}{}
	}

	setB := make(map[string]struct{}, len(b))
	for _, v := range b {
		setB[v] = struct{}{}
	}

	intersection := 0
	for k := range setA {
		if _, ok := setB[k]; ok {
			intersection++
		}
	}

	// Union = |A| + |B| - |A ∩ B|
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

const matchThreshold = 0.8

// MatchCompositions groups compositions from multiple sources by champion overlap.
// Compositions with >= 80% Jaccard similarity are considered the same comp.
func MatchCompositions(bySource map[string][]NormalizedComp) []MatchGroup {
	var groups []MatchGroup

	// Process each source in a deterministic order
	sourceOrder := []string{"mobalytics", "tacticstools", "metatft"}

	for _, source := range sourceOrder {
		comps, ok := bySource[source]
		if !ok {
			continue
		}

		for _, comp := range comps {
			bestIdx := -1
			bestSim := 0.0

			for i, group := range groups {
				sim := JaccardSimilarity(comp.ChampionIDs, group.ChampionIDs)
				if sim >= matchThreshold && sim > bestSim {
					bestIdx = i
					bestSim = sim
				}
			}

			if bestIdx >= 0 {
				// Add to existing group
				groups[bestIdx].Sources = append(groups[bestIdx].Sources, comp)
				// Merge champion IDs (union)
				groups[bestIdx].ChampionIDs = unionStrings(groups[bestIdx].ChampionIDs, comp.ChampionIDs)
				// Merge core items
				for k, v := range comp.CoreItems {
					if _, exists := groups[bestIdx].CoreItems[k]; !exists {
						groups[bestIdx].CoreItems[k] = v
					}
				}
			} else {
				// Create new group
				items := make(map[string][]string)
				for k, v := range comp.CoreItems {
					items[k] = v
				}
				groups = append(groups, MatchGroup{
					Name:        comp.Name,
					Sources:     []NormalizedComp{comp},
					ChampionIDs: append([]string{}, comp.ChampionIDs...),
					CoreItems:   items,
				})
			}
		}
	}

	return groups
}

// unionStrings returns the union of two string slices with no duplicates.
func unionStrings(a, b []string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	var result []string

	for _, v := range a {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	for _, v := range b {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
