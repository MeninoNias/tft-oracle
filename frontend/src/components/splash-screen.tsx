import { useState, useEffect, useCallback } from "react";

interface LogLine {
  timestamp: string;
  message: string;
  status: "success" | "ok" | "running" | "pending" | "waiting";
}

const statusColors: Record<string, string> = {
  success: "text-green-400",
  ok: "text-green-400",
  running: "text-yellow-400",
  pending: "text-lofi-muted",
  waiting: "text-lofi-muted",
};

const statusLabels: Record<string, string> = {
  success: "SUCCESS",
  ok: "OK",
  running: "RUNNING",
  pending: "PENDING",
  waiting: "WAITING",
};

const bootSequence: LogLine[] = [
  { timestamp: "", message: "INITIALIZING_HEXCORE_PROTOCOL", status: "success" },
  { timestamp: "", message: "LOADING_SYSTEM_MODULES", status: "success" },
  { timestamp: "", message: "FETCHING_LIVE_META_DATA [CURRENT_SET]", status: "ok" },
  { timestamp: "", message: "SCRAPING_RIOT_API", status: "ok" },
  { timestamp: "", message: "CALIBRATING_AUGMENTS", status: "running" },
  { timestamp: "", message: "COMPUTING_PROBABILITIES", status: "running" },
  { timestamp: "", message: "SYNCING_CLOUD_REPOSITORY", status: "pending" },
  { timestamp: "", message: "ESTABLISHING_WEBSOCKET_HANDSHAKE", status: "waiting" },
];

function getTimestamp(): string {
  const now = new Date();
  return `${String(now.getHours()).padStart(2, "0")}:${String(now.getMinutes()).padStart(2, "0")}:${String(now.getSeconds()).padStart(2, "0")}`;
}

interface SplashScreenProps {
  onComplete: () => void;
}

export function SplashScreen({ onComplete }: SplashScreenProps) {
  const [visibleLines, setVisibleLines] = useState<LogLine[]>([]);
  const [progress, setProgress] = useState(0);
  const [phase, setPhase] = useState<"boot" | "done">("boot");

  const addLine = useCallback((index: number) => {
    if (index >= bootSequence.length) {
      setPhase("done");
      return;
    }
    const line = bootSequence[index]!;
    setVisibleLines((prev) => [
      ...prev,
      { timestamp: getTimestamp(), message: line.message, status: line.status },
    ]);
  }, []);

  useEffect(() => {
    const timers: ReturnType<typeof setTimeout>[] = [];

    bootSequence.forEach((_, i) => {
      timers.push(setTimeout(() => addLine(i), 200 + i * 280));
    });

    // Trigger "done" after last line
    const doneDelay = 200 + bootSequence.length * 280;
    timers.push(setTimeout(() => setPhase("done"), doneDelay));

    // Progress bar animation
    const progressInterval = setInterval(() => {
      setProgress((p) => {
        if (p >= 100) {
          clearInterval(progressInterval);
          return 100;
        }
        return p + 2;
      });
    }, 40);

    timers.push(progressInterval as unknown as ReturnType<typeof setTimeout>);

    return () => {
      timers.forEach(clearTimeout);
      clearInterval(progressInterval);
    };
  }, [addLine]);

  useEffect(() => {
    if (phase === "done") {
      const timer = setTimeout(onComplete, 500);
      return () => clearTimeout(timer);
    }
  }, [phase, onComplete]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-lofi-black">
      {/* Dot grid background */}
      <div
        className="pointer-events-none fixed inset-0"
        style={{
          backgroundImage:
            "radial-gradient(circle, #27272a 1px, transparent 1px)",
          backgroundSize: "24px 24px",
        }}
      />

      <div className="relative z-10 w-full max-w-lg px-6">
        {/* Logo */}
        <div className="mb-8 text-center">
          <div className="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-lg border border-lofi-accent/30 bg-lofi-surface">
            <span className="text-3xl font-bold text-lofi-accent">//</span>
          </div>
          <h1 className="font-mono text-lg font-black uppercase tracking-widest text-white">
            TFT_ORACLE.sys
          </h1>
        </div>

        {/* Terminal window */}
        <div className="rounded-lg border border-lofi-border bg-lofi-surface overflow-hidden">
          {/* Title bar */}
          <div className="flex items-center justify-between border-b border-lofi-border px-4 py-2">
            <div className="flex gap-1.5">
              <span className="h-2.5 w-2.5 rounded-full bg-red-500/70" />
              <span className="h-2.5 w-2.5 rounded-full bg-yellow-500/70" />
              <span className="h-2.5 w-2.5 rounded-full bg-green-500/70" />
            </div>
            <span className="font-mono text-[10px] text-lofi-muted">
              SYSTEM_KERNEL v0.2.0
            </span>
          </div>

          {/* Log lines */}
          <div className="space-y-1 px-4 py-4 font-mono text-xs">
            {visibleLines.map((line, i) => (
              <div key={i} className="flex items-center gap-2 animate-in fade-in duration-200">
                <span className="text-lofi-muted">[{line.timestamp}]</span>
                <span className="text-lofi-secondary">{line.message}</span>
                <span className="ml-auto text-[10px]">
                  ...{" "}
                  <span className={statusColors[line.status]}>
                    {statusLabels[line.status]}
                  </span>
                </span>
              </div>
            ))}

            {visibleLines.length === 0 && (
              <div className="flex items-center gap-2 text-lofi-muted">
                <span className="animate-pulse">_</span>
                <span>awaiting boot sequence...</span>
              </div>
            )}
          </div>

          {/* Progress bar */}
          <div className="border-t border-lofi-border px-4 py-3">
            <div className="mb-1 flex items-center justify-between">
              <span className="font-mono text-[10px] uppercase text-lofi-muted">
                {progress < 100 ? "calibrating_augments" : "boot_complete"}
              </span>
              <span className="font-mono text-[10px] font-bold text-lofi-accent">
                {Math.min(progress, 100)}%
              </span>
            </div>
            <div className="h-1.5 w-full overflow-hidden rounded-full bg-lofi-black">
              <div
                className="h-full rounded-full bg-lofi-accent transition-all duration-100 ease-linear"
                style={{ width: `${Math.min(progress, 100)}%` }}
              />
            </div>
          </div>
        </div>

        {/* Faded status text */}
        {progress < 100 && (
          <p className="mt-4 text-center font-mono text-[10px] italic text-lofi-muted/50 animate-pulse">
            syncing unit_stats from cloud repository...
          </p>
        )}

        {/* System stats */}
        <div className="mt-6 flex items-center justify-center gap-6 font-mono text-[10px] text-lofi-muted">
          <div className="text-center">
            <p className="uppercase tracking-wider">status</p>
            <p className="font-bold text-green-400">ENCRYPTED</p>
          </div>
          <div className="text-center">
            <p className="uppercase tracking-wider">region</p>
            <p className="font-bold text-white">GLOBAL</p>
          </div>
          <div className="text-center">
            <p className="uppercase tracking-wider">latency</p>
            <p className="font-bold text-white">3.8ms</p>
          </div>
        </div>
      </div>
    </div>
  );
}
