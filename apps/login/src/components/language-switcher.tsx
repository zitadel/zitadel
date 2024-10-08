"use client";

import { Listbox, Transition } from "@headlessui/react";
import { CheckIcon, ChevronUpDownIcon } from "@heroicons/react/24/outline";
import { usePathname, useRouter } from "next/navigation";
import { Fragment, useState } from "react";
import { useTranslation } from "react-i18next";
import i18nConfig from "../../i18nConfig";

interface Lang {
  id: number;
  name: string;
  img: string;
  code: string;
}

const LANGS: Lang[] = [
  {
    id: 1,
    name: "Deutsch",
    code: "de",
    img: "/images/flags/de.png",
  },
  {
    id: 2,
    name: "Italiano",
    code: "it",
    img: "/images/flags/it.png",
  },
  {
    id: 3,
    name: "English",
    code: "en",
    img: "/images/flags/us.png",
  },
];

type Props = {
  locale: string;
};

export default function LanguageSwitcher({ locale }: Props) {
  const { i18n } = useTranslation();

  const currentLocale = locale || i18n.language || i18nConfig.defaultLocale;

  const [selected, setSelected] = useState(
    LANGS.find((l) => l.code === currentLocale) || LANGS[0],
  );

  const router = useRouter();
  const currentPathname = usePathname();

  const handleChange = (language: Lang) => {
    setSelected(language);
    const newLocale = language.code;
    // set cookie for next-i18n-router
    const days = 30;
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = date.toUTCString();
    document.cookie = `NEXT_LOCALE=${newLocale};expires=${expires};path=/`;

    // redirect to the new locale path
    if (
      currentLocale === i18nConfig.defaultLocale &&
      !i18nConfig.prefixDefault
    ) {
      router.push("/" + newLocale + currentPathname);
    } else {
      router.push(
        currentPathname.replace(`/${currentLocale}`, `/${newLocale}`),
      );
    }

    router.refresh();
  };

  return (
    <div className="mb-4 w-32">
      <Listbox value={selected} onChange={handleChange}>
        <div className="relative mt-1">
          <Listbox.Button className="relative w-full cursor-default rounded-lg border border-divider-light bg-white py-2 pl-3 pr-10 text-left focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-orange-300 sm:text-sm">
            <span className="block truncate">{selected.name}</span>
            <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <ChevronUpDownIcon
                className="h-5 w-5 text-gray-400"
                aria-hidden="true"
              />
            </span>
          </Listbox.Button>
          <Transition
            as={Fragment}
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Listbox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm">
              {LANGS.map((lang, index) => (
                <Listbox.Option
                  key={lang.code}
                  className={({ active }) =>
                    `relative cursor-default select-none py-2 pl-10 pr-4 ${
                      active ? "bg-amber-100 text-amber-900" : "text-gray-900"
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
                        <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-amber-600">
                          <CheckIcon className="h-5 w-5" aria-hidden="true" />
                        </span>
                      ) : null}
                    </>
                  )}
                </Listbox.Option>
              ))}
            </Listbox.Options>
          </Transition>
        </div>
      </Listbox>
    </div>
  );
}
