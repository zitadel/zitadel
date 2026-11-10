"use client";

import { clsx } from "clsx";

// Generic theme-aware properties interface (kept for backward compatibility)
export interface ThemeableProps {
  roundness?: string;
  spacing?: string;
  padding?: string;
  typography?: string;
}

// Utility for conditional roundness classes (avoids Tailwind conflicts)
export function getRoundnessClasses(roundness: string, baseClasses: string = "") {
  return clsx(baseClasses, {
    "rounded-none": roundness === "rounded-none",
    "rounded-md": roundness === "rounded-md",
    "rounded-full": roundness === "rounded-full",
    "rounded-lg": roundness === "rounded-lg",
    "rounded-3xl": roundness === "rounded-3xl",
  });
}

// Utility for button-specific roundness
export function getButtonRoundnessClasses(roundness: string) {
  return clsx({
    "rounded-none": roundness === "rounded-none",
    "rounded-md": roundness === "rounded-md",
    "rounded-full": roundness === "rounded-full",
  });
}

// Utility for input-specific roundness
export function getInputRoundnessClasses(roundness: string) {
  return clsx({
    "rounded-none": roundness === "rounded-none",
    "rounded-md": roundness === "rounded-md",
    "rounded-full": roundness === "rounded-full",
  });
}

// Utility for card-specific roundness
export function getCardRoundnessClasses(roundness: string) {
  return clsx({
    "rounded-none": roundness === "rounded-none",
    "rounded-lg": roundness === "rounded-lg",
    "rounded-3xl": roundness === "rounded-3xl",
  });
}
