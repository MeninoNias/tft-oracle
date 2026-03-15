import type { PatchLookup } from "@/hooks/use-patch-lookup";
import type { MatchTrait } from "@/gen/tft/v1/player_pb";

interface TraitBadgeProps {
  trait: MatchTrait;
  lookup: PatchLookup;
}

const styleColors: Record<number, string> = {
  1: "bg-amber-900/30 border-amber-700/50 text-amber-500",
  2: "bg-gray-400/20 border-gray-400/50 text-gray-300",
  3: "bg-yellow-500/20 border-yellow-500/50 text-yellow-400",
  4: "bg-yellow-400/25 border-yellow-300/50 text-yellow-300",
};

export function TraitBadge({ trait, lookup }: TraitBadgeProps) {
  const traitData = lookup.traitByApiName.get(trait.apiName);
  const color =
    styleColors[trait.style] ??
    "bg-lofi-surface border-lofi-border text-lofi-muted";

  return (
    <div
      className={`flex items-center gap-1 rounded border px-1 py-0.5 ${color}`}
      title={traitData?.name ?? trait.apiName}
    >
      {traitData?.iconUrl ? (
        <img
          src={traitData.iconUrl}
          alt={traitData.name}
          className="h-4 w-4"
          loading="lazy"
        />
      ) : (
        <span className="text-[9px] leading-none">
          {formatTraitName(trait.apiName)}
        </span>
      )}
      <span className="text-xs font-bold leading-none">{trait.numUnits}</span>
    </div>
  );
}

function formatTraitName(apiName: string): string {
  return apiName.replace(/^Set\d+_/, "").replace(/TFT\d+_/, "");
}
