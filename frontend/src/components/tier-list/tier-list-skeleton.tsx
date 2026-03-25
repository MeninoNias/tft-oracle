import { Skeleton } from "@/components/ui/skeleton";

export function TierListSkeleton() {
  return (
    <div className="space-y-6">
      {["S", "A", "B", "C"].map((tier) => (
        <div key={tier} className="border-l-4 border-l-lofi-border pl-3">
          <Skeleton className="mb-3 h-6 w-16" />
          <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
            {Array.from({ length: tier === "S" ? 2 : 3 }).map((_, i) => (
              <Skeleton key={i} className="h-36" />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
