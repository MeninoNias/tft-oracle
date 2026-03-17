import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";

interface ErrorScreenProps {
  code: string;
  title: string;
  lines: string[];
  fadedLine?: string;
  statusLabel?: string;
  transmissionId?: string;
}

export function ErrorScreen({
  code,
  title,
  lines,
  fadedLine,
  statusLabel = "DISCONNECTED",
  transmissionId = "TFT-ORCL-ERR-09X",
}: ErrorScreenProps) {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col items-center justify-center py-16">
      {/* Big faded code behind title */}
      <div className="relative mb-8 text-center">
        <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 select-none font-mono text-[100px] font-black leading-none text-lofi-accent/5">
          {code}
        </span>
        <h1 className="relative font-mono text-4xl font-black uppercase tracking-widest text-white">
          ERROR
        </h1>
      </div>

      {/* Error card */}
      <div className="mb-8 w-full max-w-xl rounded-lg border border-lofi-accent/20 bg-lofi-surface/80 p-6 shadow-lg shadow-lofi-accent/5 backdrop-blur-sm">
        <div className="mb-4 flex flex-col gap-2">
          <p className="font-mono text-2xl font-black leading-tight tracking-tighter text-lofi-accent">
            {code}: {title}
          </p>
          <div className="flex items-center gap-2">
            <span className="h-2 w-2 animate-pulse rounded-full bg-red-500" />
            <p className="font-mono text-xs uppercase tracking-widest text-lofi-muted">
              status: {statusLabel}
            </p>
          </div>
        </div>

        <div className="mb-4 h-px w-full bg-lofi-accent/20" />

        <div className="space-y-2 font-mono text-sm leading-relaxed opacity-80">
          {lines.map((line, i) => (
            <p key={i} className="flex gap-2 text-lofi-secondary">
              <span className="text-lofi-accent">&gt;</span>
              {line}
            </p>
          ))}
          {fadedLine && (
            <p className="flex gap-2 italic text-lofi-accent/30">
              <span>&gt;</span>
              {fadedLine}
            </p>
          )}
        </div>
      </div>

      {/* Action buttons */}
      <div className="flex gap-3">
        <Button onClick={() => navigate("/")}>home: return to base</Button>
        <Button variant="secondary" onClick={() => navigate(-1 as never)}>
          back: previous page
        </Button>
      </div>

      {/* Footer */}
      <div className="mt-10 flex flex-col items-center gap-3">
        <div className="flex gap-3 opacity-40">
          <span className="h-2 w-2 rounded-full bg-lofi-accent" />
          <span className="h-2 w-2 rounded-full bg-lofi-accent" />
          <span className="h-2 w-2 rounded-full bg-lofi-accent" />
        </div>
        <p className="font-mono text-[10px] uppercase tracking-[0.2em] text-lofi-muted">
          end of transmission // {transmissionId}
        </p>
      </div>
    </div>
  );
}
