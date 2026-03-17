import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useRegister, useLogin } from "@/hooks/use-auth";
import { SERVERS } from "@/lib/servers";
import { friendlyError } from "@/lib/error";

export function AuthPage() {
  const [searchParams] = useSearchParams();
  const isExpired = searchParams.get("expired") === "1";
  const [tab, setTab] = useState<"register" | "login">(
    isExpired ? "login" : "register",
  );

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

      <div className="relative z-10 w-full max-w-md px-6">
        {/* Header */}
        <div className="mb-8 text-center">
          <div className="mb-4 inline-flex h-14 w-14 items-center justify-center rounded-lg border border-lofi-border bg-lofi-surface">
            <span className="text-2xl font-bold text-lofi-accent">//</span>
          </div>
          <h1 className="text-lg font-bold tracking-wider text-white">
            TFT_ORACLE
          </h1>
          <p className="mt-1 text-xs text-lofi-muted">
            tactical intelligence system
          </p>
        </div>

        {/* Tab switcher */}
        <div className="mb-6 flex rounded-sm border border-lofi-border bg-lofi-black">
          <button
            onClick={() => setTab("register")}
            className={`flex-1 py-2 text-xs font-medium transition-colors ${
              tab === "register"
                ? "bg-lofi-surface text-white"
                : "text-lofi-muted hover:text-lofi-secondary"
            }`}
          >
            register
          </button>
          <button
            onClick={() => setTab("login")}
            className={`flex-1 py-2 text-xs font-medium transition-colors ${
              tab === "login"
                ? "bg-lofi-surface text-white"
                : "text-lofi-muted hover:text-lofi-secondary"
            }`}
          >
            login
          </button>
        </div>

        {/* Expired session banner */}
        {isExpired && (
          <div className="mb-4 rounded-sm border border-yellow-500/30 bg-yellow-500/10 px-4 py-3">
            <p className="text-xs font-medium text-yellow-400">
              session expired or invalid access key.
            </p>
            <p className="mt-0.5 text-[10px] text-yellow-400/70">
              please log in again to continue.
            </p>
          </div>
        )}

        {/* Form card */}
        <div className="rounded-lg border border-lofi-border bg-lofi-surface p-6">
          {tab === "register" ? <RegisterForm /> : <LoginForm />}
        </div>

        {/* Footer */}
        <div className="mt-6 text-center">
          <p className="text-[10px] text-lofi-muted">
            tft oracle is not affiliated with riot games
          </p>
          <div className="mt-2 flex items-center justify-center gap-3 text-[10px] text-lofi-muted">
            <span className="flex items-center gap-1">
              <span className="h-1.5 w-1.5 rounded-full bg-green-500" />
              system: operational
            </span>
            <span>v0.2.0</span>
          </div>
        </div>
      </div>
    </div>
  );
}

function RegisterForm() {
  const [gameName, setGameName] = useState("");
  const [tagLine, setTagLine] = useState("");
  const [region, setRegion] = useState("br");
  const [accessKey, setAccessKey] = useState("");
  const [copied, setCopied] = useState(false);
  const register = useRegister();
  const navigate = useNavigate();

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await register.mutateAsync({
      gameName: gameName.trim(),
      tagLine: tagLine.trim(),
      region,
    });
    setAccessKey(res.accessKey);
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(accessKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (accessKey) {
    return (
      <div>
        <div className="mb-4 flex items-center gap-2">
          <span className="flex h-6 w-6 items-center justify-center rounded-full bg-green-500/20 text-xs text-green-400">
            ✓
          </span>
          <h3 className="text-sm font-bold text-white">
            registration complete
          </h3>
        </div>

        <p className="mb-4 text-xs leading-relaxed text-lofi-secondary">
          Save this access key — it will only be shown{" "}
          <span className="font-bold text-lofi-accent">once</span>. Use it to
          log in on any device.
        </p>

        <div className="mb-4 rounded-sm border border-lofi-accent/30 bg-lofi-black p-4">
          <p className="mb-1 text-[10px] uppercase tracking-wider text-lofi-muted">
            your access key
          </p>
          <code className="block break-all text-sm leading-relaxed text-lofi-accent">
            {accessKey}
          </code>
        </div>

        <div className="flex gap-2">
          <Button onClick={handleCopy} className="flex-1">
            {copied ? "copied!" : "copy key"}
          </Button>
          <Button
            variant="secondary"
            className="flex-1"
            onClick={() => navigate("/player")}
          >
            continue →
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h3 className="mb-1 text-sm font-bold text-white">link your riot id</h3>
      <p className="mb-5 text-xs text-lofi-secondary">
        verify your identity to generate a unique access key.
      </p>

      <form onSubmit={handleRegister} className="space-y-4">
        <div>
          <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
            game name
          </label>
          <Input
            value={gameName}
            onChange={(e) => setGameName(e.target.value)}
            placeholder="pai de olivia"
            disabled={register.isPending}
          />
        </div>

        <div className="flex gap-3">
          <div className="flex-1">
            <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
              tag line
            </label>
            <Input
              value={tagLine}
              onChange={(e) => setTagLine(e.target.value)}
              placeholder="NIAS"
              disabled={register.isPending}
            />
          </div>
          <div className="w-28">
            <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
              server
            </label>
            <select
              value={region}
              onChange={(e) => setRegion(e.target.value)}
              disabled={register.isPending}
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

        {register.error && (
          <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
            {friendlyError(register.error)}
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          disabled={register.isPending || !gameName.trim() || !tagLine.trim()}
        >
          {register.isPending ? "verifying riot id..." : "register"}
        </Button>
      </form>

      <div className="mt-4 flex items-center justify-center gap-4 text-[10px] text-lofi-muted">
        <span className="flex items-center gap-1">
          <span className="h-1 w-1 rounded-full bg-green-500" />
          encrypted
        </span>
        <span className="flex items-center gap-1">
          <span className="h-1 w-1 rounded-full bg-green-500" />
          riot api verified
        </span>
      </div>
    </div>
  );
}

function LoginForm() {
  const [accessKey, setAccessKey] = useState("");
  const login = useLogin();
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    await login.mutateAsync(accessKey.trim());
    navigate("/player");
  };

  return (
    <div>
      <h3 className="mb-1 text-sm font-bold text-white">enter access key</h3>
      <p className="mb-5 text-xs text-lofi-secondary">
        paste the access key you received during registration.
      </p>

      <form onSubmit={handleLogin} className="space-y-4">
        <div>
          <label className="mb-1.5 block text-[10px] uppercase tracking-wider text-lofi-muted">
            access key
          </label>
          <Input
            value={accessKey}
            onChange={(e) => setAccessKey(e.target.value)}
            placeholder="paste your access key..."
            disabled={login.isPending}
          />
        </div>

        {login.error && (
          <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
            {friendlyError(login.error)}
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          disabled={login.isPending || !accessKey.trim()}
        >
          {login.isPending ? "authenticating..." : "login"}
        </Button>
      </form>

      <p className="mt-4 text-center text-[10px] text-lofi-muted">
        don't have a key?{" "}
        <span className="text-lofi-accent">switch to register tab above</span>
      </p>
    </div>
  );
}
