import Image from "next/image";
import Link from "next/link";
import type { Doctor } from "@/lib/types";

export function DoctorCard({ doctor }: { doctor: Doctor }) {
  return (
    <Link
      href={`/doctors/${doctor.id}`}
      className="group flex flex-col rounded-lg border border-hairline bg-white p-4 transition-colors hover:border-brand"
    >
      <div className="flex gap-4">
        <div className="relative h-16 w-16 shrink-0 overflow-hidden rounded-full bg-secondary">
          {doctor.image_url ? (
            <Image
              src={doctor.image_url}
              alt={doctor.name}
              fill
              sizes="64px"
              className="object-cover"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center text-lg font-semibold text-slate">
              {doctor.name.charAt(0)}
            </div>
          )}
        </div>

        <div className="min-w-0 flex-1">
          <h3 className="line-clamp-2 font-display text-base font-semibold text-ink group-hover:text-brand">
            {doctor.name}
          </h3>
          {doctor.degrees && (
            <p className="mt-0.5 truncate text-xs text-slate">{doctor.degrees}</p>
          )}
        </div>
      </div>

      <div className="mt-3 border-t border-hairline pt-3">
        {doctor.specialties && (
          <p className="line-clamp-2 text-sm text-ink">{doctor.specialties}</p>
        )}
        {doctor.workplace && (
          <p className="mt-1 truncate text-xs text-slate">{doctor.workplace}</p>
        )}
      </div>

      {doctor.bmdc_reg_no && (
        <div className="mt-3 flex items-center gap-1.5">
          <span className="flex h-5 w-5 items-center justify-center rounded-full bg-seal/10 text-[10px] text-seal">
            ✓
          </span>
          <span className="font-mono text-[11px] text-slate">
            BMDC {doctor.bmdc_reg_no}
          </span>
        </div>
      )}
    </Link>
  );
}