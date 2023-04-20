"use client";

import { useTheme } from "next-themes";

const ThemeWrapper = ({ children }: any) => {
  const { resolvedTheme } = useTheme();

  const isDark = resolvedTheme && resolvedTheme === "dark";

  return (
    <div className={`${isDark ? "ui-dark" : "ui-light"} `}>{children}</div>
  );
};

export default ThemeWrapper;
