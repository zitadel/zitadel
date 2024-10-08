import { i18nRouter } from "next-i18n-router";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import i18nConfig from "../i18nConfig";

const INSTANCE = process.env.ZITADEL_API_URL;
const SERVICE_USER_ID = process.env.ZITADEL_SERVICE_USER_ID as string;

export const config = {
  matcher: "/((?!api|static|.*\\..*|_next).*)",
};

export function middleware(request: NextRequest) {
  // OIDC specific routes
  if (
    request.nextUrl.pathname.startsWith("/.well-known") ||
    request.nextUrl.pathname.startsWith("/oauth") ||
    request.nextUrl.pathname.startsWith("/oidc") ||
    request.nextUrl.pathname.startsWith("/idps/callback")
  ) {
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set("x-zitadel-login-client", SERVICE_USER_ID);

    // this is a workaround for the next.js server not forwarding the host header
    // requestHeaders.set("x-zitadel-forwarded", `host="${request.nextUrl.host}"`);
    requestHeaders.set("x-zitadel-public-host", `${request.nextUrl.host}`);

    // this is a workaround for the next.js server not forwarding the host header
    requestHeaders.set(
      "x-zitadel-instance-host",
      `${INSTANCE}`.replace("https://", ""),
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

  return i18nRouter(request, i18nConfig);
}
