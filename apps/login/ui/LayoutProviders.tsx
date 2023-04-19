"use client";

import { ThemeProvider } from "next-themes";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      storageKey="cp-theme"
      value={{ dark: "dark" }}
    >
      {children}
    </ThemeProvider>
  );
}
