"use client";

import Link from "next/link";
import { BadgeState, StateBadge } from "./StateBadge";
import clsx from "clsx";
import { CheckIcon } from "@heroicons/react/24/outline";
import { EMAIL, SMS, TOTP, U2F } from "./AuthMethods";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2beta/user_service_pb";

type Props = {
  loginName?: string;
  sessionId?: string;
  authRequestId?: string;
  organization?: string;
  userMethods: AuthenticationMethodType[];
};

export default function ChooseSecondFactor({
  loginName,
  sessionId,
  authRequestId,
  organization,
  userMethods,
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

  return (
    <div className="grid grid-cols-1 gap-5 w-full pt-4">
      {userMethods.map((method, i) => {
        return (
          <div key={"method-" + i}>
            {method === AuthenticationMethodType.TOTP &&
              TOTP(false, "/otp/time-based?" + params)}
            {method === AuthenticationMethodType.U2F &&
              U2F(false, "/u2f?" + params)}
            {method === AuthenticationMethodType.OTP_EMAIL &&
              EMAIL(false, "/otp/email?" + params)}
            {method === AuthenticationMethodType.OTP_SMS &&
              SMS(false, "/otp/sms?" + params)}
          </div>
        );
      })}
    </div>
  );
}
