import { ReadonlyHeaders } from "next/dist/server/web/spec-extension/adapters/headers";
import { NextRequest } from "next/server";
import { ServiceConfig } from "./zitadel";
import { getInstanceHost, getPublicHost } from "./server/host";

/**
 * Extracts the service URL based on deployment mode and configuration.
 *
 * Priority:
 * 1. ZITADEL_API_URL (required) - Used by both self-hosted and multi-tenant
 * 2. x-zitadel-forward-host (multi-tenant only) - Set by Zitadel proxy
 * 3. host header (multi-tenant fallback) - For dynamic host resolution
 *
 * @param headers - Request headers
 * @returns Object containing the service Configuration
 * @throws Error if the service Configuration could not be determined
 */

function stripProtocol(url: string): string {
  return url.replace(/^https?:\/\//, "");
}

export function getServiceConfig(headers: ReadonlyHeaders): { serviceConfig: ServiceConfig } {
  if (!process.env.ZITADEL_API_URL) {
    throw new Error("ZITADEL_API_URL is not set");
  }

  // use forwarded host from proxy - headers are forwarded to the APIs.
  const instanceHost = getInstanceHost(headers);
  const publicHost = getPublicHost(headers);

  return {
    serviceConfig: {
      baseUrl: process.env.ZITADEL_API_URL,
      ...(instanceHost && { instanceHost: stripProtocol(instanceHost) }),
      ...(publicHost && { publicHost: stripProtocol(publicHost) }),
    },
  };
}

export function constructUrl(request: NextRequest, path: string) {
  const protocol = request.nextUrl.protocol;

  const forwardedHost = getPublicHost(request.headers);
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH || "";
  return new URL(`${basePath}${path}`, `${protocol}//${forwardedHost}`);
}
