import { useEffect } from "react";
import { usePlayerProfile } from "@/hooks/use-player-profile";
import { useMatchHistory } from "@/hooks/use-match-history";
import { usePatchLookup } from "@/hooks/use-patch-lookup";
import { usePlayerStore } from "@/stores/player-store";
import { SectionHeader } from "@/components/layout/section-header";
import { RiotIdForm } from "@/components/player/riot-id-form";
import { PlayerData } from "@/components/player/player-data";

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

  return (
    <div>
      <SectionHeader number="04" title="player search" />

      <div className="mb-4">
        <RiotIdForm
          onSubmit={(gn, tl, r) => setSearch(gn, tl, r)}
          isLoading={profile.isLoading}
        />
      </div>

      <PlayerData
        profile={profile}
        matchHistory={matchHistory}
        puuid={puuid}
        patchLookup={patchLookup}
        sectionPrefix="04"
      />
    </div>
  );
}
