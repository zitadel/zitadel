"use client";

import { usePathname } from "next/navigation";
import { Fragment } from "react";

type Props = {
  domain: string;
};

export function AddressBar({ domain }: Props) {
  const pathname = usePathname();

  return (
    <div className="flex items-center space-x-2 p-3.5 lg:px-5 lg:py-3 overflow-hidden">
      <div className="text-gray-600">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-4"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fillRule="evenodd"
            d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
            clipRule="evenodd"
          />
        </svg>
      </div>
      <div className="flex space-x-1 text-sm font-medium">
        <div className="max-w-[150px] px-2 overflow-hidden text-gray-500  text-ellipsis">
          <span className="whitespace-nowrap">{domain}</span>
        </div>
        {pathname ? (
          <>
            <span className="text-gray-600">/</span>
            {pathname
              .split("/")
              .slice(1)
              .filter((s) => !!s)
              .map((segment) => {
                return (
                  <Fragment key={segment}>
                    <span>
                      <span
                        key={segment}
                        className="animate-[highlight_1s_ease-in-out_1] rounded-full px-1.5 py-0.5 text-gray-800 dark:text-gray-100"
                      >
                        {segment}
                      </span>
                    </span>

                    <span className="text-gray-600">/</span>
                  </Fragment>
                );
              })}
          </>
        ) : null}
      </div>
    </div>
  );
}
