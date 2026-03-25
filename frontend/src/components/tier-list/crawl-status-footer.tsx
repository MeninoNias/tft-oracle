interface CrawlStatusFooterProps {
  lastUpdated: string;
}

export function CrawlStatusFooter({ lastUpdated }: CrawlStatusFooterProps) {
  const formatted = lastUpdated
    ? formatRelativeTime(lastUpdated)
    : "never";

  return (
    <div className="mt-6 border-t border-lofi-border pt-3 text-[10px] text-lofi-muted">
      last updated: {formatted} // sources: metatft + tftactics + mobalytics
    </div>
  );
}

function formatRelativeTime(iso: string): string {
  const date = new Date(iso);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMin = Math.floor(diffMs / 60000);

  if (diffMin < 1) return "just now";
  if (diffMin < 60) return `${diffMin}m ago`;

  const diffH = Math.floor(diffMin / 60);
  if (diffH < 24) return `${diffH}h ago`;

  const diffD = Math.floor(diffH / 24);
  return `${diffD}d ago`;
}
