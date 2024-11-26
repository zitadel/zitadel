import { NextRequest, NextResponse } from "next/server";

export const config = {
  matcher: [
    "/.well-known/:path*",
    "/oauth/:path*",
    "/oidc/:path*",
    "/idps/callback/:path*",
  ],
};

const INSTANCE = process.env.ZITADEL_API_URL;
const SERVICE_USER_ID = process.env.ZITADEL_SERVICE_USER_ID as string;

export function middleware(request: NextRequest) {
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set("x-zitadel-login-client", SERVICE_USER_ID);

  // this is a workaround for the next.js server not forwarding the host header
  // requestHeaders.set("x-zitadel-forwarded", `host="${request.nextUrl.host}"`);
  requestHeaders.set("x-zitadel-public-host", `${request.nextUrl.host}`);

  // this is a workaround for the next.js server not forwarding the host header
  requestHeaders.set(
    "x-zitadel-instance-host",
    `${INSTANCE}`.replace(/^https?:\/\//, ""),
  );

  const responseHeaders = new Headers();
  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  request.nextUrl.href = `${INSTANCE}${request.nextUrl.pathname}${request.nextUrl.search}`;
  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
    headers: responseHeaders,
  });
}
