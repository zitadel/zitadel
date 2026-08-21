import { isSafeRedirectUri } from "@/lib/client-utils";
import { Cookie } from "@/lib/cookies";
import { isClassifiedError } from "@/lib/grpc/interceptors/error-classification";
import { sendLoginname, SendLoginnameCommand } from "@/lib/server/loginname";
import { createResponse, getLoginSettings, ServiceConfig } from "@/lib/zitadel";
import { Code, create } from "@zitadel/client";
import { CreateResponseRequestSchema } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { isSessionValid } from "./session";

type LoginWithSAMLAndSession = {
  serviceConfig: ServiceConfig;
  samlRequest: string;
  sessionId: string;
  sessions: Session[];
  sessionCookies: Cookie[];
};

export async function loginWithSAMLAndSession({
  serviceConfig,
  samlRequest,
  sessionId,
  sessions,
  sessionCookies,
}: LoginWithSAMLAndSession): Promise<
  { error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }
> {
  console.log(`Login with session: ${sessionId} and samlRequest: ${samlRequest}`);

  const selectedSession = sessions.find((s) => s.id === sessionId);

  if (selectedSession && selectedSession.id) {
    const isValid = await isSessionValid({ serviceConfig, session: selectedSession });

    if (!isValid && selectedSession.factors?.user) {
      // if the session is not valid anymore, we need to redirect the user to re-authenticate /
      // TODO: handle IDP intent direcly if available
      const command: SendLoginnameCommand = {
        loginName: selectedSession.factors.user?.loginName,
        organization: selectedSession.factors?.user?.organizationId,
        requestId: `saml_${samlRequest}`,
      };

      const res = await sendLoginname(command);

      if (res && "redirect" in res && res?.redirect) {
        return { redirect: res.redirect };
      }

      if (res && "samlData" in res && res?.samlData) {
        return { samlData: res.samlData };
      }
    }

    const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession?.id);

    if (cookie && cookie.id && cookie.token) {
      const session = {
        sessionId: cookie?.id,
        sessionToken: cookie?.token,
      };

      // works not with _rsc request
      try {
        const { url, binding } = await createResponse({
          serviceConfig,
          req: create(CreateResponseRequestSchema, {
            samlRequestId: samlRequest,
            responseKind: {
              case: "session",
              value: session,
            },
          }),
        });
        if (url && binding.case === "redirect") {
          if (!isSafeRedirectUri(url)) {
            console.warn("loginWithSAMLAndSession: Blocked unsafe SAML redirect URL:", url);
            return { error: "Unsafe redirect URI was blocked" };
          }
          return { redirect: url };
        } else if (url && binding.case === "post") {
          if (!isSafeRedirectUri(url)) {
            console.warn("loginWithSAMLAndSession: Blocked unsafe SAML post URL:", url);
            return { error: "Unsafe redirect URI was blocked" };
          }
          return {
            samlData: {
              url,
              fields: {
                RelayState: binding.value.relayState,
                SAMLResponse: binding.value.samlResponse,
              },
            },
          };
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
            console.warn("loginWithSAMLAndSession: Unsafe defaultRedirectUri prevented:", loginSettings.defaultRedirectUri);
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
        } else if (selectedSession.factors?.user) {
          // The session could not be used to complete this authentication request
          // (e.g. the target organization enforces a different authentication
          // policy than the one the session was created with). Instead of
          // surfacing a generic error and leaving the user stuck, guide them into
          // re-authentication for the current request - the same outcome as
          // selecting "Add another account".
          const command: SendLoginnameCommand = {
            loginName: selectedSession.factors.user.loginName,
            organization: selectedSession.factors.user.organizationId,
            requestId: `saml_${samlRequest}`,
          };

          const res = await sendLoginname(command).catch(() => undefined);

          if (res && "redirect" in res && res?.redirect) {
            return { redirect: res.redirect };
          }

          return { error: "This session can't be reused for this request. Please authenticate again." };
        } else {
          return { error: "Unknown error occurred" };
        }
      }
    }
  }

  // If no session found or no valid cookie, return error
  return { error: "Session not found or invalid" };
}
