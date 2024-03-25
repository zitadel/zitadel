import { listAuthenticationMethodTypes, listUsers } from "#/lib/zitadel";
import {
  createSessionAndUpdateCookie,
  createSessionForUserIdAndUpdateCookie,
} from "#/utils/session";
import { U } from "@zitadel/server/dist/index-79b5dba4";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName, authRequestId, organization } = body;
    // TODO - search for users with org
    console.log(
      "loginName",
      loginName,
      "authRequestId",
      authRequestId,
      "organization",
      organization
    );
    return listUsers(loginName, organization).then((users) => {
      console.log("users", users);
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
          authRequestId
        )
          .then((session) => {
            console.log(session);
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
            return NextResponse.json(error, { status: 500 });
          });
      } else {
        return NextResponse.json(
          { message: "Could not find user" },
          { status: 404 }
        );
      }
    });
  } else {
    return NextResponse.error();
  }
}
