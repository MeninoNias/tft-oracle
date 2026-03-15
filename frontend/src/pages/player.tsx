import { useEffect } from "react";
import { usePlayerProfile } from "@/hooks/use-player-profile";
import { useMatchHistory } from "@/hooks/use-match-history";
import { usePatchLookup } from "@/hooks/use-patch-lookup";
import { usePlayerStore } from "@/stores/player-store";
import { SectionHeader } from "@/components/layout/section-header";
import { RiotIdForm } from "@/components/player/riot-id-form";
import { ProfileHeader } from "@/components/player/profile-header";
import { MatchCard } from "@/components/player/match-card";
import { AnalyticsPanel } from "@/components/player/analytics-panel";
import { Skeleton } from "@/components/ui/skeleton";

export function PlayerPage() {
  const { gameName, tagLine, region, puuid, setSearch, setPUUID } =
    usePlayerStore();

  const profile = usePlayerProfile(gameName, tagLine, region);
  const matchHistory = useMatchHistory(puuid, region);
  const patchLookup = usePatchLookup();

  useEffect(() => {
    if (profile.data?.account?.puuid) {
      setPUUID(profile.data.account.puuid);
    }
  }, [profile.data?.account?.puuid, setPUUID]);

  const handleSearch = (gn: string, tl: string, r: string) => {
    setSearch(gn, tl, r);
  };

  return (
    <div>
      <SectionHeader number="03" title="player" />

      <div className="mb-4">
        <RiotIdForm onSubmit={handleSearch} isLoading={profile.isLoading} />
      </div>

      {profile.error && (
        <div className="mb-4 rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
          {profile.error.message}
        </div>
      )}

      {profile.isLoading && (
        <div className="space-y-3">
          <Skeleton className="h-20" />
          <div className="grid grid-cols-1 gap-3 lg:grid-cols-3">
            <Skeleton className="h-48" />
            <Skeleton className="h-48" />
            <Skeleton className="h-48" />
          </div>
        </div>
      )}

      {profile.data?.account && (
        <div className="space-y-4">
          <ProfileHeader
            account={profile.data.account}
            summoner={profile.data.summoner}
            rankedEntries={profile.data.rankedEntries}
          />

          {matchHistory.data?.matches &&
            matchHistory.data.matches.length > 0 && (
              <>
                <SectionHeader number="03.1" title="analytics" />
                <AnalyticsPanel
                  matches={matchHistory.data.matches}
                  puuid={puuid}
                />
              </>
            )}

          <SectionHeader number="03.2" title="match history" />

          {matchHistory.isLoading && (
            <div className="space-y-2">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-20" />
              ))}
            </div>
          )}

          {matchHistory.error && (
            <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
              {matchHistory.error.message}
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
      )}
    </div>
  );
}
