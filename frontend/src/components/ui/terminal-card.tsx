import { type ReactNode } from "react";

interface TerminalCardProps {
  children: ReactNode;
  className?: string;
}

export function TerminalCard({ children, className = "" }: TerminalCardProps) {
  return (
    <div
      className={`rounded-sm border border-lofi-border bg-lofi-surface p-4 ${className}`}
    >
      {children}
    </div>
  );
}
