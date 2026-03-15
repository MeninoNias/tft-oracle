interface SectionHeaderProps {
  number: string;
  title: string;
}

export function SectionHeader({ number, title }: SectionHeaderProps) {
  return (
    <div className="mb-6 flex items-center gap-3">
      <span className="text-xs text-lofi-muted">{number}</span>
      <h2 className="text-sm font-medium text-white">{title}</h2>
      <div className="h-px flex-1 bg-lofi-border" />
    </div>
  );
}
