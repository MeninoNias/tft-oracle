import { useBoardStore } from "@/stores/board-store";
import type { Champion } from "@/gen/tft/v1/patch_pb";
import { BoardCell } from "./board-cell";

const ROWS = 4;
const COLS = 7;

interface BoardGridProps {
  champions: Champion[];
  selectedPosition: number | null;
  onSelectPosition: (position: number) => void;
}

export function BoardGrid({
  champions,
  selectedPosition,
  onSelectPosition,
}: BoardGridProps) {
  const board = useBoardStore((s) => s.getBoard());
  const champMap = new Map(champions.map((c) => [c.apiName, c]));

  return (
    <div className="flex flex-col items-center gap-1">
      <p className="mb-1 text-[10px] uppercase tracking-wider text-lofi-muted">
        front row
      </p>
      {Array.from({ length: ROWS })
        .map((_, row) => ROWS - 1 - row) // render front row (3) at top
        .map((row) => (
          <div
            key={row}
            className={`flex gap-1 ${row % 2 === 1 ? "ml-[34px]" : ""}`}
          >
            {Array.from({ length: COLS }).map((_, col) => {
              const position = row * COLS + col;
              const placed = board[position];
              return (
                <BoardCell
                  key={position}
                  position={position}
                  placed={placed}
                  championData={
                    placed ? champMap.get(placed.championApiName) : undefined
                  }
                  onSelect={onSelectPosition}
                  selectedPosition={selectedPosition}
                />
              );
            })}
          </div>
        ))}
      <p className="mt-1 text-[10px] uppercase tracking-wider text-lofi-muted">
        back row
      </p>
    </div>
  );
}
