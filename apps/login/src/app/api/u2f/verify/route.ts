import { getSession, verifyU2FRegistration } from "@/lib/zitadel";
import { getSessionCookieById } from "@/utils/cookies";
import { NextRequest, NextResponse, userAgent } from "next/server";
import { VerifyU2FRegistrationRequest } from "@zitadel/proto/zitadel/user/v2beta/user_service_pb";
import { PlainMessage } from "@zitadel/client2";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    let { u2fId, passkeyName, publicKeyCredential, sessionId } = body;

    if (!!!passkeyName) {
      const { browser, device, os } = userAgent(request);
      passkeyName = `${device.vendor ?? ""} ${device.model ?? ""}${
        device.vendor || device.model ? ", " : ""
      }${os.name}${os.name ? ", " : ""}${browser.name}`;
    }
    const sessionCookie = await getSessionCookieById(sessionId);

    const session = await getSession(sessionCookie.id, sessionCookie.token);

    const userId = session?.session?.factors?.user?.id;

    if (userId) {
      const req: PlainMessage<VerifyU2FRegistrationRequest> = {
        publicKeyCredential,
        u2fId,
        userId,
        tokenName: passkeyName,
      };
      return verifyU2FRegistration(req)
        .then((resp) => {
          return NextResponse.json(resp);
        })
        .catch((error) => {
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
