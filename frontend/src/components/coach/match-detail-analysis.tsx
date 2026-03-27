import { TerminalCard } from "@/components/ui/terminal-card";
import { PlacementBadge } from "@/components/player/placement-badge";
import { UnitIcon } from "@/components/player/unit-icon";
import { TraitBadge } from "@/components/player/trait-badge";
import type { Match, Participant } from "@/gen/tft/v1/player_pb";
import type { AnalyzeMatchResponse } from "@/gen/tft/v1/coach_pb";
import type { PatchLookup } from "@/hooks/use-patch-lookup";

interface MatchDetailAnalysisProps {
  match: Match;
  playerPuuid: string;
  patchLookup: PatchLookup;
  analysis: AnalyzeMatchResponse;
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

export function MatchDetailAnalysis({
  match,
  playerPuuid,
  patchLookup,
  analysis,
}: MatchDetailAnalysisProps) {
  const info = match.info;
  if (!info) return null;

  const player = info.participants.find((p) => p.puuid === playerPuuid);
  if (!player) return null;

  const gradeColor =
    gradeColors[analysis.overallGrade] ?? "text-lofi-muted border-lofi-muted";

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

  return (
    <div className="mt-2 space-y-4 border-l-2 border-lofi-accent/30 pl-4">
      {/* ── Header: Grade + Summary ── */}
      <div className="flex items-start gap-4">
        <div
          className={`flex h-14 w-14 shrink-0 items-center justify-center rounded-sm border-2 text-xl font-bold ${gradeColor}`}
        >
          {analysis.overallGrade}
        </div>
        <div className="min-w-0 flex-1">
          <p className="text-sm text-lofi-secondary">{analysis.summary}</p>
          <div className="mt-1 flex flex-wrap gap-3 text-[10px] text-lofi-muted">
            <span>placement: {player.placement}/8</span>
            <span>level: {player.level}</span>
            <span>gold left: {player.goldLeft}</span>
            <span>round: {formatRound(player.lastRound)}</span>
            <span>damage: {player.totalDamageToPlayers}</span>
            <span>duration: {gameLengthMin}m</span>
          </div>
        </div>
      </div>

      {/* ── Your Board ── */}
      <TerminalCard className="!p-3">
        <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
          your board
        </p>

        {/* Champions with full detail */}
        <div className="flex flex-wrap gap-2">
          {sortedUnits.map((unit, i) => (
            <div key={i} className="flex flex-col items-center">
              <UnitIcon unit={unit} lookup={patchLookup} />
              <span className="mt-0.5 text-[9px] text-lofi-muted">
                {patchLookup.championByApiName.get(unit.characterId)?.name ??
                  unit.characterId.replace(/^TFT\d+_/, "")}
              </span>
            </div>
          ))}
        </div>

        {/* Augments */}
        {player.augments.length > 0 && (
          <div className="mt-3">
            <p className="mb-1 text-[10px] uppercase text-lofi-muted">
              augments
            </p>
            <div className="flex flex-wrap gap-1">
              {player.augments.map((aug) => {
                const item = patchLookup.itemByApiName.get(aug);
                return (
                  <div
                    key={aug}
                    className="flex items-center gap-1 rounded border border-lofi-border bg-lofi-black px-1.5 py-0.5"
                  >
                    {item?.iconUrl && (
                      <img
                        src={item.iconUrl}
                        alt={item.name}
                        className="h-4 w-4"
                        loading="lazy"
                      />
                    )}
                    <span className="text-[10px] text-lofi-secondary">
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

      {/* ── Coaching Insights ── */}
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

      {/* ── Meta Comparison ── */}
      {analysis.metaComparison && (
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
                missing units:{" "}
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

      {/* ── Lobby Overview ── */}
      <div>
        <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
          lobby ({analysis.lobbyContext
            ? `${analysis.lobbyContext.contestedCount} contested`
            : "8 players"})
        </p>
        {analysis.lobbyContext && (
          <p className="mb-2 text-xs text-lofi-secondary">
            {analysis.lobbyContext.lobbyStrengthAssessment}
          </p>
        )}
        <TerminalCard className="!p-3">
          <div className="space-y-1.5">
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

      {/* ── Suggestions ── */}
      {analysis.suggestions.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            what to improve next game
          </p>
          <div className="space-y-2">
            {analysis.suggestions.map((s, i) => (
              <div
                key={i}
                className="flex items-start gap-3 text-sm"
              >
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
