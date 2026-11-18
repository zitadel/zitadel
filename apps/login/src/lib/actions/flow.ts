"use server";

import { completeAuthFlow } from "@/lib/server/auth-flow";
import { createLogger } from "@/lib/server/logger";
import { FinishFlowCommand, getNextUrl } from "@/lib/server/next-url";

const logger = createLogger("flow");

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
  logger.debug("completeFlowOrGetUrl called", { command, defaultRedirectUri });

  // Complete OIDC/SAML flows directly with server action
  if (
    "sessionId" in command &&
    "requestId" in command &&
    (command.requestId.startsWith("saml_") || command.requestId.startsWith("oidc_"))
  ) {
    logger.debug("completeFlowOrGetUrl: OIDC/SAML flow detected");
    // This completes the flow and returns a redirect URL or error
    const result = await completeAuthFlow({
      sessionId: command.sessionId,
      requestId: command.requestId,
    });
    logger.debug("completeFlowOrGetUrl: got OIDC/SAML flow result");
    return result;
  }

  logger.debug("completeFlowOrGetUrl: Regular flow, getting next URL");
  // For all other cases, return URL for navigation
  const url = await getNextUrl(command, defaultRedirectUri);
  logger.debug("completeFlowOrGetUrl: got Next URL", { url });
  const result = { redirect: url };
  logger.debug("completeFlowOrGetUrl: got final result");
  return result;
}
