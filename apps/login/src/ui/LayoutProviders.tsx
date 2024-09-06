"use client";

import { useTheme } from "next-themes";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  const { resolvedTheme } = useTheme();
  const isDark = resolvedTheme === "dark";

  return (
    <div className={`${isDark ? "ui-dark" : "ui-light"} `}>{children}</div>
  );
}
