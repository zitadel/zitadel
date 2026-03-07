"use client";

import { useEffect } from "react";
import { useThemeConfig } from "@/lib/theme-hooks";

export function CustomCssWrapper({ children }: { children: React.ReactNode }) {
  const themeConfig = useThemeConfig();

  useEffect(() => {
    if (!themeConfig.customCssFile || typeof document === "undefined") {
      return;
    }

    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = themeConfig.customCssFile;
    link.setAttribute("data-custom-css", "true");
    document.head.appendChild(link);

    return () => {
      const existingLink = document.querySelector('link[data-custom-css="true"]');
      if (existingLink) {
        document.head.removeChild(existingLink);
      }
    };
  }, [themeConfig.customCssFile]);

  return <>{children}</>;
}
