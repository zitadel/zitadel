import { server, deleteSession } from "#/lib/zitadel";
import {
  SessionCookie,
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "#/utils/cookies";
import {
  createSessionAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "#/utils/session";
import { RequestChallenges } from "@zitadel/server";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName, password } = body;

    // const domain: string = request.nextUrl.hostname;

    return createSessionAndUpdateCookie(
      loginName,
      password,
      undefined,
      undefined
    ).then((session) => {
      return NextResponse.json(session);
    });
  } else {
    return NextResponse.json(
      { details: "Session could not be created" },
      { status: 500 }
    );
  }
}

/**
 *
 * @param request password for the most recent session
 * @returns the updated most recent Session with the added password
 */
export async function PUT(request: NextRequest) {
  const body = await request.json();

  if (body) {
    const { loginName, password, passkey, authRequestId } = body;
    const challenges: RequestChallenges = body.challenges;

    const recentPromise: Promise<SessionCookie> = loginName
      ? getSessionCookieByLoginName(loginName).catch((error) => {
          return Promise.reject(error);
        })
      : getMostRecentSessionCookie().catch((error) => {
          return Promise.reject(error);
        });

    const domain: string = request.nextUrl.hostname;

    if (challenges.webAuthN && !challenges.webAuthN.domain) {
      challenges.webAuthN.domain = domain;
    }

    return recentPromise
      .then((recent) => {
        return setSessionAndUpdateCookie(
          recent.id,
          recent.token,
          recent.loginName,
          password,
          passkey,
          challenges,
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

/**
 *
 * @param request id of the session to be deleted
 */
export async function DELETE(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const id = searchParams.get("id");
  if (id) {
    const session = await getSessionCookieById(id);

    return deleteSession(server, session.id, session.token)
      .then(() => {
        return removeSessionFromCookie(session)
          .then(() => {
            return NextResponse.json({});
          })
          .catch((error) => {
            return NextResponse.json(
              { details: "could not set cookie" },
              { status: 500 }
            );
          });
      })
      .catch((error) => {
        return NextResponse.json(
          { details: "could not delete session" },
          { status: 500 }
        );
      });
  } else {
    return NextResponse.error();
  }
}
