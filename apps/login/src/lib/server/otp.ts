"use server";

import { setSessionAndUpdateCookie } from "@/utils/session";
import {
  CheckOTPSchema,
  ChecksSchema,
  CheckTOTPSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { create } from "@zitadel/client";
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
  const { loginName, sessionId, organization, authRequestId, code, method } =
    command;

  const recentPromise = sessionId
    ? getSessionCookieById({ sessionId }).catch((error) => {
        return Promise.reject(error);
      })
    : loginName
      ? getSessionCookieByLoginName({ loginName, organization }).catch(
          (error) => {
            return Promise.reject(error);
          },
        )
      : getMostRecentSessionCookie().catch((error) => {
          return Promise.reject(error);
        });

  return recentPromise.then((recent) => {
    const checks = create(ChecksSchema, {});

    if (method === "time-based") {
      checks.totp = create(CheckTOTPSchema, {
        code,
      });
    } else if (method === "sms") {
      checks.otpSms = create(CheckOTPSchema, {
        code,
      });
    } else if (method === "email") {
      checks.otpEmail = create(CheckOTPSchema, {
        code,
      });
    }

    return setSessionAndUpdateCookie(
      recent,
      checks,
      undefined,
      authRequestId,
    ).then((session) => {
      return {
        sessionId: session.id,
        factors: session.factors,
        challenges: session.challenges,
      };
    });
  });
}
