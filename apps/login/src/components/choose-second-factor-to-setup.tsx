"use client";

import { skipMFAAndContinueWithNextUrl } from "@/lib/server/session";
import { LoginSettings, SecondFactorType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useRouter } from "next/navigation";
import { EMAIL, SMS, TOTP, U2F } from "./auth-methods";
import { Translated } from "./translated";
import { useState } from "react";
import { Alert } from "./alert";

type Props = {
  userId: string;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
  loginSettings: LoginSettings;
  userMethods: AuthenticationMethodType[];
  checkAfter: boolean;
  phoneVerified: boolean;
  emailVerified: boolean;
  force: boolean;
};

export function ChooseSecondFactorToSetup({
  userId,
  loginName,
  sessionId,
  requestId,
  organization,
  loginSettings,
  userMethods,
  checkAfter,
  phoneVerified,
  emailVerified,
  force,
}: Props) {
  const router = useRouter();
  const params = new URLSearchParams({});

  const [error, setError] = useState<string>("");

  if (loginName) {
    params.append("loginName", loginName);
  }
  if (sessionId) {
    params.append("sessionId", sessionId);
  }
  if (requestId) {
    params.append("requestId", requestId);
  }
  if (organization) {
    params.append("organization", organization);
  }
  if (checkAfter) {
    params.append("checkAfter", "true");
  }

  return (
    <>
      <div className="grid w-full grid-cols-1 gap-5 pt-4">
        {loginSettings.secondFactors.map((factor) => {
          switch (factor) {
            case SecondFactorType.OTP:
              return TOTP(userMethods.includes(AuthenticationMethodType.TOTP), "/otp/time-based/set?" + params);
            case SecondFactorType.U2F:
              return U2F(userMethods.includes(AuthenticationMethodType.U2F), "/u2f/set?" + params);
            case SecondFactorType.OTP_EMAIL:
              return (
                emailVerified && EMAIL(userMethods.includes(AuthenticationMethodType.OTP_EMAIL), "/otp/email/set?" + params)
              );
            case SecondFactorType.OTP_SMS:
              return phoneVerified && SMS(userMethods.includes(AuthenticationMethodType.OTP_SMS), "/otp/sms/set?" + params);
            default:
              return null;
          }
        })}
      </div>
      {!force && (
        <button
          className="text-sm transition-all hover:text-primary-light-500 dark:hover:text-primary-dark-500"
          onClick={async () => {
            const skipResponse = await skipMFAAndContinueWithNextUrl({
              userId,
              loginName,
              sessionId,
              organization,
              requestId,
            });

            if (skipResponse && "error" in skipResponse && skipResponse.error) {
              setError(skipResponse.error);
              return;
            }

            // For regular flows (non-OIDC/SAML), return URL for client-side navigation
            if ("redirect" in skipResponse && skipResponse.redirect) {
              router.push(skipResponse.redirect);
            }
          }}
          type="button"
          data-testid="reset-button"
        >
          <Translated i18nKey="set.skip" namespace="mfa" />
        </button>
      )}
      {error && (
        <div className="py-4" data-testid="error">
          <Alert>{error}</Alert>
        </div>
      )}
    </>
  );
}
