interface PositioningTipProps {
  tip: string;
}

export function PositioningTip({ tip }: PositioningTipProps) {
  return (
    <div className="rounded-sm border-l-2 border-lofi-accent bg-lofi-surface px-4 py-3">
      <p className="mb-0.5 text-[10px] uppercase tracking-wider text-lofi-accent">
        positioning tip
      </p>
      <p className="text-xs text-lofi-text">{tip}</p>
    </div>
  );
}
