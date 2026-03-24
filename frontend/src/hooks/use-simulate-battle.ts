import { useMutation } from "@tanstack/react-query";
import { simulationClient } from "@/lib/transport";
import { useBoardStore } from "@/stores/board-store";

export function useSimulateBattle() {
  return useMutation({
    mutationFn: () => {
      const state = useBoardStore.getState();

      const playerChampions = Object.entries(state.playerBoard).map(
        ([pos, champ]) => ({
          championApiName: champ.championApiName,
          position: parseInt(pos, 10),
          starLevel: champ.starLevel,
          itemApiNames: champ.itemApiNames,
        }),
      );

      const opponentChampions = Object.entries(state.opponentBoard).map(
        ([pos, champ]) => ({
          championApiName: champ.championApiName,
          position: parseInt(pos, 10),
          starLevel: champ.starLevel,
          itemApiNames: champ.itemApiNames,
        }),
      );

      return simulationClient.simulateBattle({
        playerBoard: {
          champions: playerChampions,
          augments: state.playerAugments,
          level: state.playerLevel,
        },
        opponentBoard: {
          champions: opponentChampions,
          augments: state.opponentAugments,
          level: state.opponentLevel,
        },
      });
    },
  });
}
