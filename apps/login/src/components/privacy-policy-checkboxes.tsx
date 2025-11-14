"use client";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import Link from "next/link";
import { useState } from "react";
import { Checkbox } from "./checkbox";
import { Translated } from "./translated";

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

  // Helper function to check if all required checkboxes are accepted
  const checkAllAccepted = (newState: AcceptanceState) => {
    const hasTosLink = !!legal?.tosLink;
    const hasPrivacyLink = !!legal?.privacyPolicyLink;

    // Check that all required checkboxes are accepted
    return (
      (!hasTosLink || newState.tosAccepted) &&
      (!hasPrivacyLink || newState.privacyPolicyAccepted)
    );
  };

  return (
    <>
      <p className="mt-4 flex flex-row items-center text-sm text-text-light-secondary-500 dark:text-text-dark-secondary-500">
        <Translated i18nKey="agreeTo" namespace="register" />
        {legal?.helpLink && (
          <span>
            <Link href={legal.helpLink} target="_blank">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="ml-1 h-5 w-5"
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
            checked={acceptanceState.tosAccepted}
            value={"tos"}
            onChangeVal={(checked: boolean) => {
              const newState = {
                ...acceptanceState,
                tosAccepted: checked,
              };
              setAcceptanceState(newState);
              onChange(checkAllAccepted(newState));
            }}
            data-testid="tos-checkbox"
          />

          <div className="mr-4 w-[28rem]">
            <p className="text-sm text-text-light-500 dark:text-text-dark-500">
              <Link href={legal.tosLink} className="underline" target="_blank">
                <Translated i18nKey="termsOfService" namespace="register" />
              </Link>
            </p>
          </div>
        </div>
      )}
      {legal?.privacyPolicyLink && (
        <div className="mt-4 flex items-center">
          <Checkbox
            className="mr-4"
            checked={acceptanceState.privacyPolicyAccepted}
            value={"privacypolicy"}
            onChangeVal={(checked: boolean) => {
              const newState = {
                ...acceptanceState,
                privacyPolicyAccepted: checked,
              };
              setAcceptanceState(newState);
              onChange(checkAllAccepted(newState));
            }}
            data-testid="privacy-policy-checkbox"
          />

          <div className="mr-4 w-[28rem]">
            <p className="text-sm text-text-light-500 dark:text-text-dark-500">
              <Link href={legal.privacyPolicyLink} className="underline" target="_blank">
                <Translated i18nKey="privacyPolicy" namespace="register" />
              </Link>
            </p>
          </div>
        </div>
      )}
    </>
  );
}
