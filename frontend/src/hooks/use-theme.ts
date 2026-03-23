import { useEffect } from "react";
import { useSettingsStore } from "@/stores/settings-store";

/**
 * Applies the current theme from settings store to the <html> element.
 * Call once at the app root.
 */
export function useTheme() {
  const theme = useSettingsStore((s) => s.theme);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);
}
