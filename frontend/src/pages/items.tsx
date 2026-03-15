import { useState, useMemo } from "react";
import { usePatchData } from "@/hooks/use-patch-data";
import { SectionHeader } from "@/components/layout/section-header";
import { SearchFilterBar } from "@/components/search-filter-bar";
import { ItemCard } from "@/components/item-card";
import { Skeleton } from "@/components/ui/skeleton";

const typeFilters = [
  { label: "component", value: "component" },
  { label: "crafted", value: "crafted" },
];

export function ItemsPage() {
  const { data, isLoading, error } = usePatchData();
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState("all");

  const filtered = useMemo(() => {
    if (!data?.items) return [];
    return data.items.filter((item) => {
      const matchesSearch = item.name
        .toLowerCase()
        .includes(search.toLowerCase());

      let matchesType = true;
      if (typeFilter === "component") {
        matchesType = item.composition.length === 0;
      } else if (typeFilter === "crafted") {
        matchesType = item.composition.length === 2;
      }

      return matchesSearch && matchesType;
    });
  }, [data?.items, search, typeFilter]);

  if (error) {
    return (
      <div className="text-xs text-red-400">
        error: {error.message}
      </div>
    );
  }

  return (
    <div>
      <SectionHeader number="02" title="items" />

      <SearchFilterBar
        searchValue={search}
        onSearchChange={setSearch}
        searchPlaceholder="search items..."
        filters={typeFilters}
        activeFilter={typeFilter}
        onFilterChange={setTypeFilter}
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
            {filtered.length} item{filtered.length !== 1 ? "s" : ""}
            {data?.set ? ` // set ${data.set.number}` : ""}
          </p>
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
            {filtered.map((item) => (
              <ItemCard key={item.apiName} item={item} />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
