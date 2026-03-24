import { ConnectError, Code } from "@connectrpc/connect";
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
      <div className="rounded-sm border border-lofi-accent/30 bg-lofi-surface p-4">
        <div className="mb-2 flex items-center gap-2">
          <span className="flex h-6 w-6 items-center justify-center rounded-full bg-red-500/20 text-xs text-red-400">
            !
          </span>
          <h3 className="text-xs font-bold text-white">
            simulation unavailable
          </h3>
        </div>
        <p className="text-xs leading-relaxed text-lofi-secondary">
          {getSimulationErrorMessage(error)}
        </p>
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

function getSimulationErrorMessage(error: Error): string {
  if (error instanceof ConnectError) {
    switch (error.code) {
      case Code.Unimplemented:
        return "the simulation service is not available on this server. make sure the backend is running with the SimulationService enabled.";
      case Code.Unavailable:
        return "the AI engine is offline — OPENAI_API_KEY may not be configured on the server.";
      case Code.InvalidArgument:
        return "invalid board state. place at least one champion before simulating.";
      case Code.Internal:
        return "the AI analysis failed. this may be a temporary issue — try again in a moment.";
      case Code.ResourceExhausted:
        return "too many simulation requests. please wait a moment and try again.";
    }
  }

  // Fallback for network errors
  if (error.message.includes("fetch") || error.message.includes("network")) {
    return "cannot reach the server. make sure the backend is running on localhost:8080.";
  }

  return error.message;
}
