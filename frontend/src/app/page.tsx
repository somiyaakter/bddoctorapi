import Link from "next/link";
import { listDoctors, listLocations, listSpecialties } from "@/lib/api";
import { SearchBar } from "../components/ui/search/search-bar";
import { DoctorGrid } from "../components/ui/doctor/doctor-grid";

export default async function HomePage() {
  const [doctorsRes, locationsRes, specialtiesRes] = await Promise.all([
    listDoctors({ page: 1, pageSize: 6 }),
    listLocations(),
    listSpecialties(),
  ]);

  const featuredDoctors = doctorsRes.data ?? [];
  const totalDoctors = doctorsRes.pagination?.total_items ?? 0;
  const totalLocations = locationsRes.data?.length ?? 0;
  const totalSpecialties = specialtiesRes.data?.length ?? 0;

  return (
    <div>
      {/* Hero */}
      <section className="border-b border-hairline bg-white">
        <div className="mx-auto max-w-4xl px-4 py-16 text-center sm:px-6 sm:py-24">
          <h1 className="font-display text-3xl font-semibold text-ink sm:text-5xl">
            Find a verified doctor near you
          </h1>
          <p className="mx-auto mt-4 max-w-xl text-base text-slate">
            Search {totalDoctors.toLocaleString()} doctors across Bangladesh by name,
            specialty, hospital, or BMDC number.
          </p>

          <div className="mx-auto mt-8 max-w-xl">
            <SearchBar />
          </div>
        </div>
      </section>

      {/* Stats strip */}
      <section className="border-b border-hairline bg-paper">
        <div className="mx-auto flex max-w-4xl flex-wrap justify-center gap-x-10 gap-y-3 px-4 py-6 font-mono text-sm text-slate sm:px-6">
          <span>{totalDoctors.toLocaleString()} Doctors</span>
          <span>{totalSpecialties} Specialties</span>
          <span>{totalLocations} Locations</span>
        </div>
      </section>

      {/* Featured doctors */}
      {featuredDoctors.length > 0 && (
        <section className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
          <div className="flex items-baseline justify-between">
            <h2 className="font-display text-xl font-semibold text-ink sm:text-2xl">
              Recently listed doctors
            </h2>
            <Link href="/doctors" className="text-sm font-medium text-brand hover:text-brand-hover">
              View all →
            </Link>
          </div>
          <div className="mt-6">
            <DoctorGrid doctors={featuredDoctors} />
          </div>
        </section>
      )}

      {/* CTA band */}
      <section className="border-t border-hairline bg-white">
        <div className="mx-auto max-w-4xl px-4 py-12 text-center sm:px-6">
          <h2 className="font-display text-xl font-semibold text-ink sm:text-2xl">
            Browse the full directory
          </h2>
          <p className="mt-2 text-sm text-slate">
            Filter by location and specialty to find the right doctor for you.
          </p>
          <Link
            href="/doctors"
            className="mt-6 inline-block rounded-md bg-brand px-6 py-3 text-sm font-medium text-white hover:bg-brand-hover"
          >
            Browse all doctors
          </Link>
        </div>
      </section>
    </div>
  );
}