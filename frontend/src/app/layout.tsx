import type { Metadata } from "next";
import {
  Fraunces,
  Hind_Siliguri,
  IBM_Plex_Mono,
} from "next/font/google";

import "./globals.css";
import { Navbar } from "../components/ui/layout/nabvar";
import { Footer } from "../components/ui/layout/footer";

const fraunces = Fraunces({
  subsets: ["latin"],
  variable: "--font-display",
  weight: ["500", "600"],
});

const hindSiliguri = Hind_Siliguri({
  subsets: ["latin", "bengali"],
  variable: "--font-body",
  weight: ["400", "500", "600"],
});

const plexMono = IBM_Plex_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  weight: ["400", "500"],
});

export const metadata: Metadata = {
  title: {
    default: "MediDirectory — Find Trusted Doctors",
    template: "%s | DataLab",
  },
  description:
    "A directory of 7,000+ doctors in Bangladesh. Find doctors by specialty, location, and hospital.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${fraunces.variable} ${hindSiliguri.variable} ${plexMono.variable}`}
    >
      <body className="font-body bg-paper text-ink antialiased">
        <Navbar />

        <main className="min-h-screen">{children}</main>

        <Footer />
      </body>
    </html>
  );
}