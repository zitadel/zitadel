import {
  createCallback,
  getAuthRequest,
  listSessions,
  server,
} from "#/lib/zitadel";
import { SessionCookie, getAllSessions } from "#/utils/cookies";
import { Session, AuthRequest, Prompt } from "@zitadel/server";
import { NextRequest, NextResponse } from "next/server";

async function loadSessions(ids: string[]): Promise<Session[]> {
  const response = await listSessions(
    server,
    ids.filter((id: string | undefined) => !!id)
  );
  return response?.sessions ?? [];
}

function findSession(
  sessions: Session[],
  authRequest: AuthRequest
): Session | undefined {
  if (authRequest.hintUserId) {
    console.log(`find session for hintUserId: ${authRequest.hintUserId}`);
    return sessions.find((s) => s.factors?.user?.id === authRequest.hintUserId);
  }
  if (authRequest.loginHint) {
    console.log(`find session for loginHint: ${authRequest.loginHint}`);
    return sessions.find(
      (s) => s.factors?.user?.loginName === authRequest.loginHint
    );
  }
  return undefined;
}

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const authRequestId = searchParams.get("authRequest");

  if (authRequestId) {
    const { authRequest } = await getAuthRequest(server, { authRequestId });
    const sessionCookies: SessionCookie[] = await getAllSessions();
    const ids = sessionCookies.map((s) => s.id);

    let sessions: Session[] = [];
    if (ids && ids.length) {
      sessions = await loadSessions(ids);
    } else {
      console.info("No session cookie found.");
      return [];
    }

    // use existing session and hydrate it for oidc
    if (authRequest && sessions.length) {
      // if some accounts are available for selection and select_account is set
      if (authRequest && authRequest.prompt === Prompt.PROMPT_SELECT_ACCOUNT) {
        const accountsUrl = new URL("/accounts", request.url);
        if (authRequest?.id) {
          accountsUrl.searchParams.set("authRequestId", authRequest?.id);
        }

        return NextResponse.redirect(accountsUrl);
      } else {
        // check for loginHint, userId hint sessions
        let selectedSession = findSession(sessions, authRequest);

        // if (!selectedSession) {
        //   selectedSession = sessions[0]; // TODO: remove
        // }

        if (selectedSession && selectedSession.id) {
          const cookie = sessionCookies.find(
            (cookie) => cookie.id === selectedSession?.id
          );

          if (cookie && cookie.id && cookie.token) {
            const session = {
              sessionId: cookie?.id,
              sessionToken: cookie?.token,
            };
            const { callbackUrl } = await createCallback(server, {
              authRequestId,
              session,
            });
            return NextResponse.redirect(callbackUrl);
          } else {
            const accountsUrl = new URL("/accounts", request.url);
            if (authRequest?.id) {
              accountsUrl.searchParams.set("authRequestId", authRequest?.id);
            }

            return NextResponse.redirect(accountsUrl);
          }
        } else {
          const accountsUrl = new URL("/accounts", request.url);
          if (authRequest?.id) {
            accountsUrl.searchParams.set("authRequestId", authRequest?.id);
          }

          return NextResponse.redirect(accountsUrl);
          // return NextResponse.error();
        }
      }
    } else {
      const loginNameUrl = new URL("/loginname", request.url);
      if (authRequest?.id) {
        loginNameUrl.searchParams.set("authRequestId", authRequest?.id);
      }

      return NextResponse.redirect(loginNameUrl);
    }
  } else {
    return NextResponse.error();
  }
}
