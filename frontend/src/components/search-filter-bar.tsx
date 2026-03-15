import { Input } from "./ui/input";
import { Button } from "./ui/button";

interface FilterOption {
  label: string;
  value: string;
}

interface SearchFilterBarProps {
  searchValue: string;
  onSearchChange: (value: string) => void;
  searchPlaceholder?: string;
  filters?: FilterOption[];
  activeFilter: string;
  onFilterChange: (value: string) => void;
}

export function SearchFilterBar({
  searchValue,
  onSearchChange,
  searchPlaceholder = "search...",
  filters = [],
  activeFilter,
  onFilterChange,
}: SearchFilterBarProps) {
  return (
    <div className="mb-4 flex flex-wrap items-center gap-2">
      <Input
        value={searchValue}
        onChange={(e) => onSearchChange(e.target.value)}
        placeholder={searchPlaceholder}
        className="max-w-xs"
      />
      {filters.length > 0 && (
        <div className="flex gap-1">
          <Button
            variant={activeFilter === "all" ? "primary" : "secondary"}
            onClick={() => onFilterChange("all")}
          >
            all
          </Button>
          {filters.map((f) => (
            <Button
              key={f.value}
              variant={activeFilter === f.value ? "primary" : "secondary"}
              onClick={() => onFilterChange(f.value)}
            >
              {f.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}
