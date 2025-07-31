"use client";

import { MoonIcon, SunIcon } from "@heroicons/react/24/outline";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

export function Theme() {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState<boolean>(false);

  const isDark = resolvedTheme === "dark";

  // useEffect only runs on the client, so now we can safely show the UI
  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <div
      className={`relative grid h-fit grid-cols-2 rounded-full border border-divider-light p-1 dark:border-divider-dark`}
    >
      <button
        className={`flex h-8 w-8 flex-row items-center justify-center rounded-full transition-all hover:opacity-100 ${
          isDark ? "bg-black/10 dark:bg-white/10" : "opacity-60"
        }`}
        onClick={() => setTheme("dark")}
      >
        <MoonIcon className="h-4 w-4 flex-shrink-0 rounded-full text-xl" />
      </button>
      <button
        className={`flex h-8 w-8 flex-row items-center justify-center rounded-full transition-all hover:opacity-100 ${
          !isDark ? "bg-black/10 dark:bg-white/10" : "opacity-60"
        }`}
        onClick={() => setTheme("light")}
      >
        <SunIcon className="h-6 w-6 flex-shrink-0 rounded-full text-xl" />
      </button>
    </div>
  );
}
