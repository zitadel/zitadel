"use client";

import { useState, useEffect } from "react";
import { getThemeConfig } from "./theme";

/**
 * Custom hook that returns the effective layout mode, taking into account
 * both the theme configuration and responsive breakpoints.
 *
 * On medium screens and below (md: max-width 767px), it will automatically
 * switch to top-to-bottom layout regardless of the theme setting.
 *
 * NOTE: This is a client-side hook and requires "use client" directive.
 */
export function useResponsiveLayout(): { isSideBySide: boolean; isResponsiveOverride: boolean } {
  const themeConfig = getThemeConfig();
  const [isMdOrSmaller, setIsMdOrSmaller] = useState(false);
  const [isHydrated, setIsHydrated] = useState(false);

  useEffect(() => {
    // Mark as hydrated on client side
    setIsHydrated(true);

    // Check if we're in a browser environment
    if (typeof window === "undefined") {
      return;
    }

    const mediaQuery = window.matchMedia("(max-width: 767px)"); // md breakpoint is 768px in Tailwind

    // Set initial value
    setIsMdOrSmaller(mediaQuery.matches);

    // Listen for changes
    const handleChange = (e: MediaQueryListEvent) => {
      setIsMdOrSmaller(e.matches);
    };

    mediaQuery.addEventListener("change", handleChange);

    // Cleanup
    return () => mediaQuery.removeEventListener("change", handleChange);
  }, []);

  const configuredSideBySide = themeConfig.layout === "side-by-side";

  // During SSR or before hydration, assume desktop (side-by-side if configured)
  // This prevents hydration mismatches
  const isSideBySide = configuredSideBySide && (isHydrated ? !isMdOrSmaller : true);
  const isResponsiveOverride = configuredSideBySide && isHydrated && isMdOrSmaller;

  return { isSideBySide, isResponsiveOverride };
}

/**
 * Custom hook that returns the theme configuration for client-side usage.
 *
 * NOTE: This is a client-side hook and requires "use client" directive.
 */
export function useThemeConfig() {
  return getThemeConfig();
}
