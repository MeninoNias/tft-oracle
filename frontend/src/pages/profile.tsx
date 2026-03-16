import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { usePlayerProfile } from "@/hooks/use-player-profile";
import { useMatchHistory } from "@/hooks/use-match-history";
import { usePatchLookup } from "@/hooks/use-patch-lookup";
import { useAuthStore } from "@/stores/auth-store";
import { useCurrentUser } from "@/hooks/use-auth";
import { SectionHeader } from "@/components/layout/section-header";
import { PlayerData } from "@/components/player/player-data";
import { Skeleton } from "@/components/ui/skeleton";

export function ProfilePage() {
  const navigate = useNavigate();
  const { token, user } = useAuthStore();
  const currentUser = useCurrentUser();

  // Use auth user data for auto-search
  const authUser = currentUser.data?.user ?? user;
  const gameName = authUser?.gameName ?? "";
  const tagLine = authUser?.tagLine ?? "";
  const region = authUser?.region ?? "";
  const puuid = authUser?.puuid ?? "";

  const profile = usePlayerProfile(gameName, tagLine, region);
  const matchHistory = useMatchHistory(puuid, region);
  const patchLookup = usePatchLookup();

  // Redirect to auth if not logged in
  useEffect(() => {
    if (!token) {
      navigate("/auth");
    }
  }, [token, navigate]);

  if (!token) return null;

  if (!authUser || (!gameName && !currentUser.isLoading)) {
    return (
      <div>
        <SectionHeader number="03" title="profile" />
        <div className="rounded-sm border border-lofi-border bg-lofi-surface p-6 text-center">
          <p className="text-xs text-lofi-secondary">
            no riot id linked to your account.
          </p>
          <p className="mt-1 text-[10px] text-lofi-muted">
            use the player search to look up any player, or re-register with
            your riot id.
          </p>
        </div>
      </div>
    );
  }

  if (currentUser.isLoading) {
    return (
      <div>
        <SectionHeader number="03" title="profile" />
        <div className="space-y-3">
          <Skeleton className="h-20" />
          <Skeleton className="h-48" />
        </div>
      </div>
    );
  }

  return (
    <div>
      <SectionHeader number="03" title="profile" />

      <PlayerData
        profile={profile}
        matchHistory={matchHistory}
        puuid={puuid}
        patchLookup={patchLookup}
        sectionPrefix="03"
      />
    </div>
  );
}
