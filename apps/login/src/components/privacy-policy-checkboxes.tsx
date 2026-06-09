"use client";
import { resolveLocalizedLegalLink } from "@/lib/legal-links";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { useLocale } from "next-intl";
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
  const locale = useLocale();
  const [acceptanceState, setAcceptanceState] = useState<AcceptanceState>({
    tosAccepted: false,
    privacyPolicyAccepted: false,
  });
  const helpLink = resolveLocalizedLegalLink(legal?.helpLink, locale);
  const tosLink = resolveLocalizedLegalLink(legal?.tosLink, locale);
  const privacyPolicyLink = resolveLocalizedLegalLink(legal?.privacyPolicyLink, locale);

  // Helper function to check if all required checkboxes are accepted
  const checkAllAccepted = (newState: AcceptanceState) => {
    const hasTosLink = !!tosLink;
    const hasPrivacyLink = !!privacyPolicyLink;

    // Check that all required checkboxes are accepted
    return (!hasTosLink || newState.tosAccepted) && (!hasPrivacyLink || newState.privacyPolicyAccepted);
  };

  return (
    <>
      <p className="text-text-light-secondary-500 dark:text-text-dark-secondary-500 mt-4 flex flex-row items-center text-sm">
        <Translated i18nKey="agreeTo" namespace="register" />
        {helpLink && (
          <span>
            <Link href={helpLink} target="_blank" aria-label="Open help in a new tab" data-testid="help-link">
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
      {tosLink && (
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
            <p className="text-text-light-500 dark:text-text-dark-500 text-sm">
              <Link href={tosLink} className="underline" target="_blank" data-testid="tos-link">
                <Translated i18nKey="termsOfService" namespace="register" />
              </Link>
            </p>
          </div>
        </div>
      )}
      {privacyPolicyLink && (
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
            <p className="text-text-light-500 dark:text-text-dark-500 text-sm">
              <Link href={privacyPolicyLink} className="underline" target="_blank" data-testid="privacy-policy-link">
                <Translated i18nKey="privacyPolicy" namespace="register" />
              </Link>
            </p>
          </div>
        </div>
      )}
    </>
  );
}
