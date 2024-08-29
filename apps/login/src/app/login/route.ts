import {
  CreateCallbackRequestSchema,
  SessionSchema,
} from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";

export const dynamic = "force-dynamic";
export const revalidate = false;
export const fetchCache = "default-no-store";

import {
  createCallback,
  getActiveIdentityProviders,
  getAuthRequest,
  getOrgsByDomain,
  listSessions,
  startIdentityProviderFlow,
} from "@/lib/zitadel";
import { getAllSessions } from "@zitadel/next";
import { NextRequest, NextResponse } from "next/server";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import {
  AuthRequest,
  Prompt,
} from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { idpTypeToSlug } from "@/lib/idp";
import { create } from "@zitadel/client";

async function loadSessions(ids: string[]): Promise<Session[]> {
  const response = await listSessions(
    ids.filter((id: string | undefined) => !!id),
  );

  return response?.sessions ?? [];
}

const ORG_SCOPE_REGEX = /urn:zitadel:iam:org:id:([0-9]+)/;
const ORG_DOMAIN_SCOPE_REGEX = /urn:zitadel:iam:org:domain:primary:(.+)/; // TODO: check regex for all domain character options
const IDP_SCOPE_REGEX = /urn:zitadel:iam:org:idp:id:(.+)/;

function findSession(
  sessions: Session[],
  authRequest: AuthRequest,
): Session | undefined {
  if (authRequest.hintUserId) {
    console.log(`find session for hintUserId: ${authRequest.hintUserId}`);
    return sessions.find((s) => s.factors?.user?.id === authRequest.hintUserId);
  }
  if (authRequest.loginHint) {
    console.log(`find session for loginHint: ${authRequest.loginHint}`);
    return sessions.find(
      (s) => s.factors?.user?.loginName === authRequest.loginHint,
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

  // TODO: find a better way to handle _rsc (react server components) requests and block them to avoid conflicts when creating oidc callback
  const _rsc = searchParams.get("_rsc");
  if (_rsc) {
    return NextResponse.json({ error: "No _rsc supported" }, { status: 500 });
  }

  const sessionCookies = await getAllSessions();
  const ids = sessionCookies.map((s) => s.id);
  let sessions: Session[] = [];
  if (ids && ids.length) {
    sessions = await loadSessions(ids);
  }

  /**
   * TODO: before automatically redirecting to the callbackUrl, check if the session is still valid
   * possible scenaio:
   * mfa is required, session is not valid anymore (e.g. session expired, user logged out, etc.)
   * to check for mfa for automatically selected session -> const response = await listAuthenticationMethodTypes(userId);
   **/

  if (authRequestId && sessionId) {
    console.log(
      `Login with session: ${sessionId} and authRequest: ${authRequestId}`,
    );

    let selectedSession = sessions.find((s) => s.id === sessionId);

    if (selectedSession && selectedSession.id) {
      console.log(`Found session ${selectedSession.id}`);
      const cookie = sessionCookies.find(
        (cookie) => cookie.id === selectedSession?.id,
      );

      if (cookie && cookie.id && cookie.token) {
        const session = {
          sessionId: cookie?.id,
          sessionToken: cookie?.token,
        };

        // works not with _rsc request
        try {
          const { callbackUrl } = await createCallback(
            create(CreateCallbackRequestSchema, {
              authRequestId,
              callbackKind: {
                case: "session",
                value: create(SessionSchema, session),
              },
            }),
          );
          if (callbackUrl) {
            return NextResponse.redirect(callbackUrl);
          } else {
            return NextResponse.json(
              { error: "An error occurred!" },
              { status: 500 },
            );
          }
        } catch (error) {
          return NextResponse.json({ error }, { status: 500 });
        }
      }
    }
  }

  if (authRequestId) {
    const { authRequest } = await getAuthRequest({ authRequestId });

    let organization = "";
    let idpId = "";

    if (authRequest?.scope) {
      const orgScope = authRequest.scope.find((s: string) =>
        ORG_SCOPE_REGEX.test(s),
      );

      const idpScope = authRequest.scope.find((s: string) =>
        IDP_SCOPE_REGEX.test(s),
      );

      if (orgScope) {
        const matched = ORG_SCOPE_REGEX.exec(orgScope);
        organization = matched?.[1] ?? "";
      } else {
        const orgDomainScope = authRequest.scope.find((s: string) =>
          ORG_DOMAIN_SCOPE_REGEX.test(s),
        );

        if (orgDomainScope) {
          const matched = ORG_DOMAIN_SCOPE_REGEX.exec(orgDomainScope);
          const orgDomain = matched?.[1] ?? "";
          if (orgDomain) {
            const orgs = await getOrgsByDomain(orgDomain);
            if (orgs.result && orgs.result.length === 1) {
              organization = orgs.result[0].id ?? "";
            }
          }
        }
      }

      if (idpScope) {
        const matched = IDP_SCOPE_REGEX.exec(idpScope);
        idpId = matched?.[1] ?? "";

        const identityProviders = await getActiveIdentityProviders(
          organization ? organization : undefined,
        ).then((resp) => {
          return resp.identityProviders;
        });

        const idp = identityProviders.find((idp) => idp.id === idpId);

        if (idp) {
          const host = request.nextUrl.origin;

          const identityProviderType = identityProviders[0].type;
          let provider = idpTypeToSlug(identityProviderType);

          const params = new URLSearchParams();

          if (authRequestId) {
            params.set("authRequestId", authRequestId);
          }

          if (organization) {
            params.set("organization", organization);
          }

          return startIdentityProviderFlow({
            idpId,
            urls: {
              successUrl:
                `${host}/idp/${provider}/success?` +
                new URLSearchParams(params),
              failureUrl:
                `${host}/idp/${provider}/failure?` +
                new URLSearchParams(params),
            },
          }).then((resp) => {
            if (
              resp.nextStep.value &&
              typeof resp.nextStep.value === "string"
            ) {
              return NextResponse.redirect(resp.nextStep.value);
            }
          });
        }
      }
    }

    const gotoAccounts = (): NextResponse<unknown> => {
      const accountsUrl = new URL("/accounts", request.url);
      if (authRequest?.id) {
        accountsUrl.searchParams.set("authRequestId", authRequest?.id);
      }
      if (organization) {
        accountsUrl.searchParams.set("organization", organization);
      }

      return NextResponse.redirect(accountsUrl);
    };

    if (authRequest && authRequest.prompt.includes(Prompt.CREATE)) {
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
      if (authRequest.prompt.includes(Prompt.SELECT_ACCOUNT)) {
        return gotoAccounts();
      } else if (authRequest.prompt.includes(Prompt.LOGIN)) {
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
      } else if (authRequest.prompt.includes(Prompt.NONE)) {
        // NONE prompt - silent authentication

        let selectedSession = findSession(sessions, authRequest);

        if (selectedSession && selectedSession.id) {
          const cookie = sessionCookies.find(
            (cookie) => cookie.id === selectedSession?.id,
          );

          if (cookie && cookie.id && cookie.token) {
            const session = {
              sessionId: cookie?.id,
              sessionToken: cookie?.token,
            };
            const { callbackUrl } = await createCallback(
              create(CreateCallbackRequestSchema, {
                authRequestId,
                callbackKind: {
                  case: "session",
                  value: create(SessionSchema, session),
                },
              }),
            );
            return NextResponse.redirect(callbackUrl);
          } else {
            return NextResponse.json(
              { error: "No active session found" },
              { status: 400 }, // TODO: check for correct status code
            );
          }
        } else {
          return NextResponse.json(
            { error: "No active session found" },
            { status: 400 }, // TODO: check for correct status code
          );
        }
      } else {
        // check for loginHint, userId hint sessions
        let selectedSession = findSession(sessions, authRequest);

        if (selectedSession && selectedSession.id) {
          const cookie = sessionCookies.find(
            (cookie) => cookie.id === selectedSession?.id,
          );

          if (cookie && cookie.id && cookie.token) {
            const session = {
              sessionId: cookie?.id,
              sessionToken: cookie?.token,
            };
            try {
              const { callbackUrl } = await createCallback(
                create(CreateCallbackRequestSchema, {
                  authRequestId,
                  callbackKind: {
                    case: "session",
                    value: create(SessionSchema, session),
                  },
                }),
              );
              if (callbackUrl) {
                return NextResponse.redirect(callbackUrl);
              } else {
                console.log(
                  "could not create callback, redirect user to choose other account",
                );
                return gotoAccounts();
              }
            } catch (error) {
              console.error(error);
              return gotoAccounts();
            }
          } else {
            return gotoAccounts();
          }
        } else {
          return gotoAccounts();
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
      { status: 500 },
    );
  }
}
