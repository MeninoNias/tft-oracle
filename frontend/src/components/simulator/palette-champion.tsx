import { useDraggable } from "@dnd-kit/core";
import type { Champion } from "@/gen/tft/v1/patch_pb";

const costColors: Record<number, string> = {
  1: "text-cost-1",
  2: "text-cost-2",
  3: "text-cost-3",
  4: "text-cost-4",
  5: "text-cost-5",
};

interface PaletteChampionProps {
  champion: Champion;
}

export function PaletteChampion({ champion }: PaletteChampionProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } =
    useDraggable({ id: `palette-${champion.apiName}`, data: { champion } });

  const style = transform
    ? { transform: `translate(${transform.x}px, ${transform.y}px)` }
    : undefined;

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className={`flex cursor-grab items-center gap-2 rounded-sm border border-lofi-border bg-lofi-surface px-2 py-1.5 transition-colors hover:border-lofi-muted ${
        isDragging ? "opacity-50" : ""
      }`}
    >
      {champion.squareIconUrl ? (
        <img
          src={champion.squareIconUrl}
          alt={champion.name}
          className="h-7 w-7 rounded-sm object-cover"
          loading="lazy"
        />
      ) : (
        <div className="flex h-7 w-7 items-center justify-center rounded-sm bg-lofi-black text-[8px] text-lofi-muted">
          ?
        </div>
      )}
      <span className="text-[10px] text-lofi-text">{champion.name}</span>
      <span className={`ml-auto text-[10px] font-bold ${costColors[champion.cost] ?? ""}`}>
        {champion.cost}g
      </span>
    </div>
  );
}
