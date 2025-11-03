"use client";

import { setLanguageCookie } from "@/lib/cookies";
import { Lang, LANGS } from "@/lib/i18n";
import { getThemeConfig, getComponentRoundness, APPEARANCE_STYLES } from "@/lib/theme";
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from "@headlessui/react";
import { CheckIcon, ChevronDownIcon } from "@heroicons/react/24/outline";
import clsx from "clsx";
import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";

// Helper function to get language switcher roundness from theme
function getLanguageSwitcherRoundness(): string {
  return getComponentRoundness("button");
}

// Helper function to get card appearance styles for the language switcher
function getLanguageSwitcherCardAppearance(): string {
  const themeConfig = getThemeConfig();
  const appearance = APPEARANCE_STYLES[themeConfig.appearance];
  return appearance?.card || "bg-black/5 dark:bg-white/5"; // Fallback to current styling
}

export function LanguageSwitcher() {
  const currentLocale = useLocale();
  const switcherRoundness = getLanguageSwitcherRoundness();
  const cardAppearance = getLanguageSwitcherCardAppearance();

  const [selected, setSelected] = useState(LANGS.find((l) => l.code === currentLocale) || LANGS[0]);

  const router = useRouter();

  const handleChange = async (language: Lang) => {
    setSelected(language);
    const newLocale = language.code;

    await setLanguageCookie(newLocale);

    router.refresh();
  };

  return (
    <div className="w-32">
      <Listbox value={selected} onChange={handleChange}>
        <ListboxButton
          className={clsx(
            `relative block w-full py-1.5 pl-3 pr-8 text-left text-sm/6 text-black dark:text-white ${switcherRoundness}`,
            cardAppearance,
            "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25",
          )}
        >
          {selected.name}
          <ChevronDownIcon className="group pointer-events-none absolute right-2.5 top-2.5 size-4" aria-hidden="true" />
        </ListboxButton>
        <ListboxOptions
          anchor="bottom"
          transition
          className={clsx(
            `w-[var(--button-width)] border border-black/5 bg-background-light-500 p-1 [--anchor-gap:var(--spacing-1)] focus:outline-none dark:border-white/5 dark:bg-background-dark-500 rounded-md`,
            "transition duration-100 ease-in data-[leave]:data-[closed]:opacity-0",
          )}
        >
          {LANGS.map((lang) => (
            <ListboxOption
              key={lang.code}
              value={lang}
              className={`group flex cursor-default select-none items-center gap-2 px-3 py-1.5 data-[focus]:bg-black/10 dark:data-[focus]:bg-white/10 ${switcherRoundness}`}
            >
              <CheckIcon className="invisible size-4 group-data-[selected]:visible" />
              <div className="text-sm/6 text-black dark:text-white">{lang.name}</div>
            </ListboxOption>
          ))}
        </ListboxOptions>
      </Listbox>
    </div>
  );
}
