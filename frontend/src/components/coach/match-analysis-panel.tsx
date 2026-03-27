import { TerminalCard } from "@/components/ui/terminal-card";
import type { AnalyzeMatchResponse } from "@/gen/tft/v1/coach_pb";

interface MatchAnalysisPanelProps {
  analysis: AnalyzeMatchResponse;
}

const gradeColors: Record<string, string> = {
  "S+": "text-yellow-400 border-yellow-400",
  S: "text-yellow-400 border-yellow-400",
  "A+": "text-green-400 border-green-400",
  A: "text-green-400 border-green-400",
  "B+": "text-blue-400 border-blue-400",
  B: "text-blue-400 border-blue-400",
  "C+": "text-orange-400 border-orange-400",
  C: "text-orange-400 border-orange-400",
  D: "text-red-400 border-red-400",
  F: "text-red-500 border-red-500",
};

const insightGradeColors: Record<string, string> = {
  good: "bg-green-500",
  okay: "bg-yellow-500",
  poor: "bg-red-500",
};

export function MatchAnalysisPanel({ analysis }: MatchAnalysisPanelProps) {
  const gradeColor =
    gradeColors[analysis.overallGrade] ?? "text-lofi-muted border-lofi-muted";

  return (
    <div className="mt-2 space-y-3 border-l-2 border-lofi-accent/30 pl-4">
      {/* Grade + Summary */}
      <div className="flex items-start gap-3">
        <div
          className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-sm border-2 text-lg font-bold ${gradeColor}`}
        >
          {analysis.overallGrade}
        </div>
        <p className="text-sm text-lofi-secondary">{analysis.summary}</p>
      </div>

      {/* Insights Grid */}
      <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3">
        {analysis.insights.map((insight) => (
          <TerminalCard key={insight.category} className="!p-3">
            <div className="flex items-center gap-2">
              <span
                className={`h-2 w-2 rounded-full ${insightGradeColors[insight.grade] ?? "bg-lofi-muted"}`}
              />
              <span className="text-xs font-medium uppercase text-lofi-muted">
                {insight.category}
              </span>
            </div>
            <p className="mt-1 text-sm font-medium text-lofi-text">
              {insight.title}
            </p>
            <p className="mt-0.5 text-xs text-lofi-secondary">
              {insight.detail}
            </p>
          </TerminalCard>
        ))}
      </div>

      {/* Meta Comparison */}
      {analysis.metaComparison && (
        <TerminalCard className="!p-3">
          <p className="text-xs font-medium uppercase text-lofi-muted">
            meta comparison
          </p>
          <p className="mt-1 text-sm text-lofi-text">
            Closest comp:{" "}
            <span className="font-medium text-lofi-accent">
              {analysis.metaComparison.closestMetaComp}
            </span>{" "}
            ({analysis.metaComparison.metaTier}-tier)
          </p>
          {analysis.metaComparison.missingUnits.length > 0 && (
            <p className="mt-1 text-xs text-lofi-secondary">
              Missing: {analysis.metaComparison.missingUnits.join(", ")}
            </p>
          )}
          {analysis.metaComparison.suboptimalItems.length > 0 && (
            <p className="text-xs text-lofi-secondary">
              Suboptimal items:{" "}
              {analysis.metaComparison.suboptimalItems.join(", ")}
            </p>
          )}
          <p className="mt-1 text-xs text-lofi-secondary">
            {analysis.metaComparison.assessment}
          </p>
        </TerminalCard>
      )}

      {/* Lobby Context */}
      {analysis.lobbyContext && (
        <TerminalCard className="!p-3">
          <p className="text-xs font-medium uppercase text-lofi-muted">
            lobby context
          </p>
          <p className="mt-1 text-sm text-lofi-text">
            {analysis.lobbyContext.contestedCount > 0
              ? `${analysis.lobbyContext.contestedCount} player(s) contested your comp`
              : "Your comp was uncontested"}
          </p>
          <p className="mt-0.5 text-xs text-lofi-secondary">
            {analysis.lobbyContext.lobbyStrengthAssessment}
          </p>
          {analysis.lobbyContext.contestedDetails.map((detail, i) => (
            <p key={i} className="text-xs text-lofi-muted">
              {detail}
            </p>
          ))}
        </TerminalCard>
      )}

      {/* Suggestions */}
      {analysis.suggestions.length > 0 && (
        <div>
          <p className="mb-1 text-xs font-medium uppercase text-lofi-muted">
            suggestions
          </p>
          <div className="space-y-1">
            {analysis.suggestions.map((s, i) => (
              <div
                key={i}
                className="flex items-start gap-2 text-xs text-lofi-secondary"
              >
                <span
                  className={`mt-0.5 shrink-0 rounded-sm px-1 py-0.5 text-[10px] font-medium uppercase ${
                    s.priority === "high"
                      ? "bg-red-500/20 text-red-400"
                      : s.priority === "medium"
                        ? "bg-yellow-500/20 text-yellow-400"
                        : "bg-blue-500/20 text-blue-400"
                  }`}
                >
                  {s.priority}
                </span>
                <span>{s.description}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
