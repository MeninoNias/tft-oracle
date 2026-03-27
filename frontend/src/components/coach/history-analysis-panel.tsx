import { TerminalCard } from "@/components/ui/terminal-card";
import { SkillRadarChart } from "@/components/coach/skill-radar-chart";
import type { AnalyzeHistoryResponse } from "@/gen/tft/v1/coach_pb";

interface HistoryAnalysisPanelProps {
  analysis: AnalyzeHistoryResponse;
}

const sentimentColors: Record<string, string> = {
  positive: "bg-green-500",
  neutral: "bg-yellow-500",
  negative: "bg-red-500",
};

const directionIcons: Record<string, string> = {
  improving: "↑",
  stable: "→",
  declining: "↓",
};

const directionColors: Record<string, string> = {
  improving: "text-green-400",
  stable: "text-yellow-400",
  declining: "text-red-400",
};

export function HistoryAnalysisPanel({
  analysis,
}: HistoryAnalysisPanelProps) {
  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center gap-2">
        <span className="text-xs font-medium uppercase text-lofi-accent">
          ai coach
        </span>
        <span className="text-xs text-lofi-muted">
          {analysis.matchesAnalyzed} matches analyzed
        </span>
      </div>

      {/* Summary */}
      <TerminalCard>
        <p className="text-sm text-lofi-secondary">
          {analysis.overallSummary}
        </p>
      </TerminalCard>

      {/* Skill Radar */}
      {analysis.skillRadar && (
        <TerminalCard className="flex flex-col items-center">
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            skill radar
          </p>
          <SkillRadarChart radar={analysis.skillRadar} size={220} />
        </TerminalCard>
      )}

      {/* Patterns */}
      {analysis.patterns.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            patterns detected
          </p>
          <div className="space-y-2">
            {analysis.patterns.map((p, i) => (
              <TerminalCard key={i} className="!p-3">
                <div className="flex items-center gap-2">
                  <span
                    className={`h-2 w-2 rounded-full ${sentimentColors[p.sentiment] ?? "bg-lofi-muted"}`}
                  />
                  <span className="text-xs uppercase text-lofi-muted">
                    {p.category}
                  </span>
                </div>
                <p className="mt-1 text-sm font-medium text-lofi-text">
                  {p.title}
                </p>
                <p className="mt-0.5 text-xs text-lofi-secondary">
                  {p.detail}
                </p>
              </TerminalCard>
            ))}
          </div>
        </div>
      )}

      {/* Trends */}
      {analysis.trends.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            trends
          </p>
          <div className="space-y-1">
            {analysis.trends.map((t, i) => (
              <div
                key={i}
                className="flex items-center gap-2 text-xs text-lofi-secondary"
              >
                <span
                  className={`text-sm font-bold ${directionColors[t.direction] ?? "text-lofi-muted"}`}
                >
                  {directionIcons[t.direction] ?? "?"}
                </span>
                <span className="font-medium text-lofi-text">{t.metric}</span>
                <span>{t.detail}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Improvement Plan */}
      {analysis.improvementPlan.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-medium uppercase text-lofi-muted">
            improvement plan
          </p>
          <div className="space-y-2">
            {analysis.improvementPlan.map((s, i) => (
              <div
                key={i}
                className="flex items-start gap-3 text-sm text-lofi-secondary"
              >
                <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full border border-lofi-accent text-xs font-bold text-lofi-accent">
                  {i + 1}
                </span>
                <div>
                  <p className="text-lofi-text">{s.description}</p>
                  <span className="text-xs text-lofi-muted">{s.category}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
