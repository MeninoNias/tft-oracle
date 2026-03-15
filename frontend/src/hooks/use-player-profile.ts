import { useQuery } from "@tanstack/react-query";
import { playerClient } from "@/lib/transport";

export function usePlayerProfile(
  gameName: string,
  tagLine: string,
  region: string,
) {
  return useQuery({
    queryKey: ["playerProfile", gameName, tagLine, region],
    queryFn: () => playerClient.getPlayerProfile({ gameName, tagLine, region }),
    enabled: !!gameName && !!tagLine,
    staleTime: 10 * 60 * 1000,
    retry: false,
  });
}
