import { Button } from "@/components/ui/button";

interface InlineErrorProps {
  message: string;
  onRetry?: () => void;
}

export function InlineError({ message, onRetry }: InlineErrorProps) {
  return (
    <div className="flex items-center justify-between rounded-sm border border-red-500/30 bg-red-500/10 px-4 py-3">
      <div className="flex items-center gap-2">
        <span className="text-xs text-red-400">!</span>
        <span className="text-xs text-red-400">{message}</span>
      </div>
      {onRetry && (
        <Button variant="secondary" onClick={onRetry} className="ml-4 text-red-400 border-red-500/30 hover:bg-red-500/10">
          retry
        </Button>
      )}
    </div>
  );
}
