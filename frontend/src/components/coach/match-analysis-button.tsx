import { useAnalyzeMatch } from "@/hooks/use-coach";
import type { AnalyzeMatchResponse } from "@/gen/tft/v1/coach_pb";

interface MatchAnalysisButtonProps {
  matchId: string;
  puuid: string;
  onAnalysis: (data: AnalyzeMatchResponse) => void;
}

export function MatchAnalysisButton({
  matchId,
  puuid,
  onAnalysis,
}: MatchAnalysisButtonProps) {
  const { mutate, isPending } = useAnalyzeMatch();

  return (
    <button
      onClick={(e) => {
        e.stopPropagation();
        mutate(
          { matchId, puuid },
          { onSuccess: (data) => onAnalysis(data) },
        );
      }}
      disabled={isPending}
      className="shrink-0 rounded-sm border border-lofi-border px-2 py-0.5 text-xs text-lofi-muted transition-colors hover:border-lofi-accent hover:text-lofi-accent disabled:opacity-50"
      title="AI Coach Analysis"
    >
      {isPending ? (
        <span className="animate-pulse">analyzing...</span>
      ) : (
        "analyze"
      )}
    </button>
  );
}
