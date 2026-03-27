import { useState, useEffect, useCallback, createContext, useContext } from "react";

type ToastType = "success" | "error" | "info" | "warning";

interface Toast {
  id: number;
  message: string;
  type: ToastType;
  duration: number;
  action?: { label: string; onClick: () => void };
}

interface ToastOptions {
  type?: ToastType;
  duration?: number;
  action?: { label: string; onClick: () => void };
}

interface ToastContextValue {
  toast: (message: string, options?: ToastType | ToastOptions) => void;
}

const ToastContext = createContext<ToastContextValue>({
  toast: () => {},
});

export function useToast() {
  return useContext(ToastContext);
}

const DEFAULT_DURATION = 3000;
let nextId = 0;

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback(
    (message: string, options?: ToastType | ToastOptions) => {
      const id = nextId++;
      const opts: ToastOptions =
        typeof options === "string" ? { type: options } : options ?? {};
      setToasts((prev) => [
        ...prev,
        {
          id,
          message,
          type: opts.type ?? "info",
          duration: opts.duration ?? DEFAULT_DURATION,
          action: opts.action,
        },
      ]);
    },
    [],
  );

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ toast: addToast }}>
      {children}
      <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
        {toasts.map((t) => (
          <ToastItem key={t.id} toast={t} onDismiss={removeToast} />
        ))}
      </div>
    </ToastContext.Provider>
  );
}

const typeStyles: Record<ToastType, string> = {
  success: "border-green-500/30 bg-green-500/10 text-green-400",
  error: "border-red-500/30 bg-red-500/10 text-red-400",
  info: "border-lofi-accent/30 bg-lofi-accent/10 text-lofi-accent",
  warning: "border-yellow-500/30 bg-yellow-500/10 text-yellow-400",
};

function ToastItem({
  toast,
  onDismiss,
}: {
  toast: Toast;
  onDismiss: (id: number) => void;
}) {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    // Trigger enter animation
    requestAnimationFrame(() => setVisible(true));

    const timer = setTimeout(() => {
      setVisible(false);
      setTimeout(() => onDismiss(toast.id), 200);
    }, toast.duration);

    return () => clearTimeout(timer);
  }, [toast.id, toast.duration, onDismiss]);

  return (
    <div
      className={`flex items-center gap-3 rounded-sm border px-4 py-2 text-xs font-medium shadow-lg backdrop-blur-sm transition-all duration-200 ${typeStyles[toast.type]} ${
        visible
          ? "translate-x-0 opacity-100"
          : "translate-x-4 opacity-0"
      }`}
    >
      <span>{toast.message}</span>
      {toast.action && (
        <button
          onClick={() => {
            toast.action!.onClick();
            onDismiss(toast.id);
          }}
          className="ml-auto shrink-0 underline underline-offset-2 hover:opacity-80"
        >
          {toast.action.label}
        </button>
      )}
    </div>
  );
}
