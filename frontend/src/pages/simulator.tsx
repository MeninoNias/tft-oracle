import { useState, useCallback } from "react";
import {
  DndContext,
  type DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { usePatchData } from "@/hooks/use-patch-data";
import { useBoardStore } from "@/stores/board-store";
import { useBoardSynergies } from "@/hooks/use-board-synergies";
import { SectionHeader } from "@/components/layout/section-header";
import { Skeleton } from "@/components/ui/skeleton";
import { BoardGrid } from "@/components/simulator/board-grid";
import { BoardControls } from "@/components/simulator/board-controls";
import { ChampionPalette } from "@/components/simulator/champion-palette";
import { SynergyPanel } from "@/components/simulator/synergy-panel";
import { ItemSelector } from "@/components/simulator/item-selector";
import { SimulationResults } from "@/components/simulator/simulation-results";
import { useSimulateBattle } from "@/hooks/use-simulate-battle";
import { Button } from "@/components/ui/button";

export function SimulatorPage() {
  const { data, isLoading, error } = usePatchData();
  const board = useBoardStore((s) => s.getBoard());
  const placeChampion = useBoardStore((s) => s.placeChampion);
  const removeChampion = useBoardStore((s) => s.removeChampion);
  const [selectedPosition, setSelectedPosition] = useState<number | null>(null);
  const simulation = useSimulateBattle();

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  );

  const synergies = useBoardSynergies(
    board,
    data?.champions ?? [],
    data?.traits ?? [],
  );

  const champMap = new Map(
    (data?.champions ?? []).map((c) => [c.apiName, c]),
  );

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { over, active } = event;
      if (!over) return;

      const cellId = String(over.id);
      if (!cellId.startsWith("cell-")) return;

      const position = parseInt(cellId.replace("cell-", ""), 10);
      const champion = active.data.current?.champion;
      if (champion) {
        placeChampion(position, champion.apiName);
      }
    },
    [placeChampion],
  );

  const handleSelectPosition = useCallback(
    (position: number) => {
      setSelectedPosition((prev) => (prev === position ? null : position));
    },
    [],
  );

  const handleRemoveChampion = useCallback(() => {
    if (selectedPosition !== null) {
      removeChampion(selectedPosition);
      setSelectedPosition(null);
    }
  }, [selectedPosition, removeChampion]);

  if (error) {
    return (
      <div className="text-xs text-red-400">error: {error.message}</div>
    );
  }

  return (
    <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
      <div>
        <SectionHeader number="07" title="simulator" />

        <BoardControls />

        {isLoading ? (
          <div className="space-y-3">
            <Skeleton className="h-72" />
            <Skeleton className="h-48" />
          </div>
        ) : (
          <>
            <div className="mb-4 flex gap-4">
              {/* Board */}
              <div className="flex-1">
                <BoardGrid
                  champions={data?.champions ?? []}
                  selectedPosition={selectedPosition}
                  onSelectPosition={handleSelectPosition}
                />

                {/* Remove button when champion selected */}
                {selectedPosition !== null && board[selectedPosition] && (
                  <div className="mt-2 text-center">
                    <button
                      onClick={handleRemoveChampion}
                      className="text-[10px] text-red-400 hover:text-red-300"
                    >
                      remove champion
                    </button>
                  </div>
                )}
              </div>

              {/* Synergy panel */}
              <div className="w-48 shrink-0">
                <SynergyPanel synergies={synergies} />
              </div>
            </div>

            {/* Item selector for selected champion */}
            {selectedPosition !== null && board[selectedPosition] && (
              <div className="mb-4">
                <ItemSelector
                  position={selectedPosition}
                  placed={board[selectedPosition]}
                  items={data?.items ?? []}
                  championName={
                    champMap.get(board[selectedPosition].championApiName)
                      ?.name ?? board[selectedPosition].championApiName
                  }
                  onClose={() => setSelectedPosition(null)}
                />
              </div>
            )}

            {/* Champion palette */}
            <div className="rounded-sm border border-lofi-border bg-lofi-surface p-3">
              <p className="mb-2 text-[10px] uppercase tracking-wider text-lofi-muted">
                champion palette — drag to board
              </p>
              <ChampionPalette champions={data?.champions ?? []} />
            </div>

            {/* Simulate button */}
            <div className="mt-4">
              <Button
                onClick={() => simulation.mutate()}
                disabled={
                  Object.keys(board).length === 0 || simulation.isPending
                }
                className="w-full"
              >
                {simulation.isPending
                  ? "analyzing battle..."
                  : "simulate battle"}
              </Button>
            </div>

            {/* Simulation results */}
            {(simulation.data || simulation.isPending || simulation.error) && (
              <div className="mt-4">
                <SimulationResults
                  data={simulation.data}
                  isPending={simulation.isPending}
                  error={simulation.error}
                />
              </div>
            )}

            {/* Board stats */}
            <div className="mt-3 flex items-center gap-4 text-[10px] text-lofi-muted">
              <span>
                {Object.keys(board).length} champion
                {Object.keys(board).length !== 1 ? "s" : ""} placed
              </span>
              <span>{synergies.length} active synergies</span>
              {data?.set && <span>set {data.set.number}</span>}
            </div>
          </>
        )}
      </div>
    </DndContext>
  );
}
