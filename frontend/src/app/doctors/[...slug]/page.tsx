import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { getDoctor, listDoctors, listLocations, listSpecialties } from "@/lib/api";
import { isDoctorId, findLocationBySlug, findSpecialtyBySlug } from "@/lib/slug";
import { DoctorDetails } from "../../../components/ui/doctor/doctor-details";
import { EmptyState } from "../../../components/ui/shared/empty-state";
import { DoctorGrid } from "../../../components/ui/doctor/doctor-grid";
import { Pagination } from "../../../components/ui/shared/pagination";


interface PageProps {
  params: Promise<{ slug: string[] }>;
  searchParams: Promise<{ page?: string }>;
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;

  if (slug.length === 1 && isDoctorId(slug[0])) {
    const res = await getDoctor(parseInt(slug[0], 10));
    if (!res) return { title: "Doctor Not Found" };
    const d = res.data;
    return {
      title: d.name,
      description: `${d.name}${d.specialties ? ` — ${d.specialties}` : ""}${
        d.workplace ? ` at ${d.workplace}` : ""
      }.`,
    };
  }

  return { title: "Browse Doctors" };
}

export default async function DoctorSlugPage({ params, searchParams }: PageProps) {
  const { slug } = await params;
  const { page: pageParam } = await searchParams;
  const page = pageParam ? parseInt(pageParam, 10) : 1;

  // Case: /doctors/123 -> doctor detail
  if (slug.length === 1 && isDoctorId(slug[0])) {
    const res = await getDoctor(parseInt(slug[0], 10));
    if (!res) notFound();
    return <DoctorDetails doctor={res.data} />;
  }

  if (slug.length > 2) notFound();

  const [locationsRes, specialtiesRes] = await Promise.all([listLocations(), listSpecialties()]);
  const locations = locationsRes.data ?? [];
  const specialties = specialtiesRes.data ?? [];

  let locationId: number | undefined;
  let locationName: string | undefined;

  // Case: /doctors/cardiology/dhaka -> specialty + location
  if (slug.length === 2) {
    const loc = findLocationBySlug(locations, slug[1]);
    if (!loc) notFound();
    locationId = loc.id;
    locationName = loc.name;
  }

  // Case: /doctors/cardiology -> specialty only
  const specialty = findSpecialtyBySlug(specialties, locations, slug[0], locationId);
  if (!specialty) notFound();

  const doctorsRes = await listDoctors({
    page,
    pageSize: 12,
    specialtyId: specialty.id,
    locationId,
  });
  const doctors = doctorsRes.data ?? [];
  const pagination = doctorsRes.pagination;

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6">
      <h1 className="font-display text-2xl font-semibold text-ink sm:text-3xl">
        {specialty.name}
        {locationName ? ` in ${locationName}` : ""}
      </h1>
      <p className="mt-1 text-sm text-slate">
        {pagination?.total_items ?? doctors.length} doctor
        {(pagination?.total_items ?? doctors.length) === 1 ? "" : "s"} found
      </p>

      {doctors.length === 0 ? (
        <div className="mt-8">
          <EmptyState
            title="No doctors found"
            description="Try browsing all doctors instead."
          />
        </div>
      ) : (
        <>
          <div className="mt-6">
            <DoctorGrid doctors={doctors} />
          </div>
          {pagination && pagination.total_pages > 1 && (
            <div className="mt-8">
              <Pagination  currentPage={pagination.page} totalPages={pagination.total_pages} />
            </div>
          )}
        </>
      )}
    </div>
  );
}