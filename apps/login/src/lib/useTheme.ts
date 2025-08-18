"use client";

import { useMemo } from "react";
import { getThemeConfig, getThemeClasses, type ThemeConfig } from "./theme";

/**
 * Hook to get the current theme configuration and CSS classes
 * This hook can be used in any component that needs theme-aware styling
 */
export function useTheme() {
  const theme = useMemo(() => getThemeConfig(), []);
  const classes = useMemo(() => getThemeClasses(theme), [theme]);

  return {
    theme,
    classes,
    // Helper functions for structural/visual styling (colors handled by API)
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
}
