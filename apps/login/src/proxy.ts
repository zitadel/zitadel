import { NextRequest, NextResponse } from "next/server";
import { randomBytes } from "node:crypto";
import { statSync } from "node:fs";
import { resolve } from "node:path";
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
  // A caller must never choose the nonce used by the rendering tree.
  requestHeaders.delete("x-zitadel-csp-nonce");
  const nonce = process.env.CSP_NONCE_ENABLED === "true" ? randomBytes(16).toString("base64") : undefined;

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

  // Next's synthetic fallback documents are prerendered without a request
  // nonce. They are internal implementation paths, not supplied Login routes.
  if (nonce && ["/_not-found", "/_global-error"].includes(request.nextUrl.pathname)) {
    return new NextResponse(null, { status: 404, headers: { "Cache-Control": "private, no-store" } });
  }

  // Next's missing static-chunk error uses a prerendered shell, bypassing the
  // dynamic root layout. Keep real chunk caching but never serve that shell.
  if (nonce && request.nextUrl.pathname.startsWith("/_next/static/")) {
    const root = resolve(process.cwd(), ".next/static");
    const file = resolve(root, request.nextUrl.pathname.slice("/_next/static/".length));
    let found = false;
    try {
      found = file.startsWith(root + "/") && !!statSync(file, { throwIfNoEntry: false })?.isFile();
    } catch {}
    if (!found) return new NextResponse(null, { status: 404, headers: { "Cache-Control": "private, no-store" } });
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

      responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl, iframeOrigins, nonce }));

      if (!iframeOrigins) {
        responseHeaders.set("X-Frame-Options", "deny");
      }
    } catch (err) {
      logger.error("Failed to load security settings for CSP, using default CSP", {
        error: err instanceof Error ? err.message : String(err),
      });
      responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl, nonce }));
      responseHeaders.set("X-Frame-Options", "deny");
    }
  } else {
    responseHeaders.set("Content-Security-Policy", buildCSP({ serviceUrl: baseUrl, nonce }));
    responseHeaders.set("X-Frame-Options", "deny");
  }

  // Only proxy paths need to be rewritten to the ZITADEL backend
  const proxyPaths = ["/.well-known/", "/oauth/", "/oidc/", "/idps/callback/", "/saml/", "/assets/"];
  const isMatched = proxyPaths.some((prefix) => request.nextUrl.pathname.startsWith(prefix));

  if (!isMatched) {
    if (nonce) {
      // Next reads this policy when noncing bootstrap, chunk and streamed
      // scripts. Its middleware response policy must be identical.
      requestHeaders.set("x-zitadel-csp-nonce", nonce);
      requestHeaders.set("Content-Security-Policy", responseHeaders.get("Content-Security-Policy")!);
      requestHeaders.delete("Content-Security-Policy-Report-Only");
      // Static files keep Next's asset caching; dynamic fallback documents still
      // use the root layout's no-store rendering, including missing assets.
      const path = request.nextUrl.pathname;
      const staticAsset =
        path.startsWith("/_next/static/") ||
        path.startsWith("/favicon/") ||
        path.startsWith("/logo/") ||
        [
          "/favicon.ico",
          "/checkbox.svg",
          "/grid-light.svg",
          "/grid-dark.svg",
          "/zitadel-logo-light.svg",
          "/zitadel-logo-dark.svg",
        ].includes(path);
      if (!staticAsset) responseHeaders.set("Cache-Control", "private, no-store");
    }
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
