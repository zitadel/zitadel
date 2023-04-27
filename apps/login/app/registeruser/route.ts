import { addHumanUser, server } from "#/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { email, password, firstName, lastName } = body;

    const userId = await addHumanUser(server, {
      email: email,
      firstName,
      lastName,
      password: password,
    });
    return NextResponse.json({ userId });
  } else {
    return NextResponse.error();
  }
}
