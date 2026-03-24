"use client";

import { useSyncExternalStore } from "react";

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

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

function getSnapshot() {
  return currentThemeMode;
}

/**
 * Hook to subscribe to themeMode changes using React 18+ useSyncExternalStore.
 * Ensures ThemeSwitch always sees a consistent snapshot during concurrent renders.
 */
export function useThemeMode(): number {
  return useSyncExternalStore(subscribe, getSnapshot, getSnapshot);
}
