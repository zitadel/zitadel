"use client";

import { MoonIcon, SunIcon } from "@heroicons/react/24/outline";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { getThemeConfig, ROUNDNESS_CLASSES } from "@/lib/theme";

function getThemeToggleRoundness() {
  const theme = getThemeConfig();
  // For theme toggle, use specific roundness classes
  switch (theme.roundness) {
    case "edgy":
      return "rounded-none";
    case "mid":
      return "rounded-md";
    case "full":
      return "rounded-full";
    default:
      return "rounded-md";
  }
}

export default function ThemeSwitch() {
  const [mounted, setMounted] = useState(false);
  const { theme, setTheme } = useTheme();
  const toggleRoundness = getThemeToggleRoundness();

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  return (
    <div className={`flex space-x-1 bg-black/5 dark:bg-white/5 p-1 ${toggleRoundness}`}>
      <button
        className={`p-2 ${toggleRoundness} transition-colors ${
          theme === "light" ? "bg-white text-gray-900 shadow-sm" : "text-gray-400 hover:text-gray-300"
        }`}
        onClick={() => setTheme("light")}
        aria-label="Switch to light mode"
      >
        <SunIcon className="h-4 w-4" />
      </button>
      <button
        className={`p-2 ${toggleRoundness} transition-colors ${
          theme === "dark" ? "bg-gray-800 text-white shadow-sm" : "text-gray-600 hover:text-gray-700"
        }`}
        onClick={() => setTheme("dark")}
        aria-label="Switch to dark mode"
      >
        <MoonIcon className="h-4 w-4" />
      </button>
    </div>
  );
}
