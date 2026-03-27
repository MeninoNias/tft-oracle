import { NavLink } from "react-router-dom";
import { useAuthStore } from "@/stores/auth-store";

const navItems = [
  { to: "/coach", label: "ai coach" },
  { to: "/simulator", label: "simulator" },
  { to: "/tier-list", label: "tier list" },
  { to: "/champions", label: "champions" },
  { to: "/items", label: "items" },
  { to: "/augments", label: "augments" },
  { to: "/player", label: "player search" },
  { to: "/settings", label: "settings" },
];

export function Sidebar() {
  const { token, user, logout } = useAuthStore();
  const isLoggedIn = !!token;

  return (
    <aside className="flex h-screen w-64 flex-col border-r border-lofi-border bg-lofi-black">
      {/* Logo */}
      <div className="border-b border-lofi-border px-4 py-4">
        <h1 className="text-sm font-bold tracking-wider text-white">
          TFT_ORACLE
        </h1>
        <p className="mt-0.5 text-[10px] text-lofi-muted">
          v0.4.0 // phase 4
        </p>
      </div>

      {/* Active User Card */}
      {isLoggedIn && user ? (
        <div className="border-b border-lofi-border px-3 py-3">
          <NavLink to="/profile" className="block">
            <div className="rounded-sm border border-lofi-border bg-lofi-surface px-3 py-2.5 transition-colors hover:border-lofi-muted">
              <p className="text-[10px] font-medium uppercase tracking-wider text-lofi-muted">
                active user
              </p>
              <p className="mt-1 text-sm font-bold tracking-wide text-white">
                {user.gameName}
              </p>
              <div className="mt-0.5 flex items-center gap-2">
                <span className="text-[10px] text-lofi-muted">
                  #{user.tagLine}
                </span>
                <span className="text-[10px] text-lofi-muted">
                  {user.region.toUpperCase()}
                </span>
              </div>
            </div>
          </NavLink>
        </div>
      ) : (
        <div className="border-b border-lofi-border px-3 py-3">
          <NavLink to="/auth" className="block">
            <div className="rounded-sm border border-dashed border-lofi-muted px-3 py-3 text-center transition-colors hover:border-lofi-accent">
              <p className="text-xs font-medium text-lofi-accent">
                login / register
              </p>
              <p className="mt-0.5 text-[10px] text-lofi-muted">
                sign in to use ai coach
              </p>
            </div>
          </NavLink>
        </div>
      )}

      {/* Navigation */}
      <nav className="flex-1 px-2 py-3">
        <p className="mb-1 px-3 text-[9px] font-medium uppercase tracking-widest text-lofi-muted">
          navigation
        </p>
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `flex items-center gap-2 rounded-sm px-3 py-2 text-xs transition-colors ${
                isActive
                  ? "bg-lofi-surface text-white"
                  : "text-lofi-secondary hover:bg-lofi-surface hover:text-lofi-text"
              }`
            }
          >
            {({ isActive }) => (
              <>
                <span
                  className={`h-1.5 w-1.5 rounded-full ${
                    isActive ? "bg-lofi-accent" : "bg-lofi-muted"
                  }`}
                />
                {item.label}
              </>
            )}
          </NavLink>
        ))}
      </nav>

      <div className="border-t border-lofi-border px-4 py-3">
        <p className="text-[10px] text-lofi-muted">
          data: communitydragon + riot api + metatft + tftactics + mobalytics
        </p>
        {isLoggedIn && (
          <button
            onClick={logout}
            className="mt-2 text-[10px] text-lofi-muted hover:text-red-400"
          >
            logout
          </button>
        )}
      </div>
    </aside>
  );
}
