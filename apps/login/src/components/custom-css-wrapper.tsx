"use client";

import { useEffect, useRef } from "react";
import { useThemeConfig } from "@/lib/theme-hooks";

export function CustomCssWrapper({ children }: { children: React.ReactNode }) {
  const themeConfig = useThemeConfig();
  const linkRef = useRef<HTMLLinkElement | null>(null);

  useEffect(() => {
    if (!themeConfig.customCssFile || typeof document === "undefined") {
      return;
    }

    let link = linkRef.current;

    if (!link) {
      link = document.createElement("link");
      link.rel = "stylesheet";
      link.setAttribute("data-custom-css", "true");
      document.head.appendChild(link);
      linkRef.current = link;
    }

    link.href = themeConfig.customCssFile;

    return () => {
      if (linkRef.current && linkRef.current.parentNode === document.head) {
        document.head.removeChild(linkRef.current);
      }
      linkRef.current = null;
    };
  }, [themeConfig.customCssFile]);

  return <>{children}</>;
}
