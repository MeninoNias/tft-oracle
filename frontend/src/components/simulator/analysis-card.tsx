import { TerminalCard } from "@/components/ui/terminal-card";

interface AnalysisCardProps {
  analysis: string;
}

export function AnalysisCard({ analysis }: AnalysisCardProps) {
  return (
    <TerminalCard>
      <p className="mb-1 text-[10px] uppercase tracking-wider text-lofi-muted">
        tactical analysis
      </p>
      <p className="text-xs leading-relaxed text-lofi-secondary">{analysis}</p>
    </TerminalCard>
  );
}
