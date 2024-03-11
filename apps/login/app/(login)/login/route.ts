import {
  createCallback,
  getAuthRequest,
  listSessions,
  server,
} from "#/lib/zitadel";
import { SessionCookie, getAllSessions } from "#/utils/cookies";
import { Session, AuthRequest, Prompt, login } from "@zitadel/server";
import { NextRequest, NextResponse } from "next/server";

async function loadSessions(ids: string[]): Promise<Session[]> {
  const response = await listSessions(
    server,
    ids.filter((id: string | undefined) => !!id)
  );
  return response?.sessions ?? [];
}

const ORG_SCOPE_REGEX = /urn:zitadel:iam:org:id:([0-9]*)/g;

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
  if (sessions.length) {
    return sessions[0];
  }
  return undefined;
}
export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const authRequestId = searchParams.get("authRequest");
  const sessionId = searchParams.get("sessionId");

  const sessionCookies: SessionCookie[] = await getAllSessions();
  const ids = sessionCookies.map((s) => s.id);
  let sessions: Session[] = [];
  if (ids && ids.length) {
    sessions = await loadSessions(ids);
  }

  if (authRequestId && sessionId) {
    console.log(
      `Login with session: ${sessionId} and authRequest: ${authRequestId}`
    );

    let selectedSession = sessions.find((s) => s.id === sessionId);

    if (selectedSession && selectedSession.id) {
      console.log(`Found session ${selectedSession.id}`);
      const cookie = sessionCookies.find(
        (cookie) => cookie.id === selectedSession?.id
      );

      if (cookie && cookie.id && cookie.token) {
        console.log(`Found sessioncookie ${cookie.id}`);

        const session = {
          sessionId: cookie?.id,
          sessionToken: cookie?.token,
        };

        const { callbackUrl } = await createCallback(server, {
          authRequestId,
          session,
        });
        return NextResponse.redirect(callbackUrl);
      }
    }
  }

  if (authRequestId) {
    console.log(`Login with authRequest: ${authRequestId}`);
    const { authRequest } = await getAuthRequest(server, { authRequestId });
    let organization;

    if (
      authRequest?.scope &&
      authRequest.scope.find((s) => ORG_SCOPE_REGEX.test(s))
    ) {
      const orgId = authRequest.scope.find((s) => ORG_SCOPE_REGEX.test(s));

      if (orgId) {
        const matched = orgId.replace("urn:zitadel:iam:org:id:", "");
        organization = matched;
      }
    }

    if (authRequest && authRequest.prompt.includes(Prompt.PROMPT_CREATE)) {
      const registerUrl = new URL("/register", request.url);
      if (authRequest?.id) {
        registerUrl.searchParams.set("authRequestId", authRequest?.id);
      }
      if (organization) {
        registerUrl.searchParams.set("organization", organization);
      }

      return NextResponse.redirect(registerUrl);
    }

    // use existing session and hydrate it for oidc
    if (authRequest && sessions.length) {
      // if some accounts are available for selection and select_account is set
      if (authRequest.prompt.includes(Prompt.PROMPT_SELECT_ACCOUNT)) {
        const accountsUrl = new URL("/accounts", request.url);
        if (authRequest?.id) {
          accountsUrl.searchParams.set("authRequestId", authRequest?.id);
        }
        if (organization) {
          accountsUrl.searchParams.set("organization", organization);
        }

        return NextResponse.redirect(accountsUrl);
      } else if (authRequest.prompt.includes(Prompt.PROMPT_LOGIN)) {
        // if prompt is login
        const loginNameUrl = new URL("/loginname", request.url);
        if (authRequest?.id) {
          loginNameUrl.searchParams.set("authRequestId", authRequest?.id);
        }
        if (authRequest.loginHint) {
          loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
        }
        if (organization) {
          loginNameUrl.searchParams.set("organization", organization);
        }
        return NextResponse.redirect(loginNameUrl);
      } else if (authRequest.prompt.includes(Prompt.PROMPT_NONE)) {
        // NONE prompt - silent authentication

        let selectedSession = findSession(sessions, authRequest);

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
            return NextResponse.json(
              { error: "No active session found" },
              { status: 500 } // TODO: check for correct status code
            );
          }
        } else {
          return NextResponse.json(
            { error: "No active session found" },
            { status: 500 } // TODO: check for correct status code
          );
        }
      } else {
        // check for loginHint, userId hint sessions
        let selectedSession = findSession(sessions, authRequest);

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
            accountsUrl.searchParams.set("authRequestId", authRequestId);
            if (organization) {
              accountsUrl.searchParams.set("organization", organization);
            }
            return NextResponse.redirect(accountsUrl);
          }
        } else {
          const accountsUrl = new URL("/accounts", request.url);
          accountsUrl.searchParams.set("authRequestId", authRequestId);
          if (organization) {
            accountsUrl.searchParams.set("organization", organization);
          }
          return NextResponse.redirect(accountsUrl);
        }
      }
    } else {
      const loginNameUrl = new URL("/loginname", request.url);

      loginNameUrl.searchParams.set("authRequestId", authRequestId);
      if (authRequest?.loginHint) {
        loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
        loginNameUrl.searchParams.set("submit", "true"); // autosubmit
      }

      if (organization) {
        loginNameUrl.searchParams.set("organization", organization);
      }

      return NextResponse.redirect(loginNameUrl);
    }
  } else {
    return NextResponse.json(
      { error: "No authRequestId provided" },
      { status: 500 }
    );
  }
}
