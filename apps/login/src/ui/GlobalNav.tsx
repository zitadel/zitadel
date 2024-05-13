"use client";

import { demos, type Item } from "@/lib/demos";
import { ZitadelLogo } from "@/ui/ZitadelLogo";
import Link from "next/link";
import { useSelectedLayoutSegment, usePathname } from "next/navigation";
import clsx from "clsx";
import { Bars3Icon, XMarkIcon } from "@heroicons/react/24/solid";
import { useState } from "react";
import Theme from "./Theme";

export function GlobalNav() {
  const [isOpen, setIsOpen] = useState(false);
  const close = () => setIsOpen(false);

  return (
    <div className="fixed top-0 z-10 flex w-full flex-col border-b border-divider-light dark:border-divider-dark bg-white/80 dark:bg-black/80 lg:bottom-0 lg:z-auto lg:w-72 lg:border-r">
      <div className="flex h-14 items-center py-4 px-4 lg:h-auto">
        <Link
          href="/"
          className="group flex w-full items-center space-x-2.5"
          onClick={close}
        >
          <div className="">
            <ZitadelLogo />
          </div>

          <h2 className="text-blue-500 font-bold uppercase transform translate-y-2 text-sm">
            Login
          </h2>
        </Link>
      </div>

      <div className="absolute right-0 top-0 flex flex-row items-center lg:hidden">
        <Theme />
        <button
          type="button"
          className="group flex h-14 items-center space-x-2 px-4"
          onClick={() => setIsOpen(!isOpen)}
        >
          <div className="font-medium text-text-light-secondary-500 group-hover:text-text-light-500 dark:text-text-dark-secondary-500 dark:group-hover:text-text-dark-500">
            Menu
          </div>
          {isOpen ? (
            <XMarkIcon className="block w-6 " />
          ) : (
            <Bars3Icon className="block w-6 " />
          )}
        </button>
      </div>

      <div
        className={clsx(
          "overflow-y-auto lg:static lg:flex lg:flex-col justify-between h-full",
          {
            "fixed inset-x-0 bottom-0 top-14 mt-px bg-white/80 dark:bg-black/80 backdrop-blur-lg":
              isOpen,
            hidden: !isOpen,
          },
        )}
      >
        <nav
          className={`space-y-6 px-4 py-5 ${
            isOpen ? "text-center lg:text-left" : ""
          }`}
        >
          {demos.map((section) => {
            return (
              <div key={section.name}>
                <div className="mb-2 px-3 text-[11px] font-bold uppercase tracking-wider text-black/40 dark:text-white/40">
                  <div>{section.name}</div>
                </div>

                <div className="space-y-1">
                  {section.items.map((item) => (
                    <GlobalNavItem key={item.slug} item={item} close={close} />
                  ))}
                </div>
              </div>
            );
          })}
        </nav>

        <div className="flex flex-row p-4">
          <Theme />
        </div>
      </div>
    </div>
  );
}

function GlobalNavItem({
  item,
  close,
}: {
  item: Item;
  close: () => false | void;
}) {
  const segment = useSelectedLayoutSegment();
  const pathname = usePathname();

  const isActive = `/${item.slug}` === pathname;

  return (
    <Link
      onClick={close}
      href={`/${item.slug}`}
      className={clsx(
        "block rounded-md px-3 py-2 text-[15px] font-medium text-text-light-500 dark:text-text-dark-500 opacity-60 dark:opacity-60",
        {
          "hover:opacity-100 hover:dark:opacity-100": !isActive,
          "text-text-light-500 dark:text-text-dark-500 opacity-100 dark:opacity-100 font-semibold":
            isActive,
        },
      )}
    >
      {item.name}
    </Link>
  );
}
