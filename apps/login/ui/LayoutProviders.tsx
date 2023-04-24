"use client";

import { ColorService } from "#/utils/colors";
import { ThemeProvider, useTheme } from "next-themes";
import { useEffect } from "react";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  const { resolvedTheme } = useTheme();
  const isDark = resolvedTheme && resolvedTheme === "dark";

  useEffect(() => {
    new ColorService(document);
  });

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
