import { useMemo } from "react";
import { usePatchData } from "@/hooks/use-patch-data";
import type { Champion, Item, Trait } from "@/gen/tft/v1/patch_pb";

export interface PatchLookup {
  championByApiName: Map<string, Champion>;
  itemByApiName: Map<string, Item>;
  traitByApiName: Map<string, Trait>;
  isLoading: boolean;
}

export function usePatchLookup(): PatchLookup {
  const { data, isLoading } = usePatchData();

  const championByApiName = useMemo(() => {
    const map = new Map<string, Champion>();
    if (data?.champions) {
      for (const c of data.champions) map.set(c.apiName, c);
    }
    return map;
  }, [data?.champions]);

  const itemByApiName = useMemo(() => {
    const map = new Map<string, Item>();
    if (data?.items) {
      for (const i of data.items) map.set(i.apiName, i);
    }
    return map;
  }, [data?.items]);

  const traitByApiName = useMemo(() => {
    const map = new Map<string, Trait>();
    if (data?.traits) {
      for (const t of data.traits) map.set(t.apiName, t);
    }
    return map;
  }, [data?.traits]);

  return { championByApiName, itemByApiName, traitByApiName, isLoading };
}
