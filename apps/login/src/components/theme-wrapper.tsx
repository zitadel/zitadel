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

  // Apply custom font from branding settings before paint to avoid FOUC.
  // When a custom font is uploaded via the label/branding policy, fontUrl
  // contains a fully-resolved URL to the font file served by the assets API.
  // We inject a @font-face rule and set a CSS custom property so the entire
  // login UI picks up the custom font with the existing font as fallback.
  useLayoutEffect(() => {
    const STYLE_ID = "zitadel-custom-font";

    if (branding?.fontUrl) {
      let fontSrc: string;
      try {
        fontSrc = new URL(branding.fontUrl).href;
      } catch {
        // Malformed URL — skip custom font
        return;
      }

      let styleEl = document.getElementById(STYLE_ID) as HTMLStyleElement | null;
      if (!styleEl) {
        styleEl = document.createElement("style");
        styleEl.id = STYLE_ID;
        document.head.appendChild(styleEl);
      }
      // Capture the current font-family (Lato from next/font) before overriding,
      // so it serves as fallback if the custom font fails to load.
      const existingFont =
        getComputedStyle(document.documentElement).fontFamily || "sans-serif";
      const fontStack = `'ZitadelCustomFont', ${existingFont}`;

      styleEl.textContent = `
        @font-face {
          font-family: 'ZitadelCustomFont';
          font-style: normal;
          font-display: swap;
          src: url('${fontSrc}');
        }
      `;

      document.documentElement.style.setProperty("--zitadel-font-family", fontStack);
      // Inline style overrides the class-based Lato from next/font
      document.documentElement.style.setProperty("font-family", fontStack);
    } else {
      // No custom font — remove injected style and let Lato class take over
      const existing = document.getElementById(STYLE_ID);
      if (existing) {
        existing.remove();
      }
      document.documentElement.style.removeProperty("--zitadel-font-family");
      document.documentElement.style.removeProperty("font-family");
    }

    return () => {
      const existing = document.getElementById(STYLE_ID);
      if (existing) {
        existing.remove();
      }
      document.documentElement.style.removeProperty("--zitadel-font-family");
      document.documentElement.style.removeProperty("font-family");
    };
  }, [branding?.fontUrl]);

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
