import { setEmail, server } from "@/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { userId, code } = body;

    // replace with resend Mail method once its implemented
    return setEmail(server, userId)
      .then((resp) => {
        return NextResponse.json(resp);
      })
      .catch((error) => {
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.error();
  }
}
