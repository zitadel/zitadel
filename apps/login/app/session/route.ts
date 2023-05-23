import { createSession, getSession, server, setSession } from "#/lib/zitadel";
import {
  SessionCookie,
  addSessionToCookie,
  getMostRecentSessionCookie,
  updateSessionCookie,
} from "#/utils/cookies";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName } = body;

    const createdSession = await createSession(server, loginName);

    return getSession(
      server,
      createdSession.sessionId,
      createdSession.sessionToken
    ).then(({ session }) => {
      const sessionCookie: SessionCookie = {
        id: createdSession.sessionId,
        token: createdSession.sessionToken,
        changeDate: session.changeDate,
        loginName: session.factors.user.loginName,
      };
      return addSessionToCookie(sessionCookie).then(() => {
        return NextResponse.json({ factors: session.factors });
      });
    });
  } else {
    return NextResponse.error();
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
    const { password } = body;

    const recent = await getMostRecentSessionCookie();

    return setSession(server, recent.id, recent.token, password)
      .then((session) => {
        const sessionCookie: SessionCookie = {
          id: recent.id,
          token: session.sessionToken,
          changeDate: session.details.changeDate,
          loginName: recent.loginName,
        };

        return getSession(server, sessionCookie.id, sessionCookie.token).then(
          ({ session }) => {
            const newCookie: SessionCookie = {
              id: sessionCookie.id,
              token: sessionCookie.token,
              changeDate: session.changeDate,
              loginName: session.factors.user.loginName,
            };

            return updateSessionCookie(sessionCookie.id, newCookie)
              .then(() => {
                return NextResponse.json({ factors: session.factors });
              })
              .catch((error) => {
                return NextResponse.json(
                  { details: "could not set cookie" },
                  { status: 500 }
                );
              });
          }
        );
      })
      .catch((error) => {
        console.error("erasd", error);
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.error();
  }
}
