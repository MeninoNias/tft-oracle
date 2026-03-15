import { create } from "zustand";

interface PlayerState {
  gameName: string;
  tagLine: string;
  region: string;
  puuid: string;
  setSearch: (gameName: string, tagLine: string, region: string) => void;
  setPUUID: (puuid: string) => void;
  clear: () => void;
}

export const usePlayerStore = create<PlayerState>((set) => ({
  gameName: "",
  tagLine: "",
  region: "br",
  puuid: "",
  setSearch: (gameName, tagLine, region) => set({ gameName, tagLine, region }),
  setPUUID: (puuid) => set({ puuid }),
  clear: () => set({ gameName: "", tagLine: "", region: "americas", puuid: "" }),
}));
