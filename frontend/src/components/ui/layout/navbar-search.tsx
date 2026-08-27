"use client";

import { useRouter } from "next/navigation";
import { useRef, useState } from "react";
import { Search, X } from "lucide-react";

export function NavbarSearch() {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  function openSearch() {
    setOpen(true);
    // wait for the input to mount before focusing
    requestAnimationFrame(() => inputRef.current?.focus());
  }

  function closeSearch() {
    setOpen(false);
    setValue("");
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!value.trim()) return;
    router.push(`/doctors?q=${encodeURIComponent(value.trim())}`);
    closeSearch();
  }

  if (!open) {
    return (
      <button
        type="button"
        onClick={openSearch}
        className="inline-flex items-center gap-2 rounded-md bg-brand px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand focus-visible:ring-offset-2 focus-visible:ring-offset-paper"
      >
        <Search className="h-4 w-4" aria-hidden="true" />
        Search Doctors
      </button>
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      role="search"
      className="flex items-center gap-2 rounded-md border border-hairline bg-white pl-3 pr-1.5 py-1.5"
    >
      <Search className="h-4 w-4 shrink-0 text-slate" aria-hidden="true" />
      <input
        ref={inputRef}
        type="search"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onKeyDown={(e) => e.key === "Escape" && closeSearch()}
        placeholder="Name, specialty, hospital…"
        aria-label="Search doctors"
        className="w-56 bg-transparent text-sm text-ink placeholder:text-slate focus:outline-none"
      />
      <button
        type="button"
        onClick={closeSearch}
        aria-label="Close search"
        className="shrink-0 rounded p-1 text-slate hover:text-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand"
      >
        <X className="h-4 w-4" aria-hidden="true" />
      </button>
    </form>
  );
}