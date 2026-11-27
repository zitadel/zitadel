import "server-only";

import { isSafeRedirectUri } from "@/lib/client-utils";
import { Cookie } from "@/lib/server/cookies";
import { isClassifiedError } from "@/lib/server/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/server/logger";
import { reauthenticate } from "@/lib/server/reauthenticate";
import { isSessionValid } from "@/lib/server/session";
import { createCallback, getLoginSettings, ServiceConfig } from "@/lib/server/zitadel";
import { Code, create } from "@zitadel/client";
import { CreateCallbackRequestSchema, SessionSchema } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

const logger = createLogger("oidc");

type LoginWithOIDCAndSession = {
  serviceConfig: ServiceConfig;
  authRequest: string;
  sessionId: string;
  sessions: Session[];
  sessionCookies: Cookie[];
};
export async function loginWithOIDCAndSession({
  serviceConfig,
  authRequest,
  sessionId,
  sessions,
  sessionCookies,
}: LoginWithOIDCAndSession): Promise<
  { error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }
> {
  const selectedSession = sessions.find((s) => s.id === sessionId);

  if (selectedSession && selectedSession.id) {
    const isValid = await isSessionValid({ serviceConfig, session: selectedSession });

    logger.debug("Session is valid", { isValid });

    if (!isValid) {
      // An invalid session (e.g. no verified factor, MFA missing or expired) must never be
      // linked to the auth request. Send the user to re-authenticate, or fail closed.
      return (await reauthenticate(selectedSession, `oidc_${authRequest}`)) ?? { error: "Session not found or invalid" };
    }

    const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession?.id);

    if (cookie && cookie.id && cookie.token) {
      const session = {
        sessionId: cookie?.id,
        sessionToken: cookie?.token,
      };

      try {
        const { callbackUrl } = await createCallback({
          serviceConfig,
          req: create(CreateCallbackRequestSchema, {
            authRequestId: authRequest,
            callbackKind: {
              case: "session",
              value: create(SessionSchema, session),
            },
          }),
        });
        if (callbackUrl) {
          if (!isSafeRedirectUri(callbackUrl)) {
            console.warn("loginWithOIDCAndSession: Blocked unsafe OIDC callback URL:", callbackUrl);
            return { error: "Unsafe redirect URI was blocked" };
          }
          return { redirect: callbackUrl };
        } else {
          return { error: "An error occurred!" };
        }
      } catch (error: unknown) {
        // handle already handled gracefully as these could come up if old emails with requestId are used (reset password, register emails etc.)
        console.error(error);
        if (isClassifiedError(error) && error.code === Code.FailedPrecondition) {
          const loginSettings = await getLoginSettings({
            serviceConfig,
            organization: selectedSession.factors?.user?.organizationId,
          });

          if (loginSettings?.defaultRedirectUri && isSafeRedirectUri(loginSettings.defaultRedirectUri)) {
            return { redirect: loginSettings.defaultRedirectUri };
          } else if (loginSettings?.defaultRedirectUri) {
            console.warn("loginWithOIDCAndSession: Unsafe defaultRedirectUri prevented:", loginSettings.defaultRedirectUri);
          }

          const signedinUrl = "/signedin";

          const params = new URLSearchParams();
          if (selectedSession.factors?.user?.loginName) {
            params.append("loginName", selectedSession.factors?.user?.loginName);
          }
          if (selectedSession.factors?.user?.organizationId) {
            params.append("organization", selectedSession.factors?.user?.organizationId);
          }
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
