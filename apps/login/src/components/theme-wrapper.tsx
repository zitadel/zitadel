"use client";

import { setTheme } from "@/helpers/colors";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { ReactNode, useEffect } from "react";
import { useTheme } from "next-themes";

type Props = {
  branding: BrandingSettings | undefined;
  children: ReactNode;
};

export const ThemeWrapper = ({ children, branding }: Props) => {
  const { setTheme: setNextTheme } = useTheme();

  useEffect(() => {
    setTheme(document, branding);
  }, [branding]);

  // Handle branding themeMode to force specific theme
  useEffect(() => {
    if (branding?.themeMode !== undefined) {
      // Based on the proto definition:
      // THEME_MODE_UNSPECIFIED = 0
      // THEME_MODE_AUTO = 1
      // THEME_MODE_LIGHT = 2
      // THEME_MODE_DARK = 3
      switch (branding.themeMode) {
        case 2: // THEME_MODE_LIGHT
          setNextTheme("light");
          break;
        case 3: // THEME_MODE_DARK
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
