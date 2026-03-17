import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useRegisterAndLogin } from "@/hooks/use-auth";
import { SERVERS } from "@/lib/servers";
import { friendlyError } from "@/lib/error";

type Step = "identify" | "connect" | "augment";

interface OnboardingFlowProps {
  onComplete: () => void;
}

export function OnboardingFlow({ onComplete }: OnboardingFlowProps) {
  const [step, setStep] = useState<Step>("identify");
  const [gameName, setGameName] = useState("");
  const [tagLine, setTagLine] = useState("");
  const [region, setRegion] = useState("br");
  const [accessKey, setAccessKey] = useState("");
  const [copied, setCopied] = useState(false);

  const registerAndLogin = useRegisterAndLogin();
  const navigate = useNavigate();

  const handleSync = async (e: React.FormEvent) => {
    e.preventDefault();
    setStep("connect");
    try {
      const res = await registerAndLogin.mutateAsync({
        gameName: gameName.trim(),
        tagLine: tagLine.trim(),
        region,
      });
      setAccessKey(res.accessKey);
    } catch {
      // Error is available via registerAndLogin.error
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(accessKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const steps: { key: Step; label: string }[] = [
    { key: "identify", label: "IDENTIFY" },
    { key: "connect", label: "CONNECT" },
    { key: "augment", label: "AUGMENT" },
  ];

  const currentIndex = steps.findIndex((s) => s.key === step);

  return (
    <div className="mx-auto max-w-lg py-10">
      {/* Step indicator */}
      <div className="mb-8 flex items-center justify-center gap-6">
        {steps.map((s, i) => (
          <div key={s.key} className="flex items-center gap-2">
            <span
              className={`flex h-6 w-6 items-center justify-center rounded-full text-[10px] font-bold ${
                i <= currentIndex
                  ? "bg-lofi-accent text-lofi-black"
                  : "border border-lofi-border text-lofi-muted"
              }`}
            >
              {String(i + 1).padStart(2, "0")}
            </span>
            <span
              className={`text-[10px] uppercase tracking-wider ${
                i <= currentIndex ? "text-white" : "text-lofi-muted"
              }`}
            >
              {s.label}
            </span>
          </div>
        ))}
      </div>

      {/* Step content */}
      <div className="rounded-lg border border-lofi-border bg-lofi-surface p-6">
        {step === "identify" && (
          <IdentifyStep
            gameName={gameName}
            tagLine={tagLine}
            region={region}
            onGameNameChange={setGameName}
            onTagLineChange={setTagLine}
            onRegionChange={setRegion}
            onSubmit={handleSync}
            isPending={false}
          />
        )}

        {step === "connect" && (
          <ConnectStep
            isPending={registerAndLogin.isPending}
            error={registerAndLogin.error}
            accessKey={accessKey}
            copied={copied}
            onCopy={handleCopy}
            onContinue={() => setStep("augment")}
            onRetry={() => {
              registerAndLogin.reset();
              setStep("identify");
            }}
          />
        )}

        {step === "augment" && (
          <AugmentStep
            gameName={gameName}
            tagLine={tagLine}
            onViewProfile={onComplete}
            onPlayerSearch={() => navigate("/player")}
            onChampions={() => navigate("/champions")}
          />
        )}
      </div>

      {/* Already registered link */}
      {step === "identify" && (
        <p className="mt-4 text-center text-[10px] text-lofi-muted">
          already registered?{" "}
          <button
            onClick={() => navigate("/auth")}
            className="text-lofi-accent hover:underline"
          >
            login with access key
          </button>
        </p>
      )}
    </div>
  );
}

function IdentifyStep({
  gameName,
  tagLine,
  region,
  onGameNameChange,
  onTagLineChange,
  onRegionChange,
  onSubmit,
  isPending,
}: {
  gameName: string;
  tagLine: string;
  region: string;
  onGameNameChange: (v: string) => void;
  onTagLineChange: (v: string) => void;
  onRegionChange: (v: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  isPending: boolean;
}) {
  return (
    <div>
      <h3 className="mb-1 text-sm font-bold text-white">
        Welcome, Tactician.
      </h3>
      <p className="mb-5 text-xs text-lofi-secondary">
        link your riot id to activate the oracle. your identity is verified
        through the riot api — no password needed.
      </p>

      <form onSubmit={onSubmit} className="space-y-4">
        <div>
          <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
            game name
          </label>
          <Input
            value={gameName}
            onChange={(e) => onGameNameChange(e.target.value)}
            placeholder="pai de olivia"
            disabled={isPending}
          />
        </div>

        <div className="flex gap-3">
          <div className="flex-1">
            <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
              tag line
            </label>
            <Input
              value={tagLine}
              onChange={(e) => onTagLineChange(e.target.value)}
              placeholder="NIAS"
              disabled={isPending}
            />
          </div>
          <div className="w-28">
            <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
              server
            </label>
            <select
              value={region}
              onChange={(e) => onRegionChange(e.target.value)}
              disabled={isPending}
              className="w-full rounded-sm border border-lofi-border bg-lofi-black px-3 py-1.5 text-xs text-lofi-text focus:border-lofi-accent focus:outline-none focus:ring-1 focus:ring-lofi-accent"
            >
              {SERVERS.map((s) => (
                <option key={s.value} value={s.value}>
                  {s.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={isPending || !gameName.trim() || !tagLine.trim()}
        >
          SYNC_ACCOUNT
        </Button>
      </form>
    </div>
  );
}

function ConnectStep({
  isPending,
  error,
  accessKey,
  copied,
  onCopy,
  onContinue,
  onRetry,
}: {
  isPending: boolean;
  error: Error | null;
  accessKey: string;
  copied: boolean;
  onCopy: () => void;
  onContinue: () => void;
  onRetry: () => void;
}) {
  if (isPending) {
    return (
      <div className="flex flex-col items-center py-8">
        <div className="mb-4 h-8 w-8 animate-spin rounded-full border-2 border-lofi-border border-t-lofi-accent" />
        <p className="text-xs text-lofi-secondary">
          verifying riot id and generating access key...
        </p>
      </div>
    );
  }

  if (error) {
    return (
      <div>
        <div className="mb-4 flex items-center gap-2">
          <span className="flex h-6 w-6 items-center justify-center rounded-full bg-red-500/20 text-xs text-red-400">
            !
          </span>
          <h3 className="text-sm font-bold text-white">
            connection failed
          </h3>
        </div>
        <div className="mb-4 rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
          {friendlyError(error)}
        </div>
        <Button onClick={onRetry} className="w-full">
          back to step 1
        </Button>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-4 flex items-center gap-2">
        <span className="flex h-6 w-6 items-center justify-center rounded-full bg-green-500/20 text-xs text-green-400">
          ✓
        </span>
        <h3 className="text-sm font-bold text-white">connection established</h3>
      </div>

      <p className="mb-4 text-xs leading-relaxed text-lofi-secondary">
        save this access key — it will only be shown{" "}
        <span className="font-bold text-lofi-accent">once</span>. use it to log
        in on any device.
      </p>

      <div className="mb-4 rounded-sm border border-lofi-accent/30 bg-lofi-black p-4">
        <p className="mb-1 text-[10px] uppercase tracking-wider text-lofi-muted">
          your access key
        </p>
        <code className="block break-all text-sm leading-relaxed text-lofi-accent">
          {accessKey}
        </code>
      </div>

      <div className="mb-4 flex gap-2">
        <Button onClick={onCopy} variant="secondary" className="flex-1">
          {copied ? "copied!" : "copy key"}
        </Button>
      </div>

      {/* Status badges */}
      <div className="mb-4 flex items-center justify-center gap-4 text-[10px] text-lofi-muted">
        <span className="flex items-center gap-1">
          <span className="h-1 w-1 rounded-full bg-green-500" />
          ENCRYPTION_ACTIVE
        </span>
        <span className="flex items-center gap-1">
          <span className="h-1 w-1 rounded-full bg-green-500" />
          RIOT_API_VERIFIED
        </span>
      </div>

      <Button onClick={onContinue} className="w-full">
        CONTINUE
      </Button>
    </div>
  );
}

function AugmentStep({
  gameName,
  tagLine,
  onViewProfile,
  onPlayerSearch,
  onChampions,
}: {
  gameName: string;
  tagLine: string;
  onViewProfile: () => void;
  onPlayerSearch: () => void;
  onChampions: () => void;
}) {
  return (
    <div className="text-center">
      <div className="mb-4 flex items-center justify-center gap-2">
        <span className="flex h-6 w-6 items-center justify-center rounded-full bg-green-500/20 text-xs text-green-400">
          ✓
        </span>
        <h3 className="text-sm font-bold text-white">oracle activated</h3>
      </div>

      <p className="mb-2 text-xs text-lofi-secondary">
        linked as{" "}
        <span className="font-bold text-lofi-accent">
          {gameName}#{tagLine}
        </span>
      </p>

      <p className="mb-6 text-xs text-lofi-muted">
        tactical intelligence is now online. choose your next move.
      </p>

      <div className="flex flex-col gap-2">
        <Button onClick={onViewProfile} className="w-full">
          view profile
        </Button>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            className="flex-1"
            onClick={onPlayerSearch}
          >
            player search
          </Button>
          <Button
            variant="secondary"
            className="flex-1"
            onClick={onChampions}
          >
            champions
          </Button>
        </div>
      </div>
    </div>
  );
}
