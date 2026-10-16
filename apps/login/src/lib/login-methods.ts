import { LoginSettings, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getActiveIdentityProviders, getDefaultOrg, getLoginSettings, listIDPLinks, ServiceConfig } from "./zitadel";

export async function getLoginMethods({
  serviceConfig,
  userId,
  organization,
  authMethods,
  loginSettings,
}: {
  serviceConfig: ServiceConfig;
  userId?: string;
  organization?: string;
  authMethods?: AuthenticationMethodType[];
  loginSettings: LoginSettings;
}) {
  // Match /loginname: provider availability without an explicit scope uses the default organization.
  if (!userId && !organization) {
    const defaultOrg = await getDefaultOrg({ serviceConfig });
    if (defaultOrg) {
      organization = defaultOrg.id;
      loginSettings = (await getLoginSettings({ serviceConfig, organization })) ?? loginSettings;
    }
  }
  const password =
    !!loginSettings.allowLocalAuthentication && (!authMethods || authMethods.includes(AuthenticationMethodType.PASSWORD));
  const passkey =
    !!loginSettings.allowLocalAuthentication &&
    loginSettings.passkeysType === PasskeysType.ALLOWED &&
    (!authMethods || authMethods.includes(AuthenticationMethodType.PASSKEY));
  const identityProviders = [];
  if (loginSettings.allowExternalIdp && (!authMethods || authMethods.includes(AuthenticationMethodType.IDP))) {
    const [active, links] = await Promise.all([
      getActiveIdentityProviders({ serviceConfig, orgId: organization }),
      userId ? listIDPLinks({ serviceConfig, userId }) : Promise.resolve(undefined),
    ]);
    for (const provider of active.identityProviders) {
      const link = links?.result.find((link) => link.idpId === provider.id);
      if (!userId || link) identityProviders.push({ ...provider, loginHint: link?.userName });
    }
  }
  return { password, passkey, identityProviders, count: Number(password) + Number(passkey) + identityProviders.length };
}
