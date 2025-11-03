"use client";
import { ThemeProvider as ThemeP } from "next-themes";
import { ReactNode } from "react";

export function ThemeProvider({ children }: { children: ReactNode }) {
  return (
    <ThemeP
      attribute="class"
      defaultTheme="system"
      storageKey="cp-theme"
      value={{ dark: "dark" }}
    >
      {children}
    </ThemeP>
  );
}
