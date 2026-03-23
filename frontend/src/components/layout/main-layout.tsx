import { type ReactNode } from "react";
import { Sidebar } from "./sidebar";
import { PageTransition } from "../ui/page-transition";

interface MainLayoutProps {
  children: ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className="flex h-screen bg-lofi-black">
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-6">
        <PageTransition>{children}</PageTransition>
      </main>
    </div>
  );
}
