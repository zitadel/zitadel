import { Cookie } from "@/lib/cookies";
import { sendLoginname, SendLoginnameCommand } from "@/lib/server/loginname";
import { createCallback, getLoginSettings } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  CreateCallbackRequestSchema,
  SessionSchema,
} from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { NextRequest, NextResponse } from "next/server";
import { constructUrl } from "./service-url";
import { isSessionValid } from "./session";

type LoginWithOIDCAndSession = {
  serviceUrl: string;
  authRequest: string;
  sessionId: string;
  sessions: Session[];
  sessionCookies: Cookie[];
  request: NextRequest;
};
export async function loginWithOIDCAndSession({
  serviceUrl,
  authRequest,
  sessionId,
  sessions,
  sessionCookies,
  request,
}: LoginWithOIDCAndSession) {
  console.log(
    `Login with session: ${sessionId} and authRequest: ${authRequest}`,
  );

  const selectedSession = sessions.find((s) => s.id === sessionId);

  if (selectedSession && selectedSession.id) {
    console.log(`Found session ${selectedSession.id}`);

    const isValid = await isSessionValid({
      serviceUrl,
      session: selectedSession,
    });

    console.log("Session is valid:", isValid);

    if (!isValid && selectedSession.factors?.user) {
      // if the session is not valid anymore, we need to redirect the user to re-authenticate /
      // TODO: handle IDP intent direcly if available
      const command: SendLoginnameCommand = {
        loginName: selectedSession.factors.user?.loginName,
        organization: selectedSession.factors?.user?.organizationId,
        requestId: `oidc_${authRequest}`,
      };

      const res = await sendLoginname(command);

      if (res && "redirect" in res && res?.redirect) {
        const absoluteUrl = constructUrl(request, res.redirect);
        return NextResponse.redirect(absoluteUrl.toString());
      }
    }

    const cookie = sessionCookies.find(
      (cookie) => cookie.id === selectedSession?.id,
    );

    if (cookie && cookie.id && cookie.token) {
      const session = {
        sessionId: cookie?.id,
        sessionToken: cookie?.token,
      };

      // works not with _rsc request
      try {
        const { callbackUrl } = await createCallback({
          serviceUrl,
          req: create(CreateCallbackRequestSchema, {
            authRequestId: authRequest,
            callbackKind: {
              case: "session",
              value: create(SessionSchema, session),
            },
          }),
        });
        if (callbackUrl) {
          return NextResponse.redirect(callbackUrl);
        } else {
          return NextResponse.json(
            { error: "An error occurred!" },
            { status: 500 },
          );
        }
      } catch (error: unknown) {
        // handle already handled gracefully as these could come up if old emails with requestId are used (reset password, register emails etc.)
        console.error(error);
        if (
          error &&
          typeof error === "object" &&
          "code" in error &&
          error?.code === 9
        ) {
          const loginSettings = await getLoginSettings({
            serviceUrl,
            organization: selectedSession.factors?.user?.organizationId,
          });

          if (loginSettings?.defaultRedirectUri) {
            return NextResponse.redirect(loginSettings.defaultRedirectUri);
          }

          const signedinUrl = constructUrl(request, "/signedin");

          if (selectedSession.factors?.user?.loginName) {
            signedinUrl.searchParams.set(
              "loginName",
              selectedSession.factors?.user?.loginName,
            );
          }
          if (selectedSession.factors?.user?.organizationId) {
            signedinUrl.searchParams.set(
              "organization",
              selectedSession.factors?.user?.organizationId,
            );
          }
          return NextResponse.redirect(signedinUrl);
        } else {
          return NextResponse.json({ error }, { status: 500 });
        }
      }
    }
  }
}
