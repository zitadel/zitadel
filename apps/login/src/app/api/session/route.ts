import {
  server,
  deleteSession,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import {
  SessionCookie,
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "@/utils/cookies";
import {
  createSessionAndUpdateCookie,
  createSessionForIdpAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "@/utils/session";
import { Challenges, Checks, RequestChallenges } from "@zitadel/server";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    const {
      userId,
      idpIntent,
      loginName,
      password,
      organization,
      authRequestId,
    } = body;

    if (userId && idpIntent) {
      return createSessionForIdpAndUpdateCookie(
        userId,
        idpIntent,
        organization,
        authRequestId,
      ).then((session) => {
        return NextResponse.json(session);
      });
    } else {
      return createSessionAndUpdateCookie(
        loginName,
        password,
        undefined,
        organization,
        authRequestId,
      ).then((session) => {
        return NextResponse.json(session);
      });
    }
  } else {
    return NextResponse.json(
      { details: "Session could not be created" },
      { status: 500 },
    );
  }
}

/**
 *
 * @param request password for the most recent session
 * @returns the updated most recent Session with the added password
 */
export async function PUT(request: NextRequest) {
  const body = await request.json();

  if (body) {
    const {
      loginName,
      sessionId,
      organization,
      checks,
      authRequestId,
      challenges,
    } = body;

    const recentPromise: Promise<SessionCookie> = sessionId
      ? getSessionCookieById(sessionId).catch((error) => {
          return Promise.reject(error);
        })
      : loginName
        ? getSessionCookieByLoginName(loginName, organization).catch(
            (error) => {
              return Promise.reject(error);
            },
          )
        : getMostRecentSessionCookie().catch((error) => {
            return Promise.reject(error);
          });

    const domain: string = request.nextUrl.hostname;

    if (challenges && challenges.webAuthN && !challenges.webAuthN.domain) {
      challenges.webAuthN.domain = domain;
    }

    return recentPromise
      .then(async (recent) => {
        if (
          challenges &&
          (challenges.otpEmail === "" || challenges.otpSms === "")
        ) {
          const sessionResponse = await getSession(
            server,
            recent.id,
            recent.token,
          );
          if (sessionResponse && sessionResponse.session?.factors?.user?.id) {
            const userResponse = await getUserByID(
              sessionResponse.session.factors.user.id,
            );
            if (
              challenges.otpEmail === "" &&
              userResponse.user?.human?.email?.email
            ) {
              challenges.otpEmail = userResponse.user?.human?.email?.email;
            }

            if (
              challenges.otpSms === "" &&
              userResponse.user?.human?.phone?.phone
            ) {
              challenges.otpSms = userResponse.user?.human?.phone?.phone;
            }
          }
        }

        return setSessionAndUpdateCookie(
          recent,
          checks,
          challenges,
          authRequestId,
        ).then(async (session) => {
          // if password, check if user has MFA methods
          let authMethods;
          if (checks && checks.password && session.factors?.user?.id) {
            const response = await listAuthenticationMethodTypes(
              session.factors?.user?.id,
            );
            if (response.authMethodTypes && response.authMethodTypes.length) {
              authMethods = response.authMethodTypes;
            }
          }

          return NextResponse.json({
            sessionId: session.id,
            factors: session.factors,
            challenges: session.challenges,
            authMethods,
          });
        });
      })
      .catch((error) => {
        console.error(error);
        return NextResponse.json({ details: error }, { status: 500 });
      });
  } else {
    return NextResponse.json(
      { details: "Request body is missing" },
      { status: 400 },
    );
  }
}

/**
 *
 * @param request id of the session to be deleted
 */
export async function DELETE(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const id = searchParams.get("id");
  if (id) {
    const session = await getSessionCookieById(id);

    return deleteSession(server, session.id, session.token)
      .then(() => {
        return removeSessionFromCookie(session)
          .then(() => {
            return NextResponse.json({});
          })
          .catch((error) => {
            return NextResponse.json(
              { details: "could not set cookie" },
              { status: 500 },
            );
          });
      })
      .catch((error) => {
        return NextResponse.json(
          { details: "could not delete session" },
          { status: 500 },
        );
      });
  } else {
    return NextResponse.error();
  }
}
