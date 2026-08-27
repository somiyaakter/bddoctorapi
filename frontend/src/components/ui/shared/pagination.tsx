"use client";

import Link from "next/link";
import { useSearchParams, usePathname } from "next/navigation";

export function Pagination({
  currentPage,
  totalPages,
}: {
  currentPage: number;
  totalPages: number;
}) {
  const pathname = usePathname();
  const searchParams = useSearchParams();

  function hrefFor(page: number) {
    const params = new URLSearchParams(searchParams.toString());
    params.set("page", String(page));
    return `${pathname}?${params.toString()}`;
  }

  const prevDisabled = currentPage <= 1;
  const nextDisabled = currentPage >= totalPages;

  return (
    <nav aria-label="Pagination" className="flex items-center justify-center gap-2">
      <Link
        href={hrefFor(currentPage - 1)}
        aria-disabled={prevDisabled}
        className={`rounded-md border border-hairline px-3 py-1.5 text-sm ${
          prevDisabled
            ? "pointer-events-none text-slate/40"
            : "text-ink hover:border-brand hover:text-brand"
        }`}
      >
        Previous
      </Link>
      <span className="font-mono text-sm text-slate">
        Page {currentPage} of {totalPages}
      </span>
      <Link
        href={hrefFor(currentPage + 1)}
        aria-disabled={nextDisabled}
        className={`rounded-md border border-hairline px-3 py-1.5 text-sm ${
          nextDisabled
            ? "pointer-events-none text-slate/40"
            : "text-ink hover:border-brand hover:text-brand"
        }`}
      >
        Next
      </Link>
    </nav>
  );
}