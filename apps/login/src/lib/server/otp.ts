"use server";

import { setSessionAndUpdateCookie } from "@/lib/server/cookie";
import { create } from "@zitadel/client";
import {
  CheckOTPSchema,
  ChecksSchema,
  CheckTOTPSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { headers } from "next/headers";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
} from "../cookies";
import { getServiceUrlFromHeaders } from "../service-url";
import { getLoginSettings } from "../zitadel";

export type SetOTPCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  requestId?: string;
  code: string;
  method: string;
};

export async function setOTP(command: SetOTPCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const recentSession = command.sessionId
    ? await getSessionCookieById({ sessionId: command.sessionId }).catch(
        (error) => {
          return Promise.reject(error);
        },
      )
    : command.loginName
      ? await getSessionCookieByLoginName({
          loginName: command.loginName,
          organization: command.organization,
        }).catch((error) => {
          return Promise.reject(error);
        })
      : await getMostRecentSessionCookie().catch((error) => {
          return Promise.reject(error);
        });

  const checks = create(ChecksSchema, {});

  if (command.method === "time-based") {
    checks.totp = create(CheckTOTPSchema, {
      code: command.code,
    });
  } else if (command.method === "sms") {
    checks.otpSms = create(CheckOTPSchema, {
      code: command.code,
    });
  } else if (command.method === "email") {
    checks.otpEmail = create(CheckOTPSchema, {
      code: command.code,
    });
  }

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: command.organization,
  });

  return setSessionAndUpdateCookie(
    recentSession,
    checks,
    undefined,
    command.requestId,
    loginSettings?.secondFactorCheckLifetime,
  ).then((session) => {
    return {
      sessionId: session.id,
      factors: session.factors,
      challenges: session.challenges,
    };
  });
}
