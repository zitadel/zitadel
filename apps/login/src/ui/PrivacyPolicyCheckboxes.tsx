"use client";
import React, { useState } from "react";
import Link from "next/link";
import { Checkbox } from "./Checkbox";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2beta/legal_settings_pb";

type Props = {
  legal: LegalAndSupportSettings;
  onChange: (allAccepted: boolean) => void;
};

type AcceptanceState = {
  tosAccepted: boolean;
  privacyPolicyAccepted: boolean;
};

export function PrivacyPolicyCheckboxes({ legal, onChange }: Props) {
  const [acceptanceState, setAcceptanceState] = useState<AcceptanceState>({
    tosAccepted: false,
    privacyPolicyAccepted: false,
  });

  return (
    <>
      <p className="flex flex-row items-center text-text-light-secondary-500 dark:text-text-dark-secondary-500 mt-4 text-sm">
        To register you must agree to the terms and conditions
        {legal?.helpLink && (
          <span>
            <Link href={legal.helpLink} target="_blank">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="ml-1 w-5 h-5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 5.25h.008v.008H12v-.008z"
                />
              </svg>
            </Link>
          </span>
        )}
      </p>
      {legal?.tosLink && (
        <div className="mt-4 flex items-center">
          <Checkbox
            className="mr-4"
            checked={false}
            value={"privacypolicy"}
            onChangeVal={(checked: boolean) => {
              setAcceptanceState({
                ...acceptanceState,
                tosAccepted: checked,
              });
              onChange(checked && acceptanceState.privacyPolicyAccepted);
            }}
          />

          <div className="mr-4 w-[28rem]">
            <p className="text-sm text-text-light-500 dark:text-text-dark-500">
              Agree&nbsp;
              <Link href={legal.tosLink} className="underline" target="_blank">
                Terms of Service
              </Link>
            </p>
          </div>
        </div>
      )}
      {legal?.privacyPolicyLink && (
        <div className="mt-4 flex items-center">
          <Checkbox
            className="mr-4"
            checked={false}
            value={"tos"}
            onChangeVal={(checked: boolean) => {
              setAcceptanceState({
                ...acceptanceState,
                privacyPolicyAccepted: checked,
              });
              onChange(checked && acceptanceState.tosAccepted);
            }}
          />

          <div className="mr-4 w-[28rem]">
            <p className="text-sm text-text-light-500 dark:text-text-dark-500">
              Agree&nbsp;
              <Link
                href={legal.privacyPolicyLink}
                className="underline"
                target="_blank"
              >
                Privacy Policy
              </Link>
            </p>
          </div>
        </div>
      )}
    </>
  );
}
