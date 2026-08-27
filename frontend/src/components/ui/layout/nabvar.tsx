import Link from "next/link";
import { NavbarSearch } from "./navbar-search";
import { MobileMenu } from "./mobile-menu";
import { Stethoscope } from "lucide-react";


const NAV_LINKS = [
  { href: "/doctors", label: "Doctors" },
  { href: "/api-docs", label: "API Docs" },
];

export function Navbar() {
  return (
    <header className="sticky top-0 z-40 border-b border-hairline bg-paper/95 backdrop-blur">
      <nav
        aria-label="Main navigation"
        className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-3 sm:px-6"
      >
        <Link
  href="/"
  className="flex shrink-0 items-center gap-2 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand focus-visible:ring-offset-2 focus-visible:ring-offset-paper rounded-sm"
>
  <span className="flex h-7 w-7 items-center justify-center rounded-md bg-brand text-white">
    <Stethoscope className="h-4 w-4" aria-hidden="true" />
  </span>
  <span className="font-display text-lg font-semibold text-ink">
    MediDirectory
  </span>
</Link>

        {/* Desktop: secondary nav + primary search CTA */}
        <div className="hidden md:flex md:items-center md:gap-6">
          <ul className="flex items-center gap-5">
            {NAV_LINKS.map((link) => (
              <li key={link.href}>
                <Link
                  href={link.href}
                  className="text-sm font-medium text-slate transition-colors hover:text-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand focus-visible:ring-offset-2 focus-visible:ring-offset-paper rounded-sm"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>

          <NavbarSearch />
        </div>

        {/* Mobile: hamburger only, search lives inside the panel */}
        <div className="md:hidden">
          <MobileMenu links={NAV_LINKS} />
        </div>
      </nav>
    </header>
  );
}