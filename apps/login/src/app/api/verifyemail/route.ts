import { server, verifyEmail } from "@/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { userId, code } = body;

    return verifyEmail(server, userId, code)
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
