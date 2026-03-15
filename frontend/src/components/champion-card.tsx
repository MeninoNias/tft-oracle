import { type Champion } from "../gen/tft/v1/patch_pb";
import { TerminalCard } from "./ui/terminal-card";
import { Badge } from "./ui/badge";

const costColors: Record<number, string> = {
  1: "border-cost-1 text-cost-1",
  2: "border-cost-2 text-cost-2",
  3: "border-cost-3 text-cost-3",
  4: "border-cost-4 text-cost-4",
  5: "border-cost-5 text-cost-5",
};

interface ChampionCardProps {
  champion: Champion;
  traitNames?: Record<string, string>;
}

export function ChampionCard({ champion, traitNames = {} }: ChampionCardProps) {
  const costColor = costColors[champion.cost] ?? "border-lofi-border text-lofi-secondary";

  return (
    <TerminalCard className={`border-l-2 ${costColor.split(" ")[0]}`}>
      <div className="flex items-start gap-3">
        {champion.squareIconUrl ? (
          <img
            src={champion.squareIconUrl}
            alt={champion.name}
            className="h-12 w-12 rounded-sm bg-lofi-black object-cover"
            loading="lazy"
          />
        ) : (
          <div className="flex h-12 w-12 items-center justify-center rounded-sm bg-lofi-black text-[10px] text-lofi-muted">
            ?
          </div>
        )}
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="truncate text-xs font-medium text-white">
              {champion.name}
            </span>
            <span className={`text-[10px] font-bold ${costColor.split(" ")[1]}`}>
              {champion.cost}g
            </span>
          </div>
          <div className="mt-1 flex flex-wrap gap-1">
            {champion.traitApiNames.map((trait) => (
              <Badge key={trait}>
                {traitNames[trait] ?? trait.replace(/^.*_/, "")}
              </Badge>
            ))}
          </div>
        </div>
      </div>
    </TerminalCard>
  );
}
