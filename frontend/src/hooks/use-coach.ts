import { useMutation, useQuery } from "@tanstack/react-query";
import { coachClient } from "@/lib/transport";

export function useGetMatchAnalysis(matchId: string, puuid: string) {
  return useQuery({
    queryKey: ["coachMatchAnalysis", matchId, puuid],
    queryFn: () => coachClient.getMatchAnalysis({ matchId, puuid }),
    enabled: !!matchId && !!puuid,
    staleTime: Infinity,
  });
}

export function useAnalyzeMatch() {
  return useMutation({
    mutationFn: (params: { matchId: string; puuid: string; force?: boolean }) =>
      coachClient.analyzeMatch(params),
  });
}

export function useAnalyzeHistory() {
  return useMutation({
    mutationFn: (params: { puuid: string; matchCount?: number }) =>
      coachClient.analyzeHistory(params),
  });
}
