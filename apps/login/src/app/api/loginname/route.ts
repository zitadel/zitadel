import { ProviderSlug } from "@/lib/demos";
import {
  getActiveIdentityProviders,
  getLoginSettings,
  listAuthenticationMethodTypes,
  listUsers,
  PROVIDER_NAME_MAPPING,
  startIdentityProviderFlow,
} from "@/lib/zitadel";
import { createSessionForUserIdAndUpdateCookie } from "@/utils/session";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
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
      } else {
        const loginSettings = await getLoginSettings(organization);
        // TODO: check if allowDomainDiscovery has to be allowed too, to redirect to the register page
        // user not found, check if register is enabled on organization

        if (
          loginSettings?.allowRegister &&
          !loginSettings?.allowUsernamePassword
        ) {
          // TODO redirect to loginname page with idp hint
          const identityProviders = await getActiveIdentityProviders(
            organization,
          ).then((resp) => {
            return resp.identityProviders;
          });

          if (identityProviders.length === 1) {
            const host = request.nextUrl.origin;

            const provider =
              identityProviders[0].type === IdentityProviderType.GITHUB
                ? "github"
                : identityProviders[0].type === IdentityProviderType.GOOGLE
                  ? "google"
                  : identityProviders[0].type === IdentityProviderType.AZURE_AD
                    ? "azure"
                    : identityProviders[0].type === IdentityProviderType.SAML
                      ? "saml"
                      : identityProviders[0].type === IdentityProviderType.OIDC
                        ? "oidc"
                        : "oidc";

            const params = new URLSearchParams();

            if (authRequestId) {
              params.set("authRequestId", authRequestId);
            }

            if (organization) {
              params.set("organization", organization);
            }

            return startIdentityProviderFlow({
              idpId: identityProviders[0].id,
              urls: {
                successUrl:
                  `${host}/idp/${provider}/success?` +
                  new URLSearchParams(params),
                failureUrl:
                  `${host}/idp/${provider}/failure?` +
                  new URLSearchParams(params),
              },
            }).then((resp: any) => {
              if (resp.authUrl) {
                return NextResponse.json({ nextStep: resp.authUrl });
              }
            });
          } else {
            return NextResponse.json(
              { message: "Could not find user" },
              { status: 404 },
            );
          }
        } else if (
          loginSettings?.allowRegister &&
          loginSettings?.allowUsernamePassword
        ) {
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
            nextStep: registerUrl,
            status: 200,
          });
        }

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
