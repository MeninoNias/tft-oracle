import { useDroppable } from "@dnd-kit/core";
import type { PlacedChampion } from "@/stores/board-store";
import type { Champion } from "@/gen/tft/v1/patch_pb";

const costBorderColors: Record<number, string> = {
  1: "border-cost-1",
  2: "border-cost-2",
  3: "border-cost-3",
  4: "border-cost-4",
  5: "border-cost-5",
};

interface BoardCellProps {
  position: number;
  placed?: PlacedChampion;
  championData?: Champion;
  onSelect: (position: number) => void;
  selectedPosition: number | null;
}

export function BoardCell({
  position,
  placed,
  championData,
  onSelect,
  selectedPosition,
}: BoardCellProps) {
  const { isOver, setNodeRef } = useDroppable({ id: `cell-${position}` });

  const isSelected = selectedPosition === position;
  const borderColor = placed && championData
    ? costBorderColors[championData.cost] ?? "border-lofi-border"
    : "border-lofi-border";

  return (
    <div
      ref={setNodeRef}
      onClick={() => placed && onSelect(position)}
      className={`flex h-16 w-16 items-center justify-center rounded-sm border-2 transition-all ${
        isOver
          ? "border-lofi-accent bg-lofi-accent/10"
          : placed
            ? `${borderColor} bg-lofi-surface cursor-pointer`
            : "border-dashed border-lofi-border/50 bg-lofi-black/50"
      } ${isSelected ? "ring-2 ring-lofi-accent" : ""}`}
    >
      {placed && championData ? (
        <div className="relative flex flex-col items-center">
          <img
            src={championData.squareIconUrl}
            alt={championData.name}
            className="h-10 w-10 rounded-sm object-cover"
          />
          {/* Star level indicator */}
          <div className="absolute -bottom-1 flex gap-px">
            {Array.from({ length: placed.starLevel }).map((_, i) => (
              <span
                key={i}
                className="h-1.5 w-1.5 rounded-full bg-yellow-400"
              />
            ))}
          </div>
          {/* Item dots */}
          {placed.itemApiNames.length > 0 && (
            <div className="absolute -top-1 -right-1 flex gap-px">
              {placed.itemApiNames.map((_, i) => (
                <span
                  key={i}
                  className="h-1.5 w-1.5 rounded-full bg-lofi-accent"
                />
              ))}
            </div>
          )}
        </div>
      ) : (
        <span className="text-[8px] text-lofi-muted/30">
          {Math.floor(position / 7)},{position % 7}
        </span>
      )}
    </div>
  );
}
