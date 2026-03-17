import { useAuthStore } from "@/stores/auth-store";
import { useSettingsStore } from "@/stores/settings-store";
import { SectionHeader } from "@/components/layout/section-header";
import { Button } from "@/components/ui/button";

const languages = [
  { value: "en", label: "English" },
  { value: "pt_br", label: "Portugues (BR)" },
  { value: "es", label: "Espanol" },
  { value: "ko", label: "Korean" },
  { value: "ja", label: "Japanese" },
];

export function SettingsPage() {
  const { token, user, logout } = useAuthStore();
  const { theme, language, alerts, setTheme, setLanguage, setAlert } =
    useSettingsStore();

  return (
    <div>
      <SectionHeader number="06" title="settings" />

      <div className="space-y-6">
        {/* Account */}
        <Section title="account" number="06.1">
          {token && user ? (
            <div className="space-y-3">
              <Row label="riot id">
                <span className="text-xs text-white">
                  {user.gameName}
                  <span className="text-lofi-muted">#{user.tagLine}</span>
                </span>
              </Row>
              <Row label="region">
                <span className="text-xs text-white uppercase">
                  {user.region}
                </span>
              </Row>
              <Row label="session">
                <div className="flex items-center gap-2">
                  <span className="h-1.5 w-1.5 rounded-full bg-green-500" />
                  <span className="text-xs text-green-400">active</span>
                </div>
              </Row>
              <div className="pt-2">
                <Button
                  variant="secondary"
                  onClick={logout}
                  className="text-red-400 hover:text-red-300"
                >
                  logout
                </Button>
              </div>
            </div>
          ) : (
            <p className="text-xs text-lofi-muted">
              not authenticated. visit{" "}
              <a href="/auth" className="text-lofi-accent hover:underline">
                login / register
              </a>{" "}
              to connect your account.
            </p>
          )}
        </Section>

        {/* Interface */}
        <Section title="interface" number="06.2">
          <div className="space-y-3">
            <Row label="theme">
              <div className="flex gap-2">
                <ToggleButton
                  active={theme === "dark"}
                  onClick={() => setTheme("dark")}
                >
                  dark
                </ToggleButton>
                <ToggleButton
                  active={theme === "light"}
                  onClick={() => setTheme("light")}
                >
                  light
                </ToggleButton>
              </div>
            </Row>
            <Row label="language">
              <select
                value={language}
                onChange={(e) => setLanguage(e.target.value)}
                className="rounded-sm border border-lofi-border bg-lofi-black px-3 py-1.5 text-xs text-lofi-text focus:border-lofi-accent focus:outline-none focus:ring-1 focus:ring-lofi-accent"
              >
                {languages.map((l) => (
                  <option key={l.value} value={l.value}>
                    {l.label}
                  </option>
                ))}
              </select>
            </Row>
          </div>
        </Section>

        {/* Oracle Alerts */}
        <Section title="oracle alerts" number="06.3">
          <div className="space-y-3">
            <ToggleRow
              label="meta shift alerts"
              desc="notify when win rates change >5%"
              checked={alerts.metaShift}
              onChange={(v) => setAlert("metaShift", v)}
            />
            <ToggleRow
              label="scouting assistant"
              desc="real-time lobby composition hints"
              checked={alerts.scouting}
              onChange={(v) => setAlert("scouting", v)}
            />
            <ToggleRow
              label="match ready notifications"
              desc="alert when match history updates"
              checked={alerts.matchReady}
              onChange={(v) => setAlert("matchReady", v)}
            />
          </div>
          <p className="mt-3 text-[10px] text-lofi-muted">
            alerts require backend integration — coming in a future update.
          </p>
        </Section>

        {/* Terminal Actions */}
        <Section title="terminal actions" number="06.4">
          <div className="space-y-2">
            <Button
              variant="secondary"
              onClick={() => {
                localStorage.removeItem("tft-oracle-settings");
                window.location.reload();
              }}
            >
              reset all preferences
            </Button>
            <p className="text-[10px] text-lofi-muted">
              clears local settings and reloads the app with defaults.
            </p>
          </div>
        </Section>

        {/* System Info */}
        <Section title="system" number="06.5">
          <div className="space-y-1">
            <Row label="version">
              <span className="text-xs text-lofi-text">v0.2.0</span>
            </Row>
            <Row label="phase">
              <span className="text-xs text-lofi-text">2 — analytics</span>
            </Row>
            <Row label="data sources">
              <span className="text-xs text-lofi-text">
                communitydragon + riot api
              </span>
            </Row>
          </div>
        </Section>
      </div>
    </div>
  );
}

function Section({
  title,
  number,
  children,
}: {
  title: string;
  number: string;
  children: React.ReactNode;
}) {
  return (
    <div className="rounded-sm border border-lofi-border bg-lofi-surface p-4">
      <p className="mb-3 text-[10px] uppercase tracking-wider text-lofi-muted">
        [{number}] {title}
      </p>
      {children}
    </div>
  );
}

function Row({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-xs text-lofi-secondary">{label}</span>
      {children}
    </div>
  );
}

function ToggleButton({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className={`rounded-sm border px-3 py-1 text-[10px] font-medium transition-colors ${
        active
          ? "border-lofi-accent bg-lofi-accent/10 text-lofi-accent"
          : "border-lofi-border text-lofi-muted hover:text-lofi-secondary"
      }`}
    >
      {children}
    </button>
  );
}

function ToggleRow({
  label,
  desc,
  checked,
  onChange,
}: {
  label: string;
  desc: string;
  checked: boolean;
  onChange: (value: boolean) => void;
}) {
  return (
    <div className="flex items-center justify-between">
      <div>
        <p className="text-xs text-lofi-secondary">{label}</p>
        <p className="text-[10px] text-lofi-muted">{desc}</p>
      </div>
      <button
        onClick={() => onChange(!checked)}
        className={`relative h-5 w-9 rounded-full transition-colors ${
          checked ? "bg-lofi-accent" : "bg-lofi-border"
        }`}
      >
        <span
          className={`absolute top-0.5 h-4 w-4 rounded-full bg-white transition-transform ${
            checked ? "left-[18px]" : "left-0.5"
          }`}
        />
      </button>
    </div>
  );
}
