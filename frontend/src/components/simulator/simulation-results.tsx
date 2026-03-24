import { Skeleton } from "@/components/ui/skeleton";
import { WinProbabilityGauge } from "./win-probability-gauge";
import { AnalysisCard } from "./analysis-card";
import { PositioningTip } from "./positioning-tip";
import { KeyFactors } from "./key-factors";
import { SuggestedChanges } from "./suggested-changes";

interface SimulationResultsProps {
  data?: {
    winProbability: number;
    confidence: number;
    analysis: string;
    positioningTip: string;
    keyFactors: string[];
    suggestedChanges: {
      description: string;
      priority: string;
      category: string;
    }[];
  };
  isPending: boolean;
  error: Error | null;
}

export function SimulationResults({
  data,
  isPending,
  error,
}: SimulationResultsProps) {
  if (isPending) {
    return (
      <div className="space-y-3">
        <Skeleton className="h-24" />
        <Skeleton className="h-16" />
        <Skeleton className="h-12" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
        simulation failed: {error.message}
      </div>
    );
  }

  if (!data) return null;

  return (
    <div className="space-y-3">
      <WinProbabilityGauge
        probability={data.winProbability}
        confidence={data.confidence}
      />
      <AnalysisCard analysis={data.analysis} />
      {data.positioningTip && <PositioningTip tip={data.positioningTip} />}
      <KeyFactors factors={data.keyFactors} />
      <SuggestedChanges changes={data.suggestedChanges} />
    </div>
  );
}
