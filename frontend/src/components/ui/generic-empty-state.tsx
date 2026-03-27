import { Button } from "@/components/ui/button";

interface GenericEmptyStateProps {
  title: string;
  description: string;
  actionLabel?: string;
  onAction?: () => void;
}

export function GenericEmptyState({
  title,
  description,
  actionLabel,
  onAction,
}: GenericEmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <div className="mb-6 flex h-12 w-12 items-center justify-center rounded-lg border border-lofi-border bg-lofi-surface">
        <span className="text-xl font-bold text-lofi-muted">//</span>
      </div>

      <h3 className="mb-1 text-xs font-bold uppercase tracking-widest text-lofi-secondary">
        {title}
      </h3>

      <p className="mb-6 max-w-xs text-center text-[10px] leading-relaxed text-lofi-muted">
        {description}
      </p>

      {actionLabel && onAction && (
        <Button variant="secondary" onClick={onAction}>
          {actionLabel}
        </Button>
      )}
    </div>
  );
}
