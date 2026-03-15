import { TerminalCard } from "@/components/ui/terminal-card";
import type {
  Account,
  Summoner,
  RankedEntry,
} from "@/gen/tft/v1/player_pb";

interface ProfileHeaderProps {
  account: Account;
  summoner?: Summoner;
  rankedEntries: RankedEntry[];
}

const tierColors: Record<string, string> = {
  IRON: "text-gray-400",
  BRONZE: "text-amber-700",
  SILVER: "text-gray-300",
  GOLD: "text-yellow-400",
  PLATINUM: "text-teal-400",
  EMERALD: "text-emerald-400",
  DIAMOND: "text-blue-400",
  MASTER: "text-purple-400",
  GRANDMASTER: "text-red-400",
  CHALLENGER: "text-yellow-300",
};

export function ProfileHeader({
  account,
  summoner,
  rankedEntries,
}: ProfileHeaderProps) {
  const ranked = rankedEntries.find((e) => e.queueType === "RANKED_TFT");
  const totalGames = ranked ? ranked.wins + ranked.losses : 0;
  const winRate =
    totalGames > 0 ? ((ranked!.wins / totalGames) * 100).toFixed(1) : "0";

  return (
    <TerminalCard className="border-l-2 border-l-lofi-accent">
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-baseline gap-2">
            <h2 className="text-sm font-bold text-white">
              {account.gameName}
            </h2>
            <span className="text-[10px] text-lofi-muted">
              #{account.tagLine}
            </span>
          </div>
          {summoner && (
            <p className="mt-0.5 text-[10px] text-lofi-muted">
              level {Number(summoner.summonerLevel)}
            </p>
          )}
        </div>

        {ranked && (
          <div className="text-right">
            <p
              className={`text-sm font-bold ${tierColors[ranked.tier] ?? "text-lofi-text"}`}
            >
              {ranked.tier} {ranked.rank}
            </p>
            <p className="text-[10px] text-lofi-secondary">
              {ranked.leaguePoints} LP
            </p>
            <p className="mt-0.5 text-[10px] text-lofi-muted">
              {ranked.wins}W {ranked.losses}L ({winRate}%)
            </p>
          </div>
        )}

        {!ranked && (
          <p className="text-[10px] text-lofi-muted">unranked</p>
        )}
      </div>
    </TerminalCard>
  );
}
