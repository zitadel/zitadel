"use client";

import { setLanguageCookie } from "@/lib/cookies";
import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
  Transition,
} from "@headlessui/react";
import { CheckIcon, ChevronUpDownIcon } from "@heroicons/react/24/outline";
import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import { Fragment, useState } from "react";

interface Lang {
  id: number;
  name: string;
  img: string;
  code: string;
}

const LANGS: Lang[] = [
  {
    id: 1,
    name: "English",
    code: "en",
    img: "/images/flags/us.png",
  },
  {
    id: 2,
    name: "Deutsch",
    code: "de",
    img: "/images/flags/de.png",
  },
  {
    id: 3,
    name: "Italiano",
    code: "it",
    img: "/images/flags/it.png",
  },
  {
    id: 4,
    name: "Español",
    code: "es",
    img: "/images/flags/es.png",
  },
];

export function LanguageSwitcher() {
  const currentLocale = useLocale();

  const [selected, setSelected] = useState(
    LANGS.find((l) => l.code === currentLocale) || LANGS[0],
  );

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
        <div className="relative">
          <ListboxButton className="relative w-full cursor-default rounded-lg border border-divider-light bg-background-light-500 dark:bg-background-dark-500 py-2 pl-3 pr-10 text-left focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm">
            <span className="block truncate">{selected.name}</span>
            <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <ChevronUpDownIcon
                className="h-5 w-5 text-gray-400"
                aria-hidden="true"
              />
            </span>
          </ListboxButton>
          <Transition
            as={Fragment}
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <ListboxOptions
              anchor="bottom"
              className="absolute mt-1 max-h-60 w-48 w-full overflow-auto rounded-md text-text-light-500 dark:text-text-dark-500 bg-background-light-500 dark:bg-background-dark-500 py-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm"
            >
              {LANGS.map((lang, index) => (
                <ListboxOption
                  key={lang.code}
                  className={({ active }) =>
                    `relative cursor-default select-none py-2 pl-10 pr-4 ${
                      active
                        ? "bg-background-light-300 dark:bg-background-dark-300"
                        : ""
                    }`
                  }
                  value={lang}
                >
                  {({ selected }) => (
                    <>
                      <span
                        className={`block truncate ${
                          selected ? "font-medium" : "font-normal"
                        }`}
                      >
                        {lang.name}
                      </span>
                      {selected ? (
                        <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-primary-light-500 dark:text-primary-dark-500">
                          <CheckIcon className="h-5 w-5" aria-hidden="true" />
                        </span>
                      ) : null}
                    </>
                  )}
                </ListboxOption>
              ))}
            </ListboxOptions>
          </Transition>
        </div>
      </Listbox>
    </div>
  );
}
