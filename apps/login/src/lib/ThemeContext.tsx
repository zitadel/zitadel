"use client";

import React, { createContext, useContext, useMemo } from "react";
import { getThemeConfig, getThemeClasses, type ThemeConfig } from "./theme";

interface ThemeContextValue {
  theme: ThemeConfig;
  classes: ReturnType<typeof getThemeClasses>;
  // Helper functions for common use cases
  getCardClasses: () => string;
  getButtonClasses: () => string;
  getInputClasses: () => string;
  getSpacingClasses: () => string;
  getPaddingClasses: () => string;
  getTypographyClasses: () => string;
  getContainerStyle: () => React.CSSProperties;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

interface ThemeContextProviderProps {
  children: React.ReactNode;
  // Optional: Allow injecting custom theme config
  customTheme?: Partial<ThemeConfig>;
}

export function ThemeContextProvider({ children, customTheme }: ThemeContextProviderProps) {
  const value = useMemo(() => {
    // Merge default theme with custom theme if provided
    const theme = customTheme ? { ...getThemeConfig(), ...customTheme } : getThemeConfig();
    const classes = getThemeClasses(theme);

    return {
      theme,
      classes,
      getCardClasses: () => `${classes.roundness.card} ${classes.preset.card}`,
      getButtonClasses: () => `${classes.roundness.button} ${classes.preset.typography} px-4 py-2 transition-all`,
      getInputClasses: () => `${classes.roundness.input} px-3 py-2 transition-all focus:outline-none`,
      getSpacingClasses: () => classes.preset.spacing,
      getPaddingClasses: () => classes.preset.padding,
      getTypographyClasses: () => classes.preset.typography,
      getContainerStyle: () => ({
        ...(classes.backgroundImage && theme.backgroundImage
          ? {
              backgroundImage: `url(${theme.backgroundImage})`,
              backgroundSize: "cover",
              backgroundPosition: "center",
              backgroundRepeat: "no-repeat",
            }
          : {}),
      }),
    };
  }, [customTheme]);

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeContextProvider");
  }
  return context;
}
