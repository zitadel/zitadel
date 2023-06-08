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

    const session = await getSessionCookieById(sessionId);

    return createPasskeyRegistrationLink(server, session.id, session.token)
      .then((resp) => {
        return NextResponse.json(resp);
      })
      .catch((error) => {
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.json({}, { status: 500 });
  }
}
