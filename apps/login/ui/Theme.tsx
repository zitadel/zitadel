"use client";

import React, { useEffect, useState } from "react";
import { useTheme } from "next-themes";
import { MoonIcon, SunIcon } from "@heroicons/react/24/outline";

function Theme() {
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
      className={`relative grid grid-cols-2 rounded-full border border-divider-light dark:border-divider-dark p-1`}
    >
      <button
        className={`h-8 w-8 rounded-full flex flex-row items-center justify-center hover:opacity-100 transition-all ${
          isDark ? "bg-black/10 dark:bg-white/10" : "opacity-60"
        }`}
        onClick={() => setTheme("dark")}
      >
        <MoonIcon className="h-4 w-4 flex-shrink-0 text-xl rounded-full" />
      </button>
      <button
        className={`h-8 w-8 rounded-full flex flex-row items-center justify-center hover:opacity-100 transition-all ${
          !isDark ? "bg-black/10 dark:bg-white/10" : "opacity-60"
        }`}
        onClick={() => setTheme("light")}
      >
        <SunIcon className="h-6 w-6 flex-shrink-0 text-xl rounded-full" />
      </button>
    </div>
  );
}

export default Theme;
