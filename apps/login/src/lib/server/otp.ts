"use server";

import { setSessionAndUpdateCookie } from "@/lib/server/cookie";
import { create } from "@zitadel/client";
import {
  CheckOTPSchema,
  ChecksSchema,
  CheckTOTPSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
} from "../cookies";

export type SetOTPCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  authRequestId?: string;
  code: string;
  method: string;
};

export async function setOTP(command: SetOTPCommand) {
  const recentPromise = command.sessionId
    ? getSessionCookieById({ sessionId: command.sessionId }).catch((error) => {
        return Promise.reject(error);
      })
    : command.loginName
      ? getSessionCookieByLoginName({
          loginName: command.loginName,
          organization: command.organization,
        }).catch((error) => {
          return Promise.reject(error);
        })
      : getMostRecentSessionCookie().catch((error) => {
          return Promise.reject(error);
        });

  return recentPromise.then((recent) => {
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

    return setSessionAndUpdateCookie(
      recent,
      checks,
      undefined,
      command.authRequestId,
    ).then((session) => {
      return {
        sessionId: session.id,
        factors: session.factors,
        challenges: session.challenges,
      };
    });
  });
}
