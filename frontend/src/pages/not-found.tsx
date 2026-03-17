import { ErrorScreen } from "@/components/error/error-screen";

export function NotFoundPage() {
  return (
    <ErrorScreen
      code="404"
      title="NO_DATA_FOUND"
      statusLabel="DISCONNECTED"
      lines={[
        "Critical failure in data retrieval.",
        "Remote uplink terminated by host.",
        "System pulse: OFFLINE.",
      ]}
      fadedLine="Attempting reconnection... [FAILED]"
      transmissionId="TFT-ORCL-ERR-404"
    />
  );
}
