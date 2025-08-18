"use client";

import React from "react";
import { useTheme } from "@/lib/useTheme";

interface ThemeProviderProps {
  children: React.ReactNode;
}

/**
 * Theme provider that wraps the entire application and applies global theme styles
 * This component handles the background image and global layout classes
 */
export function ThemeProvider({ children }: ThemeProviderProps) {
  const { theme, classes, getContainerStyle } = useTheme();

  return (
    <div className={classes.layout.container} style={getContainerStyle()}>
      {/* Background overlay for better text readability when using background images */}
      {classes.backgroundImage && theme.backgroundImage && <div className={classes.backgroundImage.overlay} />}

      {/* Content with proper z-index to appear above background */}
      <div className="relative z-10 w-full h-full">{children}</div>
    </div>
  );
}
