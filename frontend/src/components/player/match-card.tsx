import { useState } from "react";
import { TerminalCard } from "@/components/ui/terminal-card";
import { PlacementBadge } from "@/components/player/placement-badge";
import { UnitIcon } from "@/components/player/unit-icon";
import { TraitBadge } from "@/components/player/trait-badge";
import { MatchDetail } from "@/components/player/match-detail";
import type { Match } from "@/gen/tft/v1/player_pb";
import type { PatchLookup } from "@/hooks/use-patch-lookup";

interface MatchCardProps {
  match: Match;
  playerPuuid: string;
  patchLookup: PatchLookup;
}

export function MatchCard({
  match,
  playerPuuid,
  patchLookup,
}: MatchCardProps) {
  const [expanded, setExpanded] = useState(false);

  const info = match.info;
  const metadata = match.metadata;
  if (!info || !metadata) return null;

  const player = info.participants.find((p) => p.puuid === playerPuuid);
  if (!player) return null;

  const gameDate = new Date(Number(info.gameDatetime));
  const timeAgo = getTimeAgo(gameDate);
  const gameLengthMin = (info.gameLength / 60).toFixed(0);
  const roundStr = formatRound(player.lastRound);

  const activeTraits = player.traits
    .filter((t) => t.style > 0)
    .sort((a, b) => b.style - a.style);

  const sortedUnits = [...player.units].sort(
    (a, b) => b.rarity - a.rarity || b.tier - a.tier,
  );

  return (
    <div>
      <TerminalCard
        className={`cursor-pointer transition-colors hover:border-lofi-muted ${
          player.placement <= 4
            ? "border-l-2 border-l-green-500/40"
            : "border-l-2 border-l-red-500/30"
        }`}
      >
        <div onClick={() => setExpanded(!expanded)}>
          <div className="flex items-center gap-4">
            {/* Placement */}
            <PlacementBadge
              placement={player.placement}
              size="md"
              className="shrink-0"
            />

            {/* Meta info */}
            <div className="w-[120px] shrink-0 space-y-0.5">
              <p className="text-xs font-medium text-lofi-secondary">
                {info.tftGameType === "standard" ? "ranked" : info.tftGameType}
              </p>
              <p className="text-xs text-lofi-muted">
                {timeAgo} // {gameLengthMin}m
              </p>
              <p className="text-xs text-lofi-muted">
                lv{player.level} // {roundStr}
              </p>
              <div className="flex items-center gap-2">
                <span
                  className="text-xs text-lofi-secondary"
                  title="damage dealt"
                >
                  ⚔ {player.totalDamageToPlayers}
                </span>
                {player.playersEliminated > 0 && (
                  <span className="text-xs font-medium text-red-400">
                    {player.playersEliminated} kill
                    {player.playersEliminated > 1 ? "s" : ""}
                  </span>
                )}
              </div>
            </div>

            {/* Traits */}
            <div className="flex w-[150px] shrink-0 flex-wrap content-start gap-1">
              {activeTraits.slice(0, 8).map((trait) => (
                <TraitBadge
                  key={trait.apiName}
                  trait={trait}
                  lookup={patchLookup}
                />
              ))}
            </div>

            {/* Units */}
            <div className="flex min-w-0 flex-1 items-start gap-1.5 overflow-x-auto">
              {sortedUnits.map((unit, i) => (
                <UnitIcon key={i} unit={unit} lookup={patchLookup} />
              ))}
            </div>

            {/* Expand */}
            <button className="shrink-0 px-1 text-xs text-lofi-muted hover:text-lofi-text">
              {expanded ? "▲" : "▼"}
            </button>
          </div>
        </div>
      </TerminalCard>

      {expanded && (
        <MatchDetail
          match={match}
          playerPuuid={playerPuuid}
          patchLookup={patchLookup}
        />
      )}
    </div>
  );
}

function formatRound(lastRound: number): string {
  if (lastRound <= 3) return `1-${lastRound}`;
  const adjusted = lastRound - 3;
  const stage = Math.floor(adjusted / 7) + 2;
  const round = (adjusted % 7) + 1;
  return `${stage}-${round}`;
}

function getTimeAgo(date: Date): string {
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}
