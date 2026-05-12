import { NextRequest, NextResponse } from "next/server";
import { buildCSP } from "./lib/csp";
import { applyCustomHeaders } from "./lib/custom-headers";
import { createLogger } from "./lib/logger";
import { getIframeOrigins } from "./lib/server/security-settings";
import { getServiceConfig } from "./lib/service-url";

const logger = createLogger("middleware");

export const config = {
  matcher: ["/.well-known/:path*", "/oauth/:path*", "/oidc/:path*", "/idps/callback/:path*", "/saml/:path*", "/:path*"],
};

export async function proxy(request: NextRequest) {
  // Add the original URL as a header to all requests
  const requestHeaders = new Headers(request.headers);

  // Extract "organization" search param from the URL and set it as a header if available
  const organization = request.nextUrl.searchParams.get("organization");
  if (organization) {
    requestHeaders.set("x-zitadel-i18n-organization", organization);
  }

  // Internal infrastructure routes — skip middleware entirely.
  // /healthy and /ready are Kubernetes/Docker health probes that must respond
  // without depending on a ZITADEL backend.
  const skipPaths = ["/healthy", "/ready"];
  if (skipPaths.includes(request.nextUrl.pathname)) {
    return NextResponse.next({ request: { headers: requestHeaders } });
  }

  const { serviceConfig } = getServiceConfig(request.headers);
  const { baseUrl, publicHost, instanceHost } = serviceConfig;

  // Build CSP headers using security settings fetched directly from the
  // ZITADEL API (no self-loopback through the load balancer).
  const responseHeaders = new Headers();

  const cspFetchEnabled = process.env.CSP_FETCH_ENABLED !== "false";

  if (cspFetchEnabled) {
    try {
      const iframeOrigins = await getIframeOrigins(baseUrl, instanceHost, publicHost);

      responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl, iframeOrigins }));

      if (!iframeOrigins) {
        responseHeaders.set("X-Frame-Options", "deny");
      }
    } catch (err) {
      logger.error("Failed to load security settings for CSP, using default CSP", {
        error: err instanceof Error ? err.message : String(err),
      });
      responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl }));
      responseHeaders.set("X-Frame-Options", "deny");
    }
  } else {
    responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl }));
    responseHeaders.set("X-Frame-Options", "deny");
  }

  // Only proxy paths need to be rewritten to the ZITADEL backend
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

  // Apply headers from CUSTOM_REQUEST_HEADERS environment variable
  applyCustomHeaders({
    set: (key, value) => requestHeaders.set(key, value),
    remove: (key) => requestHeaders.delete(key),
  });

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
