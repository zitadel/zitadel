type FinishFlowCommand =
  | {
      sessionId: string;
      requestId: string;
    }
  | { loginName: string };

/**
 * for client: redirects user back to an OIDC or SAML application or to a success page when using requestId, check if a default redirect and redirect to it, or just redirect to a success page with the loginName
 * @param command
 * @returns
 */
export async function getNextUrl(
  command: FinishFlowCommand & { organization?: string },
  defaultRedirectUri?: string,
): Promise<string> {
  if ("sessionId" in command && "requestId" in command) {
    const params = new URLSearchParams({
      sessionId: command.sessionId,
      requestId: command.requestId,
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
