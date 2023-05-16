import { createSession, server, setSession } from "#/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName } = body;

    const session = await createSession(server, loginName);
    return NextResponse.json(session);
  } else {
    return NextResponse.error();
  }
}

export async function PUT(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName } = body;

    const session = await setSession(server, loginName);
    return NextResponse.json(session);
  } else {
    return NextResponse.error();
  }
}
