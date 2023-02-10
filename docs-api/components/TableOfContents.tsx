import React, { useEffect, useState } from "react";
import Link from "next/link";
import { ZitadelLogo } from "./ZitadelLogo";
import { Disclosure } from "@headlessui/react";
import { ChevronDownIcon } from "@heroicons/react/24/outline";
import Theme from "./Theme";
import { useRouter } from "next/router";

export function TableOfContents({ toc }) {
  const router = useRouter();

  const items = toc.filter(
    (item) =>
      item.id && (item.level === 2 || item.level === 3 || item.level === 4)
  );

  if (items.length <= 1) {
    return null;
  }

  return (
    <nav className="sticky bg-gray-50 dark:bg-background-dark-600 top-0 h-screen border-box overflow-y-auto bottom-0 flex-shrink-0 w-64 xl:w-72 pr-4 xl:px-6 border-r border-border-light dark:border-border-dark flex flex-col">
      <div className="flex flex-col relative">
        <div className="z-10 sticky h-16 top-0 left-0 right-0">
          <div className="pl-4 pt-4 pb-2 bg-gray-50 dark:bg-background-dark-600 flex items-center justify-between">
            <Link className="" href="/">
              <ZitadelLogo />
            </Link>

            <div className="relative">
              <Theme />
            </div>
          </div>
        </div>
        <div className="sticky top-16 left-0 h-8 bg-gradient-to-b from-gray-50 dark:from-background-dark-600 to-transparent dark:to-transparent"></div>
        <ul className="flex-1 flex flex-col pb-8">
          {items.map((item, i) => {
            const href = `#${item.id}`;

            const active = `#${router.asPath.split("#")[1]}` === href;

            if (item.level === 2 && i < items.length) {
              const remaining = items.slice(i + 1);

              const nextSection =
                i +
                (remaining.findIndex((i) => i.level === 2) ?? remaining.length);
              const subItems = items.slice(i + 1, nextSection + 1);
              items.splice(i + 1, subItems.length);

              return (
                <Disclosure key={`menu_${i}`}>
                  {({ open }) => (
                    <>
                      <Disclosure.Button className="pl-4 flex w-full justify-between rounded-lg py-2 pt-6 text-left text-sm font-medium text-gray-800 dark:text-gray-200 focus:outline-none focus-visible:ring focus-visible:ring-purple-500 focus-visible:ring-opacity-75">
                        <span className="uppercase text-xs">{item.title}</span>
                        <ChevronDownIcon
                          className={`${
                            open ? "rotate-180 transform" : ""
                          } h-4 w-4 text-gray-400 dark:text-gray-200`}
                        />
                      </Disclosure.Button>
                      <Disclosure.Panel className="pt-1 text-sm text-gray-500">
                        {subItems.map((subitem, j) => {
                          const sub_href = `#${subitem.id}`;

                          const sub_active =
                            `#${router.asPath.split("#")[1]}` === sub_href;

                          return (
                            <li
                              key={`sub_${i}_${j}_${subitem.title}`}
                              className={[
                                subitem.level === 3 ? "pl-4" : undefined,
                                subitem.level === 4 ? "pl-8" : undefined,
                                sub_active
                                  ? "bg-primary-light-500/5 dark:bg-white/5"
                                  : "",
                                "py-1 text-sm min-h-8 flex items-center rounded-r-md xl:rounded-l-md",
                              ]
                                .filter(Boolean)
                                .join(" ")}
                            >
                              <Link
                                className={[
                                  "text-sm",
                                  sub_active
                                    ? "text-primary-light-500 dark:text-primary-dark-400"
                                    : "text-gray-500 dark:text-gray-400 hover:text-black hover:dark:text-white",
                                ]
                                  .filter(Boolean)
                                  .join(" ")}
                                href={sub_href}
                              >
                                {subitem.title}
                              </Link>
                            </li>
                          );
                        })}
                      </Disclosure.Panel>
                    </>
                  )}
                </Disclosure>
              );
            } else {
              return (
                <li
                  key={`sub_${i}_${item.title}`}
                  className={[
                    "pl-4 rounded-r-md xl:rounded-l-md",
                    active ? "bg-primary-light-500/5 dark:bg-white/5" : "",
                    item.level === 3 ? "py-1" : undefined,
                  ]
                    .filter(Boolean)
                    .join(" ")}
                >
                  <Link
                    className={[
                      "text-sm ",
                      active
                        ? "text-primary-light-500 dark:text-primary-dark-400"
                        : "text-gray-500 dark:text-gray-400 hover:text-black hover:dark:text-white",
                    ]
                      .filter(Boolean)
                      .join("")}
                    href={href}
                  >
                    {item.title}
                  </Link>
                </li>
              );
            }
          })}
        </ul>
      </div>
    </nav>
  );
}
