import { type ReactNode } from "react";

interface TerminalCardProps {
  children: ReactNode;
  className?: string;
}

export function TerminalCard({ children, className = "" }: TerminalCardProps) {
  return (
    <div
      className={`rounded-sm border border-lofi-border bg-lofi-surface p-4 transition-colors duration-150 hover:border-lofi-muted ${className}`}
    >
      {children}
    </div>
  );
}
