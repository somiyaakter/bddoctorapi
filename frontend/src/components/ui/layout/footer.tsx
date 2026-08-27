import Link from "next/link";

const footerLinks = [
  { label: "Find Doctors", href: "/doctors" },
  { label: "API Documentation", href: "/api-docs" },
];

export function Footer() {
  return (
    <footer className="border-t border-hairline bg-paper">
      <div className="mx-auto max-w-6xl px-4 sm:px-6">
        <div className="flex flex-col gap-8 py-10 sm:py-12">
          {/* Top */}
          <div className="flex flex-col gap-6 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <Link
                href="/"
                className="font-display text-2xl font-semibold tracking-tight text-ink"
              >
                MediDirectory
              </Link>

              <p className="mt-2 max-w-md text-sm leading-6 text-slate">
                A trusted doctor directory for finding specialists,
                hospitals, and chamber information across Bangladesh.
              </p>
            </div>

            <nav
              aria-label="Footer navigation"
              className="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm"
            >
              {footerLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  className="font-medium text-slate transition-colors hover:text-brand"
                >
                  {link.label}
                </Link>
              ))}
            </nav>
          </div>

          {/* Bottom */}
          <div className="flex flex-col gap-3 border-t border-hairline pt-6 text-xs text-slate sm:flex-row sm:items-center sm:justify-between">
            <p>
              © {new Date().getFullYear()} MediDirectory
            </p>

          </div>
        </div>
      </div>
    </footer>
  );
}