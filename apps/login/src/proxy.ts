import { NextRequest, NextResponse } from "next/server";
import { buildCSP } from "./lib/csp";
import { applyCustomHeaders } from "./lib/custom-headers";
import { createLogger } from "./lib/logger";
import { getIframeOrigins } from "./lib/server/security-settings";
import { getServiceConfig } from "./lib/service-url";

const logger = createLogger("middleware");

const PROXY_PATHS = ["/.well-known/", "/oauth/", "/oidc/", "/idps/callback/", "/saml/"];
const SKIP_PATHS = ["/healthy", "/ready"];

function getBasePath(): string {
  return process.env.NEXT_PUBLIC_BASE_PATH || "";
}

export function normalizeProxyPathname(pathname: string, basePath = getBasePath()): string {
  if (!basePath) {
    return pathname;
  }
  if (pathname === basePath) {
    return "/";
  }
  if (pathname.startsWith(`${basePath}/`)) {
    return pathname.slice(basePath.length) || "/";
  }
  return pathname;
}

export function isProxyPath(pathname: string, basePath = getBasePath()): boolean {
  const normalizedPath = normalizeProxyPathname(pathname, basePath);
  return PROXY_PATHS.some((prefix) => normalizedPath.startsWith(prefix));
}

export const PROXY_MATCHER = [
  "/.well-known/:path*",
  "/oauth/:path*",
  "/oidc/:path*",
  "/idps/callback",
  "/idps/callback/:path*",
  "/saml/:path*",
  "/:path*",
];

export const config = {
  matcher: PROXY_MATCHER,
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
  const normalizedPathname = normalizeProxyPathname(request.nextUrl.pathname);

  if (SKIP_PATHS.includes(normalizedPathname)) {
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
  if (!isProxyPath(request.nextUrl.pathname)) {
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

  request.nextUrl.href = `${baseUrl}${normalizedPathname}${request.nextUrl.search}`;

  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
    headers: responseHeaders,
  });
}
