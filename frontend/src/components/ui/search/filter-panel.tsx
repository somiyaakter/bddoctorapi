"use client";

import { useRouter, useSearchParams, usePathname } from "next/navigation";
import type { Location, Specialty } from "@/lib/types";

interface FilterPanelProps {
  locations: Location[];
  specialties: Specialty[];
  selectedLocationId?: number;
  selectedSpecialtyId?: number;
}

export function FilterPanel({
  locations,
  specialties,
  selectedLocationId,
  selectedSpecialtyId,
}: FilterPanelProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  function updateFilter(key: "location_id" | "specialty_id", value: string) {
    const params = new URLSearchParams(searchParams.toString());
    if (value) {
      params.set(key, value);
    } else {
      params.delete(key);
    }
    // Changing location invalidates any previously selected specialty from a different location
    if (key === "location_id") {
      params.delete("specialty_id");
    }
    params.delete("page");
    router.push(`${pathname}?${params.toString()}`);
  }

  return (
    <div className="space-y-6 rounded-lg border border-hairline bg-white p-4">
      <div>
        <label htmlFor="location-filter" className="text-sm font-medium text-ink">
          Location
        </label>
        <select
          id="location-filter"
          value={selectedLocationId ?? ""}
          onChange={(e) => updateFilter("location_id", e.target.value)}
          className="mt-1.5 w-full rounded-md border border-hairline bg-white px-3 py-2 text-sm text-ink focus:border-brand focus:outline-none"
        >
          <option value="">All locations</option>
          {locations.map((loc) => (
            <option key={loc.id} value={loc.id}>
              {loc.name}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label htmlFor="specialty-filter" className="text-sm font-medium text-ink">
          Specialty
        </label>
        <select
          id="specialty-filter"
          value={selectedSpecialtyId ?? ""}
          onChange={(e) => updateFilter("specialty_id", e.target.value)}
          className="mt-1.5 w-full rounded-md border border-hairline bg-white px-3 py-2 text-sm text-ink focus:border-brand focus:outline-none"
        >
          <option value="">All specialties</option>
          {specialties.map((sp) => (
            <option key={sp.id} value={sp.id}>
              {sp.name}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
}