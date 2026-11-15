import "server-only";

import { isSafeRedirectUri } from "@/lib/client-utils";
import { getPublicHostWithProtocol } from "@/lib/server/host";
import { headers } from "next/headers";

export type FinishFlowCommand =
  | {
      sessionId: string;
      requestId: string;
    }
  | { loginName: string };

function goToPage(
  path: string,
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

  return `${path}?` + params;
}

function goToSignedInPage(props: Parameters<typeof goToPage>[1]) {
  return goToPage("/signedin", props);
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

  // ask for consent before the Device Authorization Flow is completed
  if (
    "requestId" in command &&
    command.requestId.startsWith("device_") &&
    ("loginName" in command || "sessionId" in command)
  ) {
    const result = goToPage("/device/consent", {
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
  // 1. Environment variable override
  const envOverride = process.env.DEFAULT_REDIRECT_URI;
  if (envOverride) {
    if (envOverride.startsWith("/")) {
      // Special state: trigger absolute host-based redirect with provided path
      try {
        const _headers = await headers();
        const host = getPublicHostWithProtocol(_headers);
        const result = `${host}${envOverride}`;
        console.log("resolveRedirectUri: Using host-based redirect from override:", result);
        return result;
      } catch (error) {
        console.warn("resolveRedirectUri: Could not determine host for override, falling back", error);
      }
    } else {
      console.log("resolveRedirectUri: Using DEFAULT_REDIRECT_URI override:", envOverride);
      return envOverride;
    }
  }

  // 2. Default redirect URI from settings
  if (defaultRedirectUri) {
    if (isSafeRedirectUri(defaultRedirectUri)) {
      console.log("resolveRedirectUri: Using defaultRedirectUri from settings:", defaultRedirectUri);
      return defaultRedirectUri;
    } else {
      console.warn("resolveRedirectUri: Unsafe defaultRedirectUri prevented:", defaultRedirectUri);
    }
  }

  // 3. Default signed-in page (relative)
  const result = goToSignedInPage(command);
  console.log("resolveRedirectUri: Using relative goToSignedInPage result:", result);
  return result;
}
