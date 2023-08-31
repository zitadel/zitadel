import { listAuthenticationMethodTypes } from "#/lib/zitadel";
import { createSessionAndUpdateCookie } from "#/utils/session";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName, authRequestId } = body;

    return createSessionAndUpdateCookie(
      loginName,
      undefined,
      undefined,
      authRequestId
    )
      .then((session) => {
        if (session.factors?.user?.id) {
          return listAuthenticationMethodTypes(session.factors?.user?.id)
            .then((methods) => {
              return NextResponse.json({
                authMethodTypes: methods.authMethodTypes,
                sessionId: session.id,
                factors: session.factors,
              });
            })
            .catch((error) => {
              return NextResponse.json(error, { status: 500 });
            });
        } else {
          throw "No user id found in session";
        }
      })
      .catch((error) => {
        return NextResponse.json(
          {
            details: "could not add session to cookie",
          },
          { status: 500 }
        );
      });
  } else {
    return NextResponse.error();
  }
}
