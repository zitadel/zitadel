"use client";

import { RadioGroup } from "@headlessui/react";
import { Translated } from "./translated";

export enum AuthenticationMethod {
  Passkey = "passkey",
  Password = "password",
}

export const methods = [
  AuthenticationMethod.Passkey,
  AuthenticationMethod.Password,
];

export function AuthenticationMethodRadio({
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
          <div className="flex flex-row space-x-4">
            {methods.map((method) => (
              <RadioGroup.Option
                key={method}
                value={method}
                data-testid={method + "-radio"}
                className={({ active, checked }) =>
                  `${
                    active
                      ? "ring-2 ring-primary-light-500 ring-opacity-60 dark:ring-white/20"
                      : ""
                  } ${
                    checked
                      ? "bg-background-light-400 ring-2 ring-primary-light-500 dark:bg-background-dark-400 dark:ring-primary-dark-500"
                      : "bg-background-light-400 dark:bg-background-dark-400"
                  } boder-divider-light relative flex h-full flex-1 cursor-pointer rounded-lg border px-5 py-4 hover:shadow-lg focus:outline-none dark:border-divider-dark dark:hover:bg-white/10`
                }
              >
                {({ checked }) => (
                  <>
                    <div className="flex w-full flex-col items-center text-sm">
                      {method === "passkey" && (
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth="1.5"
                          stroke="currentColor"
                          className="mb-3 h-8 w-8"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M7.864 4.243A7.5 7.5 0 0119.5 10.5c0 2.92-.556 5.709-1.568 8.268M5.742 6.364A7.465 7.465 0 004.5 10.5a7.464 7.464 0 01-1.15 3.993m1.989 3.559A11.209 11.209 0 008.25 10.5a3.75 3.75 0 117.5 0c0 .527-.021 1.049-.064 1.565M12 10.5a14.94 14.94 0 01-3.6 9.75m6.633-4.596a18.666 18.666 0 01-2.485 5.33"
                          />
                        </svg>
                      )}
                      {method === "password" && (
                        <svg
                          className="mb-3 h-8 w-8 fill-current"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 24 24"
                        >
                          <title>form-textbox-password</title>
                          <path d="M17,7H22V17H17V19A1,1 0 0,0 18,20H20V22H17.5C16.95,22 16,21.55 16,21C16,21.55 15.05,22 14.5,22H12V20H14A1,1 0 0,0 15,19V5A1,1 0 0,0 14,4H12V2H14.5C15.05,2 16,2.45 16,3C16,2.45 16.95,2 17.5,2H20V4H18A1,1 0 0,0 17,5V7M2,7H13V9H4V15H13V17H2V7M20,15V9H17V15H20M8.5,12A1.5,1.5 0 0,0 7,10.5A1.5,1.5 0 0,0 5.5,12A1.5,1.5 0 0,0 7,13.5A1.5,1.5 0 0,0 8.5,12M13,10.89C12.39,10.33 11.44,10.38 10.88,11C10.32,11.6 10.37,12.55 11,13.11C11.55,13.63 12.43,13.63 13,13.11V10.89Z" />
                        </svg>
                      )}
                      <RadioGroup.Label
                        as="p"
                        className={`font-medium ${checked ? "" : ""}`}
                      >
                        {method === AuthenticationMethod.Passkey && (
                          <Translated
                            i18nKey="methods.passkey"
                            namespace="register"
                          />
                        )}
                        {method === AuthenticationMethod.Password && (
                          <Translated
                            i18nKey="methods.password"
                            namespace="register"
                          />
                        )}
                      </RadioGroup.Label>
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
