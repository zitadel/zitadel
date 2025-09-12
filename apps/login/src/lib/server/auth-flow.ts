'use server';

import { getAllSessions } from "@/lib/cookies";
import { loginWithOIDCAndSession } from "@/lib/oidc";
import { loginWithSAMLAndSession } from "@/lib/saml";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { listSessions } from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { headers } from "next/headers";
import { NextRequest } from "next/server";
import { redirect } from 'next/navigation';
import { validateAuthRequest, isRSCRequest } from '../auth-utils';

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
 * Complete OIDC/SAML authentication flow with session
 * This is the shared logic for flow completion
 */
export async function completeAuthFlow(params: AuthFlowParams, request?: NextRequest) {
  const { sessionId, requestId } = params;
  
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
    return await loginWithOIDCAndSession({
      serviceUrl,
      authRequest: requestId.replace("oidc_", ""),
      sessionId,
      sessions,
      sessionCookies,
      request: request as NextRequest,
    });
    
  } else if (requestId.startsWith("saml_")) {
    // Complete SAML flow
    return await loginWithSAMLAndSession({
      serviceUrl,
      samlRequest: requestId.replace("saml_", ""),
      sessionId,
      sessions,
      sessionCookies,
      request: request as NextRequest,
    });
  }

  throw new Error('Invalid request ID format');
}

/**
 * Server Action to complete authentication flow
 * This replaces client-side navigation to /login
 */
export async function completeAuthFlowAction(params: AuthFlowParams) {
  // For server actions, we need to extract the redirect URL and call redirect()
  const result = await completeAuthFlow(params);
  
  // Extract redirect URL from the response
  if (result && result.status === 302) {
    const location = result.headers.get('location');
    if (location) {
      redirect(location);
      return;
    }
  }
  
  throw new Error('Authentication flow completion failed');
}