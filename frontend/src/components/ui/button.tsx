import { type ButtonHTMLAttributes } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary";
}

export function Button({
  variant = "primary",
  className = "",
  children,
  ...props
}: ButtonProps) {
  const base =
    "inline-flex items-center justify-center rounded-sm px-3 py-1.5 text-xs font-medium transition-all focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-lofi-accent disabled:pointer-events-none disabled:opacity-50 active:scale-95";

  const variants = {
    primary: "bg-lofi-text text-lofi-black hover:opacity-80",
    secondary:
      "border border-lofi-border bg-transparent text-lofi-text hover:bg-lofi-surface",
  };

  return (
    <button className={`${base} ${variants[variant]} ${className}`} {...props}>
      {children}
    </button>
  );
}
