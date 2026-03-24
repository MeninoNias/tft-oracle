interface WinProbabilityGaugeProps {
  probability: number; // 0.0 to 1.0
  confidence: number;  // 0.0 to 1.0
}

function getColor(p: number): string {
  if (p < 0.35) return "bg-red-500";
  if (p < 0.55) return "bg-yellow-500";
  return "bg-green-500";
}

function getTextColor(p: number): string {
  if (p < 0.35) return "text-red-400";
  if (p < 0.55) return "text-yellow-400";
  return "text-green-400";
}

function getConfidenceLabel(c: number): string {
  if (c >= 0.8) return "high";
  if (c >= 0.5) return "medium";
  return "low";
}

export function WinProbabilityGauge({
  probability,
  confidence,
}: WinProbabilityGaugeProps) {
  const pct = Math.round(probability * 100);

  return (
    <div className="rounded-sm border border-lofi-border bg-lofi-surface p-4">
      <div className="mb-2 flex items-center justify-between">
        <span className="text-[10px] uppercase tracking-wider text-lofi-muted">
          win probability
        </span>
        <span className="text-[10px] text-lofi-muted">
          confidence:{" "}
          <span className="font-bold text-lofi-secondary">
            {getConfidenceLabel(confidence)}
          </span>
        </span>
      </div>

      <div className="mb-2 flex items-baseline gap-2">
        <span className={`text-3xl font-black ${getTextColor(probability)}`}>
          {pct}%
        </span>
      </div>

      <div className="h-2 w-full overflow-hidden rounded-full bg-lofi-black">
        <div
          className={`h-full rounded-full transition-all duration-500 ${getColor(probability)}`}
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  );
}
