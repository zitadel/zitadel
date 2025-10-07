"use client";

import { MoonIcon, SunIcon } from "@heroicons/react/24/outline";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { getThemeConfig, getComponentRoundness, APPEARANCE_STYLES } from "@/lib/theme";

function getThemeToggleRoundness() {
  return getComponentRoundness("themeSwitch");
}

// Helper function to get card appearance styles for the theme switch wrapper
function getThemeSwitchCardAppearance(): string {
  const themeConfig = getThemeConfig();
  const appearance = APPEARANCE_STYLES[themeConfig.appearance];
  return appearance?.card || "bg-black/5 dark:bg-white/5"; // Fallback to current styling
}

// Helper function to get selected button styling for clear visibility
function getSelectedButtonStyle(isSelected: boolean): string {
  const themeConfig = getThemeConfig();

  if (!isSelected) {
    return "text-gray-400 hover:text-gray-300 dark:text-gray-500 dark:hover:text-gray-400";
  }

  // Selected state styling based on appearance theme
  switch (themeConfig.appearance) {
    case "glass":
      return "bg-white/30 dark:bg-black/30 text-gray-900 dark:text-white shadow-lg backdrop-blur-sm border border-white/40 dark:border-white/20";
    case "material":
      return "bg-white dark:bg-gray-800 text-gray-900 dark:text-white shadow-md";
    case "flat":
    default:
      return "bg-white dark:bg-gray-800 text-gray-900 dark:text-white border border-gray-200 dark:border-gray-700";
  }
}

export default function ThemeSwitch() {
  const [mounted, setMounted] = useState(false);
  const { theme, setTheme } = useTheme();
  const toggleRoundness = getThemeToggleRoundness();
  const cardAppearance = getThemeSwitchCardAppearance();

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  return (
    <div className={`flex space-x-1 p-1 ${toggleRoundness} ${cardAppearance}`}>
      <button
        className={`w-8 h-8 flex flex-row items-center justify-center ${toggleRoundness} transition-colors ${getSelectedButtonStyle(theme === "light")}`}
        onClick={() => setTheme("light")}
        aria-label="Switch to light mode"
      >
        <SunIcon className="h-5 w-5" />
      </button>
      <button
        className={`w-8 h-8 flex flex-row items-center justify-center ${toggleRoundness} transition-colors ${getSelectedButtonStyle(theme === "dark")}`}
        onClick={() => setTheme("dark")}
        aria-label="Switch to dark mode"
      >
        <MoonIcon className="h-4 w-4" />
      </button>
    </div>
  );
}
