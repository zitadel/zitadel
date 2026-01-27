"use server";

import { getAllSessions } from "@/lib/cookies";
import { loginWithOIDCAndSession } from "@/lib/oidc";
import { loginWithSAMLAndSession } from "@/lib/saml";
import { getServiceConfig } from "@/lib/service-url";
import { listSessions, ServiceConfig } from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { headers } from "next/headers";

export interface AuthFlowParams {
  sessionId: string;
  requestId: string;
  organization?: string;
}

async function loadSessions({ serviceConfig, ids }: { serviceConfig: ServiceConfig; ids: string[] }): Promise<Session[]> {
  const response = await listSessions({ serviceConfig, ids: ids.filter((id: string | undefined) => !!id),
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
  const { serviceConfig } = getServiceConfig(_headers);

  const sessionCookies = await getAllSessions();
  const ids = sessionCookies.map((s) => s.id);
  let sessions: Session[] = [];

  if (ids && ids.length) {
    sessions = await loadSessions({ serviceConfig, ids });
  }

  if (requestId.startsWith("oidc_")) {
    // Complete OIDC flow
    const result = await loginWithOIDCAndSession({ serviceConfig, authRequest: requestId.replace("oidc_", ""),
      sessionId,
      sessions,
      sessionCookies,
    });

    // Safety net - ensure we always return a valid object
    if (!result || typeof result !== "object" || (!("redirect" in result) && !("error" in result))) {
      console.error("Auth flow: Invalid result from loginWithOIDCAndSession:", result);
      return { error: "Authentication completed but navigation failed" };
    }

    return result;
  } else if (requestId.startsWith("saml_")) {
    // Complete SAML flow
    const result = await loginWithSAMLAndSession({ serviceConfig, samlRequest: requestId.replace("saml_", ""),
      sessionId,
      sessions,
      sessionCookies,
    });

    // Safety net - ensure we always return a valid object
    if (!result || typeof result !== "object" || (!("redirect" in result) && !("error" in result))) {
      console.error("Auth flow: Invalid result from loginWithSAMLAndSession:", result);
      return { error: "Authentication completed but navigation failed" };
    }

    return result;
  }

  return { error: "Invalid request ID format" };
}
