import { useMemo } from "react";
import type { Champion, Trait } from "@/gen/tft/v1/patch_pb";
import type { PlacedChampion } from "@/stores/board-store";

export interface ActiveSynergy {
  apiName: string;
  name: string;
  count: number;
  threshold: number;
  style: number; // 0=inactive, 1=bronze, 2=silver, 3=gold, 4=chromatic
}

export function useBoardSynergies(
  board: Record<number, PlacedChampion>,
  champions: Champion[],
  traits: Trait[],
): ActiveSynergy[] {
  return useMemo(() => {
    const champMap = new Map(champions.map((c) => [c.apiName, c]));
    const traitMap = new Map(traits.map((t) => [t.apiName, t]));

    // Count units per trait
    const traitCounts = new Map<string, number>();
    for (const placed of Object.values(board)) {
      const champ = champMap.get(placed.championApiName);
      if (!champ) continue;
      for (const traitApi of champ.traitApiNames) {
        traitCounts.set(traitApi, (traitCounts.get(traitApi) ?? 0) + 1);
      }
    }

    // Compute active synergies
    const synergies: ActiveSynergy[] = [];
    for (const [apiName, count] of traitCounts) {
      const trait = traitMap.get(apiName);
      if (!trait) continue;

      let bestThreshold = 0;
      let bestStyle = 0;
      for (const effect of trait.effects) {
        if (count >= effect.minUnits && effect.style > 0) {
          bestThreshold = effect.minUnits;
          bestStyle = effect.style;
        }
      }

      if (bestStyle > 0) {
        synergies.push({
          apiName,
          name: trait.name,
          count,
          threshold: bestThreshold,
          style: bestStyle,
        });
      }
    }

    // Sort: higher style first, then count desc
    synergies.sort((a, b) => {
      if (a.style !== b.style) return b.style - a.style;
      return b.count - a.count;
    });

    return synergies;
  }, [board, champions, traits]);
}
