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
import { SettingsPage } from "@/pages/settings";

export function App() {
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
              <Route path="/settings" element={<SettingsPage />} />
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
