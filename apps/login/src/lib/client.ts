type FinishFlowCommand =
  | {
      sessionId: string;
      authRequestId: string;
    }
  | { loginName: string };

/**
 * for client: redirects user back to OIDC application or to a success page when using authRequestId, check if a default redirect and redirect to it, or just redirect to a success page with the loginName
 * @param command
 * @returns
 */
export async function getNextUrl(
  command: FinishFlowCommand & { organization?: string },
  defaultRedirectUri?: string,
): Promise<string> {
  if ("sessionId" in command && "authRequestId" in command) {
    const url =
      `/login?` +
      new URLSearchParams({
        sessionId: command.sessionId,
        authRequest: command.authRequestId,
      });
    return url;
  }

  if (defaultRedirectUri) {
    return defaultRedirectUri;
  }

  const signedInUrl =
    `/signedin?` +
    new URLSearchParams({
      loginName: command.loginName,
    });

  return signedInUrl;
}
