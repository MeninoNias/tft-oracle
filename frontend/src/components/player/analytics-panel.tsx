import { useMemo } from "react";
import { TerminalCard } from "@/components/ui/terminal-card";
import type { Match } from "@/gen/tft/v1/player_pb";

interface AnalyticsPanelProps {
  matches: Match[];
  puuid: string;
}

interface TraitStat {
  apiName: string;
  games: number;
  avgPlacement: number;
  top4Count: number;
}

interface ChampionStat {
  characterId: string;
  games: number;
  avgPlacement: number;
}

export function AnalyticsPanel({ matches, puuid }: AnalyticsPanelProps) {
  const stats = useMemo(() => {
    const placements: number[] = [];
    const traitMap = new Map<
      string,
      { games: number; totalPlacement: number; top4: number }
    >();
    const champMap = new Map<
      string,
      { games: number; totalPlacement: number }
    >();

    for (const match of matches) {
      const p = match.info?.participants.find((p) => p.puuid === puuid);
      if (!p) continue;

      placements.push(p.placement);

      for (const trait of p.traits) {
        if (trait.style === 0) continue;
        const existing = traitMap.get(trait.apiName) ?? {
          games: 0,
          totalPlacement: 0,
          top4: 0,
        };
        existing.games++;
        existing.totalPlacement += p.placement;
        if (p.placement <= 4) existing.top4++;
        traitMap.set(trait.apiName, existing);
      }

      for (const unit of p.units) {
        const existing = champMap.get(unit.characterId) ?? {
          games: 0,
          totalPlacement: 0,
        };
        existing.games++;
        existing.totalPlacement += p.placement;
        champMap.set(unit.characterId, existing);
      }
    }

    const avgPlacement =
      placements.length > 0
        ? placements.reduce((a, b) => a + b, 0) / placements.length
        : 0;
    const top4Rate =
      placements.length > 0
        ? (placements.filter((p) => p <= 4).length / placements.length) * 100
        : 0;
    const winRate =
      placements.length > 0
        ? (placements.filter((p) => p === 1).length / placements.length) * 100
        : 0;

    const placementDist = Array.from({ length: 8 }, (_, i) => ({
      placement: i + 1,
      count: placements.filter((p) => p === i + 1).length,
    }));

    const topTraits: TraitStat[] = Array.from(traitMap.entries())
      .map(([apiName, data]) => ({
        apiName,
        games: data.games,
        avgPlacement: data.totalPlacement / data.games,
        top4Count: data.top4,
      }))
      .sort((a, b) => b.games - a.games)
      .slice(0, 5);

    const topChampions: ChampionStat[] = Array.from(champMap.entries())
      .map(([characterId, data]) => ({
        characterId,
        games: data.games,
        avgPlacement: data.totalPlacement / data.games,
      }))
      .sort((a, b) => b.games - a.games)
      .slice(0, 5);

    return {
      avgPlacement,
      top4Rate,
      winRate,
      placementDist,
      topTraits,
      topChampions,
      totalGames: placements.length,
    };
  }, [matches, puuid]);

  if (stats.totalGames === 0) return null;

  const maxCount = Math.max(...stats.placementDist.map((d) => d.count), 1);

  return (
    <div className="grid grid-cols-1 gap-3 lg:grid-cols-3">
      {/* Overview */}
      <TerminalCard>
        <h3 className="mb-2 text-[10px] uppercase tracking-wider text-lofi-muted">
          overview // {stats.totalGames} games
        </h3>
        <div className="space-y-1.5">
          <StatRow label="avg placement" value={stats.avgPlacement.toFixed(2)} />
          <StatRow label="top 4 rate" value={`${stats.top4Rate.toFixed(1)}%`} />
          <StatRow label="win rate" value={`${stats.winRate.toFixed(1)}%`} />
        </div>

        <h3 className="mb-2 mt-4 text-[10px] uppercase tracking-wider text-lofi-muted">
          placement distribution
        </h3>
        <div className="space-y-1">
          {stats.placementDist.map((d) => (
            <div key={d.placement} className="flex items-center gap-2">
              <span className="w-4 text-right text-[10px] text-lofi-secondary">
                {d.placement}
              </span>
              <div className="h-3 flex-1 overflow-hidden rounded-sm bg-lofi-black">
                <div
                  className={`h-full rounded-sm ${d.placement <= 4 ? "bg-green-500/40" : "bg-red-500/30"}`}
                  style={{ width: `${(d.count / maxCount) * 100}%` }}
                />
              </div>
              <span className="w-5 text-right text-[10px] text-lofi-muted">
                {d.count}
              </span>
            </div>
          ))}
        </div>
      </TerminalCard>

      {/* Top Traits */}
      <TerminalCard>
        <h3 className="mb-2 text-[10px] uppercase tracking-wider text-lofi-muted">
          top traits
        </h3>
        {stats.topTraits.length === 0 && (
          <p className="text-[10px] text-lofi-muted">no data</p>
        )}
        <div className="space-y-2">
          {stats.topTraits.map((trait) => (
            <div
              key={trait.apiName}
              className="flex items-center justify-between"
            >
              <div>
                <span className="text-xs text-lofi-text">
                  {formatName(trait.apiName)}
                </span>
                <span className="ml-2 text-[10px] text-lofi-muted">
                  {trait.games} games
                </span>
              </div>
              <div className="text-right">
                <span className="text-[10px] text-lofi-secondary">
                  avg {trait.avgPlacement.toFixed(1)}
                </span>
                <span className="ml-2 text-[10px] text-lofi-muted">
                  top4{" "}
                  {((trait.top4Count / trait.games) * 100).toFixed(0)}%
                </span>
              </div>
            </div>
          ))}
        </div>
      </TerminalCard>

      {/* Top Champions */}
      <TerminalCard>
        <h3 className="mb-2 text-[10px] uppercase tracking-wider text-lofi-muted">
          top champions
        </h3>
        {stats.topChampions.length === 0 && (
          <p className="text-[10px] text-lofi-muted">no data</p>
        )}
        <div className="space-y-2">
          {stats.topChampions.map((champ) => (
            <div
              key={champ.characterId}
              className="flex items-center justify-between"
            >
              <div>
                <span className="text-xs text-lofi-text">
                  {formatName(champ.characterId)}
                </span>
                <span className="ml-2 text-[10px] text-lofi-muted">
                  {champ.games} games
                </span>
              </div>
              <span className="text-[10px] text-lofi-secondary">
                avg {champ.avgPlacement.toFixed(1)}
              </span>
            </div>
          ))}
        </div>
      </TerminalCard>
    </div>
  );
}

function StatRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-[10px] text-lofi-secondary">{label}</span>
      <span className="text-xs font-bold text-white">{value}</span>
    </div>
  );
}

function formatName(apiName: string): string {
  return apiName
    .replace(/^Set\d+_/, "")
    .replace(/^TFT\d+_/, "")
    .replace(/^TFT_/, "");
}
