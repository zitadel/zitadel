import { completeAuthFlow } from "./server/auth-flow";

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
 * Complete authentication flow or get next URL for navigation
 * - For OIDC/SAML flows with sessionId+requestId: completes flow directly via server action
 * - For device flows: returns URL to signed-in page
 * - For other cases: returns default redirect or fallback URL
 */
export async function completeFlowOrGetUrl(
  command: FinishFlowCommand & { organization?: string },
  defaultRedirectUri?: string,
): Promise<{ redirect: string } | { error: string }> {
  console.debug("[client] completeFlowOrGetUrl called", { command, defaultRedirectUri });

  // Complete OIDC/SAML flows directly with server action
  if (
    "sessionId" in command &&
    "requestId" in command &&
    (command.requestId.startsWith("saml_") || command.requestId.startsWith("oidc_"))
  ) {
    console.debug("[client] OIDC/SAML flow detected");
    // This completes the flow and returns a redirect URL or error
    const result = await completeAuthFlow({
      sessionId: command.sessionId,
      requestId: command.requestId,
    });
    console.debug("[client] OIDC/SAML flow result", { result });
    return result;
  }

  console.debug("[client] Regular flow, getting next URL");
  // For all other cases, return URL for navigation
  const url = await getNextUrl(command, defaultRedirectUri);
  console.debug("[client] Next URL", { url });
  const result = { redirect: url };
  console.debug("[client] Final result", { result });
  return result;
}

/**
 * for client: redirects user back to device flow completion, default redirect, or success page
 * Note: OIDC/SAML flows now use completeAuthFlowAction() instead of URL navigation
 * @param command
 * @returns
 */
export async function getNextUrl(
  command: FinishFlowCommand & { organization?: string },
  defaultRedirectUri?: string,
): Promise<string> {
  console.debug("[client] getNextUrl called", { command, defaultRedirectUri });

  // finish Device Authorization Flow
  if (
    "requestId" in command &&
    command.requestId.startsWith("device_") &&
    ("loginName" in command || "sessionId" in command)
  ) {
    const result = goToSignedInPage({
      ...command,
      organization: command.organization,
    });
    console.debug("[client] Device flow result", { result });
    return result;
  }

  // OIDC/SAML flows are now handled by completeAuthFlowAction() server action
  // This function only handles device flows and fallback navigation

  if (defaultRedirectUri) {
    console.debug("[client] Using defaultRedirectUri", { defaultRedirectUri });
    return defaultRedirectUri;
  }

  const result = goToSignedInPage(command);
  console.debug("[client] Using goToSignedInPage result", { result });
  return result;
}
