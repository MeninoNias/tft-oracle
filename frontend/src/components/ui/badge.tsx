import { type ReactNode } from "react";

interface BadgeProps {
  children: ReactNode;
  className?: string;
}

export function Badge({ children, className = "" }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center rounded-sm border border-lofi-border px-1.5 py-0.5 text-[10px] uppercase tracking-wider text-lofi-secondary ${className}`}
    >
      {children}
    </span>
  );
}
