import type {
  ApiListResponse,
  ApiSingleResponse,
  Doctor,
  Location,
  Specialty,
} from "./types";

const API_BASE_URL = process.env.API_BASE_URL;
const API_KEY = process.env.API_KEY;

export class ApiRequestError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

async function apiFetch<T>(path: string, revalidateSeconds: number): Promise<T> {
  const res = await fetch(`${API_BASE_URL}/api/v1${path}`, {
    headers: { "X-API-Key": API_KEY ?? "" },
    next: { revalidate: revalidateSeconds },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: "Unknown error" }));
    throw new ApiRequestError(body.error ?? `Request failed (${res.status})`, res.status);
  }

  return res.json();
}

export interface ListDoctorsParams {
  page?: number;
  pageSize?: number;
  q?: string;
  locationId?: number;
  specialtyId?: number;
}

export async function listDoctors(params: ListDoctorsParams = {}): Promise<ApiListResponse<Doctor>> {
  const search = new URLSearchParams();
  if (params.page) search.set("page", String(params.page));
  if (params.pageSize) search.set("page_size", String(params.pageSize));
  if (params.q) search.set("q", params.q);
  if (params.locationId) search.set("location_id", String(params.locationId));
  if (params.specialtyId) search.set("specialty_id", String(params.specialtyId));

  return apiFetch<ApiListResponse<Doctor>>(`/doctors?${search}`, 60);
}

export async function getDoctor(id: number): Promise<ApiSingleResponse<Doctor> | null> {
  try {
    return await apiFetch<ApiSingleResponse<Doctor>>(`/doctors/${id}`, 3600);
  } catch (err) {
    if (err instanceof ApiRequestError && err.status === 404) return null;
    throw err;
  }
}

export async function listLocations(): Promise<ApiListResponse<Location>> {
  return apiFetch<ApiListResponse<Location>>(`/locations`, 3600);
}

export async function listSpecialties(locationId?: number): Promise<ApiListResponse<Specialty>> {
  const search = locationId ? `?location_id=${locationId}` : "";
  return apiFetch<ApiListResponse<Specialty>>(`/specialties${search}`, 3600);
}