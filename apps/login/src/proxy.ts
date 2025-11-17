import { SecuritySettings } from "@zitadel/proto/zitadel/settings/v2/security_settings_pb";

import { headers } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { DEFAULT_CSP } from "../constants/csp";
import { getServiceUrlFromHeaders } from "./lib/service-url";
export const config = {
  matcher: [
    "/.well-known/:path*",
    "/oauth/:path*",
    "/oidc/:path*",
    "/idps/callback/:path*",
    "/saml/:path*",
    "/:path*",
  ],
};

async function loadSecuritySettings(
  request: NextRequest,
): Promise<SecuritySettings | null> {
  const securityResponse = await fetch(`${request.nextUrl.origin}/security`);

  if (!securityResponse.ok) {
    console.error(
      "Failed to fetch security settings:",
      securityResponse.statusText,
    );
    return null;
  }

  const response = await securityResponse.json();

  if (!response || !response.settings) {
    console.error("No security settings found in the response.");
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

  // Only run the rest of the logic for the original matcher paths
  const proxyPaths = [
    "/.well-known/",
    "/oauth/",
    "/oidc/",
    "/idps/callback/",
    "/saml/",
  ];

  const isMatched = proxyPaths.some((prefix) =>
    request.nextUrl.pathname.startsWith(prefix),
  );

  // escape proxy if the environment is setup for multitenancy
  if (
    !isMatched ||
    !process.env.ZITADEL_API_URL ||
    !process.env.ZITADEL_SERVICE_USER_TOKEN
  ) {
    // For all other routes, just add the header and continue
    return NextResponse.next({
      request: { headers: requestHeaders },
    });
  }

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const instanceHost = `${serviceUrl}`
    .replace("https://", "")
    .replace("http://", "");

  // Add additional headers as before
  requestHeaders.set("x-zitadel-public-host", `${request.nextUrl.host}`);
  requestHeaders.set("x-zitadel-instance-host", instanceHost);

  const responseHeaders = new Headers();
  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  const securitySettings = await loadSecuritySettings(request);

  if (securitySettings?.embeddedIframe?.enabled) {
    responseHeaders.set(
      "Content-Security-Policy",
      `${DEFAULT_CSP} frame-ancestors ${securitySettings.embeddedIframe.allowedOrigins.join(" ")};`,
    );
    responseHeaders.delete("X-Frame-Options");
  }

  request.nextUrl.href = `${serviceUrl}${request.nextUrl.pathname}${request.nextUrl.search}`;

  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
    headers: responseHeaders,
  });
}
