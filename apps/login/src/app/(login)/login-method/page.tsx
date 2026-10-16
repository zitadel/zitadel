import { PASSKEYS, PASSWORD } from "@/components/auth-methods";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getLoginMethods } from "@/lib/login-methods";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getLoginSettings, listAuthenticationMethodTypes } from "@/lib/zitadel";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { redirect } from "next/navigation";

export async function generateMetadata() {
  const t = await getTranslations("authenticator");
  return { title: t("title") };
}

export default async function Page(props: { searchParams: Promise<Record<string, string | undefined>> }) {
  const { loginName, organization, requestId } = await props.searchParams;
  const { serviceConfig } = getServiceConfig(await headers());
  const params = new URLSearchParams();
  if (loginName) params.set("loginName", loginName);
  if (organization) params.set("organization", organization);
  if (requestId) params.set("requestId", requestId);

  if (!loginName) redirect("/loginname?" + params);
  const contextSettings = await getLoginSettings({ serviceConfig, organization });
  if (!contextSettings) redirect("/loginname?" + params);

  let user: { id?: string; loginName?: string; displayName?: string; organizationId?: string } = {
    loginName,
    organizationId: organization,
  };
  let loginSettings = contextSettings;
  let authMethods;
  // Protected flows offer policy-level choices without looking up whether this user exists.
  if (!contextSettings.ignoreUnknownUsernames) {
    const session = await loadMostRecentSession({ serviceConfig, sessionParams: { loginName, organization } });
    const sessionUser = session?.factors?.user;
    if (
      !sessionUser?.id ||
      sessionUser.loginName !== loginName ||
      (organization && organization !== sessionUser.organizationId)
    )
      redirect("/loginname?" + params);
    user = sessionUser;
    const settings = await getLoginSettings({ serviceConfig, organization: user.organizationId });
    if (!settings || settings.ignoreUnknownUsernames) redirect("/loginname?" + params);
    loginSettings = settings;
    authMethods = (await listAuthenticationMethodTypes({ serviceConfig, userId: sessionUser.id })).authMethodTypes;
  }
  if (user.organizationId) params.set("organization", user.organizationId);
  const options = await getLoginMethods({
    serviceConfig,
    userId: user.id,
    organization: user.organizationId,
    authMethods,
    loginSettings,
  });
  if (!options.count) redirect("/loginname?" + params);
  const branding = await getBrandingSettings({ serviceConfig, organization: user.organizationId });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="authenticator" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="authenticator" />
        </p>
        <UserAvatar
          loginName={user.loginName}
          displayName={user.displayName}
          showDropdown
          searchParams={{ loginName, organization: user.organizationId, requestId }}
        />
      </div>
      <div className="flex w-full flex-col gap-5 pt-6">
        {options.password && PASSWORD(false, "/password?" + params)}
        {options.passkey && PASSKEYS(false, "/passkey?" + params)}
        {options.identityProviders.map((provider) => (
          <SignInWithIdp
            key={provider.id}
            identityProviders={[provider]}
            requestId={requestId}
            organization={user.organizationId}
            loginHint={provider.loginHint}
            postErrorRedirectUrl="/loginname"
            showLabel={false}
          />
        ))}
      </div>
    </DynamicTheme>
  );
}
