"use client";
import { ThemeProvider as ThemeP } from "next-themes";
import { ReactNode } from "react";

export function ThemeProvider({ children, nonce }: { children: ReactNode; nonce?: string }) {
  return (
    <ThemeP nonce={nonce} attribute="class" defaultTheme="system" storageKey="cp-theme" value={{ dark: "dark" }}>
      {children}
    </ThemeP>
  );
}
