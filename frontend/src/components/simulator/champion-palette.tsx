import { useState, useMemo } from "react";
import type { Champion } from "@/gen/tft/v1/patch_pb";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { PaletteChampion } from "./palette-champion";

const costFilters = [1, 2, 3, 4, 5];

interface ChampionPaletteProps {
  champions: Champion[];
}

export function ChampionPalette({ champions }: ChampionPaletteProps) {
  const [search, setSearch] = useState("");
  const [costFilter, setCostFilter] = useState<number | null>(null);

  const filtered = useMemo(() => {
    return champions.filter((c) => {
      const matchesSearch = c.name
        .toLowerCase()
        .includes(search.toLowerCase());
      const matchesCost = costFilter === null || c.cost === costFilter;
      return matchesSearch && matchesCost;
    });
  }, [champions, search, costFilter]);

  return (
    <div>
      <div className="mb-2 flex flex-wrap items-center gap-2">
        <Input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="search champions..."
          className="max-w-[200px]"
        />
        <div className="flex gap-1">
          <Button
            variant={costFilter === null ? "primary" : "secondary"}
            onClick={() => setCostFilter(null)}
          >
            all
          </Button>
          {costFilters.map((cost) => (
            <Button
              key={cost}
              variant={costFilter === cost ? "primary" : "secondary"}
              onClick={() => setCostFilter(cost)}
            >
              {cost}g
            </Button>
          ))}
        </div>
      </div>

      <div className="grid max-h-48 grid-cols-2 gap-1 overflow-y-auto pr-1 lg:grid-cols-3">
        {filtered.map((c) => (
          <PaletteChampion key={c.apiName} champion={c} />
        ))}
      </div>

      {filtered.length === 0 && (
        <p className="py-4 text-center text-[10px] text-lofi-muted">
          no champions found
        </p>
      )}
    </div>
  );
}
