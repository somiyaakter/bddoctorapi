import type { Doctor } from "@/lib/types";
import { DoctorCard } from "./doctor-card";


export function DoctorGrid({ doctors }: { doctors: Doctor[] }) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
      {doctors.map((doctor) => (
        <DoctorCard key={doctor.id} doctor={doctor} />
      ))}
    </div>
  );
}