import { addHumanUser, server } from "@/lib/zitadel";
import {
  createSessionAndUpdateCookie,
  createSessionForUserIdAndUpdateCookie,
} from "@/utils/session";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const {
      email,
      password,
      firstName,
      lastName,
      organization,
      authRequestId,
    } = body;

    return addHumanUser(server, {
      email: email,
      firstName,
      lastName,
      password: password ? password : undefined,
      organization,
    })
      .then((user) => {
        return createSessionForUserIdAndUpdateCookie(
          user.userId,
          password,
          undefined,
          authRequestId,
        ).then((session) => {
          return NextResponse.json({
            userId: user.userId,
            sessionId: session.id,
            factors: session.factors,
          });
        });
      })
      .catch((error) => {
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.error();
  }
}
