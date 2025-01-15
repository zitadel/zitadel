import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { ChooseAuthenticatorToSetup } from "@/components/choose-authenticator-to-setup";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import {
  getActiveIdentityProviders,
  getBrandingSettings,
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "authenticator" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { loginName, authRequestId, organization, sessionId } = searchParams;

  const host = (await headers()).get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  const sessionWithData = sessionId
    ? await loadSessionById(host, sessionId, organization)
    : await loadSessionByLoginname(host, loginName, organization);

  async function getAuthMethodsAndUser(host: string, session?: Session) {
    const userId = session?.factors?.user?.id;

    if (!userId) {
      throw Error("Could not get user id from session");
    }

    return listAuthenticationMethodTypes({ host, userId }).then((methods) => {
      return getUserByID({ host, userId }).then((user) => {
        const humanUser =
          user.user?.type.case === "human" ? user.user?.type.value : undefined;

        return {
          factors: session?.factors,
          authMethods: methods.authMethodTypes ?? [],
          phoneVerified: humanUser?.phone?.isVerified ?? false,
          emailVerified: humanUser?.email?.isVerified ?? false,
          expirationDate: session?.expirationDate,
        };
      });
    });
  }

  async function loadSessionByLoginname(
    host: string,
    loginName?: string,
    organization?: string,
  ) {
    return loadMostRecentSession({
      host,
      sessionParams: {
        loginName,
        organization,
      },
    }).then((session) => {
      return getAuthMethodsAndUser(host, session);
    });
  }

  async function loadSessionById(
    host: string,
    sessionId: string,
    organization?: string,
  ) {
    const recent = await getSessionCookieById({ sessionId, organization });
    return getSession({
      host,
      sessionId: recent.id,
      sessionToken: recent.token,
    }).then((sessionResponse) => {
      return getAuthMethodsAndUser(host, sessionResponse.session);
    });
  }

  if (!sessionWithData) {
    return <Alert>{tError("unknownContext")}</Alert>;
  }

  const branding = await getBrandingSettings({
    host,
    organization: sessionWithData.factors?.user?.organizationId,
  });

  const loginSettings = await getLoginSettings({
    host,
    organization: sessionWithData.factors?.user?.organizationId,
  });

  const identityProviders = await getActiveIdentityProviders({
    host,
    orgId: sessionWithData.factors?.user?.organizationId,
    linking_allowed: true,
  }).then((resp) => {
    return resp.identityProviders;
  });

  const params = new URLSearchParams({
    initial: "true", // defines that a code is not required and is therefore not shown in the UI
  });

  if (sessionWithData.factors?.user?.loginName) {
    params.set("loginName", sessionWithData.factors?.user?.loginName);
  }

  if (sessionWithData.factors?.user?.organizationId) {
    params.set("organization", sessionWithData.factors?.user?.organizationId);
  }

  if (authRequestId) {
    params.set("authRequestId", authRequestId);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>

        <p className="ztdl-p">{t("description")}</p>

        <UserAvatar
          loginName={sessionWithData.factors?.user?.loginName}
          displayName={sessionWithData.factors?.user?.displayName}
          showDropdown
          searchParams={searchParams}
        ></UserAvatar>

        {loginSettings && (
          <ChooseAuthenticatorToSetup
            authMethods={sessionWithData.authMethods}
            loginSettings={loginSettings}
            params={params}
          ></ChooseAuthenticatorToSetup>
        )}

        <div className="py-3 flex flex-col">
          <p className="ztdl-p text-center">{t("linkWithIDP")}</p>
        </div>

        {loginSettings?.allowExternalIdp && identityProviders && (
          <SignInWithIdp
            identityProviders={identityProviders}
            authRequestId={authRequestId}
            organization={sessionWithData.factors?.user?.organizationId}
            linkOnly={true} // tell the callback function to just link the IDP and not login, to get an error when user is already available
          ></SignInWithIdp>
        )}

        <div className="mt-8 flex w-full flex-row items-center">
          <BackButton />
          <span className="flex-grow"></span>
        </div>
      </div>
    </DynamicTheme>
  );
}
