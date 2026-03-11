import { type NextRequest, NextResponse } from "next/server";
import { decrypt } from "./crypto.js";

export interface MiddlewareOptions {
  /** Paths that require authentication (supports glob-like prefix matching). */
  protectedPaths?: string[];
  /** The sign-in URL to redirect unauthenticated users to. Defaults to `/api/auth/signin`. */
  signInUrl?: string;
  /**
   * Cookie secret for decrypting the session cookie.
   * Falls back to the `ZITADEL_COOKIE_SECRET` environment variable.
   */
  cookieSecret?: string;
}

const SESSION_COOKIE = "zitadel.auth.oidc.session";

/**
 * Creates a Next.js middleware function that protects routes based on session state.
 *
 * @example
 * ```ts
 * // middleware.ts
 * import { createZitadelMiddleware } from "@zitadel/nextjs";
 *
 * export default createZitadelMiddleware({
 *   protectedPaths: ["/dashboard", "/settings"],
 *   signInUrl: "/api/auth/signin",
 * });
 *
 * export const config = {
 *   matcher: ["/dashboard/:path*", "/settings/:path*"],
 * };
 * ```
 */
export function createZitadelMiddleware(options?: MiddlewareOptions) {
  const protectedPaths = options?.protectedPaths ?? ["/"];
  const signInUrl = options?.signInUrl ?? "/api/auth/signin";

  return async function middleware(request: NextRequest): Promise<NextResponse> {
    const { pathname } = request.nextUrl;

    // Check if the current path is protected
    const isProtected = protectedPaths.some((p) => pathname.startsWith(p));
    if (!isProtected) {
      return NextResponse.next();
    }

    const secret =
      options?.cookieSecret ?? process.env.ZITADEL_COOKIE_SECRET;
    if (!secret || secret.length < 32) {
      // Misconfigured — redirect to sign-in as a safety fallback
      return NextResponse.redirect(new URL(signInUrl, request.url));
    }

    // Read the session cookie
    const sessionCookie = request.cookies.get(SESSION_COOKIE);
    if (!sessionCookie?.value) {
      return NextResponse.redirect(new URL(signInUrl, request.url));
    }

    // Decrypt and validate the session
    const json = await decrypt(sessionCookie.value, secret);
    if (!json) {
      return NextResponse.redirect(new URL(signInUrl, request.url));
    }

    const session = JSON.parse(json) as {
      accessToken: string;
      expiresAt: number;
    };

    // Check expiry
    const now = Math.floor(Date.now() / 1000);
    if (session.expiresAt <= now) {
      return NextResponse.redirect(new URL(signInUrl, request.url));
    }

    // Attach session info to request headers for downstream server components
    const response = NextResponse.next();
    response.headers.set("x-zitadel-access-token", session.accessToken);
    return response;
  };
}
