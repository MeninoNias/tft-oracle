import { useState, useMemo } from "react";
import type { Item } from "@/gen/tft/v1/patch_pb";
import type { PlacedChampion } from "@/stores/board-store";
import { useBoardStore } from "@/stores/board-store";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface ItemSelectorProps {
  position: number;
  placed: PlacedChampion;
  items: Item[];
  championName: string;
  onClose: () => void;
}

export function ItemSelector({
  position,
  placed,
  items,
  championName,
  onClose,
}: ItemSelectorProps) {
  const [search, setSearch] = useState("");
  const { addItem, removeItem, setStarLevel } = useBoardStore();

  const filtered = useMemo(() => {
    return items.filter((i) =>
      i.name.toLowerCase().includes(search.toLowerCase()),
    );
  }, [items, search]);

  const itemMap = new Map(items.map((i) => [i.apiName, i]));

  return (
    <div className="rounded-sm border border-lofi-accent/30 bg-lofi-surface p-4">
      <div className="mb-3 flex items-center justify-between">
        <div>
          <p className="text-xs font-medium text-white">{championName}</p>
          <p className="text-[10px] text-lofi-muted">
            star level + items (max 3)
          </p>
        </div>
        <Button variant="secondary" onClick={onClose}>
          close
        </Button>
      </div>

      {/* Star level */}
      <div className="mb-3 flex items-center gap-2">
        <span className="text-[10px] text-lofi-muted">stars:</span>
        {([1, 2, 3] as const).map((star) => (
          <button
            key={star}
            onClick={() => setStarLevel(position, star)}
            className={`rounded-sm border px-2 py-0.5 text-[10px] transition-colors ${
              placed.starLevel === star
                ? "border-yellow-500 bg-yellow-500/10 text-yellow-400"
                : "border-lofi-border text-lofi-muted hover:text-lofi-secondary"
            }`}
          >
            {"★".repeat(star)}
          </button>
        ))}
      </div>

      {/* Equipped items */}
      {placed.itemApiNames.length > 0 && (
        <div className="mb-3">
          <p className="mb-1 text-[10px] text-lofi-muted">equipped:</p>
          <div className="flex flex-wrap gap-1">
            {placed.itemApiNames.map((apiName, i) => {
              const item = itemMap.get(apiName);
              return (
                <button
                  key={i}
                  onClick={() => removeItem(position, i)}
                  className="flex items-center gap-1 rounded-sm border border-lofi-border bg-lofi-black px-2 py-1 text-[10px] text-lofi-text hover:border-red-500/50 hover:text-red-400"
                >
                  {item?.name ?? apiName}
                  <span className="text-lofi-muted">×</span>
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* Item search + grid */}
      {placed.itemApiNames.length < 3 && (
        <>
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="search items..."
            className="mb-2"
          />
          <div className="grid max-h-32 grid-cols-2 gap-1 overflow-y-auto">
            {filtered.slice(0, 30).map((item) => (
              <button
                key={item.apiName}
                onClick={() => {
                  addItem(position, item.apiName);
                  setSearch("");
                }}
                className="flex items-center gap-1.5 rounded-sm border border-lofi-border px-2 py-1 text-left text-[10px] text-lofi-text transition-colors hover:border-lofi-muted hover:bg-lofi-surface"
              >
                {item.iconUrl && (
                  <img
                    src={item.iconUrl}
                    alt={item.name}
                    className="h-5 w-5 rounded-sm object-cover"
                  />
                )}
                {item.name}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
