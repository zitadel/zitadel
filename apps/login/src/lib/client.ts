type FinishFlowCommand =
  | {
      sessionId: string;
      requestId: string;
    }
  | { loginName: string };

function goToSignedInPage(
  props:
    | { sessionId: string; organization?: string; requestId?: string }
    | { organization?: string; loginName: string; requestId?: string },
) {
  const params = new URLSearchParams({});

  if ("loginName" in props && props.loginName) {
    params.append("loginName", props.loginName);
  }

  if ("sessionId" in props && props.sessionId) {
    params.append("sessionId", props.sessionId);
  }

  if (props.organization) {
    params.append("organization", props.organization);
  }

  // required to show conditional UI for device flow
  if (props.requestId) {
    params.append("requestId", props.requestId);
  }

  return `/signedin?` + params;
}

/**
 * for client: redirects user back to an OIDC or SAML application or to a success page when using requestId, check if a default redirect and redirect to it, or just redirect to a success page with the loginName
 * @param command
 * @returns
 */
export async function getNextUrl(
  command: FinishFlowCommand & { organization?: string },
  defaultRedirectUri?: string,
): Promise<string> {
  // finish Device Authorization Flow
  if (
    "requestId" in command &&
    command.requestId.startsWith("device_") &&
    ("loginName" in command || "sessionId" in command)
  ) {
    return goToSignedInPage({
      ...command,
      organization: command.organization,
    });
  }

  // finish SAML or OIDC flow
  if (
    "sessionId" in command &&
    "requestId" in command &&
    (command.requestId.startsWith("saml_") ||
      command.requestId.startsWith("oidc_"))
  ) {
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

  return goToSignedInPage(command);
}
