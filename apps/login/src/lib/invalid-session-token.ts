import { Code } from "@zitadel/client";

/**
 * Zitadel rotates the session token on every password check (session.token.set).
 * A second check a moment later invalidates the token still stored in the cookie.
 * CreateCallback then answers COMMAND-sGr42. The account button must send the
 * user back to the password form, which already creates a fresh session when
 * the stored token no longer updates.
 */
export function isStaleSessionToken(code: number, message: string): boolean {
  return code === Code.PermissionDenied && message.includes("COMMAND-sGr42");
}

export function invalidSessionTokenRedirect(input: {
  loginName?: string;
  organization?: string;
  requestId: string;
}): { redirect: string } | { error: string } {
  if (!input.loginName) {
    return { error: "Session not found or invalid" };
  }

  const params = new URLSearchParams();
  params.set("loginName", input.loginName);
  if (input.organization) {
    params.set("organization", input.organization);
  }
  const requestId =
    input.requestId.startsWith("oidc_") || input.requestId.startsWith("saml_") ? input.requestId : `oidc_${input.requestId}`;
  params.set("requestId", requestId);
  return { redirect: `/password?${params.toString()}` };
}
