import { type InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {}

export function Input({ className = "", ...props }: InputProps) {
  return (
    <input
      className={`w-full rounded-sm border border-lofi-border bg-lofi-black px-3 py-1.5 text-xs text-lofi-text placeholder:text-lofi-muted focus:border-lofi-accent focus:outline-none focus:ring-1 focus:ring-lofi-accent ${className}`}
      {...props}
    />
  );
}
