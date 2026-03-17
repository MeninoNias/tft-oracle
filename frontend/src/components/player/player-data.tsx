import { SectionHeader } from "@/components/layout/section-header";
import { ProfileHeader } from "@/components/player/profile-header";
import { MatchCard } from "@/components/player/match-card";
import { AnalyticsPanel } from "@/components/player/analytics-panel";
import { Skeleton } from "@/components/ui/skeleton";
import { friendlyError } from "@/lib/error";
import type { PatchLookup } from "@/hooks/use-patch-lookup";
import type { GetPlayerProfileResponse } from "@/gen/tft/v1/player_pb";
import type { GetMatchHistoryResponse } from "@/gen/tft/v1/player_pb";

interface PlayerDataProps {
  profile: {
    data?: GetPlayerProfileResponse;
    isLoading: boolean;
    error: Error | null;
  };
  matchHistory: {
    data?: GetMatchHistoryResponse;
    isLoading: boolean;
    error: Error | null;
  };
  puuid: string;
  patchLookup: PatchLookup;
  sectionPrefix?: string;
}

export function PlayerData({
  profile,
  matchHistory,
  puuid,
  patchLookup,
  sectionPrefix = "03",
}: PlayerDataProps) {
  if (profile.error) {
    return (
      <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
        {friendlyError(profile.error)}
      </div>
    );
  }

  if (profile.isLoading) {
    return (
      <div className="space-y-3">
        <Skeleton className="h-20" />
        <div className="grid grid-cols-1 gap-3 lg:grid-cols-3">
          <Skeleton className="h-48" />
          <Skeleton className="h-48" />
          <Skeleton className="h-48" />
        </div>
      </div>
    );
  }

  if (!profile.data?.account) return null;

  return (
    <div className="space-y-4">
      <ProfileHeader
        account={profile.data.account}
        summoner={profile.data.summoner}
        rankedEntries={profile.data.rankedEntries}
      />

      {matchHistory.data?.matches &&
        matchHistory.data.matches.length > 0 && (
          <>
            <SectionHeader
              number={`${sectionPrefix}.1`}
              title="analytics"
            />
            <AnalyticsPanel
              matches={matchHistory.data.matches}
              puuid={puuid}
            />
          </>
        )}

      <SectionHeader number={`${sectionPrefix}.2`} title="match history" />

      {matchHistory.isLoading && (
        <div className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={i} className="h-20" />
          ))}
        </div>
      )}

      {matchHistory.error && (
        <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
          {friendlyError(matchHistory.error)}
        </div>
      )}

      {matchHistory.data?.matches && (
        <>
          <p className="mb-2 text-xs text-lofi-muted">
            {matchHistory.data.matches.length} matches
          </p>
          <div className="space-y-2">
            {matchHistory.data.matches.map((match) => (
              <MatchCard
                key={match.metadata?.matchId}
                match={match}
                playerPuuid={puuid}
                patchLookup={patchLookup}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
