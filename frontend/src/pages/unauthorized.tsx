import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";

export function UnauthorizedPage() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col items-center justify-center py-16">
      {/* Big faded code behind title */}
      <div className="relative mb-8 text-center">
        <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 select-none font-mono text-[100px] font-black leading-none text-lofi-accent/5">
          401
        </span>
        <h1 className="relative font-mono text-4xl font-black uppercase tracking-widest text-white">
          ERROR
        </h1>
      </div>

      {/* Error card */}
      <div className="mb-8 w-full max-w-xl rounded-lg border border-lofi-accent/20 bg-lofi-surface/80 p-6 shadow-lg shadow-lofi-accent/5 backdrop-blur-sm">
        <div className="mb-4 flex flex-col gap-2">
          <p className="font-mono text-2xl font-black leading-tight tracking-tighter text-lofi-accent">
            401: ACCESS_DENIED
          </p>
          <div className="flex items-center gap-2">
            <span className="h-2 w-2 animate-pulse rounded-full bg-red-500" />
            <p className="font-mono text-xs uppercase tracking-widest text-lofi-muted">
              status: UNAUTHORIZED
            </p>
          </div>
        </div>

        <div className="mb-4 h-px w-full bg-lofi-accent/20" />

        <div className="space-y-2 font-mono text-sm leading-relaxed opacity-80">
          <p className="flex gap-2 text-lofi-secondary">
            <span className="text-lofi-accent">&gt;</span>
            Access key invalid or expired.
          </p>
          <p className="flex gap-2 text-lofi-secondary">
            <span className="text-lofi-accent">&gt;</span>
            Identity verification failed.
          </p>
          <p className="flex gap-2 text-lofi-secondary">
            <span className="text-lofi-accent">&gt;</span>
            Clearance level: NONE.
          </p>
          <p className="flex gap-2 italic text-lofi-accent/30">
            <span>&gt;</span>
            Requesting re-authentication... [AWAITING]
          </p>
        </div>
      </div>

      {/* Action buttons */}
      <div className="flex gap-3">
        <Button onClick={() => navigate("/auth")}>
          authenticate: login
        </Button>
        <Button variant="secondary" onClick={() => navigate("/")}>
          home: return to base
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
          end of transmission // TFT-ORCL-ERR-401
        </p>
      </div>
    </div>
  );
}
