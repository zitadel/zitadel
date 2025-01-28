import { headers } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { getApiUrlOfHeaders } from "./lib/service";

export const config = {
  matcher: [
    "/.well-known/:path*",
    "/oauth/:path*",
    "/oidc/:path*",
    "/idps/callback/:path*",
  ],
};

export async function middleware(request: NextRequest) {
  // escape proxy if the environment is setup for multitenancy
  // if (
  //   !process.env.ZITADEL_API_URL ||
  //   !process.env.ZITADEL_USER_ID ||
  //   !process.env.ZITADEL_USER_TOKEN
  // ) {
  //   return NextResponse.next();
  // }
  const _headers = await headers();

  const instanceUrl = getApiUrlOfHeaders(_headers);

  const instanceHost = `${instanceUrl}`.replace("https://", "");

  const requestHeaders = new Headers(request.headers);
  requestHeaders.set(
    "x-zitadel-login-client",
    process.env.ZITADEL_SERVICE_USER_ID,
  );

  // this is a workaround for the next.js server not forwarding the host header
  // requestHeaders.set("x-zitadel-forwarded", `host="${request.nextUrl.host}"`);
  requestHeaders.set("x-zitadel-public-host", `${request.nextUrl.host}`);

  // this is a workaround for the next.js server not forwarding the host header
  requestHeaders.set("x-zitadel-instance-host", instanceHost);

  const responseHeaders = new Headers();
  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  request.nextUrl.href = `${instanceUrl}${request.nextUrl.pathname}${request.nextUrl.search}`;
  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
    headers: responseHeaders,
  });
}
