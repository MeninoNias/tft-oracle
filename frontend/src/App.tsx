import { useState, useCallback } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { MainLayout } from "@/components/layout/main-layout";
import { ChampionsPage } from "@/pages/champions";
import { ItemsPage } from "@/pages/items";
import { ProfilePage } from "@/pages/profile";
import { PlayerPage } from "@/pages/player";
import { AuthPage } from "@/pages/auth";
import { NotFoundPage } from "@/pages/not-found";
import { InternalErrorPage } from "@/pages/internal-error";
import { UnauthorizedPage } from "@/pages/unauthorized";
import { SplashScreen } from "@/components/splash-screen";

const SPLASH_KEY = "tft-oracle-splash-seen";

export function App() {
  const [showSplash, setShowSplash] = useState(
    () => !sessionStorage.getItem(SPLASH_KEY),
  );

  const handleSplashComplete = useCallback(() => {
    sessionStorage.setItem(SPLASH_KEY, "1");
    setShowSplash(false);
  }, []);

  if (showSplash) {
    return <SplashScreen onComplete={handleSplashComplete} />;
  }

  return (
    <Routes>
      {/* Auth page — full-screen, no sidebar */}
      <Route path="/auth" element={<AuthPage />} />

      {/* App pages — with sidebar layout */}
      <Route
        path="*"
        element={
          <MainLayout>
            <Routes>
              <Route path="/" element={<Navigate to="/champions" replace />} />
              <Route path="/champions" element={<ChampionsPage />} />
              <Route path="/items" element={<ItemsPage />} />
              <Route path="/profile" element={<ProfilePage />} />
              <Route path="/player" element={<PlayerPage />} />
              <Route path="/error/500" element={<InternalErrorPage />} />
              <Route path="/error/401" element={<UnauthorizedPage />} />
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </MainLayout>
        }
      />
    </Routes>
  );
}
