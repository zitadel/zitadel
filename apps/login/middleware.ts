import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export const config = {
  matcher: ["/.well-known/:path*", "/oauth/:path*", "/oidc/:path*"],
};

const INSTANCE = process.env.ZITADEL_API_URL;
const SERVICE_USER_ID = process.env.ZITADEL_SERVICE_USER_ID as string;

export function middleware(request: NextRequest) {
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set("x-zitadel-login-client", SERVICE_USER_ID);

  const proto = request.nextUrl.protocol.replace(":", "");
  requestHeaders.set(
    "Forwarded",
    `host=${request.nextUrl.host};proto=${proto}`
  );

  console.log(
    "intercept",
    request.nextUrl.basePath,
    request.nextUrl.host,
    proto
  );

  request.nextUrl.href = `${INSTANCE}${request.nextUrl.pathname}${request.nextUrl.search}`;
  return NextResponse.rewrite(request.nextUrl, {
    request: {
      headers: requestHeaders,
    },
  });
}
