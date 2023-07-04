"use client";

import { RadioGroup } from "@headlessui/react";

export const methods = [
  {
    name: "Passkeys",
    description: "Authenticate with your device.",
  },
  {
    name: "Password",
    description: "Authenticate with a password",
  },
];

export default function AuthenticationMethodRadio({
  selected,
  selectionChanged,
}: {
  selected: any;
  selectionChanged: (value: any) => void;
}) {
  return (
    <div className="w-full">
      <div className="mx-auto w-full max-w-md">
        <RadioGroup value={selected} onChange={selectionChanged}>
          <RadioGroup.Label className="sr-only">Server size</RadioGroup.Label>
          <div className="grid grid-cols-2 space-x-2">
            {methods.map((method) => (
              <RadioGroup.Option
                key={method.name}
                value={method}
                className={({ active, checked }) =>
                  `${
                    active
                      ? "h-full ring-2 ring-opacity-60 ring-primary-light-500 dark:ring-white/20"
                      : "h-full "
                  }
                    ${
                      checked
                        ? "bg-background-light-400 dark:bg-background-dark-400"
                        : "bg-background-light-400 dark:bg-background-dark-400"
                    }
                      relative border boder-divider-light dark:border-divider-dark flex cursor-pointer rounded-lg px-5 py-4 focus:outline-none hover:shadow-lg dark:hover:bg-white/10`
                }
              >
                {({ active, checked }) => (
                  <>
                    <div className="flex w-full items-center justify-between">
                      <div className="flex items-center">
                        <div className="text-sm">
                          <RadioGroup.Label
                            as="p"
                            className={`font-medium  ${checked ? "" : ""}`}
                          >
                            {method.name}
                          </RadioGroup.Label>
                          <RadioGroup.Description
                            as="span"
                            className={`text-xs text-opacity-80 dark:text-opacity-80 inline ${
                              checked ? "" : ""
                            }`}
                          >
                            {method.description}
                            <span aria-hidden="true">&middot;</span>{" "}
                          </RadioGroup.Description>
                        </div>
                      </div>
                      {checked && (
                        <div className="shrink-0 text-white">
                          <CheckIcon className="h-6 w-6" />
                        </div>
                      )}
                    </div>
                  </>
                )}
              </RadioGroup.Option>
            ))}
          </div>
        </RadioGroup>
      </div>
    </div>
  );
}

function CheckIcon(props: any) {
  return (
    <svg viewBox="0 0 24 24" fill="none" {...props}>
      <circle
        className="fill-current text-black/50 dark:text-white/50"
        cx={12}
        cy={12}
        r={12}
        opacity="0.2"
      />
      <path
        d="M7 13l3 3 7-7"
        className="stroke-black dark:stroke-white"
        strokeWidth={1.5}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
