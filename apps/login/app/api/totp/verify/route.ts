import {
  SessionCookie,
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
} from "#/utils/cookies";
import { setSessionAndUpdateCookie } from "#/utils/session";
import { NextRequest, NextResponse, userAgent } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();

  if (body) {
    const { loginName, sessionId, organization, authRequestId, code } = body;

    const recentPromise: Promise<SessionCookie> = sessionId
      ? getSessionCookieById(sessionId).catch((error) => {
          return Promise.reject(error);
        })
      : loginName
      ? getSessionCookieByLoginName(loginName, organization).catch((error) => {
          return Promise.reject(error);
        })
      : getMostRecentSessionCookie().catch((error) => {
          return Promise.reject(error);
        });

    return recentPromise
      .then((recent) => {
        return setSessionAndUpdateCookie(
          recent,
          undefined,
          undefined,
          undefined,
          code,
          authRequestId
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
      { status: 400 }
    );
  }
}
