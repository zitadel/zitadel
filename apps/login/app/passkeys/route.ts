import {
  createPasskeyRegistrationLink,
  getSession,
  server,
} from "#/lib/zitadel";
import { getSessionCookieById } from "#/utils/cookies";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { sessionId } = body;

    const sessionCookie = await getSessionCookieById(sessionId);
    console.log(sessionCookie);

    const session = await getSession(
      server,
      sessionCookie.id,
      sessionCookie.token
    );

    if (session?.session && session.session?.factors?.user?.id) {
      console.log(session.session.factors.user.id, sessionCookie.token);
      return createPasskeyRegistrationLink(
        session.session.factors.user.id,
        sessionCookie.token
      )
        .then((resp) => {
          return NextResponse.json(resp);
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
    return NextResponse.json({}, { status: 500 });
  }
}
