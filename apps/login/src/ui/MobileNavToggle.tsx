"use client";

import { Bars3Icon, XMarkIcon } from "@heroicons/react/24/solid";
import clsx from "clsx";
import React from "react";

const MobileNavContext = React.createContext<
  [boolean, React.Dispatch<React.SetStateAction<boolean>>] | undefined
>(undefined);

export function MobileNavContextProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isOpen, setIsOpen] = React.useState(false);
  return (
    <MobileNavContext.Provider value={[isOpen, setIsOpen]}>
      {children}
    </MobileNavContext.Provider>
  );
}

export function useMobileNavToggle() {
  const context = React.useContext(MobileNavContext);
  if (context === undefined) {
    throw new Error(
      "useMobileNavToggle must be used within a MobileNavContextProvider",
    );
  }
  return context;
}

export function MobileNavToggle({ children }: { children: React.ReactNode }) {
  const [isOpen, setIsOpen] = useMobileNavToggle();

  return (
    <>
      <button
        type="button"
        className="group absolute right-0 top-0 flex h-14 items-center space-x-2 px-4 lg:hidden"
        onClick={() => setIsOpen(!isOpen)}
      >
        <div className="font-medium text-text-light-secondary-500 dark:text-text-dark-secondary-500 group-hover:text-text-light-500 dark:group-hover:text-text-dark-500">
          Menu
        </div>
        {isOpen ? (
          <XMarkIcon className="block w-6" />
        ) : (
          <Bars3Icon className="block w-6" />
        )}
      </button>

      <div
        className={clsx("overflow-y-auto lg:static lg:block", {
          "fixed inset-x-0 bottom-0 top-14 bg-gray-900": isOpen,
          hidden: !isOpen,
        })}
      >
        {children}
      </div>
    </>
  );
}
