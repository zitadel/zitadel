"use client";

import { skipMFAAndContinueWithNextUrl } from "@/lib/server/session";
import { LoginSettings, SecondFactorType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useRouter } from "next/navigation";
import { EMAIL, SMS, TOTP, U2F } from "./auth-methods";
import { Translated } from "./translated";
import { useState } from "react";
import { handleServerActionResponse } from "@/lib/client";
import { AutoSubmitForm } from "./auto-submit-form";
import { Alert } from "./alert";
import { trackEvent, MixpanelEvents } from "@/lib/mixpanel";

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
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

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
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <div className="grid w-full grid-cols-1 gap-5 pt-4">
        {loginSettings.secondFactors.map((factor) => {
          const trackSetupSelection = () => trackEvent(MixpanelEvents.mfa_setup_method_selected, { factor: String(factor) });
          switch (factor) {
            case SecondFactorType.OTP:
              return <div onClick={trackSetupSelection} key={factor}>{TOTP(userMethods.includes(AuthenticationMethodType.TOTP), "/otp/time-based/set?" + params)}</div>;
            case SecondFactorType.U2F:
              return <div onClick={trackSetupSelection} key={factor}>{U2F(userMethods.includes(AuthenticationMethodType.U2F), "/u2f/set?" + params)}</div>;
            case SecondFactorType.OTP_EMAIL:
              return (
                emailVerified && <div onClick={trackSetupSelection} key={factor}>{EMAIL(userMethods.includes(AuthenticationMethodType.OTP_EMAIL), "/otp/email/set?" + params)}</div>
              );
            case SecondFactorType.OTP_SMS:
              return phoneVerified && <div onClick={trackSetupSelection} key={factor}>{SMS(userMethods.includes(AuthenticationMethodType.OTP_SMS), "/otp/sms/set?" + params)}</div>;
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

            handleServerActionResponse(skipResponse, router, setSamlData, setError);
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
