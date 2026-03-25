import type { TierListEntry } from "@/gen/tft/v1/tierlist_pb";
import type { Champion, Item } from "@/gen/tft/v1/patch_pb";
import { CompositionCard } from "./composition-card";

const tierStyles: Record<string, string> = {
  S: "border-l-4 border-l-yellow-500",
  A: "border-l-4 border-l-green-500",
  B: "border-l-4 border-l-blue-500",
  C: "border-l-4 border-l-lofi-muted",
  D: "border-l-4 border-l-lofi-border",
};

const tierLabelStyles: Record<string, string> = {
  S: "text-yellow-500",
  A: "text-green-500",
  B: "text-blue-500",
  C: "text-lofi-muted",
  D: "text-lofi-border",
};

interface TierSectionProps {
  tier: string;
  entries: TierListEntry[];
  championMap: Map<string, Champion>;
  itemMap: Map<string, Item>;
}

export function TierSection({
  tier,
  entries,
  championMap,
  itemMap,
}: TierSectionProps) {
  if (entries.length === 0) return null;

  return (
    <div className={`mb-6 rounded-sm pl-3 ${tierStyles[tier] ?? ""}`}>
      <div className="mb-3 flex items-center gap-2">
        <span
          className={`text-lg font-black ${tierLabelStyles[tier] ?? "text-lofi-muted"}`}
        >
          {tier}
        </span>
        <span className="text-[10px] text-lofi-muted">
          {entries.length} comp{entries.length !== 1 ? "s" : ""}
        </span>
      </div>

      <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
        {entries.map((entry) => (
          <CompositionCard
            key={entry.compositionName}
            entry={entry}
            championMap={championMap}
            itemMap={itemMap}
          />
        ))}
      </div>
    </div>
  );
}
