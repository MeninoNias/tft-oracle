import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { SectionHeader } from "@/components/layout/section-header";
import { TerminalCard } from "@/components/ui/terminal-card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useRegister, useLogin } from "@/hooks/use-auth";

const servers = [
  { value: "br", label: "BR" },
  { value: "na", label: "NA" },
  { value: "euw", label: "EUW" },
  { value: "eune", label: "EUNE" },
  { value: "kr", label: "KR" },
  { value: "jp", label: "JP" },
  { value: "oce", label: "OCE" },
  { value: "lan", label: "LAN" },
  { value: "las", label: "LAS" },
  { value: "tr", label: "TR" },
  { value: "ru", label: "RU" },
];

export function AuthPage() {
  const [tab, setTab] = useState<"register" | "login">("register");

  return (
    <div>
      <SectionHeader number="00" title="authenticate" />

      <div className="mb-4 flex gap-2">
        <Button
          variant={tab === "register" ? "primary" : "secondary"}
          onClick={() => setTab("register")}
        >
          register
        </Button>
        <Button
          variant={tab === "login" ? "primary" : "secondary"}
          onClick={() => setTab("login")}
        >
          login
        </Button>
      </div>

      {tab === "register" ? <RegisterForm /> : <LoginForm />}
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
      <TerminalCard className="max-w-lg border-l-2 border-l-lofi-accent">
        <h3 className="mb-2 text-sm font-bold text-white">
          registration complete
        </h3>
        <p className="mb-3 text-xs text-lofi-secondary">
          Save this access key — it will only be shown once. Use it to log in
          on any device.
        </p>

        <div className="mb-3 rounded-sm border border-lofi-border bg-lofi-black p-3">
          <code className="break-all text-xs text-lofi-accent">{accessKey}</code>
        </div>

        <div className="flex gap-2">
          <Button onClick={handleCopy}>
            {copied ? "copied!" : "copy key"}
          </Button>
          <Button
            variant="secondary"
            onClick={() => {
              // Auto-login after registration
              navigate("/player");
            }}
          >
            continue
          </Button>
        </div>

        <p className="mt-3 text-[10px] text-lofi-muted">
          you can now use this key to log in. keep it safe.
        </p>
      </TerminalCard>
    );
  }

  return (
    <TerminalCard className="max-w-lg">
      <h3 className="mb-1 text-sm font-bold text-white">
        link your riot id
      </h3>
      <p className="mb-4 text-xs text-lofi-secondary">
        verify your identity by entering your Riot ID. we'll generate a unique
        access key for you.
      </p>

      <form onSubmit={handleRegister} className="space-y-3">
        <div>
          <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
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
            <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
              tag line
            </label>
            <Input
              value={tagLine}
              onChange={(e) => setTagLine(e.target.value)}
              placeholder="NIAS"
              disabled={register.isPending}
            />
          </div>
          <div className="w-24">
            <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
              server
            </label>
            <select
              value={region}
              onChange={(e) => setRegion(e.target.value)}
              disabled={register.isPending}
              className="w-full rounded-sm border border-lofi-border bg-lofi-black px-3 py-1.5 text-xs text-lofi-text focus:border-lofi-accent focus:outline-none focus:ring-1 focus:ring-lofi-accent"
            >
              {servers.map((s) => (
                <option key={s.value} value={s.value}>
                  {s.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        {register.error && (
          <div className="rounded-sm border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-400">
            {register.error.message}
          </div>
        )}

        <Button
          type="submit"
          disabled={
            register.isPending || !gameName.trim() || !tagLine.trim()
          }
        >
          {register.isPending ? "verifying..." : "register"}
        </Button>
      </form>
    </TerminalCard>
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
    <TerminalCard className="max-w-lg">
      <h3 className="mb-1 text-sm font-bold text-white">
        enter access key
      </h3>
      <p className="mb-4 text-xs text-lofi-secondary">
        paste the access key you received during registration.
      </p>

      <form onSubmit={handleLogin} className="space-y-3">
        <div>
          <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
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
            {login.error.message}
          </div>
        )}

        <Button
          type="submit"
          disabled={login.isPending || !accessKey.trim()}
        >
          {login.isPending ? "authenticating..." : "login"}
        </Button>
      </form>
    </TerminalCard>
  );
}
