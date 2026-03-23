"use client";

import { useEffect, useState } from "react";

/**
 * Module-level store for themeMode from branding settings.
 * This is used instead of React Context because ThemeSwitch in the layout
 * is a sibling of (not a descendant of) the ThemeWrapper/DynamicTheme tree.
 */
let currentThemeMode = 0; // Default to UNSPECIFIED
const listeners = new Set<() => void>();

function notifyListeners() {
  listeners.forEach((listener) => listener());
}

export function setThemeMode(mode: number) {
  if (currentThemeMode !== mode) {
    currentThemeMode = mode;
    notifyListeners();
  }
}

export function getThemeMode() {
  return currentThemeMode;
}

/**
 * Hook to subscribe to themeMode changes.
 * Works across the component tree regardless of React Context boundaries.
 */
export function useThemeMode(): number {
  const [themeMode, setThemeModeState] = useState(currentThemeMode);

  useEffect(() => {
    // Sync in case value changed between render and effect
    setThemeModeState(currentThemeMode);

    const listener = () => setThemeModeState(currentThemeMode);
    listeners.add(listener);
    return () => {
      listeners.delete(listener);
    };
  }, []);

  return themeMode;
}
