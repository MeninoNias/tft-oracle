import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User } from "@/gen/tft/v1/auth_pb";

interface AuthState {
  token: string;
  user: User | null;
  setAuth: (token: string, user: User) => void;
  logout: () => void;
  isAuthenticated: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: "",
      user: null,
      setAuth: (token, user) => set({ token, user }),
      logout: () => set({ token: "", user: null }),
      isAuthenticated: () => !!get().token,
    }),
    {
      name: "tft-oracle-auth",
      partialize: (state) => ({ token: state.token }),
    },
  ),
);
