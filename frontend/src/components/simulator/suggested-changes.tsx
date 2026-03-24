import { Badge } from "@/components/ui/badge";

interface Change {
  description: string;
  priority: string;
  category: string;
}

const priorityColors: Record<string, string> = {
  high: "border-red-500 text-red-400",
  medium: "border-yellow-500 text-yellow-400",
  low: "border-lofi-border text-lofi-muted",
};

interface SuggestedChangesProps {
  changes: Change[];
}

export function SuggestedChanges({ changes }: SuggestedChangesProps) {
  if (changes.length === 0) return null;

  return (
    <div>
      <p className="mb-1.5 text-[10px] uppercase tracking-wider text-lofi-muted">
        suggested changes
      </p>
      <div className="space-y-1.5">
        {changes.map((change, i) => (
          <div
            key={i}
            className="flex items-start gap-2 rounded-sm border border-lofi-border bg-lofi-surface px-3 py-2"
          >
            <Badge className={priorityColors[change.priority] ?? ""}>
              {change.priority}
            </Badge>
            <Badge>{change.category}</Badge>
            <span className="text-xs text-lofi-secondary">
              {change.description}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
