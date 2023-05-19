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
      console.log(session);
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
    console.log("found recent cookie: ", recent);
    const session = await setSession(server, recent.id, recent.token, password);
    console.log("updatedsession", session);

    const sessionCookie: SessionCookie = {
      id: recent.id,
      token: session.sessionToken,
      changeDate: session.details.changeDate,
      loginName: recent.loginName,
    };

    return getSession(server, sessionCookie.id, sessionCookie.token).then(
      ({ session }) => {
        console.log(session);
        const newCookie: SessionCookie = {
          id: sessionCookie.id,
          token: sessionCookie.token,
          changeDate: session.changeDate,
          loginName: session.factors.user.loginName,
        };
        // return addSessionToCookie(sessionCookie).then(() => {
        //   return NextResponse.json({ factors: session.factors });
        // });
        return updateSessionCookie(sessionCookie.id, sessionCookie).then(() => {
          console.log("updatedRecent:", sessionCookie);
          return NextResponse.json({ factors: session.factors });
        });
      }
    );
  } else {
    return NextResponse.error();
  }
}

// /**
//  *
//  * @param request loginName of a session
//  * @returns the session
//  */
// export async function GET(request: NextRequest) {
//   console.log(request);
//   if (request) {
//     const { loginName } = request.params;

//     const recent = await getMostRecentCookieWithLoginname(loginName);
//     console.log("found recent cookie: ", recent);

//     return getSession(server, recent.id, recent.token).then(({ session }) => {
//       console.log(session);

//       return NextResponse.json({ factors: session.factors });
//     });
//   } else {
//     return NextResponse.error();
//   }
// }
