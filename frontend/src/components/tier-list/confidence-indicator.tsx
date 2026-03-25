const styles: Record<string, string> = {
  high: "bg-green-500",
  medium: "bg-yellow-500",
  low: "bg-red-500",
};

interface ConfidenceIndicatorProps {
  level: string;
}

export function ConfidenceIndicator({ level }: ConfidenceIndicatorProps) {
  const dotColor = styles[level] ?? styles.low;

  return (
    <span className="inline-flex items-center gap-1 text-[10px] text-lofi-secondary">
      <span className={`h-1.5 w-1.5 rounded-full ${dotColor}`} />
      {level}
    </span>
  );
}
