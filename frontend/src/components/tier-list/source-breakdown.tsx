import { useState } from "react";

interface SourceBreakdownProps {
  metatftTier: string;
  tftacticsTier: string;
  mobalyticsTier: string;
}

const sources = [
  { key: "metatftTier" as const, label: "MetaTFT" },
  { key: "tftacticsTier" as const, label: "TFTactics" },
  { key: "mobalyticsTier" as const, label: "Mobalytics" },
];

export function SourceBreakdown({
  metatftTier,
  tftacticsTier,
  mobalyticsTier,
}: SourceBreakdownProps) {
  const [open, setOpen] = useState(false);

  const tiers: Record<string, string> = {
    metatftTier,
    tftacticsTier,
    mobalyticsTier,
  };

  return (
    <div className="mt-2 border-t border-lofi-border pt-2">
      <button
        onClick={() => setOpen(!open)}
        className="text-[10px] text-lofi-muted hover:text-lofi-secondary transition-colors"
      >
        {open ? "- hide sources" : "+ view sources"}
      </button>

      {open && (
        <div className="mt-2 space-y-1">
          {sources.map((s) => (
            <div
              key={s.key}
              className="flex items-center justify-between text-[10px]"
            >
              <span className="text-lofi-muted">{s.label}</span>
              <span className="text-lofi-secondary">
                {tiers[s.key] || "—"}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
