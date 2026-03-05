"use server";

import { headers } from "next/headers";
import { completeAuthFlow } from "./server/auth-flow";
import { getPublicHostWithProtocol } from "./server/host";
import { sanitizeRedirectUri } from "./client-utils";

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
  console.log("getNextUrl called with:", command, "defaultRedirectUri:", defaultRedirectUri);

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
    console.log("getNextUrl: Got Device flow result");
    return result;
  }

  // OIDC/SAML flows are now handled by completeAuthFlowAction() server action
  // This function only handles device flows and fallback navigation

  const result = await resolveRedirectUri(command, defaultRedirectUri);
  console.log("getNextUrl: Resolved redirect URI:", result);
  return result;
}

/**
 * Resolves the redirect URI based on the following priority:
 * 1. DEFAULT_REDIRECT_URI environment variable
 * 2. defaultRedirectUri from organization settings
 * 3. Relative signed-in page fallback
 * 4. Reserved for future extensions
 */
export async function resolveRedirectUri(command: FinishFlowCommand, defaultRedirectUri?: string): Promise<string> {
  // Determine current origin for same-origin validation of redirect URIs
  let currentOrigin: string | undefined;
  try {
    const _headers = await headers();
    currentOrigin = getPublicHostWithProtocol(_headers);
  } catch {
    // Origin unavailable; only relative paths will be accepted
  }

  // 1. Environment variable override
  const envOverride = process.env.DEFAULT_REDIRECT_URI;
  if (envOverride) {
    if (envOverride.startsWith("/") && !envOverride.startsWith("//")) {
      // Relative path: construct absolute URL for proper host-based redirect (bypasses Next.js base path)
      if (currentOrigin) {
        const sanitized = sanitizeRedirectUri(`${currentOrigin}${envOverride}`, currentOrigin);
        if (sanitized) {
          return sanitized;
        }
      }
      // Fall back to relative path if origin unavailable
      const sanitized = sanitizeRedirectUri(envOverride);
      if (sanitized) {
        return sanitized;
      }
    } else {
      // Absolute env var: admin-controlled, trust the origin
      const sanitized = sanitizeRedirectUri(envOverride, undefined, true);
      if (sanitized) {
        return sanitized;
      }
    }
    console.warn("resolveRedirectUri: Invalid DEFAULT_REDIRECT_URI:", envOverride);
  }

  // 2. Default redirect URI from settings
  if (defaultRedirectUri) {
    const sanitized = sanitizeRedirectUri(defaultRedirectUri, currentOrigin);
    if (sanitized) {
      return sanitized;
    }
    console.warn("resolveRedirectUri: Invalid defaultRedirectUri:", defaultRedirectUri);
  }

  // 3. Default signed-in page (relative)
  return goToSignedInPage(command);
}

