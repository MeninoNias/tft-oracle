import { TerminalCard } from "@/components/ui/terminal-card";
import { PlacementBadge } from "@/components/player/placement-badge";
import { UnitIcon } from "@/components/player/unit-icon";
import { TraitBadge } from "@/components/player/trait-badge";
import type { Match } from "@/gen/tft/v1/player_pb";
import type { PatchLookup } from "@/hooks/use-patch-lookup";

interface MatchDetailProps {
  match: Match;
  playerPuuid: string;
  patchLookup: PatchLookup;
}

export function MatchDetail({
  match,
  playerPuuid,
  patchLookup,
}: MatchDetailProps) {
  const info = match.info;
  if (!info) return null;

  const sorted = [...info.participants].sort(
    (a, b) => a.placement - b.placement,
  );

  return (
    <TerminalCard className="mt-1 rounded-t-none border-t-0">
      <div className="space-y-2">
        {sorted.map((p) => {
          const isPlayer = p.puuid === playerPuuid;
          const activeTraits = p.traits
            .filter((t) => t.style > 0)
            .sort((a, b) => b.style - a.style);
          const sortedUnits = [...p.units].sort(
            (a, b) => b.rarity - a.rarity || b.tier - a.tier,
          );

          return (
            <div
              key={p.puuid}
              className={`flex items-center gap-3 rounded-sm p-2 ${
                isPlayer ? "bg-lofi-accent/10" : ""
              }`}
            >
              <PlacementBadge placement={p.placement} />

              <div className="w-[80px] shrink-0 space-y-0.5">
                <p className="text-xs text-lofi-muted">
                  lv{p.level} // {p.goldLeft}g
                </p>
                <p className="text-xs text-lofi-muted">
                  ⚔ {p.totalDamageToPlayers}
                </p>
              </div>

              <div className="flex w-[110px] shrink-0 flex-wrap gap-0.5">
                {activeTraits.slice(0, 5).map((trait) => (
                  <TraitBadge
                    key={trait.apiName}
                    trait={trait}
                    lookup={patchLookup}
                  />
                ))}
              </div>

              <div className="flex min-w-0 flex-1 items-start gap-1 overflow-x-auto">
                {sortedUnits.map((unit, i) => (
                  <UnitIcon key={i} unit={unit} lookup={patchLookup} />
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </TerminalCard>
  );
}
