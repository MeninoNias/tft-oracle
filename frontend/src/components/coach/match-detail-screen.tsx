import { useState, useEffect } from "react";
import { TerminalCard } from "@/components/ui/terminal-card";
import { PlacementBadge } from "@/components/player/placement-badge";
import { UnitIcon } from "@/components/player/unit-icon";
import { TraitBadge } from "@/components/player/trait-badge";
import { useGetMatchAnalysis, useAnalyzeMatch } from "@/hooks/use-coach";
import type { Match, Participant } from "@/gen/tft/v1/player_pb";
import type { AnalyzeMatchResponse } from "@/gen/tft/v1/coach_pb";
import type { PatchLookup } from "@/hooks/use-patch-lookup";

interface MatchDetailScreenProps {
  match: Match;
  playerPuuid: string;
  patchLookup: PatchLookup;
  onBack: () => void;
}

const gradeColors: Record<string, string> = {
  "S+": "text-yellow-400 border-yellow-400",
  S: "text-yellow-400 border-yellow-400",
  "A+": "text-green-400 border-green-400",
  A: "text-green-400 border-green-400",
  "B+": "text-blue-400 border-blue-400",
  B: "text-blue-400 border-blue-400",
  "C+": "text-orange-400 border-orange-400",
  C: "text-orange-400 border-orange-400",
  D: "text-red-400 border-red-400",
  F: "text-red-500 border-red-500",
};

const insightGradeColors: Record<string, string> = {
  good: "bg-green-500",
  okay: "bg-yellow-500",
  poor: "bg-red-500",
};

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

export function MatchDetailScreen({
  match,
  playerPuuid,
  patchLookup,
  onBack,
}: MatchDetailScreenProps) {
  const info = match.info;
  const metadata = match.metadata;
  const matchId = metadata?.matchId ?? "";

  const [analysis, setAnalysis] = useState<AnalyzeMatchResponse | null>(null);
  const cachedQuery = useGetMatchAnalysis(matchId, playerPuuid);
  const analyzeMutation = useAnalyzeMatch();

  // Load cached analysis on mount
  useEffect(() => {
    if (cachedQuery.data?.found && cachedQuery.data.analysis) {
      setAnalysis(cachedQuery.data.analysis);
    }
  }, [cachedQuery.data]);

  if (!info || !metadata) return null;

  const player = info.participants.find((p) => p.puuid === playerPuuid);
  if (!player) return null;

  const gameDate = new Date(Number(info.gameDatetime));
  const gameLengthMin = (info.gameLength / 60).toFixed(0);

  const activeTraits = player.traits
    .filter((t) => t.style > 0)
    .sort((a, b) => b.style - a.style);
  const sortedUnits = [...player.units].sort(
    (a, b) => b.rarity - a.rarity || b.tier - a.tier,
  );
  const lobbyParticipants = [...info.participants].sort(
    (a, b) => a.placement - b.placement,
  );

  const handleAnalyze = (force = false) => {
    analyzeMutation.mutate(
      { matchId, puuid: playerPuuid, force },
      { onSuccess: (data) => setAnalysis(data) },
    );
  };

  return (
    <div className="space-y-4">
      {/* ── Back + Header ── */}
      <div className="flex items-center gap-3">
        <button
          onClick={onBack}
          className="text-xs text-lofi-muted transition-colors hover:text-lofi-text"
        >
          &larr; back
        </button>
        <div className="h-px flex-1 bg-lofi-border" />
        <span className="text-[10px] text-lofi-muted">
          {getTimeAgo(gameDate)} // {gameLengthMin}m //
          {info.tftGameType === "standard" ? " ranked" : ` ${info.tftGameType}`}
        </span>
      </div>

      {/* ── Match Summary Bar ── */}
      <div className="flex items-center gap-4">
        <PlacementBadge placement={player.placement} size="md" />
        <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-lofi-secondary">
          <span>
            level <span className="font-medium text-lofi-text">{player.level}</span>
          </span>
          <span>
            gold <span className="font-medium text-lofi-text">{player.goldLeft}</span>
          </span>
          <span>
            round{" "}
            <span className="font-medium text-lofi-text">
              {formatRound(player.lastRound)}
            </span>
          </span>
          <span>
            damage{" "}
            <span className="font-medium text-lofi-text">
              {player.totalDamageToPlayers}
            </span>
          </span>
          {player.playersEliminated > 0 && (
            <span>
              kills{" "}
              <span className="font-medium text-red-400">
                {player.playersEliminated}
              </span>
            </span>
          )}
        </div>

        <div className="flex-1" />

        {/* Analyze / Re-analyze button + Grade */}
        {analysis ? (
          <div className="flex shrink-0 items-center gap-2">
            <button
              onClick={() => handleAnalyze(true)}
              disabled={analyzeMutation.isPending}
              className="rounded-sm px-3 py-1.5 text-[10px] text-lofi-muted transition-colors hover:text-lofi-accent disabled:opacity-50"
            >
              {analyzeMutation.isPending ? (
                <span className="animate-pulse">re-analyzing...</span>
              ) : (
                "re-analyze"
              )}
            </button>
            <div
              className={`flex h-11 w-11 items-center justify-center rounded-sm border-2 text-lg font-bold ${
                gradeColors[analysis.overallGrade] ??
                "text-lofi-muted border-lofi-muted"
              }`}
            >
              {analysis.overallGrade}
            </div>
          </div>
        ) : (
          <button
            onClick={() => handleAnalyze()}
            disabled={analyzeMutation.isPending || cachedQuery.isLoading}
            className="shrink-0 rounded-sm border border-lofi-accent px-4 py-2 text-xs font-medium text-lofi-accent transition-colors hover:bg-lofi-accent hover:text-lofi-bg disabled:opacity-50"
          >
            {analyzeMutation.isPending ? (
              <span className="animate-pulse">analyzing...</span>
            ) : cachedQuery.isLoading ? (
              <span className="animate-pulse">loading...</span>
            ) : (
              "analyze match"
            )}
          </button>
        )}
      </div>

      {analyzeMutation.isError && (
        <p className="text-xs text-red-400">
          {analyzeMutation.error?.message ?? "Analysis failed"}
        </p>
      )}

      {/* ── AI Summary (when analyzed) ── */}
      {analysis && (
        <TerminalCard className="border-lofi-accent/30">
          <p className="mb-1 text-[10px] font-medium uppercase text-lofi-accent">
            ai coach
          </p>
          <p className="text-sm text-lofi-secondary">{analysis.summary}</p>
        </TerminalCard>
      )}

      {/* ── Your Board ── */}
      <TerminalCard className="!p-3">
        <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
          your board
        </p>

        <div className="flex flex-wrap gap-3">
          {sortedUnits.map((unit, i) => {
            const champ = patchLookup.championByApiName.get(unit.characterId);
            return (
              <div key={i} className="flex flex-col items-center">
                <UnitIcon unit={unit} lookup={patchLookup} />
                <span className="mt-0.5 max-w-[52px] truncate text-center text-[9px] text-lofi-muted">
                  {champ?.name ?? unit.characterId.replace(/^TFT\d+_/, "")}
                </span>
              </div>
            );
          })}
        </div>

        {/* Augments */}
        {player.augments.length > 0 && (
          <div className="mt-3">
            <p className="mb-1 text-[10px] uppercase text-lofi-muted">
              augments
            </p>
            <div className="flex flex-wrap gap-1.5">
              {player.augments.map((aug) => {
                const item = patchLookup.itemByApiName.get(aug);
                return (
                  <div
                    key={aug}
                    className="flex items-center gap-1.5 rounded border border-lofi-border bg-lofi-black px-2 py-1"
                  >
                    {item?.iconUrl && (
                      <img
                        src={item.iconUrl}
                        alt={item.name}
                        className="h-4 w-4"
                        loading="lazy"
                      />
                    )}
                    <span className="text-xs text-lofi-secondary">
                      {item?.name ?? aug.replace(/^TFT\d+_/, "")}
                    </span>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Active Traits */}
        <div className="mt-3">
          <p className="mb-1 text-[10px] uppercase text-lofi-muted">
            active traits
          </p>
          <div className="flex flex-wrap gap-1">
            {activeTraits.map((trait) => (
              <TraitBadge
                key={trait.apiName}
                trait={trait}
                lookup={patchLookup}
              />
            ))}
          </div>
        </div>
      </TerminalCard>

      {/* ── Coaching Insights (when analyzed) ── */}
      {analysis && analysis.insights.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            coaching insights
          </p>
          <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3">
            {analysis.insights.map((insight) => (
              <TerminalCard key={insight.category} className="!p-3">
                <div className="flex items-center gap-2">
                  <span
                    className={`h-2 w-2 rounded-full ${insightGradeColors[insight.grade] ?? "bg-lofi-muted"}`}
                  />
                  <span className="text-xs font-medium uppercase text-lofi-muted">
                    {insight.category}
                  </span>
                </div>
                <p className="mt-1 text-sm font-medium text-lofi-text">
                  {insight.title}
                </p>
                <p className="mt-0.5 text-xs text-lofi-secondary">
                  {insight.detail}
                </p>
              </TerminalCard>
            ))}
          </div>
        </div>
      )}

      {/* ── Meta Comparison (when analyzed) ── */}
      {analysis?.metaComparison && (
        <TerminalCard className="!p-3">
          <p className="mb-1 text-xs font-medium uppercase text-lofi-muted">
            meta comparison
          </p>
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-lofi-accent">
              {analysis.metaComparison.closestMetaComp}
            </span>
            <span className="rounded border border-lofi-border px-1.5 py-0.5 text-[10px] font-bold text-lofi-secondary">
              {analysis.metaComparison.metaTier}-tier
            </span>
          </div>
          <p className="mt-1 text-xs text-lofi-secondary">
            {analysis.metaComparison.assessment}
          </p>
          {analysis.metaComparison.missingUnits.length > 0 && (
            <div className="mt-2">
              <span className="text-[10px] uppercase text-lofi-muted">
                missing:{" "}
              </span>
              <span className="text-xs text-red-400">
                {analysis.metaComparison.missingUnits.join(", ")}
              </span>
            </div>
          )}
          {analysis.metaComparison.suboptimalItems.length > 0 && (
            <div>
              <span className="text-[10px] uppercase text-lofi-muted">
                suboptimal items:{" "}
              </span>
              <span className="text-xs text-orange-400">
                {analysis.metaComparison.suboptimalItems.join(", ")}
              </span>
            </div>
          )}
        </TerminalCard>
      )}

      {/* ── Lobby ── */}
      <div>
        <div className="mb-2 flex items-center gap-2">
          <p className="text-xs font-medium uppercase text-lofi-muted">
            lobby
          </p>
          {analysis?.lobbyContext && (
            <span className="text-[10px] text-lofi-muted">
              // {analysis.lobbyContext.contestedCount > 0
                ? `${analysis.lobbyContext.contestedCount} contested`
                : "uncontested"}
            </span>
          )}
        </div>
        {analysis?.lobbyContext && (
          <p className="mb-2 text-xs text-lofi-secondary">
            {analysis.lobbyContext.lobbyStrengthAssessment}
          </p>
        )}
        <TerminalCard className="!p-2">
          <div className="space-y-1">
            {lobbyParticipants.map((p) => (
              <LobbyRow
                key={p.puuid}
                participant={p}
                isPlayer={p.puuid === playerPuuid}
                patchLookup={patchLookup}
              />
            ))}
          </div>
        </TerminalCard>
      </div>

      {/* ── Suggestions (when analyzed) ── */}
      {analysis && analysis.suggestions.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            what to improve next game
          </p>
          <div className="space-y-2">
            {analysis.suggestions.map((s, i) => (
              <div key={i} className="flex items-start gap-3">
                <span
                  className={`mt-0.5 shrink-0 rounded-sm px-1.5 py-0.5 text-[10px] font-bold uppercase ${
                    s.priority === "high"
                      ? "bg-red-500/20 text-red-400"
                      : s.priority === "medium"
                        ? "bg-yellow-500/20 text-yellow-400"
                        : "bg-blue-500/20 text-blue-400"
                  }`}
                >
                  {s.priority}
                </span>
                <div>
                  <p className="text-xs text-lofi-text">{s.description}</p>
                  <span className="text-[10px] text-lofi-muted">
                    {s.category}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

// ── Lobby Row ──

function LobbyRow({
  participant,
  isPlayer,
  patchLookup,
}: {
  participant: Participant;
  isPlayer: boolean;
  patchLookup: PatchLookup;
}) {
  const activeTraits = participant.traits
    .filter((t) => t.style > 0)
    .sort((a, b) => b.style - a.style);
  const sortedUnits = [...participant.units].sort(
    (a, b) => b.rarity - a.rarity || b.tier - a.tier,
  );

  return (
    <div
      className={`flex items-center gap-2 rounded-sm px-2 py-1.5 ${
        isPlayer
          ? "border border-lofi-accent/30 bg-lofi-accent/5"
          : "hover:bg-lofi-surface/50"
      }`}
    >
      <PlacementBadge placement={participant.placement} />

      <div className="w-[60px] shrink-0">
        <p className="text-[10px] text-lofi-muted">
          lv{participant.level} // {participant.goldLeft}g
        </p>
        <p className="text-[10px] text-lofi-muted">
          {participant.totalDamageToPlayers} dmg
        </p>
      </div>

      <div className="flex w-[90px] shrink-0 flex-wrap gap-0.5">
        {activeTraits.slice(0, 4).map((trait) => (
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

      {isPlayer && (
        <span className="shrink-0 text-[9px] font-bold uppercase text-lofi-accent">
          you
        </span>
      )}
    </div>
  );
}
