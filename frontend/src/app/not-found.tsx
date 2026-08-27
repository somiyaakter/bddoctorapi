import Link from "next/link";

export default function NotFound() {
  return (
    <div className="mx-auto max-w-md px-4 py-24 text-center sm:px-6">
      <p className="font-mono text-sm text-seal">404</p>

      <h2 className="mt-2 font-display text-xl font-semibold text-ink">
        Page Not Found
      </h2>

      <p className="mt-2 text-sm text-slate">
        The page you are looking for has been removed or never existed.
      </p>

      <Link
        href="/doctors"
        className="mt-6 inline-block rounded-md bg-brand px-4 py-2 text-sm font-medium text-white hover:bg-brand-hover"
      >
        Find Doctors
      </Link>
    </div>
  );
}