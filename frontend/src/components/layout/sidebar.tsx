import { NavLink } from "react-router-dom";
import { useAuthStore } from "@/stores/auth-store";

const navItems = [
  { to: "/champions", label: "champions" },
  { to: "/items", label: "items" },
  { to: "/profile", label: "profile" },
  { to: "/player", label: "player search" },
  { to: "/settings", label: "settings" },
];

export function Sidebar() {
  const { token, user, logout } = useAuthStore();
  const isLoggedIn = !!token;

  return (
    <aside className="flex h-screen w-64 flex-col border-r border-lofi-border bg-lofi-black">
      <div className="border-b border-lofi-border px-4 py-4">
        <h1 className="text-sm font-bold tracking-wider text-white">
          TFT_ORACLE
        </h1>
        <p className="mt-0.5 text-[10px] text-lofi-muted">
          v0.2.0 // phase 2
        </p>
      </div>

      <nav className="flex-1 px-2 py-3">
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
        {isLoggedIn && user ? (
          <div className="space-y-1">
            <p className="text-xs text-lofi-text">
              {user.gameName}
              <span className="text-lofi-muted">#{user.tagLine}</span>
            </p>
            <button
              onClick={logout}
              className="text-[10px] text-lofi-muted hover:text-red-400"
            >
              logout
            </button>
          </div>
        ) : (
          <NavLink
            to="/auth"
            className="text-[10px] text-lofi-accent hover:text-lofi-text"
          >
            login / register
          </NavLink>
        )}
      </div>

      <div className="border-t border-lofi-border px-4 py-3">
        <p className="text-[10px] text-lofi-muted">
          data: communitydragon + riot api
        </p>
      </div>
    </aside>
  );
}
