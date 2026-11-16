import "server-only";

import { sendLoginname } from "@/lib/actions/loginname";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

/**
 * Resolves the next authentication step for the user of an invalid session.
 * Returns undefined if no step could be resolved: the session must then never be used
 * to complete an OIDC, SAML or device authorization request.
 */
export async function reauthenticate(
  session: Session,
  requestId?: string,
): Promise<{ redirect: string } | { samlData: { url: string; fields: Record<string, string> } } | undefined> {
  if (!session.factors?.user) {
    return undefined;
  }

  // TODO: if the session was authenticated with an IDP intent, send the user directly back to that IDP.
  // Not possible yet, as the session's intent factor does not expose which IDP was used.
  const res = await sendLoginname({
    loginName: session.factors.user.loginName,
    organization: session.factors.user.organizationId,
    requestId,
  });

  if (res && "redirect" in res && res.redirect) {
    return { redirect: res.redirect };
  }

  if (res && "samlData" in res && res.samlData) {
    return { samlData: res.samlData };
  }

  return undefined;
}
