import {
  getSession,
  listAuthenticationMethodTypes,
  server,
} from "#/lib/zitadel";
import { getSessionCookieById } from "#/utils/cookies";
import { createSessionAndUpdateCookie } from "#/utils/session";
import { NextRequest, NextResponse } from "next/server";

// export async function GET(request: NextRequest) {
//   const { searchParams } = new URL(request.url);
//   const sessionId = searchParams.get("sessionId");
//   if (sessionId) {
//     const sessionCookie = await getSessionCookieById(sessionId);

//     const session = await getSession(
//       server,
//       sessionCookie.id,
//       sessionCookie.token
//     );

//     const userId = session?.session?.factors?.user?.id;

//     if (userId) {
//       return listAuthenticationMethodTypes(userId)
//         .then((methods) => {
//           return NextResponse.json(methods);
//         })
//         .catch((error) => {
//           return NextResponse.json(error, { status: 500 });
//         });
//     } else {
//       return NextResponse.json(
//         { details: "could not get session" },
//         { status: 500 }
//       );
//     }
//   } else {
//     return NextResponse.json({}, { status: 400 });
//   }
// }

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const { loginName, authRequestId } = body;

    // const domain: string = request.nextUrl.hostname;

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
