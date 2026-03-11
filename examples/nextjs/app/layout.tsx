import type { Metadata } from "next";
import { Arimo, Inter } from "next/font/google";
import type { ReactNode } from "react";
import "./globals.css";

const arimo = Arimo({
  subsets: ["latin"],
  variable: "--font-sans",
  weight: ["400", "500", "600", "700"],
  display: "swap",
});

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: "ZITADEL Next.js SDK demo",
  description: "Workspace-linked demo shell for local SDK development",
};

export default function RootLayout({
  children,
}: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="en" className={`${arimo.variable} ${inter.variable}`}>
      <body>{children}</body>
    </html>
  );
}
