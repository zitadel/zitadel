"use client";

import { setTheme } from "@/helpers/colors";
import { BrandingSettings, ThemeMode } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
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
    setThemeMode(branding?.themeMode ?? ThemeMode.UNSPECIFIED);
  }, [branding?.themeMode]);

  // Handle branding themeMode to force specific theme.
  // Uses useLayoutEffect to apply before paint and writes the forced value
  // to localStorage so next-themes doesn't fall back to system default.
  useLayoutEffect(() => {
    if (branding?.themeMode !== undefined) {
      switch (branding.themeMode) {
        case ThemeMode.LIGHT:
          document.documentElement.classList.remove("dark");
          try {
            localStorage.setItem("cp-theme", "light");
          } catch {
            /* localStorage unavailable (e.g. private mode) */
          }
          setNextTheme("light");
          break;
        case ThemeMode.DARK:
          document.documentElement.classList.add("dark");
          try {
            localStorage.setItem("cp-theme", "dark");
          } catch {
            /* localStorage unavailable (e.g. private mode) */
          }
          setNextTheme("dark");
          break;
        case ThemeMode.AUTO:
        case ThemeMode.UNSPECIFIED:
        default:
          setNextTheme("system");
          break;
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [branding?.themeMode]);

  return <div>{children}</div>;
};
