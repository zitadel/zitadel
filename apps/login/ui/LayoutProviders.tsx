"use client";

import { ThemeProvider, useTheme } from "next-themes";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  const { resolvedTheme } = useTheme();
  const isDark = resolvedTheme && resolvedTheme === "dark";

  console.log(isDark);
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      storageKey="cp-theme"
      value={{ dark: "dark" }}
    >
      <div className={`${isDark ? "ui-dark" : "ui-light"} `}>{children}</div>
    </ThemeProvider>
  );
}
