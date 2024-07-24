import {
  getLoginSettings,
  listAuthenticationMethodTypes,
  listUsers,
} from "@/lib/zitadel";
import { createSessionForUserIdAndUpdateCookie } from "@/utils/session";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName, authRequestId, organization } = body;
    return listUsers(loginName, organization).then(async (users) => {
      if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
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
      } else if (organization) {
        const loginSettings = await getLoginSettings(organization);

        // user not found, check if register is enabled on organization
        if (loginSettings?.allowRegister) {
          const params: any = { organization };
          if (authRequestId) {
            params.authRequestId = authRequestId;
          }
          if (loginName) {
            params.email = loginName;
          }

          const registerUrl = new URL(
            "/register?" + new URLSearchParams(params),
            request.url,
          );

          return NextResponse.json({
            nextUrl: registerUrl,
            status: 200,
          });
        } else {
          return NextResponse.json(
            { message: "Could not find user" },
            { status: 404 },
          );
        }
      }
    });
  } else {
    return NextResponse.error();
  }
}
