import {
  server,
  deleteSession,
  listHumanAuthFactors,
  getSession,
} from "#/lib/zitadel";
import {
  SessionCookie,
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "#/utils/cookies";
import {
  createSessionAndUpdateCookie,
  createSessionForIdpAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "#/utils/session";
import { Checks, RequestChallenges } from "@zitadel/server";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const {
      userId,
      idpIntent,
      loginName,
      password,
      organization,
      authRequestId,
    } = body;

    if (userId && idpIntent) {
      return createSessionForIdpAndUpdateCookie(
        userId,
        idpIntent,
        organization,
        authRequestId
      ).then((session) => {
        return NextResponse.json(session);
      });
    } else {
      return createSessionAndUpdateCookie(
        loginName,
        password,
        undefined,
        organization,
        authRequestId
      ).then((session) => {
        return NextResponse.json(session);
      });
    }
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
    const {
      loginName,
      sessionId,
      organization,
      checks,
      authRequestId,
      challenges,
    } = body;

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

    const domain: string = request.nextUrl.hostname;

    if (challenges && challenges.webAuthN && !challenges.webAuthN.domain) {
      challenges.webAuthN.domain = domain;
    }

    return recentPromise
      .then((recent) => {
        return setSessionAndUpdateCookie(
          recent,
          checks,
          challenges,
          authRequestId
        ).then(async (session) => {
          // if password, check if user has MFA methods
          let authFactors;
          if (checks && checks.password && session.factors?.user?.id) {
            const response = await listHumanAuthFactors(
              server,
              session.factors?.user?.id
            );
            if (response.result && response.result.length) {
              authFactors = response.result;
            }
          }
          return NextResponse.json({
            sessionId: session.id,
            factors: session.factors,
            challenges: session.challenges,
            authFactors,
          });
        });
      })
      .catch((error) => {
        console.error(error);
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
