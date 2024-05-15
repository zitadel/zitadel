"use client";

import { AuthenticationMethodType, LoginSettings } from "@zitadel/server";
import { EMAIL, SMS, TOTP, U2F } from "./AuthMethods";

type Props = {
  loginName?: string;
  sessionId?: string;
  authRequestId?: string;
  organization?: string;
  loginSettings: LoginSettings;
  userMethods: AuthenticationMethodType[];
  checkAfter: boolean;
  phoneVerified: boolean;
  emailVerified: boolean;
};

export default function ChooseSecondFactorToSetup({
  loginName,
  sessionId,
  authRequestId,
  organization,
  loginSettings,
  userMethods,
  checkAfter,
  phoneVerified,
  emailVerified,
}: Props) {
  const params = new URLSearchParams({});

  if (loginName) {
    params.append("loginName", loginName);
  }
  if (sessionId) {
    params.append("sessionId", sessionId);
  }
  if (authRequestId) {
    params.append("authRequestId", authRequestId);
  }
  if (organization) {
    params.append("organization", organization);
  }
  if (checkAfter) {
    params.append("checkAfter", "true");
  }

  return (
    <div className="grid grid-cols-1 gap-5 w-full pt-4">
      {loginSettings.secondFactors.map((factor, i) => {
        return factor === 1
          ? TOTP(userMethods.includes(4), "/otp/time-based/set?" + params)
          : factor === 2
            ? U2F(userMethods.includes(5), "/u2f/set?" + params)
            : factor === 3 && emailVerified
              ? EMAIL(userMethods.includes(7), "/otp/email/set?" + params)
              : factor === 4 && phoneVerified
                ? SMS(userMethods.includes(6), "/otp/sms/set?" + params)
                : null;
      })}
    </div>
  );
}
