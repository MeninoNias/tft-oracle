import { useState } from "react";
import { usePlayerProfile } from "@/hooks/use-player-profile";
import { useMatchHistory } from "@/hooks/use-match-history";
import { usePatchLookup } from "@/hooks/use-patch-lookup";
import { useAuthStore } from "@/stores/auth-store";
import { useCurrentUser } from "@/hooks/use-auth";
import { SectionHeader } from "@/components/layout/section-header";
import { ProfileHeader } from "@/components/player/profile-header";
import { PlacementBadge } from "@/components/player/placement-badge";
import { UnitIcon } from "@/components/player/unit-icon";
import { TraitBadge } from "@/components/player/trait-badge";
import { MatchDetailScreen } from "@/components/coach/match-detail-screen";
import { TerminalCard } from "@/components/ui/terminal-card";
import { Skeleton } from "@/components/ui/skeleton";
import { AuthGate } from "@/components/auth/auth-gate";
import type { Match } from "@/gen/tft/v1/player_pb";
import type { PatchLookup } from "@/hooks/use-patch-lookup";

export function CoachPage() {
  const { user } = useAuthStore();
  const currentUser = useCurrentUser();
  const [selectedMatch, setSelectedMatch] = useState<Match | null>(null);

  const authUser = currentUser.data?.user ?? user;
  const gameName = authUser?.gameName ?? "";
  const tagLine = authUser?.tagLine ?? "";
  const region = authUser?.region ?? "";
  const puuid = authUser?.puuid ?? "";

  const profile = usePlayerProfile(gameName, tagLine, region);
  const matchHistory = useMatchHistory(puuid, region);
  const patchLookup = usePatchLookup();

  const hasMatches =
    matchHistory.data?.matches && matchHistory.data.matches.length > 0;

  return (
    <div>
      <SectionHeader number="07" title="ai coach" />

      <AuthGate>
        {/* Loading */}
        {currentUser.isLoading && (
          <div className="space-y-3">
            <Skeleton className="h-20" />
            <Skeleton className="h-48" />
          </div>
        )}

        {/* ── Match Detail Screen ── */}
        {selectedMatch && patchLookup ? (
          <MatchDetailScreen
            match={selectedMatch}
            playerPuuid={puuid}
            patchLookup={patchLookup}
            onBack={() => setSelectedMatch(null)}
          />
        ) : (
          <>
            {/* Profile */}
            {profile.data?.account && (
              <div className="mb-4">
                <ProfileHeader
                  account={profile.data.account}
                  summoner={profile.data.summoner}
                  rankedEntries={profile.data.rankedEntries}
                />
              </div>
            )}

            {/* Match History List */}
            {matchHistory.isLoading && !currentUser.isLoading && (
              <div className="space-y-2">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-16 w-full" />
                ))}
              </div>
            )}

            {hasMatches && patchLookup && (
              <div className="space-y-1.5">
                <p className="text-xs font-medium uppercase text-lofi-muted">
                  match history — click to view details
                </p>
                {matchHistory.data!.matches.map((match) => (
                  <MatchRow
                    key={match.metadata?.matchId}
                    match={match}
                    playerPuuid={puuid}
                    patchLookup={patchLookup}
                    onClick={() => setSelectedMatch(match)}
                  />
                ))}
              </div>
            )}
          </>
        )}
      </AuthGate>
    </div>
  );
}

// ── Clickable match row ──

function MatchRow({
  match,
  playerPuuid,
  patchLookup,
  onClick,
}: {
  match: Match;
  playerPuuid: string;
  patchLookup: PatchLookup;
  onClick: () => void;
}) {
  const info = match.info;
  if (!info) return null;

  const player = info.participants.find((p) => p.puuid === playerPuuid);
  if (!player) return null;

  const gameDate = new Date(Number(info.gameDatetime));
  const timeAgo = getTimeAgo(gameDate);
  const gameLengthMin = (info.gameLength / 60).toFixed(0);

  const activeTraits = player.traits
    .filter((t) => t.style > 0)
    .sort((a, b) => b.style - a.style);
  const sortedUnits = [...player.units].sort(
    (a, b) => b.rarity - a.rarity || b.tier - a.tier,
  );

  return (
    <TerminalCard
      className={`cursor-pointer transition-colors hover:border-lofi-muted ${
        player.placement <= 4
          ? "border-l-2 border-l-green-500/40"
          : "border-l-2 border-l-red-500/30"
      }`}
    >
      <div onClick={onClick} className="space-y-2">
        {/* Top row: placement + info + button */}
        <div className="flex items-center gap-4">
          <PlacementBadge placement={player.placement} size="md" className="shrink-0" />

          <div className="min-w-0 flex-1 space-y-0.5">
            <div className="flex items-center gap-3 text-xs text-lofi-muted">
              <span>{timeAgo}</span>
              <span>{gameLengthMin}m</span>
              <span>lv{player.level}</span>
              <span>{player.goldLeft}g</span>
              <span>⚔ {player.totalDamageToPlayers}</span>
              {player.playersEliminated > 0 && (
                <span className="font-medium text-red-400">
                  {player.playersEliminated} kill{player.playersEliminated > 1 ? "s" : ""}
                </span>
              )}
            </div>

            {/* Traits */}
            <div className="flex flex-wrap gap-0.5">
              {activeTraits.slice(0, 8).map((trait) => (
                <TraitBadge
                  key={trait.apiName}
                  trait={trait}
                  lookup={patchLookup}
                />
              ))}
            </div>
          </div>

          <button className="shrink-0 rounded-sm border border-lofi-accent px-4 py-2 text-xs font-medium text-lofi-accent transition-colors hover:bg-lofi-accent hover:text-lofi-bg">
            view details
          </button>
        </div>

        {/* Bottom row: units */}
        <div className="flex items-start gap-1.5 overflow-x-auto pl-[60px]">
          {sortedUnits.map((unit, i) => (
            <UnitIcon key={i} unit={unit} lookup={patchLookup} />
          ))}
        </div>
      </div>
    </TerminalCard>
  );
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
