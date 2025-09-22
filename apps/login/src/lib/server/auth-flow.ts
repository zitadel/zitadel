"use server";

import { getAllSessions } from "@/lib/cookies";
import { loginWithOIDCAndSession } from "@/lib/oidc";
import { loginWithSAMLAndSession } from "@/lib/saml";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { listSessions } from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { headers } from "next/headers";

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
 * Returns either an error or a redirect URL for client-side navigation
 */
export async function completeAuthFlow(command: AuthFlowParams): Promise<{ error: string } | { redirect: string }> {
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
    return (
 loginWithOIDCAndSession({
        serviceUrl,
        authRequest: requestId.replace("oidc_", ""),
        sessionId,
        sessions,
        sessionCookies,
      })) ?? { error: "Unknown error occurred in OIDC flow" }
    );
  } else if (requestId.startsWith("saml_")) {
    // Complete SAML flow
    return (
      (await loginWithSAMLAndSession({
        serviceUrl,
        samlRequest: requestId.replace("saml_", ""),
        sessionId,
        sessions,
        sessionCookies,
      })) ?? { error: "Unknown error occurred in SAML flow" }
    );
  }

  return { error: "Invalid request ID format" };
}
