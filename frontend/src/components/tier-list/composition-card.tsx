import { TerminalCard } from "@/components/ui/terminal-card";
import { Badge } from "@/components/ui/badge";
import { ConfidenceIndicator } from "./confidence-indicator";
import { SourceBreakdown } from "./source-breakdown";
import type { TierListEntry } from "@/gen/tft/v1/tierlist_pb";
import type { Champion, Item } from "@/gen/tft/v1/patch_pb";

const consensusStyles: Record<string, string> = {
  unanimous: "border-green-500/30 bg-green-500/10 text-green-400",
  majority: "border-yellow-500/30 bg-yellow-500/10 text-yellow-400",
  split: "border-red-500/30 bg-red-500/10 text-red-400",
  single_source: "border-lofi-border bg-lofi-surface text-lofi-muted",
};

const placementColor = (p: number) =>
  p > 0 && p <= 4 ? "text-green-400" : "text-red-400";

interface CompositionCardProps {
  entry: TierListEntry;
  championMap: Map<string, Champion>;
  itemMap: Map<string, Item>;
}

export function CompositionCard({
  entry,
  championMap,
  itemMap,
}: CompositionCardProps) {
  return (
    <TerminalCard>
      {/* Header */}
      <div className="mb-3 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Badge className="font-bold">{entry.consolidatedTier}</Badge>
          <span className="text-xs font-medium text-white">
            {entry.compositionName}
          </span>
        </div>
        <ConfidenceIndicator level={entry.confidence} />
      </div>

      {/* Champion portraits */}
      <div className="mb-3 flex flex-wrap gap-1.5">
        {entry.coreChampions.map((apiName) => {
          const champ = championMap.get(apiName);
          return (
            <div key={apiName} className="flex flex-col items-center gap-0.5">
              {champ?.squareIconUrl ? (
                <img
                  src={champ.squareIconUrl}
                  alt={champ.name}
                  className="h-8 w-8 rounded-sm border border-lofi-border"
                />
              ) : (
                <div className="flex h-8 w-8 items-center justify-center rounded-sm border border-lofi-border bg-lofi-black text-[8px] text-lofi-muted">
                  {apiName.split("_").pop()?.slice(0, 3)}
                </div>
              )}
              {/* Items for this champion */}
              <div className="flex gap-0.5">
                {(entry.recommendedItems[apiName]?.itemApiNames ?? []).map(
                  (itemApi, idx) => {
                    const item = itemMap.get(itemApi);
                    return item?.iconUrl ? (
                      <img
                        key={idx}
                        src={item.iconUrl}
                        alt={item.name}
                        className="h-3.5 w-3.5 rounded-sm"
                      />
                    ) : (
                      <div
                        key={idx}
                        className="h-3.5 w-3.5 rounded-sm bg-lofi-border"
                      />
                    );
                  },
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Stats row */}
      <div className="mb-2 flex items-center gap-4 text-[10px]">
        {entry.avgWinRate > 0 && (
          <span className="text-lofi-secondary">
            win{" "}
            <span className="text-green-400">
              {entry.avgWinRate.toFixed(1)}%
            </span>
          </span>
        )}
        {entry.avgPlayRate > 0 && (
          <span className="text-lofi-secondary">
            play{" "}
            <span className="text-lofi-text">
              {entry.avgPlayRate.toFixed(1)}%
            </span>
          </span>
        )}
        {entry.avgPlacement > 0 && (
          <span className="text-lofi-secondary">
            avg{" "}
            <span className={placementColor(entry.avgPlacement)}>
              {entry.avgPlacement.toFixed(2)}
            </span>
          </span>
        )}
      </div>

      {/* Consensus badge */}
      <div className="flex items-center gap-2">
        <span
          className={`inline-flex items-center rounded-sm border px-1.5 py-0.5 text-[10px] ${
            consensusStyles[entry.consensus] ?? consensusStyles.split
          }`}
        >
          {entry.consensus}
        </span>
      </div>

      {/* Source breakdown */}
      <SourceBreakdown
        metatftTier={entry.metatftTier}
        tftacticsTier={entry.tftacticsTier}
        mobalyticsTier={entry.mobalyticsTier}
      />
    </TerminalCard>
  );
}
