import { useQuery } from "@tanstack/react-query";
import { playerClient } from "@/lib/transport";

export function useMatchHistory(
  puuid: string,
  region: string,
  count = 20,
) {
  return useQuery({
    queryKey: ["matchHistory", puuid, region, count],
    queryFn: () => playerClient.getMatchHistory({ puuid, region, count }),
    enabled: !!puuid,
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}
