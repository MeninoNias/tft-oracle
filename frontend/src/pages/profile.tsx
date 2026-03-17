import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { usePlayerProfile } from "@/hooks/use-player-profile";
import { useMatchHistory } from "@/hooks/use-match-history";
import { usePatchLookup } from "@/hooks/use-patch-lookup";
import { useAuthStore } from "@/stores/auth-store";
import { useCurrentUser } from "@/hooks/use-auth";
import { SectionHeader } from "@/components/layout/section-header";
import { PlayerData } from "@/components/player/player-data";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/onboarding/empty-state";
import { OnboardingFlow } from "@/components/onboarding/onboarding-flow";

export function ProfilePage() {
  const navigate = useNavigate();
  const { token, user, logout } = useAuthStore();
  const currentUser = useCurrentUser();
  const [showOnboarding, setShowOnboarding] = useState(false);

  // Stale/invalid token — logout and redirect to auth with expired flag
  useEffect(() => {
    if (currentUser.isError && token) {
      logout();
      navigate("/auth?expired=1", { replace: true });
    }
  }, [currentUser.isError, token, logout, navigate]);

  const isAuthenticated = !!token && !currentUser.isError;

  // Use auth user data for auto-search
  const authUser = currentUser.data?.user ?? user;
  const gameName = authUser?.gameName ?? "";
  const tagLine = authUser?.tagLine ?? "";
  const region = authUser?.region ?? "";
  const puuid = authUser?.puuid ?? "";

  const profile = usePlayerProfile(gameName, tagLine, region);
  const matchHistory = useMatchHistory(puuid, region);
  const patchLookup = usePatchLookup();

  // Not authenticated — show onboarding flow
  if (!isAuthenticated && !showOnboarding) {
    return (
      <div>
        <SectionHeader number="03" title="profile" />
        <EmptyState onStartOnboarding={() => setShowOnboarding(true)} />
      </div>
    );
  }

  if (!isAuthenticated && showOnboarding) {
    return (
      <div>
        <SectionHeader number="03" title="profile" />
        <OnboardingFlow onComplete={() => setShowOnboarding(false)} />
      </div>
    );
  }

  // Authenticated but no gameName — show onboarding at identify step
  if (!gameName && !currentUser.isLoading) {
    return (
      <div>
        <SectionHeader number="03" title="profile" />
        <OnboardingFlow onComplete={() => setShowOnboarding(false)} />
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
