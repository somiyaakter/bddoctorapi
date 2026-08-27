import type { Location, Specialty } from "./types";

export function locationSlugFromUrl(url: string): string {
  const trimmed = url.replace(/\/$/, "");
  const idx = trimmed.lastIndexOf("doctors-");
  if (idx === -1) return "";
  return trimmed.slice(idx + "doctors-".length).toLowerCase();
}

export function specialtySlugFromUrl(url: string, locationSlug: string): string {
  const trimmed = url.replace(/\/$/, "");
  const idx = trimmed.lastIndexOf("/");
  if (idx === -1) return "";
  const slug = trimmed.slice(idx + 1);
  return slug.endsWith(`-${locationSlug}`) ? slug.slice(0, -(`-${locationSlug}`.length)) : slug;
}

export function findLocationBySlug(locations: Location[], slug: string): Location | undefined {
  return locations.find((loc) => locationSlugFromUrl(loc.url) === slug.toLowerCase());
}

export function findSpecialtyBySlug(
  specialties: Specialty[],
  locations: Location[],
  slug: string,
  locationId?: number
): Specialty | undefined {
  const pool = locationId ? specialties.filter((s) => s.location_id === locationId) : specialties;
  return pool.find((sp) => {
    const loc = locations.find((l) => l.id === sp.location_id);
    if (!loc) return false;
    return specialtySlugFromUrl(sp.url, locationSlugFromUrl(loc.url)) === slug.toLowerCase();
  });
}

export function isDoctorId(segment: string): boolean {
  return /^\d+$/.test(segment);
}