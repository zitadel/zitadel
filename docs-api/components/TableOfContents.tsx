import React, { useEffect, useState } from "react";
import Link from "next/link";
import { ZitadelLogo } from "./ZitadelLogo";
import { Disclosure } from "@headlessui/react";
import { ChevronDownIcon } from "@heroicons/react/24/outline";

export function TableOfContents({ toc }) {
  const items = toc.filter(
    (item) => item.id && (item.level === 2 || item.level === 3)
  );

  if (items.length <= 1) {
    return null;
  }

  return (
    <nav className="sticky top-0 h-screen border-box overflow-y-auto bottom-0 flex-shrink-0 w-60 px-4 border-r border-border-light dark:border-border-dark flex flex-col">
      <div className="flex flex-col relative">
        <div className="z-10 sticky h-16 top-0 left-0 right-0 px-4 pt-4 pb-2 bg-white dark:bg-background-dark-500">
          <Link className="mb-4" href="/">
            <ZitadelLogo />
          </Link>
        </div>
        <div className="sticky top-16 left-0 h-8 bg-gradient-to-b from-white dark:from-background-dark-500 to-transparent dark:to-transparent"></div>
        <ul className="flex-1 flex flex-col pb-8">
          {items.map((item, i) => {
            const href = `#${item.id}`;
            const active = false;
            // const active =
            //   typeof window !== "undefined" && window.location.hash === href;

            if (item.level === 2 && i < items.length) {
              const remaining = items.slice(i + 1);

              const nextSection =
                i +
                (remaining.findIndex((i) => i.level === 2) ?? remaining.length);
              const subItems = items.slice(i + 1, nextSection + 1);
              items.splice(i + 1, subItems.length);
              console.log(subItems);

              return (
                <Disclosure key={`menu_${i}`}>
                  {({ open }) => (
                    <>
                      <Disclosure.Button className="flex w-full justify-between rounded-lg py-2 pt-6 text-left text-sm font-medium text-gray-500 dark:text-gray-200 focus:outline-none focus-visible:ring focus-visible:ring-purple-500 focus-visible:ring-opacity-75">
                        <span className="uppercase text-xs">{item.title}</span>
                        <ChevronDownIcon
                          className={`${
                            open ? "rotate-180 transform" : ""
                          } h-4 w-4 text-gray-500 dark:text-gray-200`}
                        />
                      </Disclosure.Button>
                      <Disclosure.Panel className="pt-1 text-sm text-gray-500">
                        {subItems.map((subitem, j) => {
                          return (
                            <li
                              key={`sub_${i}_${j}_${subitem.title}`}
                              className={[
                                subitem.level === 3 ? "" : undefined,
                                subitem.level === 4 ? "pl-4" : undefined,
                                "py-1 text-sm min-h-8",
                              ]
                                .filter(Boolean)
                                .join(" ")}
                            >
                              <Link
                                className="text-sm  dark:text-gray-400 text-gray-500 dark:bg-background-dark-500 hover:text-black hover:dark:text-white"
                                href={href}
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
                    active
                      ? "text-black dark:bg-background-dark-500 dark:text-primary-dark-500"
                      : "",
                    item.level === 3 ? "py-1" : undefined,
                  ]
                    .filter(Boolean)
                    .join(" ")}
                >
                  <Link
                    className="text-sm text-gray-500 dark:text-gray-400"
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
