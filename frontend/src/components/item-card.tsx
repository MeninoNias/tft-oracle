import { type Item } from "../gen/tft/v1/patch_pb";
import { TerminalCard } from "./ui/terminal-card";
import { Badge } from "./ui/badge";

interface ItemCardProps {
  item: Item;
}

export function ItemCard({ item }: ItemCardProps) {
  const isComponent = item.composition.length === 0;
  const isCraftable = item.composition.length === 2;

  return (
    <TerminalCard>
      <div className="flex items-start gap-3">
        {item.iconUrl ? (
          <img
            src={item.iconUrl}
            alt={item.name}
            className="h-10 w-10 rounded-sm bg-lofi-black object-cover"
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
              {item.name}
            </span>
            {isComponent && <Badge>component</Badge>}
            {isCraftable && <Badge>crafted</Badge>}
            {item.unique && <Badge className="border-lofi-accent text-lofi-accent">unique</Badge>}
          </div>
          {item.desc && (
            <p className="mt-1 line-clamp-2 text-[10px] leading-relaxed text-lofi-secondary">
              {item.desc.replace(/<[^>]*>/g, "")}
            </p>
          )}
        </div>
      </div>
    </TerminalCard>
  );
}
