import {
  createSession,
  getSession,
  listAuthenticationMethodTypes,
  server,
} from "#/lib/zitadel";
import {
  SessionCookie,
  addSessionToCookie,
  getSessionCookieById,
} from "#/utils/cookies";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const sessionId = searchParams.get("sessionId");
  if (sessionId) {
    const sessionCookie = await getSessionCookieById(sessionId);

    const session = await getSession(
      server,
      sessionCookie.id,
      sessionCookie.token
    );

    const userId = session?.session?.factors?.user?.id;

    if (userId) {
      return listAuthenticationMethodTypes(userId)
        .then((methods) => {
          return NextResponse.json(methods);
        })
        .catch((error) => {
          return NextResponse.json(error, { status: 500 });
        });
    } else {
      return NextResponse.json(
        { details: "could not get session" },
        { status: 500 }
      );
    }
  } else {
    return NextResponse.json({}, { status: 400 });
  }
}

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName } = body;

    const domain: string = request.nextUrl.hostname;

    const createdSession = await createSession(
      server,
      loginName,
      undefined,
      domain
    );

    if (createdSession) {
      return getSession(
        server,
        createdSession.sessionId,
        createdSession.sessionToken
      ).then((response) => {
        if (response?.session && response.session?.factors?.user?.loginName) {
          const userId = response?.session?.factors?.user?.id;

          const sessionCookie: SessionCookie = {
            id: createdSession.sessionId,
            token: createdSession.sessionToken,
            changeDate: response.session.changeDate?.toString() ?? "",
            loginName: response.session?.factors?.user?.loginName ?? "",
          };
          return addSessionToCookie(sessionCookie)
            .then(() => {
              return listAuthenticationMethodTypes(userId)
                .then((methods) => {
                  return NextResponse.json({
                    authMethodTypes: methods.authMethodTypes,
                    sessionId: createdSession.sessionId,
                    factors: response?.session?.factors,
                  });
                })
                .catch((error) => {
                  return NextResponse.json(error, { status: 500 });
                });
            })
            .catch((error) => {
              return NextResponse.json(
                {
                  details: "could not add session to cookie",
                },
                { status: 500 }
              );
            });
        } else {
          return NextResponse.json(
            {
              details:
                "could not get session or session does not have loginName",
            },
            { status: 500 }
          );
        }
      });
    } else {
      return NextResponse.error();
    }
  } else {
    return NextResponse.error();
  }
}
