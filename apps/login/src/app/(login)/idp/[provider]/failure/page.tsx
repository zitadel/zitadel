import { ChooseAuthenticatorToLogin } from "@/components/choose-authenticator-to-login";
import { DynamicTheme } from "@/components/dynamic-theme";
import {
  getBrandingSettings,
  getLoginSettings,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  const { organization, userId } = searchParams;

  const branding = await getBrandingSettings(organization);

  const loginSettings = await getLoginSettings(organization);

  let authMethods: AuthenticationMethodType[] = [];
  if (userId) {
    const authMethodsResponse = await listAuthenticationMethodTypes(userId);
    if (authMethodsResponse.authMethodTypes) {
      authMethods = authMethodsResponse.authMethodTypes;
    }
  }

  const params = new URLSearchParams({});
  if (organization) {
    params.set("organization", organization);
  }
  if (userId) {
    params.set("userId", userId);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("loginError.title")}</h1>
        <p className="ztdl-p">{t("loginError.description")}</p>

        {userId && authMethods.length && (
          <ChooseAuthenticatorToLogin
            authMethods={sessionWithData.authMethods}
            loginSettings={loginSettings}
            params={params}
          ></ChooseAuthenticatorToLogin>
        )}
      </div>
    </DynamicTheme>
  );
}
