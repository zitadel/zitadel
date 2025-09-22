"use server";

import { getAllSessions } from "@/lib/cookies";
import { loginWithOIDCAndSession } from "@/lib/oidc";
import { loginWithSAMLAndSession } from "@/lib/saml";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { listSessions } from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { headers } from "next/headers";
import { redirect } from "next/navigation";

export interface AuthFlowParams {
  sessionId: string;
  requestId: string;
  organization?: string;
}

async function loadSessions({ serviceUrl, ids }: { serviceUrl: string; ids: string[] }): Promise<Session[]> {
  const response = await listSessions({
    serviceUrl,
    ids: ids.filter((id: string | undefined) => !!id),
  });

  return response?.sessions ?? [];
}

/**
 * Server Action to complete authentication flow
 * Complete OIDC/SAML authentication flow with session
 * This is the shared logic for flow completion
 */
export async function completeAuthFlow(command: AuthFlowParams): Promise<void | { error: string }> {
  const { sessionId, requestId } = command;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const sessionCookies = await getAllSessions();
  const ids = sessionCookies.map((s) => s.id);
  let sessions: Session[] = [];

  if (ids && ids.length) {
    sessions = await loadSessions({ serviceUrl, ids });
  }

  if (requestId.startsWith("oidc_")) {
    // Complete OIDC flow
    const result = await loginWithOIDCAndSession({
      serviceUrl,
      authRequest: requestId.replace("oidc_", ""),
      sessionId,
      sessions,
      sessionCookies,
    });

    // Handle redirect response from loginWithOIDCAndSession
    if (result && "redirect" in result && result.redirect) {
      redirect(result.redirect);
    }

    // Only return error, not redirect (since redirect throws)
    if (result && "error" in result && result.error) {
      return { error: result.error };
    }

    return; // No error, redirect was called
  } else if (requestId.startsWith("saml_")) {
    // Complete SAML flow
    const result = await loginWithSAMLAndSession({
      serviceUrl,
      samlRequest: requestId.replace("saml_", ""),
      sessionId,
      sessions,
      sessionCookies,
    });

    // Handle redirect response from loginWithSAMLAndSession
    if (result && "redirect" in result && result.redirect) {
      redirect(result.redirect);
    }

    // Only return error, not redirect (since redirect throws)
    if (result && "error" in result && result.error) {
      return { error: result.error };
    }

    return; // No error, redirect was called
  }

  return { error: "Invalid request ID format" };
}
