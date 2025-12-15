import { constructUrl } from "@/lib/service-url";
import { findValidSession } from "@/lib/session";
import {
  createCallback,
  createResponse,
  getActiveIdentityProviders,
  getAuthRequest,
  getOrgsByDomain,
  getSAMLRequest,
  getSecuritySettings,
  startIdentityProviderFlow,
  ServiceConfig,
} from "@/lib/zitadel";
import { sendLoginname, SendLoginnameCommand } from "@/lib/server/loginname";
import { idpTypeToSlug } from "@/lib/idp";
import { create } from "@zitadel/client";
import { Prompt } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { CreateCallbackRequestSchema, SessionSchema } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { CreateResponseRequestSchema } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { NextRequest, NextResponse } from "next/server";
import { DEFAULT_CSP } from "../../../constants/csp";

const ORG_SCOPE_REGEX = /urn:zitadel:iam:org:id:([0-9]+)/;
const ORG_DOMAIN_SCOPE_REGEX = /urn:zitadel:iam:org:domain:primary:(.+)/;
const IDP_SCOPE_REGEX = /urn:zitadel:iam:org:idp:id:(.+)/;

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

export interface FlowInitiationParams {
  serviceConfig: ServiceConfig;
  requestId: string;
  sessions: Session[];
  sessionCookies: any[];
  request: NextRequest;
}

/**
 * Handle OIDC flow initiation
 */
export async function handleOIDCFlowInitiation(params: FlowInitiationParams): Promise<NextResponse> {
  const { serviceConfig, requestId, sessions, sessionCookies, request } = params;

  const { authRequest } = await getAuthRequest({ serviceConfig, authRequestId: requestId.replace("oidc_", ""),
  });

  let organization = "";
  let suffix = "";
  let idpId = "";

  if (authRequest?.scope) {
    const orgScope = authRequest.scope.find((s: string) => ORG_SCOPE_REGEX.test(s));
    const idpScope = authRequest.scope.find((s: string) => IDP_SCOPE_REGEX.test(s));

    if (orgScope) {
      const matched = ORG_SCOPE_REGEX.exec(orgScope);
      organization = matched?.[1] ?? "";
    } else {
      const orgDomainScope = authRequest.scope.find((s: string) => ORG_DOMAIN_SCOPE_REGEX.test(s));

      if (orgDomainScope) {
        const matched = ORG_DOMAIN_SCOPE_REGEX.exec(orgDomainScope);
        const orgDomain = matched?.[1] ?? "";

        console.log("Extracted org domain:", orgDomain);
        if (orgDomain) {
          const orgs = await getOrgsByDomain({ serviceConfig, domain: orgDomain,
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

      const identityProviders = await getActiveIdentityProviders({ serviceConfig, orgId: organization ? organization : undefined,
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

        const params = new URLSearchParams({
          requestId: requestId,
        });

        if (organization) {
          params.set("organization", organization);
        }

        let url: string | null = await startIdentityProviderFlow({ serviceConfig, idpId,
          urls: {
            successUrl: `${origin}/idp/${provider}/process?` + new URLSearchParams(params),
            failureUrl: `${origin}/idp/${provider}/failure?` + new URLSearchParams(params),
          },
        });

        if (!url) {
          return NextResponse.json({ error: "Could not start IDP flow" }, { status: 500 });
        }

        if (url.startsWith("/")) {
          url = constructUrl(request, url).toString();
        }

        return NextResponse.redirect(url);
      }
    }
  }

  if (authRequest && authRequest.prompt.includes(Prompt.CREATE)) {
    const registerUrl = constructUrl(request, "/register");
    registerUrl.searchParams.set("requestId", requestId);

    if (organization) {
      registerUrl.searchParams.set("organization", organization);
    }

    return NextResponse.redirect(registerUrl);
  }

  // use existing session and hydrate it for oidc
  if (authRequest && sessions.length) {
    if (authRequest.prompt.includes(Prompt.SELECT_ACCOUNT)) {
      return gotoAccounts({
        request,
        requestId: `oidc_${authRequest.id}`,
        organization,
      });
    } else if (authRequest.prompt.includes(Prompt.LOGIN)) {
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
      loginNameUrl.searchParams.set("requestId", requestId);

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
      const securitySettings = await getSecuritySettings({ serviceConfig, });

      const selectedSession = await findValidSession({ serviceConfig, sessions,
        authRequest,
      });

      const noSessionResponse = NextResponse.json({ error: "No active session found" }, { status: 400 });

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

      const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession.id);

      if (!cookie || !cookie.id || !cookie.token) {
        return noSessionResponse;
      }

      const session = {
        sessionId: cookie.id,
        sessionToken: cookie.token,
      };

      const { callbackUrl } = await createCallback({ serviceConfig, req: create(CreateCallbackRequestSchema, {
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
      let selectedSession = await findValidSession({ serviceConfig, sessions,
        authRequest,
      });

      if (!selectedSession || !selectedSession.id) {
        return gotoAccounts({
          request,
          requestId: `oidc_${authRequest.id}`,
          organization,
        });
      }

      const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession.id);

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
        const { callbackUrl } = await createCallback({ serviceConfig, req: create(CreateCallbackRequestSchema, {
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
          console.log("could not create callback, redirect user to choose other account");
          return gotoAccounts({
            request,
            organization,
            requestId,
          });
        }
      } catch (error) {
        console.error(error);
        return gotoAccounts({
          request,
          requestId,
          organization,
        });
      }
    }
  } else {
    const loginNameUrl = constructUrl(request, "/loginname");
    loginNameUrl.searchParams.set("requestId", requestId);

    if (authRequest?.loginHint) {
      loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
      loginNameUrl.searchParams.set("submit", "true");
    }

    if (organization) {
      loginNameUrl.searchParams.append("organization", organization);
    }

    if (suffix) {
      loginNameUrl.searchParams.append("suffix", suffix);
    }

    return NextResponse.redirect(loginNameUrl);
  }
}

/**
 * Handle SAML flow initiation
 */
export async function handleSAMLFlowInitiation(params: FlowInitiationParams): Promise<NextResponse> {
  const { serviceConfig, requestId, sessions, sessionCookies, request } = params;

  const { samlRequest } = await getSAMLRequest({ serviceConfig, samlRequestId: requestId.replace("saml_", ""),
  });

  if (!samlRequest) {
    return NextResponse.json({ error: "No samlRequest found" }, { status: 400 });
  }

  // Early return: No sessions available - redirect to login
  if (sessions.length === 0) {
    const loginNameUrl = constructUrl(request, "/loginname");
    loginNameUrl.searchParams.set("requestId", requestId);
    return NextResponse.redirect(loginNameUrl);
  }

  // Try to find a valid session
  let selectedSession = await findValidSession({ serviceConfig, sessions,
    samlRequest,
  });

  // Early return: No valid session found - show account selection
  if (!selectedSession || !selectedSession.id) {
    return gotoAccounts({
      request,
      requestId,
    });
  }

  const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession.id);

  // Early return: No valid cookie/token found - show account selection
  // Note: We need the session token from the cookie to authenticate API calls
  if (!cookie || !cookie.id || !cookie.token) {
    return gotoAccounts({
      request,
      requestId,
    });
  }

  // Valid session and cookie found - attempt to complete SAML flow
  const session = {
    sessionId: cookie.id,
    sessionToken: cookie.token,
  };

  try {
    const { url, binding } = await createResponse({ serviceConfig, req: create(CreateResponseRequestSchema, {
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
    }
  } catch (error) {
    console.error("SAML createResponse failed:", error);
  }

  // Final fallback: SAML response creation failed - show account selection
  return gotoAccounts({
    request,
    requestId,
  });
}
