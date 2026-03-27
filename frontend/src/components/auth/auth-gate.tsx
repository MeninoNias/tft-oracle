import { type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "@/stores/auth-store";
import { GenericEmptyState } from "@/components/ui/generic-empty-state";

interface AuthGateProps {
  children: ReactNode;
}

export function AuthGate({ children }: AuthGateProps) {
  const token = useAuthStore((s) => s.token);
  const navigate = useNavigate();

  if (!token) {
    return (
      <GenericEmptyState
        title="authentication required"
        description="you need to be logged in to access this feature. sign in with your riot id to continue."
        actionLabel="login / register"
        onAction={() => navigate("/auth")}
      />
    );
  }

  return <>{children}</>;
}
