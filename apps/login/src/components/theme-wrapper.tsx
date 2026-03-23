"use client";

import { setTheme } from "@/helpers/colors";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { useTheme } from "next-themes";
import { ReactNode, useEffect, useLayoutEffect } from "react";
import { setThemeMode } from "./branding-context";

type Props = {
  branding: BrandingSettings | undefined;
  children: ReactNode;
};

export const ThemeWrapper = ({ children, branding }: Props) => {
  const { setTheme: setNextTheme } = useTheme();

  useEffect(() => {
    setTheme(document, branding);
  }, [branding]);

  // Publish themeMode to the module-level store so ThemeSwitch can read it
  useEffect(() => {
    setThemeMode(branding?.themeMode ?? 0);
  }, [branding?.themeMode]);

  // Handle branding themeMode to force specific theme.
  // Uses useLayoutEffect to apply before paint and writes the forced value
  // to localStorage so next-themes doesn't fall back to system default.
  useLayoutEffect(() => {
    if (branding?.themeMode !== undefined) {
      // Based on the proto definition:
      // THEME_MODE_UNSPECIFIED = 0
      // THEME_MODE_AUTO = 1
      // THEME_MODE_LIGHT = 2
      // THEME_MODE_DARK = 3
      switch (branding.themeMode) {
        case 2: // THEME_MODE_LIGHT
          document.documentElement.classList.remove("dark");
          try {
            localStorage.setItem("cp-theme", "light");
          } catch {}
          setNextTheme("light");
          break;
        case 3: // THEME_MODE_DARK
          document.documentElement.classList.add("dark");
          try {
            localStorage.setItem("cp-theme", "dark");
          } catch {}
          setNextTheme("dark");
          break;
        case 1: // THEME_MODE_AUTO
        case 0: // THEME_MODE_UNSPECIFIED
        default:
          setNextTheme("system");
          break;
      }
    }
  }, [branding?.themeMode, setNextTheme]);

  return <div>{children}</div>;
};
