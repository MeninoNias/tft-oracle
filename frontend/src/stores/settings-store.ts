import { create } from "zustand";
import { persist } from "zustand/middleware";

interface SettingsState {
  theme: "dark" | "light";
  language: string;
  alerts: {
    metaShift: boolean;
    scouting: boolean;
    matchReady: boolean;
  };
  setTheme: (theme: "dark" | "light") => void;
  setLanguage: (language: string) => void;
  setAlert: (key: keyof SettingsState["alerts"], value: boolean) => void;
}

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      theme: "dark",
      language: "en",
      alerts: {
        metaShift: true,
        scouting: false,
        matchReady: false,
      },
      setTheme: (theme) => set({ theme }),
      setLanguage: (language) => set({ language }),
      setAlert: (key, value) =>
        set((state) => ({
          alerts: { ...state.alerts, [key]: value },
        })),
    }),
    {
      name: "tft-oracle-settings",
    },
  ),
);
