import {
  createPasskeyRegistrationLink,
  getSession,
  registerPasskey,
} from "@/lib/zitadel";
import { getSessionCookieById } from "@/utils/cookies";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { sessionId } = body;

    const sessionCookie = await getSessionCookieById(sessionId);

    const session = await getSession(sessionCookie.id, sessionCookie.token);

    const domain: string = request.nextUrl.hostname;

    const userId = session?.session?.factors?.user?.id;

    if (userId) {
      // TODO: add org context
      return createPasskeyRegistrationLink(userId)
        .then((resp) => {
          const code = resp.code;
          if (!code) {
            throw new Error("Missing code in response");
          }
          return registerPasskey(userId, code, domain).then((resp) => {
            return NextResponse.json(resp);
          });
        })
        .catch((error) => {
          console.error("error on creating passkey registration link");
          return NextResponse.json(error, { status: 500 });
        });
    } else {
      return NextResponse.json(
        { details: "could not get session" },
        { status: 500 },
      );
    }
  } else {
    return NextResponse.json({}, { status: 400 });
  }
}
