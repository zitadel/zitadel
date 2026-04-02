import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

/**
 * Console routes that the standalone console app links to without /console prefix.
 * When a user is viewing an instance, bare paths like /users/123 need to be
 * rewritten to /console/instances/{instanceId}/users/123.
 */
const CONSOLE_ROUTES = [
  "overview",
  "users",
  "organizations",
  "projects",
  "applications",
  "actions",
  "sessions",
  "administrators",
  "activity",
  "settings",
  "getting-started",
  "account-settings",
  "feedback",
  "roles",
  "analytics",
  "billing",
  "support",
  "usage",
]

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // 1. When visiting /console/instances/{id}/* — set the active-instance cookie
  const instanceMatch = pathname.match(/^\/console\/instances\/([^/]+)/)
  if (instanceMatch) {
    const instanceId = instanceMatch[1]
    const response = NextResponse.next()
    response.cookies.set("zitadel-active-instance", instanceId, {
      path: "/",
      httpOnly: false,
      sameSite: "lax",
    })
    return response
  }

  // 2. When visiting a bare console route (e.g. /users, /users/123)
  //    Rewrite to the instance-scoped path using the cookie
  const firstSegment = pathname.split("/")[1]
  if (CONSOLE_ROUTES.includes(firstSegment)) {
    const activeInstance = request.cookies.get("zitadel-active-instance")?.value
    if (activeInstance) {
      const url = request.nextUrl.clone()
      url.pathname = `/console/instances/${activeInstance}${pathname}`
      return NextResponse.rewrite(url)
    }
    // No active instance — rewrite to /console/* (standalone fallback)
    const url = request.nextUrl.clone()
    url.pathname = `/console${pathname}`
    return NextResponse.rewrite(url)
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    // Instance-scoped paths (to set cookie)
    "/console/instances/:path*",
    // Bare console routes (to rewrite)
    "/(overview|users|organizations|projects|applications|actions|sessions|administrators|activity|settings|getting-started|account-settings|feedback|roles|analytics|billing|support|usage)(.*)",
  ],
}
