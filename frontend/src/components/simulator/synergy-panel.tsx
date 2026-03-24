import type { ActiveSynergy } from "@/hooks/use-board-synergies";
import { Badge } from "@/components/ui/badge";

const styleColors: Record<number, string> = {
  1: "border-amber-700 text-amber-600",    // bronze
  2: "border-gray-400 text-gray-300",      // silver
  3: "border-yellow-500 text-yellow-400",  // gold
  4: "border-purple-400 text-purple-300",  // chromatic
};

const styleNames: Record<number, string> = {
  1: "bronze",
  2: "silver",
  3: "gold",
  4: "chromatic",
};

interface SynergyPanelProps {
  synergies: ActiveSynergy[];
}

export function SynergyPanel({ synergies }: SynergyPanelProps) {
  if (synergies.length === 0) {
    return (
      <div className="rounded-sm border border-lofi-border bg-lofi-surface p-3">
        <p className="text-[10px] uppercase tracking-wider text-lofi-muted">
          active synergies
        </p>
        <p className="mt-2 text-[10px] text-lofi-muted">
          no active synergies — place champions to see trait activations
        </p>
      </div>
    );
  }

  return (
    <div className="rounded-sm border border-lofi-border bg-lofi-surface p-3">
      <p className="mb-2 text-[10px] uppercase tracking-wider text-lofi-muted">
        active synergies ({synergies.length})
      </p>
      <div className="space-y-1.5">
        {synergies.map((s) => (
          <div
            key={s.apiName}
            className="flex items-center justify-between gap-2"
          >
            <div className="flex items-center gap-2">
              <Badge className={styleColors[s.style] ?? ""}>
                {styleNames[s.style] ?? "?"}
              </Badge>
              <span className="text-[10px] text-lofi-text">{s.name}</span>
            </div>
            <span className="text-[10px] font-bold text-lofi-secondary">
              {s.count}/{s.threshold}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
