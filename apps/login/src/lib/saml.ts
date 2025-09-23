import { Cookie } from "@/lib/cookies";
import { sendLoginname, SendLoginnameCommand } from "@/lib/server/loginname";
import { createResponse, getLoginSettings } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { CreateResponseRequestSchema } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { cookies } from "next/headers";
import { v4 as uuidv4 } from "uuid";
import { isSessionValid } from "./session";

type LoginWithSAMLAndSession = {
  serviceUrl: string;
  samlRequest: string;
  sessionId: string;
  sessions: Session[];
  sessionCookies: Cookie[];
};

export async function getSAMLFormUID() {
  return uuidv4();
}

export async function setSAMLFormCookie(value: string): Promise<string> {
  const cookiesList = await cookies();
  const uid = await getSAMLFormUID();

  try {
    // Check cookie size limits (typical limit is 4KB)
    if (value.length > 4000) {
      console.warn(`SAML form cookie value is large (${value.length} characters), may exceed browser limits`);
    }

    // Log the attempt
    console.log(`Setting SAML form cookie with uid: ${uid}, value length: ${value.length}`);

    await cookiesList.set({
      name: uid,
      value: value,
      httpOnly: true,
      secure: process.env.NODE_ENV === "production", // Required for HTTPS in production
      sameSite: "lax", // Allows cookies with top-level navigation (needed for SAML redirects)
      path: "/",
      maxAge: 5 * 60, // 5 minutes
    });

    // Note: We can't reliably verify immediately due to Next.js cookies API behavior
    // Instead, we'll rely on the getSAMLFormCookie function to detect failures
    console.log(`Successfully set SAML form cookie with uid: ${uid}`);

    return uid;
  } catch (error) {
    console.error(`Failed to set SAML form cookie with uid: ${uid}`, {
      error,
      valueLength: value.length,
      uid,
    });
    throw new Error(`Failed to set SAML form cookie: ${error instanceof Error ? error.message : String(error)}`);
  }
}

export async function getSAMLFormCookie(uid: string): Promise<string | null> {
  const cookiesList = await cookies();

  try {
    const cookie = cookiesList.get(uid);

    if (!cookie) {
      console.warn(`SAML form cookie not found for uid: ${uid}`);
      return null;
    }

    if (!cookie.value) {
      console.warn(`SAML form cookie found but empty value for uid: ${uid}`);
      return null;
    }

    console.log(`Successfully retrieved SAML form cookie for uid: ${uid}, value length: ${cookie.value.length}`);
    return cookie.value;
  } catch (error) {
    console.error(`Error retrieving SAML form cookie for uid: ${uid}`, error);
    return null;
  }
}

export async function loginWithSAMLAndSession({
  serviceUrl,
  samlRequest,
  sessionId,
  sessions,
  sessionCookies,
}: LoginWithSAMLAndSession): Promise<{ error: string } | { redirect: string }> {
  console.log(`Login with session: ${sessionId} and samlRequest: ${samlRequest}`);

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
        requestId: `saml_${samlRequest}`,
      };

      const res = await sendLoginname(command);

      if (res && "redirect" in res && res?.redirect) {
        return { redirect: res.redirect };
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
        const { url } = await createResponse({
          serviceUrl,
          req: create(CreateResponseRequestSchema, {
            samlRequestId: samlRequest,
            responseKind: {
              case: "session",
              value: session,
            },
          }),
        });
        if (url) {
          return { redirect: url };
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
