import { useState, useMemo } from "react";
import { usePatchData } from "@/hooks/use-patch-data";
import { SectionHeader } from "@/components/layout/section-header";
import { SearchFilterBar } from "@/components/search-filter-bar";
import { ChampionCard } from "@/components/champion-card";
import { Skeleton } from "@/components/ui/skeleton";

const costFilters = [
  { label: "1g", value: "1" },
  { label: "2g", value: "2" },
  { label: "3g", value: "3" },
  { label: "4g", value: "4" },
  { label: "5g", value: "5" },
];

export function ChampionsPage() {
  const { data, isLoading, error } = usePatchData();
  const [search, setSearch] = useState("");
  const [costFilter, setCostFilter] = useState("all");

  const traitNames = useMemo(() => {
    if (!data?.traits) return {};
    const map: Record<string, string> = {};
    for (const t of data.traits) {
      map[t.apiName] = t.name;
    }
    return map;
  }, [data?.traits]);

  const filtered = useMemo(() => {
    if (!data?.champions) return [];
    return data.champions.filter((c) => {
      const matchesSearch = c.name
        .toLowerCase()
        .includes(search.toLowerCase());
      const matchesCost =
        costFilter === "all" || c.cost === Number(costFilter);
      return matchesSearch && matchesCost;
    });
  }, [data?.champions, search, costFilter]);

  if (error) {
    return (
      <div className="text-xs text-red-400">
        error: {error.message}
      </div>
    );
  }

  return (
    <div>
      <SectionHeader number="01" title="champions" />

      <SearchFilterBar
        searchValue={search}
        onSearchChange={setSearch}
        searchPlaceholder="search champions..."
        filters={costFilters}
        activeFilter={costFilter}
        onFilterChange={setCostFilter}
      />

      {isLoading ? (
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 12 }).map((_, i) => (
            <Skeleton key={i} className="h-20" />
          ))}
        </div>
      ) : (
        <>
          <p className="mb-3 text-[10px] text-lofi-muted">
            {filtered.length} champion{filtered.length !== 1 ? "s" : ""}
            {data?.set ? ` // set ${data.set.number}: ${data.set.name}` : ""}
          </p>
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
            {filtered.map((champion) => (
              <ChampionCard
                key={champion.apiName}
                champion={champion}
                traitNames={traitNames}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
