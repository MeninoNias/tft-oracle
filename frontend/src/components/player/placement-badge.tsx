interface PlacementBadgeProps {
  placement: number;
  size?: "sm" | "md";
  className?: string;
}

const placementColors: Record<number, string> = {
  1: "bg-yellow-500/20 text-yellow-400 border-yellow-500/40",
  2: "bg-gray-300/20 text-gray-300 border-gray-400/40",
  3: "bg-amber-600/20 text-amber-500 border-amber-600/40",
  4: "bg-green-500/20 text-green-400 border-green-500/40",
  5: "bg-red-400/10 text-red-400/80 border-red-400/30",
  6: "bg-red-400/10 text-red-400/80 border-red-400/30",
  7: "bg-red-500/10 text-red-500/80 border-red-500/30",
  8: "bg-red-600/10 text-red-500 border-red-600/30",
};

const sizes = {
  sm: "h-8 w-8 text-sm font-bold",
  md: "h-11 w-11 text-xl font-bold",
};

export function PlacementBadge({
  placement,
  size = "sm",
  className = "",
}: PlacementBadgeProps) {
  const color =
    placementColors[placement] ??
    "bg-lofi-surface text-lofi-muted border-lofi-border";

  return (
    <span
      className={`inline-flex items-center justify-center rounded-sm border ${sizes[size]} ${color} ${className}`}
    >
      {placement}
    </span>
  );
}
