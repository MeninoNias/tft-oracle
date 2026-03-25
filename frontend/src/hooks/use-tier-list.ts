import { useQuery } from "@tanstack/react-query";
import { tierListClient } from "@/lib/transport";

export function useTierList(patch = "", tierFilter = "") {
  return useQuery({
    queryKey: ["tierList", patch, tierFilter],
    queryFn: () =>
      tierListClient.getConsolidatedTierList({ patch, tierFilter }),
    staleTime: 10 * 60 * 1000, // 10 minutes
    retry: 1,
  });
}
