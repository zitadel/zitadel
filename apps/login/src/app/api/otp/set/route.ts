import {
  SessionCookie,
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
} from "@/utils/cookies";
import { setSessionAndUpdateCookie } from "@/utils/session";
import { NextRequest, NextResponse, userAgent } from "next/server";
import { Checks } from "@zitadel/proto/zitadel/session/v2beta/session_service_pb";
import { PlainMessage } from "@zitadel/client2";

export async function POST(request: NextRequest) {
  const body = await request.json();

  if (body) {
    const { loginName, sessionId, organization, authRequestId, code, method } =
      body;

    const recentPromise: Promise<SessionCookie> = sessionId
      ? getSessionCookieById(sessionId).catch((error) => {
          return Promise.reject(error);
        })
      : loginName
        ? getSessionCookieByLoginName(loginName, organization).catch(
            (error) => {
              return Promise.reject(error);
            },
          )
        : getMostRecentSessionCookie().catch((error) => {
            return Promise.reject(error);
          });

    return recentPromise
      .then((recent) => {
        const checks: PlainMessage<Checks> = {};

        if (method === "time-based") {
          checks.totp = {
            code,
          };
        } else if (method === "sms") {
          checks.otpSms = {
            code,
          };
        } else if (method === "email") {
          checks.otpEmail = {
            code,
          };
        }

        return setSessionAndUpdateCookie(
          recent,
          checks,
          undefined,
          authRequestId,
        ).then((session) => {
          return NextResponse.json({
            sessionId: session.id,
            factors: session.factors,
            challenges: session.challenges,
          });
        });
      })
      .catch((error) => {
        return NextResponse.json({ details: error }, { status: 500 });
      });
  } else {
    return NextResponse.json(
      { details: "Request body is missing" },
      { status: 400 },
    );
  }
}
