import { getAuthRequest, listSessions, server } from "#/lib/zitadel";
import { getAllSessionIds } from "#/utils/cookies";
import { Session } from "@zitadel/server";
import { NextRequest, NextResponse } from "next/server";

async function loadSessions(): Promise<Session[]> {
  const ids: string[] = await getAllSessionIds();

  if (ids && ids.length) {
    const response = await listSessions(
      server,
      ids.filter((id: string | undefined) => !!id)
    );
    return response?.sessions ?? [];
  } else {
    console.info("No session cookie found.");
    return [];
  }
}

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const authRequestId = searchParams.get("authRequest");

  if (authRequestId) {
    const response = await getAuthRequest(server, { authRequestId });
    const sessions = await loadSessions();
    if (sessions.length) {
      return NextResponse.json(sessions);
    } else {
      const loginNameUrl = new URL("/loginname", request.url);
      if (response.authRequest?.id) {
        loginNameUrl.searchParams.set(
          "authRequestId",
          response.authRequest?.id
        );
      }

      return NextResponse.redirect(loginNameUrl);
    }
  } else {
    return NextResponse.error();
  }
}
