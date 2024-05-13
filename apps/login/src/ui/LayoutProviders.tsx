"use client";
import { ZitadelReactProvider } from "@zitadel/react";
import { useTheme } from "next-themes";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  const { resolvedTheme } = useTheme();
  const isDark = resolvedTheme === "dark";

  return (
    <div className={`${isDark ? "ui-dark" : "ui-light"} `}>
      <ZitadelReactProvider dark={isDark}>{children}</ZitadelReactProvider>
    </div>
  );
}
