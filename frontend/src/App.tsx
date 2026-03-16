import { Routes, Route, Navigate } from "react-router-dom";
import { MainLayout } from "@/components/layout/main-layout";
import { ChampionsPage } from "@/pages/champions";
import { ItemsPage } from "@/pages/items";
import { PlayerPage } from "@/pages/player";
import { AuthPage } from "@/pages/auth";

export function App() {
  return (
    <MainLayout>
      <Routes>
        <Route path="/" element={<Navigate to="/champions" replace />} />
        <Route path="/champions" element={<ChampionsPage />} />
        <Route path="/items" element={<ItemsPage />} />
        <Route path="/player" element={<PlayerPage />} />
        <Route path="/auth" element={<AuthPage />} />
      </Routes>
    </MainLayout>
  );
}
