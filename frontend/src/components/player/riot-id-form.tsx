import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { SERVERS } from "@/lib/servers";

interface RiotIdFormProps {
  onSubmit: (gameName: string, tagLine: string, region: string) => void;
  isLoading?: boolean;
}

export function RiotIdForm({ onSubmit, isLoading }: RiotIdFormProps) {
  const [gameName, setGameName] = useState("");
  const [tagLine, setTagLine] = useState("");
  const [region, setRegion] = useState("br");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (gameName.trim() && tagLine.trim()) {
      onSubmit(gameName.trim(), tagLine.trim(), region);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-wrap items-end gap-3">
      <div className="min-w-[140px] flex-1">
        <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
          game name
        </label>
        <Input
          value={gameName}
          onChange={(e) => setGameName(e.target.value)}
          placeholder="MeninoNias"
          disabled={isLoading}
        />
      </div>
      <div className="w-24">
        <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
          tag line
        </label>
        <Input
          value={tagLine}
          onChange={(e) => setTagLine(e.target.value)}
          placeholder="BR1"
          disabled={isLoading}
        />
      </div>
      <div className="w-24">
        <label className="mb-1 block text-[10px] uppercase tracking-wider text-lofi-muted">
          server
        </label>
        <select
          value={region}
          onChange={(e) => setRegion(e.target.value)}
          disabled={isLoading}
          className="w-full rounded-sm border border-lofi-border bg-lofi-black px-3 py-1.5 text-xs text-lofi-text focus:border-lofi-accent focus:outline-none focus:ring-1 focus:ring-lofi-accent"
        >
          {SERVERS.map((r) => (
            <option key={r.value} value={r.value}>
              {r.label}
            </option>
          ))}
        </select>
      </div>
      <Button
        type="submit"
        disabled={isLoading || !gameName.trim() || !tagLine.trim()}
      >
        {isLoading ? "loading..." : "search"}
      </Button>
    </form>
  );
}
