import type { PatchLookup } from "@/hooks/use-patch-lookup";
import type { MatchUnit } from "@/gen/tft/v1/player_pb";

interface UnitIconProps {
  unit: MatchUnit;
  lookup: PatchLookup;
}

const costBorderColors: Record<number, string> = {
  1: "border-cost-1",
  2: "border-cost-2",
  3: "border-cost-3",
  4: "border-cost-4",
  5: "border-cost-5",
};

export function UnitIcon({ unit, lookup }: UnitIconProps) {
  const champion = lookup.championByApiName.get(unit.characterId);
  const cost = champion?.cost ?? unit.rarity + 1;
  const borderColor = costBorderColors[cost] ?? "border-lofi-border";
  const iconUrl = champion?.squareIconUrl;

  return (
    <div className="flex h-[62px] flex-col items-center gap-0.5">
      {/* Champion icon */}
      <div className="relative">
        {iconUrl ? (
          <img
            src={iconUrl}
            alt={champion?.name ?? unit.characterId}
            className={`h-11 w-11 rounded border-2 ${borderColor} object-cover`}
            loading="lazy"
          />
        ) : (
          <div
            className={`flex h-11 w-11 items-center justify-center rounded border-2 ${borderColor} bg-lofi-black text-[9px] text-lofi-muted`}
          >
            ?
          </div>
        )}

        {/* Star overlay */}
        {unit.tier >= 2 && (
          <div className="absolute -top-2 left-1/2 -translate-x-1/2">
            <span
              className={`text-[9px] font-bold leading-none ${
                unit.tier >= 3
                  ? "text-yellow-300 drop-shadow-[0_0_4px_rgba(253,224,71,0.8)]"
                  : "text-yellow-400 drop-shadow-[0_0_2px_rgba(250,204,21,0.5)]"
              }`}
            >
              {"★".repeat(unit.tier)}
            </span>
          </div>
        )}
      </div>

      {/* Item icons */}
      {unit.itemNames.length > 0 && (
        <div className="flex gap-0.5">
          {unit.itemNames.slice(0, 3).map((itemName, i) => {
            const item = lookup.itemByApiName.get(itemName);
            return item?.iconUrl ? (
              <img
                key={i}
                src={item.iconUrl}
                alt={item.name}
                title={item.name}
                className="h-3.5 w-3.5 rounded-sm border border-lofi-border/50"
                loading="lazy"
              />
            ) : (
              <div
                key={i}
                className="h-3.5 w-3.5 rounded-sm bg-lofi-muted/50"
                title={itemName}
              />
            );
          })}
        </div>
      )}
    </div>
  );
}
