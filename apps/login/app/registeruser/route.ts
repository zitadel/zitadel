import { addHumanUser, server } from "#/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { email, password, firstName, lastName } = body;

    return addHumanUser(server, {
      email: email,
      firstName,
      lastName,
      password: password ? password : undefined,
    })
      .then((userId) => {
        return NextResponse.json({ userId });
      })
      .catch((error) => {
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.error();
  }
}
