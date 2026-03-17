import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";

interface EmptyStateProps {
  onStartOnboarding: () => void;
}

export function EmptyState({ onStartOnboarding }: EmptyStateProps) {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col items-center justify-center py-20">
      {/* Icon */}
      <div className="mb-6 flex h-16 w-16 items-center justify-center rounded-lg border border-lofi-border bg-lofi-surface">
        <span className="text-3xl font-bold text-lofi-accent">//</span>
      </div>

      {/* Heading */}
      <h2 className="mb-2 text-sm font-bold uppercase tracking-widest text-white">
        DATA_NOT_FOUND
      </h2>

      {/* Description */}
      <p className="mb-8 max-w-sm text-center text-xs leading-relaxed text-lofi-secondary">
        the oracle is silent. link your riot id to unlock tactical intelligence
        — match history, ranked stats, and ai-powered coaching.
      </p>

      {/* Actions */}
      <div className="flex gap-3">
        <Button onClick={onStartOnboarding}>CONSULT_ORACLE</Button>
        <Button variant="secondary" onClick={() => navigate("/player")}>
          VIEW_DEMO_DATA
        </Button>
      </div>

      {/* Status footer */}
      <div className="mt-12 flex items-center gap-4 text-[10px] text-lofi-muted">
        <span className="flex items-center gap-1">
          <span className="h-1.5 w-1.5 rounded-full bg-green-500" />
          system: operational
        </span>
        <span>v0.2.0</span>
      </div>
    </div>
  );
}
