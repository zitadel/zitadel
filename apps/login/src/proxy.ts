import { SecuritySettings } from "@zitadel/proto/zitadel/settings/v2/security_settings_pb";

import { headers } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { buildCSP } from "./lib/csp";
import { createLogger } from "./lib/logger";
import { getServiceConfig } from "./lib/service-url";

const logger = createLogger("middleware");
export const config = {
  matcher: ["/.well-known/:path*", "/oauth/:path*", "/oidc/:path*", "/idps/callback/:path*", "/saml/:path*", "/:path*"],
};

async function loadSecuritySettings(request: NextRequest): Promise<SecuritySettings | null> {
  const securityResponse = await fetch(`${request.nextUrl.origin}/security`);

  if (!securityResponse.ok) {
    logger.error("Failed to fetch security settings:", { status: securityResponse.statusText });
    return null;
  }

  const response = await securityResponse.json();

  if (!response || !response.settings) {
    logger.error("No security settings found in the response.");
    return null;
  }

  return response.settings;
}

export async function proxy(request: NextRequest) {
  // Add the original URL as a header to all requests
  const requestHeaders = new Headers(request.headers);

  // Extract "organization" search param from the URL and set it as a header if available
  const organization = request.nextUrl.searchParams.get("organization");
  if (organization) {
    requestHeaders.set("x-zitadel-i18n-organization", organization);
  }

  // Internal infrastructure routes — skip middleware entirely.
  // /security is the internal API this middleware fetches (loop prevention).
  // /healthy and /ready are Kubernetes/Docker health probes that must respond
  // without depending on a ZITADEL backend.
  const skipPaths = ["/security", "/healthy", "/ready"];
  if (skipPaths.includes(request.nextUrl.pathname)) {
    return NextResponse.next({ request: { headers: requestHeaders } });
  }

  // Build security response headers (shared by all routes).
  // Wrapped in try/catch so the middleware gracefully degrades with a default
  // CSP when the ZITADEL backend is unavailable (e.g. during container startup).
  const responseHeaders = new Headers();
  let publicHost: string | undefined;
  let instanceHost: string | undefined;
  let baseUrl: string | undefined;

  try {
    const _headers = await headers();
    const { serviceConfig } = getServiceConfig(_headers);
    publicHost = serviceConfig.publicHost;
    instanceHost = serviceConfig.instanceHost;
    baseUrl = serviceConfig.baseUrl;

    const securitySettings = await loadSecuritySettings(request);
    const iframeOrigins =
      securitySettings?.embeddedIframe?.enabled && securitySettings.embeddedIframe.allowedOrigins.length > 0
        ? securitySettings.embeddedIframe.allowedOrigins
        : undefined;

    responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl, iframeOrigins }));

    if (!iframeOrigins) {
      responseHeaders.set("X-Frame-Options", "deny");
    }
  } catch (err) {
    logger.error("Failed to load security settings for CSP, using default CSP", {
      error: err instanceof Error ? err.message : String(err),
    });
    responseHeaders.set("Content-Security-Policy", buildCSP());
    responseHeaders.set("X-Frame-Options", "deny");
  }

  // Only run the rest of the logic for the original matcher paths
  const proxyPaths = ["/.well-known/", "/oauth/", "/oidc/", "/idps/callback/", "/saml/"];
  const isMatched = proxyPaths.some((prefix) => request.nextUrl.pathname.startsWith(prefix));

  if (!isMatched) {
    return NextResponse.next({
      request: { headers: requestHeaders },
      headers: responseHeaders,
    });
  }

  // Proxy-specific headers
  if (publicHost) {
    requestHeaders.set("x-zitadel-public-host", publicHost);
  }
  if (instanceHost) {
    requestHeaders.set("x-zitadel-instance-host", instanceHost);
  }

  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  request.nextUrl.href = `${baseUrl}${request.nextUrl.pathname}${request.nextUrl.search}`;

  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
    headers: responseHeaders,
  });
}
