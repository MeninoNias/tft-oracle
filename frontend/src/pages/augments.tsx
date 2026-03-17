import { useState, useMemo } from "react";
import { usePatchData } from "@/hooks/use-patch-data";
import { SectionHeader } from "@/components/layout/section-header";
import { SearchFilterBar } from "@/components/search-filter-bar";
import { AugmentCard } from "@/components/augment-card";
import { Skeleton } from "@/components/ui/skeleton";

const tierFilters = [
  { label: "silver", value: "silver" },
  { label: "gold", value: "gold" },
  { label: "prismatic", value: "prismatic" },
];

function getAugmentTier(tags: string[]): string | null {
  for (const tag of tags) {
    const lower = tag.toLowerCase();
    if (lower.includes("prismatic")) return "prismatic";
    if (lower.includes("gold")) return "gold";
    if (lower.includes("silver")) return "silver";
  }
  return null;
}

export function AugmentsPage() {
  const { data, isLoading, error } = usePatchData();
  const [search, setSearch] = useState("");
  const [tierFilter, setTierFilter] = useState("all");

  const filtered = useMemo(() => {
    if (!data?.augments) return [];
    return data.augments.filter((aug) => {
      const matchesSearch = aug.name
        .toLowerCase()
        .includes(search.toLowerCase());

      let matchesTier = true;
      if (tierFilter !== "all") {
        matchesTier = getAugmentTier(aug.tags) === tierFilter;
      }

      return matchesSearch && matchesTier;
    });
  }, [data?.augments, search, tierFilter]);

  if (error) {
    return (
      <div className="text-xs text-red-400">
        error: {error.message}
      </div>
    );
  }

  return (
    <div>
      <SectionHeader number="05" title="augments" />

      <SearchFilterBar
        searchValue={search}
        onSearchChange={setSearch}
        searchPlaceholder="search augments..."
        filters={tierFilters}
        activeFilter={tierFilter}
        onFilterChange={setTierFilter}
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
            {filtered.length} augment{filtered.length !== 1 ? "s" : ""}
            {data?.set ? ` // set ${data.set.number}` : ""}
          </p>
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
            {filtered.map((aug) => (
              <AugmentCard key={aug.apiName} augment={aug} />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
