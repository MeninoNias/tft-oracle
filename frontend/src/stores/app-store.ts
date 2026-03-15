import { create } from "zustand";

interface AppState {
  currentPage: string;
  setCurrentPage: (page: string) => void;
}

export const useAppStore = create<AppState>((set) => ({
  currentPage: "champions",
  setCurrentPage: (page) => set({ currentPage: page }),
}));
