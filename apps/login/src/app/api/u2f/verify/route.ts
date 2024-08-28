import { getSession, verifyU2FRegistration } from "@/lib/zitadel";
import { getSessionCookieById } from "@zitadel/next";
import { NextRequest, NextResponse, userAgent } from "next/server";
import {
  VerifyU2FRegistrationRequestSchema,
  VerifyU2FRegistrationResponseSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { createMessage, toJson } from "@zitadel/client";

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
    const sessionCookie = await getSessionCookieById({ sessionId });

    const session = await getSession(sessionCookie.id, sessionCookie.token);

    const userId = session?.session?.factors?.user?.id;

    if (userId) {
      // TODO: this does not make sens to me
      //  We create the object, and later on we assign another object to it.
      // let req: VerifyU2FRegistrationRequest = {
      //   publicKeyCredential,
      //   u2fId,
      //   userId,
      //   tokenName: passkeyName,
      // };

      const req = createMessage(
        VerifyU2FRegistrationRequestSchema,
        // TODO: why did we passed the request instead of body here?
        body,
      );

      return verifyU2FRegistration(req)
        .then((resp) => {
          return NextResponse.json(
            toJson(VerifyU2FRegistrationResponseSchema, resp),
          );
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
