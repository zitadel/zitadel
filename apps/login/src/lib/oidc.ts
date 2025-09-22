import { Cookie } from "@/lib/cookies";
import { sendLoginname, SendLoginnameCommand } from "@/lib/server/loginname";
import { createCallback, getLoginSettings } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { CreateCallbackRequestSchema, SessionSchema } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { isSessionValid } from "./session";

type LoginWithOIDCAndSession = {
  serviceUrl: string;
  authRequest: string;
  sessionId: string;
  sessions: Session[];
  sessionCookies: Cookie[];
};
export async function loginWithOIDCAndSession({
  serviceUrl,
  authRequest,
  sessionId,
  sessions,
  sessionCookies,
}: LoginWithOIDCAndSession): Promise<{ error: string } | { redirect: string }> {
  console.log(`Login with session: ${sessionId} and authRequest: ${authRequest}`);

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
        console.log("Redirecting to re-authenticate:", res.redirect);
        return { redirect: res.redirect };
      }
    }

    const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession?.id);

    if (cookie && cookie.id && cookie.token) {
      const session = {
        sessionId: cookie?.id,
        sessionToken: cookie?.token,
      };

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
          console.log("Redirecting to callback URL:", callbackUrl);
          return { redirect: callbackUrl };
        } else {
          return { error: "An error occurred!" };
        }
      } catch (error: unknown) {
        // handle already handled gracefully as these could come up if old emails with requestId are used (reset password, register emails etc.)
        console.error(error);
        if (error && typeof error === "object" && "code" in error && error?.code === 9) {
          const loginSettings = await getLoginSettings({
            serviceUrl,
            organization: selectedSession.factors?.user?.organizationId,
          });

          if (loginSettings?.defaultRedirectUri) {
            return { redirect: loginSettings.defaultRedirectUri };
          }

          const signedinUrl = "/signedin";

          const params = new URLSearchParams();
          if (selectedSession.factors?.user?.loginName) {
            params.append("loginName", selectedSession.factors?.user?.loginName);
          }
          if (selectedSession.factors?.user?.organizationId) {
            params.append("organization", selectedSession.factors?.user?.organizationId);
          }
          console.log("Redirecting to signed-in page:", signedinUrl + "?" + params.toString());
          return { redirect: signedinUrl + "?" + params.toString() };
        } else {
          return { error: "Unknown error occurred" };
        }
      }
    }
  }

  // If no session found or no valid cookie, return error
  return { error: "Session not found or invalid" };
}
