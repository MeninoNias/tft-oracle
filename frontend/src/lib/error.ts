import { ConnectError, Code } from "@connectrpc/connect";

const FRIENDLY_MESSAGES: Partial<Record<Code, string>> = {
  [Code.Unauthenticated]: "session expired or invalid access key. please log in again.",
  [Code.PermissionDenied]: "you don't have permission to perform this action.",
  [Code.NotFound]: "the requested data could not be found.",
  [Code.AlreadyExists]: "this riot id is already registered. use your access key to log in.",
  [Code.Unavailable]: "server is temporarily unavailable. try again in a moment.",
  [Code.Internal]: "an unexpected server error occurred. please try again.",
  [Code.InvalidArgument]: "invalid input — please check your fields and try again.",
  [Code.DeadlineExceeded]: "request timed out. please try again.",
};

/**
 * Extracts a user-friendly message from a Connect RPC error or generic Error.
 */
export function friendlyError(err: Error | null): string {
  if (!err) return "an unknown error occurred.";

  if (err instanceof ConnectError) {
    // Try to extract nested status message from the JSON payload
    const nested = parseNestedMessage(err.rawMessage);
    if (nested) return nested;

    return FRIENDLY_MESSAGES[err.code] ?? err.rawMessage ?? err.message;
  }

  return err.message;
}

function parseNestedMessage(raw: string): string | null {
  try {
    // Handle: unauthorized: {"status":{"message":"Unknown apikey","status_code":401}}
    const jsonStart = raw.indexOf("{");
    if (jsonStart === -1) return null;
    const parsed = JSON.parse(raw.slice(jsonStart));
    if (parsed?.status?.message) {
      return parsed.status.message.toLowerCase() + ".";
    }
    if (parsed?.message) {
      return parsed.message.toLowerCase() + ".";
    }
  } catch {
    // not JSON, ignore
  }
  return null;
}
