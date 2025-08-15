import { getAllSessions } from "@/lib/cookies";
import { idpTypeToSlug } from "@/lib/idp";
import { loginWithOIDCAndSession } from "@/lib/oidc";
import { loginWithSAMLAndSession } from "@/lib/saml";
import { sendLoginname, SendLoginnameCommand } from "@/lib/server/loginname";
import { constructUrl, getServiceUrlFromHeaders } from "@/lib/service-url";
import { findValidSession } from "@/lib/session";
import {
  createCallback,
  createResponse,
  getActiveIdentityProviders,
  getAuthRequest,
  getOrgsByDomain,
  getSAMLRequest,
  getSecuritySettings,
  listSessions,
  startIdentityProviderFlow,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { Prompt } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import {
  CreateCallbackRequestSchema,
  SessionSchema,
} from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { CreateResponseRequestSchema } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { headers } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { DEFAULT_CSP } from "../../../constants/csp";

export const dynamic = "force-dynamic";
export const revalidate = false;
export const fetchCache = "default-no-store";

const gotoAccounts = ({
  request,
  requestId,
  organization,
}: {
  request: NextRequest;
  requestId: string;
  organization?: string;
}): NextResponse<unknown> => {
  const accountsUrl = constructUrl(request, "/accounts");

  if (requestId) {
    accountsUrl.searchParams.set("requestId", requestId);
  }
  if (organization) {
    accountsUrl.searchParams.set("organization", organization);
  }

  return NextResponse.redirect(accountsUrl);
};

async function loadSessions({
  serviceUrl,
  ids,
}: {
  serviceUrl: string;
  ids: string[];
}): Promise<Session[]> {
  const response = await listSessions({
    serviceUrl,
    ids: ids.filter((id: string | undefined) => !!id),
  });

  return response?.sessions ?? [];
}

const ORG_SCOPE_REGEX = /urn:zitadel:iam:org:id:([0-9]+)/;
const ORG_DOMAIN_SCOPE_REGEX = /urn:zitadel:iam:org:domain:primary:(.+)/; // TODO: check regex for all domain character options
const IDP_SCOPE_REGEX = /urn:zitadel:iam:org:idp:id:(.+)/;

export async function GET(request: NextRequest) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const searchParams = request.nextUrl.searchParams;

  const oidcRequestId = searchParams.get("authRequest"); // oidc initiated request
  const samlRequestId = searchParams.get("samlRequest"); // saml initiated request

  // internal request id which combines authRequest and samlRequest with the prefix oidc_ or saml_
  let requestId =
    searchParams.get("requestId") ??
    (oidcRequestId
      ? `oidc_${oidcRequestId}`
      : samlRequestId
        ? `saml_${samlRequestId}`
        : undefined);

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
    sessions = await loadSessions({ serviceUrl, ids });
  }

  // complete flow if session and request id are provided
  if (requestId && sessionId) {
    if (requestId.startsWith("oidc_")) {
      // this finishes the login process for OIDC
      return loginWithOIDCAndSession({
        serviceUrl,
        authRequest: requestId.replace("oidc_", ""),
        sessionId,
        sessions,
        sessionCookies,
        request,
      });
    } else if (requestId.startsWith("saml_")) {
      // this finishes the login process for SAML
      return loginWithSAMLAndSession({
        serviceUrl,
        samlRequest: requestId.replace("saml_", ""),
        sessionId,
        sessions,
        sessionCookies,
        request,
      });
    }
  }

  // continue with OIDC
  if (requestId && requestId.startsWith("oidc_")) {
    const { authRequest } = await getAuthRequest({
      serviceUrl,
      authRequestId: requestId.replace("oidc_", ""),
    });

    let organization = "";
    let suffix = "";
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
            const orgs = await getOrgsByDomain({
              serviceUrl,
              domain: orgDomain,
            });
            if (orgs.result && orgs.result.length === 1) {
              organization = orgs.result[0].id ?? "";
              suffix = orgDomain;
            }
          }
        }
      }

      if (idpScope) {
        const matched = IDP_SCOPE_REGEX.exec(idpScope);
        idpId = matched?.[1] ?? "";

        const identityProviders = await getActiveIdentityProviders({
          serviceUrl,
          orgId: organization ? organization : undefined,
        }).then((resp) => {
          return resp.identityProviders;
        });

        const idp = identityProviders.find((idp) => idp.id === idpId);

        if (idp) {
          const origin = request.nextUrl.origin;

          const identityProviderType = identityProviders[0].type;

          if (identityProviderType === IdentityProviderType.LDAP) {
            const ldapUrl = constructUrl(request, "/ldap");
            if (authRequest.id) {
              ldapUrl.searchParams.set("requestId", `oidc_${authRequest.id}`);
            }
            if (organization) {
              ldapUrl.searchParams.set("organization", organization);
            }

            return NextResponse.redirect(ldapUrl);
          }

          let provider = idpTypeToSlug(identityProviderType);

          const params = new URLSearchParams();

          if (requestId) {
            params.set("requestId", requestId);
          }

          if (organization) {
            params.set("organization", organization);
          }

          let url: string | null = await startIdentityProviderFlow({
            serviceUrl,
            idpId,
            urls: {
              successUrl:
                `${origin}/idp/${provider}/success?` +
                new URLSearchParams(params),
              failureUrl:
                `${origin}/idp/${provider}/failure?` +
                new URLSearchParams(params),
            },
          });

          if (!url) {
            return NextResponse.json(
              { error: "Could not start IDP flow" },
              { status: 500 },
            );
          }

          if (url.startsWith("/")) {
            // if the url is a relative path, construct the absolute url
            url = constructUrl(request, url).toString();
          }

          return NextResponse.redirect(url);
        }
      }
    }

    if (authRequest && authRequest.prompt.includes(Prompt.CREATE)) {
      const registerUrl = constructUrl(request, "/register");
      if (authRequest.id) {
        registerUrl.searchParams.set("requestId", `oidc_${authRequest.id}`);
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
        return gotoAccounts({
          request,
          requestId: `oidc_${authRequest.id}`,
          organization,
        });
      } else if (authRequest.prompt.includes(Prompt.LOGIN)) {
        /**
         * The login prompt instructs the authentication server to prompt the user for re-authentication, regardless of whether the user is already authenticated
         */

        // if a hint is provided, skip loginname page and jump to the next page
        if (authRequest.loginHint) {
          try {
            let command: SendLoginnameCommand = {
              loginName: authRequest.loginHint,
              requestId: authRequest.id,
            };

            if (organization) {
              command = { ...command, organization };
            }

            const res = await sendLoginname(command);

            if (res && "redirect" in res && res?.redirect) {
              const absoluteUrl = constructUrl(request, res.redirect);
              return NextResponse.redirect(absoluteUrl.toString());
            }
          } catch (error) {
            console.error("Failed to execute sendLoginname:", error);
          }
        }

        const loginNameUrl = constructUrl(request, "/loginname");
        if (authRequest.id) {
          loginNameUrl.searchParams.set("requestId", `oidc_${authRequest.id}`);
        }
        if (authRequest.loginHint) {
          loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
        }
        if (organization) {
          loginNameUrl.searchParams.set("organization", organization);
        }
        if (suffix) {
          loginNameUrl.searchParams.set("suffix", suffix);
        }
        return NextResponse.redirect(loginNameUrl);
      } else if (authRequest.prompt.includes(Prompt.NONE)) {
        /**
         * With an OIDC none prompt, the authentication server must not display any authentication or consent user interface pages.
         * This means that the user should not be prompted to enter their password again.
         * Instead, the server attempts to silently authenticate the user using an existing session or other authentication mechanisms that do not require user interaction
         **/
        const securitySettings = await getSecuritySettings({
          serviceUrl,
        });

        const selectedSession = await findValidSession({
          serviceUrl,
          sessions,
          authRequest,
        });

        const noSessionResponse = NextResponse.json(
          { error: "No active session found" },
          { status: 400 },
        );

        if (securitySettings?.embeddedIframe?.enabled) {
          securitySettings.embeddedIframe.allowedOrigins;
          noSessionResponse.headers.set(
            "Content-Security-Policy",
            `${DEFAULT_CSP} frame-ancestors ${securitySettings.embeddedIframe.allowedOrigins.join(" ")};`,
          );
          noSessionResponse.headers.delete("X-Frame-Options");
        }

        if (!selectedSession || !selectedSession.id) {
          return noSessionResponse;
        }

        const cookie = sessionCookies.find(
          (cookie) => cookie.id === selectedSession.id,
        );

        if (!cookie || !cookie.id || !cookie.token) {
          return noSessionResponse;
        }

        const session = {
          sessionId: cookie.id,
          sessionToken: cookie.token,
        };

        const { callbackUrl } = await createCallback({
          serviceUrl,
          req: create(CreateCallbackRequestSchema, {
            authRequestId: requestId.replace("oidc_", ""),
            callbackKind: {
              case: "session",
              value: create(SessionSchema, session),
            },
          }),
        });

        const callbackResponse = NextResponse.redirect(callbackUrl);

        if (securitySettings?.embeddedIframe?.enabled) {
          securitySettings.embeddedIframe.allowedOrigins;
          callbackResponse.headers.set(
            "Content-Security-Policy",
            `${DEFAULT_CSP} frame-ancestors ${securitySettings.embeddedIframe.allowedOrigins.join(" ")};`,
          );
          callbackResponse.headers.delete("X-Frame-Options");
        }

        return callbackResponse;
      } else {
        // check for loginHint, userId hint and valid sessions
        let selectedSession = await findValidSession({
          serviceUrl,
          sessions,
          authRequest,
        });

        if (!selectedSession || !selectedSession.id) {
          return gotoAccounts({
            request,
            requestId: `oidc_${authRequest.id}`,
            organization,
          });
        }

        const cookie = sessionCookies.find(
          (cookie) => cookie.id === selectedSession.id,
        );

        if (!cookie || !cookie.id || !cookie.token) {
          return gotoAccounts({
            request,
            requestId: `oidc_${authRequest.id}`,
            organization,
          });
        }

        const session = {
          sessionId: cookie.id,
          sessionToken: cookie.token,
        };

        try {
          const { callbackUrl } = await createCallback({
            serviceUrl,
            req: create(CreateCallbackRequestSchema, {
              authRequestId: requestId.replace("oidc_", ""),
              callbackKind: {
                case: "session",
                value: create(SessionSchema, session),
              },
            }),
          });
          if (callbackUrl) {
            return NextResponse.redirect(callbackUrl);
          } else {
            console.log(
              "could not create callback, redirect user to choose other account",
            );
            return gotoAccounts({
              request,
              organization,
              requestId: `oidc_${authRequest.id}`,
            });
          }
        } catch (error) {
          console.error(error);
          return gotoAccounts({
            request,
            requestId: `oidc_${authRequest.id}`,
            organization,
          });
        }
      }
    } else {
      const loginNameUrl = constructUrl(request, "/loginname");

      loginNameUrl.searchParams.set("requestId", requestId);
      if (authRequest?.loginHint) {
        loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
        loginNameUrl.searchParams.set("submit", "true"); // autosubmit
      }

      if (organization) {
        loginNameUrl.searchParams.append("organization", organization);
        // loginNameUrl.searchParams.set("organization", organization);
      }

      return NextResponse.redirect(loginNameUrl);
    }
  }
  // continue with SAML
  else if (requestId && requestId.startsWith("saml_")) {
    const { samlRequest } = await getSAMLRequest({
      serviceUrl,
      samlRequestId: requestId.replace("saml_", ""),
    });

    if (!samlRequest) {
      return NextResponse.json(
        { error: "No samlRequest found" },
        { status: 400 },
      );
    }

    let selectedSession = await findValidSession({
      serviceUrl,
      sessions,
      samlRequest,
    });

    if (!selectedSession || !selectedSession.id) {
      return gotoAccounts({
        request,
        requestId: `saml_${samlRequest.id}`,
      });
    }

    const cookie = sessionCookies.find(
      (cookie) => cookie.id === selectedSession.id,
    );

    if (!cookie || !cookie.id || !cookie.token) {
      return gotoAccounts({
        request,
        requestId: `saml_${samlRequest.id}`,
        // organization,
      });
    }

    const session = {
      sessionId: cookie.id,
      sessionToken: cookie.token,
    };

    try {
      const { url, binding } = await createResponse({
        serviceUrl,
        req: create(CreateResponseRequestSchema, {
          samlRequestId: requestId.replace("saml_", ""),
          responseKind: {
            case: "session",
            value: session,
          },
        }),
      });
      if (url && binding.case === "redirect") {
        return NextResponse.redirect(url);
      } else if (url && binding.case === "post") {
        // Create HTML form that auto-submits via POST and escape the SAML cookie
        const html = `
          <html>
            <body onload="document.forms[0].submit()">
              <form action="${url}" method="post">
                <input type="hidden" name="RelayState" value="${binding.value.relayState}" />
                <input type="hidden" name="SAMLResponse" value="${binding.value.samlResponse}" />
                <noscript>
                  <button type="submit">Continue</button>
                </noscript>
              </form>
            </body>
          </html>
        `;

        return new NextResponse(html, {
          headers: { "Content-Type": "text/html" },
        });
      } else {
        console.log(
          "could not create response, redirect user to choose other account",
        );
        return gotoAccounts({
          request,
          requestId: `saml_${samlRequest.id}`,
        });
      }
    } catch (error) {
      console.error(error);
      return gotoAccounts({
        request,
        requestId: `saml_${samlRequest.id}`,
      });
    }
  }
  // Device Authorization does not need to start here as it is handled on the /device endpoint
  else {
    return NextResponse.json(
      { error: "No authRequest nor samlRequest provided" },
      { status: 500 },
    );
  }
}
