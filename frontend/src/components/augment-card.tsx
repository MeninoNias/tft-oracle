import type { Item } from "@/gen/tft/v1/patch_pb";
import { TerminalCard } from "@/components/ui/terminal-card";
import { Badge } from "@/components/ui/badge";

interface AugmentCardProps {
  augment: Item;
}

const tierColors: Record<string, string> = {
  silver: "border-gray-400 text-gray-300",
  gold: "border-yellow-500 text-yellow-400",
  prismatic: "border-purple-400 text-purple-300",
};

function getAugmentTier(tags: string[]): string {
  for (const tag of tags) {
    const lower = tag.toLowerCase();
    if (lower.includes("prismatic")) return "prismatic";
    if (lower.includes("gold")) return "gold";
    if (lower.includes("silver")) return "silver";
  }
  return "silver";
}

export function AugmentCard({ augment }: AugmentCardProps) {
  const tier = getAugmentTier(augment.tags);
  const tierColor = tierColors[tier] ?? "border-lofi-border text-lofi-secondary";

  return (
    <TerminalCard>
      <div className="flex items-start gap-3">
        {augment.iconUrl ? (
          <img
            src={augment.iconUrl}
            alt={augment.name}
            className={`h-10 w-10 rounded-sm border ${tierColor} bg-lofi-black object-cover`}
            loading="lazy"
          />
        ) : (
          <div className="flex h-10 w-10 items-center justify-center rounded-sm bg-lofi-black text-[10px] text-lofi-muted">
            ?
          </div>
        )}
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="truncate text-xs font-medium text-white">
              {augment.name}
            </span>
            <Badge className={tierColor}>{tier}</Badge>
          </div>
          {augment.desc && (
            <p className="mt-1 line-clamp-2 text-[10px] leading-relaxed text-lofi-secondary">
              {augment.desc.replace(/<[^>]*>/g, "")}
            </p>
          )}
        </div>
      </div>
    </TerminalCard>
  );
}
