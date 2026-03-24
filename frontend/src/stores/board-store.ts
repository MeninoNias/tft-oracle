import { create } from "zustand";

export interface PlacedChampion {
  championApiName: string;
  starLevel: 1 | 2 | 3;
  itemApiNames: string[];
}

interface BoardState {
  playerBoard: Record<number, PlacedChampion>;
  opponentBoard: Record<number, PlacedChampion>;
  playerAugments: string[];
  opponentAugments: string[];
  playerLevel: number;
  opponentLevel: number;
  activeTab: "player" | "opponent";

  // Board accessors
  getBoard: () => Record<number, PlacedChampion>;
  getAugments: () => string[];
  getLevel: () => number;

  // Actions
  setActiveTab: (tab: "player" | "opponent") => void;
  placeChampion: (position: number, championApiName: string) => void;
  removeChampion: (position: number) => void;
  moveChampion: (from: number, to: number) => void;
  setStarLevel: (position: number, starLevel: 1 | 2 | 3) => void;
  addItem: (position: number, itemApiName: string) => void;
  removeItem: (position: number, itemIndex: number) => void;
  setAugments: (augments: string[]) => void;
  setLevel: (level: number) => void;
  clearBoard: () => void;
  resetAll: () => void;
}

export const useBoardStore = create<BoardState>()((set, get) => ({
  playerBoard: {},
  opponentBoard: {},
  playerAugments: [],
  opponentAugments: [],
  playerLevel: 8,
  opponentLevel: 8,
  activeTab: "player",

  getBoard: () => {
    const state = get();
    return state.activeTab === "player"
      ? state.playerBoard
      : state.opponentBoard;
  },
  getAugments: () => {
    const state = get();
    return state.activeTab === "player"
      ? state.playerAugments
      : state.opponentAugments;
  },
  getLevel: () => {
    const state = get();
    return state.activeTab === "player"
      ? state.playerLevel
      : state.opponentLevel;
  },

  setActiveTab: (tab) => set({ activeTab: tab }),

  placeChampion: (position, championApiName) =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      return {
        [boardKey]: {
          ...state[boardKey],
          [position]: {
            championApiName,
            starLevel: 1 as const,
            itemApiNames: [],
          },
        },
      };
    }),

  removeChampion: (position) =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      const { [position]: _, ...rest } = state[boardKey];
      return { [boardKey]: rest };
    }),

  moveChampion: (from, to) =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      const board = { ...state[boardKey] };
      const champ = board[from];
      if (!champ) return state;
      delete board[from];
      board[to] = champ;
      return { [boardKey]: board };
    }),

  setStarLevel: (position, starLevel) =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      const champ = state[boardKey][position];
      if (!champ) return state;
      return {
        [boardKey]: {
          ...state[boardKey],
          [position]: { ...champ, starLevel },
        },
      };
    }),

  addItem: (position, itemApiName) =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      const champ = state[boardKey][position];
      if (!champ || champ.itemApiNames.length >= 3) return state;
      return {
        [boardKey]: {
          ...state[boardKey],
          [position]: {
            ...champ,
            itemApiNames: [...champ.itemApiNames, itemApiName],
          },
        },
      };
    }),

  removeItem: (position, itemIndex) =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      const champ = state[boardKey][position];
      if (!champ) return state;
      return {
        [boardKey]: {
          ...state[boardKey],
          [position]: {
            ...champ,
            itemApiNames: champ.itemApiNames.filter((_, i) => i !== itemIndex),
          },
        },
      };
    }),

  setAugments: (augments) =>
    set((state) => {
      const key =
        state.activeTab === "player"
          ? "playerAugments"
          : "opponentAugments";
      return { [key]: augments.slice(0, 3) };
    }),

  setLevel: (level) =>
    set((state) => {
      const key =
        state.activeTab === "player" ? "playerLevel" : "opponentLevel";
      return { [key]: Math.max(1, Math.min(10, level)) };
    }),

  clearBoard: () =>
    set((state) => {
      const boardKey =
        state.activeTab === "player" ? "playerBoard" : "opponentBoard";
      const augKey =
        state.activeTab === "player"
          ? "playerAugments"
          : "opponentAugments";
      return { [boardKey]: {}, [augKey]: [] };
    }),

  resetAll: () =>
    set({
      playerBoard: {},
      opponentBoard: {},
      playerAugments: [],
      opponentAugments: [],
      playerLevel: 8,
      opponentLevel: 8,
      activeTab: "player",
    }),
}));
