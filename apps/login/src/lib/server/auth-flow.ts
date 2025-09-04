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
export async function completeAuthFlow(params: AuthFlowParams, request?: NextRequest): Promise<string> {
  const { sessionId, requestId, organization } = params;
  
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
      request: request as NextRequest, // Type assertion for now
    });

    // Extract redirect URL from the response
    if (result && result.status === 302) {
      const location = result.headers.get('location');
      if (location) {
        return location;
      }
    }
    
    throw new Error('OIDC flow completion failed');
    
  } else if (requestId.startsWith("saml_")) {
    // Complete SAML flow
    const result = await loginWithSAMLAndSession({
      serviceUrl,
      samlRequest: requestId.replace("saml_", ""),
      sessionId,
      sessions,
      sessionCookies,
      request: request as NextRequest,
    });

    // Extract redirect URL from the response  
    if (result && result.status === 302) {
      const location = result.headers.get('location');
      if (location) {
        return location;
      }
    }
    
    throw new Error('SAML flow completion failed');
  }

  throw new Error('Invalid request ID format');
}

/**
 * Server Action to complete authentication flow
 * This replaces client-side navigation to /login
 */
export async function completeAuthFlowAction(params: AuthFlowParams) {
  const redirectUrl = await completeAuthFlow(params);
  redirect(redirectUrl);
}

/**
 * Validate authentication request parameters
 */
export function validateAuthRequest(searchParams: URLSearchParams): string | null {
  const oidcRequestId = searchParams.get("authRequest");
  const samlRequestId = searchParams.get("samlRequest");
  
  const requestId = searchParams.get("requestId") ??
    (oidcRequestId ? `oidc_${oidcRequestId}` : samlRequestId ? `saml_${samlRequestId}` : undefined);

  return requestId || null;
}

/**
 * Check if request is an RSC request
 */
export function isRSCRequest(searchParams: URLSearchParams): boolean {
  return searchParams.has("_rsc");
}