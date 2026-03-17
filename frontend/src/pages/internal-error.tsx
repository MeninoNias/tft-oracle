import { ErrorScreen } from "@/components/error/error-screen";

export function InternalErrorPage() {
  return (
    <ErrorScreen
      code="500"
      title="SYSTEM_FAILURE"
      statusLabel="CRITICAL"
      lines={[
        "Internal process terminated unexpectedly.",
        "Core module returned non-zero exit code.",
        "Tactical subsystem: UNRESPONSIVE.",
      ]}
      fadedLine="Initiating emergency protocol... [STANDBY]"
      transmissionId="TFT-ORCL-ERR-500"
    />
  );
}
