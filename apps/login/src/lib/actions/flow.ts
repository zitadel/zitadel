"use server";

import { completeAuthFlow } from "@/lib/server/auth-flow";
import { FinishFlowCommand, getNextUrl } from "@/lib/server/next-url";

/**
 * Complete authentication flow or get next URL for navigation
 * - For OIDC/SAML flows with sessionId+requestId: completes flow directly via server action
 * - For device flows: returns URL to signed-in page
 * - For other cases: returns default redirect or fallback URL
 */
export async function completeFlowOrGetUrl(
  command: FinishFlowCommand & { organization?: string },
  defaultRedirectUri?: string,
): Promise<{ redirect: string } | { error: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  console.log("completeFlowOrGetUrl called with:", command, "defaultRedirectUri:", defaultRedirectUri);

  // Complete OIDC/SAML flows directly with server action
  if (
    "sessionId" in command &&
    "requestId" in command &&
    (command.requestId.startsWith("saml_") || command.requestId.startsWith("oidc_"))
  ) {
    console.log("completeFlowOrGetUrl: OIDC/SAML flow detected");
    // This completes the flow and returns a redirect URL or error
    const result = await completeAuthFlow({
      sessionId: command.sessionId,
      requestId: command.requestId,
    });
    console.log("completeFlowOrGetUrl: got OIDC/SAML flow result");
    return result;
  }

  console.log("completeFlowOrGetUrl: Regular flow, getting next URL");
  // For all other cases, return URL for navigation
  const url = await getNextUrl(command, defaultRedirectUri);
  console.log("completeFlowOrGetUrl: got Next URL:", url);
  const result = { redirect: url };
  console.log("completeFlowOrGetUrl: got final result");
  return result;
}
