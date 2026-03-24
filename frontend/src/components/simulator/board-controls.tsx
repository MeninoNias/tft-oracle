import { useBoardStore } from "@/stores/board-store";
import { Button } from "@/components/ui/button";

export function BoardControls() {
  const { activeTab, setActiveTab, clearBoard, resetAll } = useBoardStore();
  const level = useBoardStore((s) => s.getLevel());
  const setLevel = useBoardStore((s) => s.setLevel);

  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      {/* Tab switcher */}
      <div className="flex rounded-sm border border-lofi-border">
        <button
          onClick={() => setActiveTab("player")}
          className={`px-3 py-1.5 text-[10px] font-medium transition-colors ${
            activeTab === "player"
              ? "bg-lofi-surface text-white"
              : "text-lofi-muted hover:text-lofi-secondary"
          }`}
        >
          your board
        </button>
        <button
          onClick={() => setActiveTab("opponent")}
          className={`px-3 py-1.5 text-[10px] font-medium transition-colors ${
            activeTab === "opponent"
              ? "bg-lofi-surface text-white"
              : "text-lofi-muted hover:text-lofi-secondary"
          }`}
        >
          opponent
        </button>
      </div>

      {/* Level selector */}
      <div className="flex items-center gap-1.5">
        <span className="text-[10px] text-lofi-muted">level:</span>
        <button
          onClick={() => setLevel(level - 1)}
          className="flex h-6 w-6 items-center justify-center rounded-sm border border-lofi-border text-xs text-lofi-secondary hover:bg-lofi-surface"
        >
          -
        </button>
        <span className="w-5 text-center text-xs font-bold text-white">
          {level}
        </span>
        <button
          onClick={() => setLevel(level + 1)}
          className="flex h-6 w-6 items-center justify-center rounded-sm border border-lofi-border text-xs text-lofi-secondary hover:bg-lofi-surface"
        >
          +
        </button>
      </div>

      {/* Actions */}
      <div className="ml-auto flex gap-2">
        <Button variant="secondary" onClick={clearBoard}>
          clear
        </Button>
        <Button variant="secondary" onClick={resetAll}>
          reset all
        </Button>
      </div>
    </div>
  );
}
