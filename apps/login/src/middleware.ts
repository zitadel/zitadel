import { headers } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { getServiceUrlFromHeaders } from "./lib/service";

export const config = {
  matcher: [
    "/.well-known/:path*",
    "/oauth/:path*",
    "/oidc/:path*",
    "/idps/callback/:path*",
    "/saml/:path*",
  ],
};

export async function middleware(request: NextRequest) {
  // escape proxy if the environment is setup for multitenancy
  if (!process.env.ZITADEL_API_URL || !process.env.ZITADEL_SERVICE_USER_TOKEN) {
    return NextResponse.next();
  }

  const _headers = await headers();

  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const instanceHost = `${serviceUrl}`
    .replace("https://", "")
    .replace("http://", "");

  const requestHeaders = new Headers(request.headers);

  // this is a workaround for the next.js server not forwarding the host header
  // requestHeaders.set("x-zitadel-forwarded", `host="${request.nextUrl.host}"`);
  requestHeaders.set("x-zitadel-public-host", `${request.nextUrl.host}`);

  // this is a workaround for the next.js server not forwarding the host header
  requestHeaders.set("x-zitadel-instance-host", instanceHost);

  const responseHeaders = new Headers();
  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  request.nextUrl.href = `${serviceUrl}${request.nextUrl.pathname}${request.nextUrl.search}`;
  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
    headers: responseHeaders,
  });
}
