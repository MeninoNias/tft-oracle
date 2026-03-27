import { useState, useMemo } from "react";
import { useTierList } from "@/hooks/use-tier-list";
import { usePatchData } from "@/hooks/use-patch-data";
import { SectionHeader } from "@/components/layout/section-header";
import { SearchFilterBar } from "@/components/search-filter-bar";
import { TierSection } from "@/components/tier-list/tier-section";
import { TierListSkeleton } from "@/components/tier-list/tier-list-skeleton";
import { CrawlStatusFooter } from "@/components/tier-list/crawl-status-footer";
import { InlineError } from "@/components/ui/inline-error";
import { GenericEmptyState } from "@/components/ui/generic-empty-state";
import { friendlyError } from "@/lib/error";
import type { TierListEntry } from "@/gen/tft/v1/tierlist_pb";
import type { Champion, Item } from "@/gen/tft/v1/patch_pb";

const tierFilters = [
  { label: "S", value: "S" },
  { label: "A", value: "A" },
  { label: "B", value: "B" },
  { label: "C", value: "C" },
];

const TIER_ORDER = ["S", "A", "B", "C", "D"];

export function TierListPage() {
  const { data, isLoading, error, refetch } = useTierList();
  const { data: patchData } = usePatchData();
  const [search, setSearch] = useState("");
  const [tierFilter, setTierFilter] = useState("all");

  // Build lookup maps from patch data
  const championMap = useMemo(() => {
    const map = new Map<string, Champion>();
    for (const c of patchData?.champions ?? []) {
      map.set(c.apiName, c);
    }
    return map;
  }, [patchData?.champions]);

  const itemMap = useMemo(() => {
    const map = new Map<string, Item>();
    for (const i of patchData?.items ?? []) {
      map.set(i.apiName, i);
    }
    return map;
  }, [patchData?.items]);

  // Filter entries
  const filtered = useMemo(() => {
    const entries = data?.entries ?? [];
    return entries.filter((e) => {
      const matchesTier =
        tierFilter === "all" || e.consolidatedTier === tierFilter;
      const matchesSearch =
        !search ||
        e.compositionName.toLowerCase().includes(search.toLowerCase()) ||
        e.coreChampions.some((id) => {
          const champ = championMap.get(id);
          return champ?.name.toLowerCase().includes(search.toLowerCase());
        });
      return matchesTier && matchesSearch;
    });
  }, [data?.entries, tierFilter, search, championMap]);

  // Group by tier
  const groupedByTier = useMemo(() => {
    const groups = new Map<string, TierListEntry[]>();
    for (const tier of TIER_ORDER) {
      groups.set(tier, []);
    }
    for (const entry of filtered) {
      const list = groups.get(entry.consolidatedTier);
      if (list) {
        list.push(entry);
      }
    }
    return groups;
  }, [filtered]);

  if (error) {
    return (
      <div>
        <SectionHeader number="08" title="tier list" />
        <InlineError message={friendlyError(error)} onRetry={() => refetch()} />
      </div>
    );
  }

  return (
    <div>
      <SectionHeader number="08" title="tier list" />

      <SearchFilterBar
        searchValue={search}
        onSearchChange={setSearch}
        searchPlaceholder="search compositions or champions..."
        filters={tierFilters}
        activeFilter={tierFilter}
        onFilterChange={setTierFilter}
      />

      {isLoading ? (
        <TierListSkeleton />
      ) : filtered.length === 0 ? (
        <GenericEmptyState
          title="NO_DATA"
          description="no compositions found. the crawler may not have run yet, or no results match your filters."
          actionLabel="clear filters"
          onAction={() => {
            setSearch("");
            setTierFilter("all");
          }}
        />
      ) : (
        <>
          <p className="mb-4 text-[10px] text-lofi-muted">
            {filtered.length} composition{filtered.length !== 1 ? "s" : ""}
            {data?.patch ? ` // patch ${data.patch}` : ""}
          </p>

          {TIER_ORDER.map((tier) => {
            const entries = groupedByTier.get(tier) ?? [];
            return (
              <TierSection
                key={tier}
                tier={tier}
                entries={entries}
                championMap={championMap}
                itemMap={itemMap}
              />
            );
          })}

          <CrawlStatusFooter lastUpdated={data?.lastUpdated ?? ""} />
        </>
      )}
    </div>
  );
}
