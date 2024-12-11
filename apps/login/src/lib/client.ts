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
    const params = new URLSearchParams({
      sessionId: command.sessionId,
      authRequest: command.authRequestId,
    });

    if (command.organization) {
      params.append("organization", command.organization);
    }

    return `/login?` + params;
  }

  if (defaultRedirectUri) {
    return defaultRedirectUri;
  }

  const params = new URLSearchParams({
    loginName: command.loginName,
  });

  if (command.organization) {
    params.append("organization", command.organization);
  }

  return `/signedin?` + params;
}
