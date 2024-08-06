import { listSessions } from "@/lib/zitadel";
import { SessionCookie, getAllSessions } from "@/utils/cookies";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { NextRequest, NextResponse } from "next/server";

async function loadSessions(ids: string[]): Promise<Session[]> {
  const response = await listSessions(
    ids.filter((id: string | undefined) => !!id),
  );

  return response?.sessions ?? [];
}

export async function GET(request: NextRequest) {
  const sessionCookies: SessionCookie[] = await getAllSessions();
  const ids = sessionCookies.map((s) => s.id);
  let sessions: Session[] = [];
  if (ids && ids.length) {
    sessions = await loadSessions(ids);
  }

  const responseHeaders = new Headers();
  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  return NextResponse.json(
    { sessions },
    { status: 200, headers: responseHeaders },
  );
}
