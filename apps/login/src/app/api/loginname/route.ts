import { listAuthenticationMethodTypes, listUsers } from "@/lib/zitadel";
import { createSessionForUserIdAndUpdateCookie } from "@/utils/session";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName, authRequestId, organization } = body;
    return listUsers(loginName, organization).then((users) => {
      if (
        users.details &&
        users.details.totalResult == 1 &&
        users.result[0].userId
      ) {
        const userId = users.result[0].userId;
        return createSessionForUserIdAndUpdateCookie(
          userId,
          undefined,
          undefined,
          authRequestId,
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
              throw { details: "No user id found in session" };
            }
          })
          .catch((error) => {
            console.error(error);
            return NextResponse.json(error, { status: 500 });
          });
      } else {
        return NextResponse.json(
          { message: "Could not find user" },
          { status: 404 },
        );
      }
    });
  } else {
    return NextResponse.error();
  }
}
