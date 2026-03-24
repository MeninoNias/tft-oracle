import { Badge } from "@/components/ui/badge";

interface KeyFactorsProps {
  factors: string[];
}

export function KeyFactors({ factors }: KeyFactorsProps) {
  if (factors.length === 0) return null;

  return (
    <div>
      <p className="mb-1.5 text-[10px] uppercase tracking-wider text-lofi-muted">
        key factors
      </p>
      <div className="flex flex-wrap gap-1.5">
        {factors.map((factor, i) => (
          <Badge key={i}>{factor}</Badge>
        ))}
      </div>
    </div>
  );
}
