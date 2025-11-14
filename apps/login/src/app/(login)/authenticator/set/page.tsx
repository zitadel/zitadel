import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { ChooseAuthenticatorToSetup } from "@/components/choose-authenticator-to-setup";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { checkUserVerification } from "@/lib/verify-helper";
import {
  getActiveIdentityProviders,
  getBrandingSettings,
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
// import { getLocale } from "next-intl/server";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { redirect } from "next/navigation";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("authenticator");
  return { title: t("title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const { loginName, requestId, organization, sessionId } = searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const sessionWithData = sessionId
    ? await loadSessionById(sessionId, organization)
    : await loadSessionByLoginname(loginName, organization);

  async function getAuthMethodsAndUser(
    serviceUrl: string,

    session?: Session,
  ) {
    const userId = session?.factors?.user?.id;

    if (!userId) {
      throw Error("Could not get user id from session");
    }

    return listAuthenticationMethodTypes({
      serviceUrl,
      userId,
    }).then((methods) => {
      return getUserByID({ serviceUrl, userId }).then((user) => {
        const humanUser = user.user?.type.case === "human" ? user.user?.type.value : undefined;

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

  async function loadSessionByLoginname(loginName?: string, organization?: string) {
    return loadMostRecentSession({
      serviceUrl,
      sessionParams: {
        loginName,
        organization,
      },
    }).then((session) => {
      return getAuthMethodsAndUser(serviceUrl, session);
    });
  }

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById({ sessionId, organization });
    return getSession({
      serviceUrl,
      sessionId: recent.id,
      sessionToken: recent.token,
    }).then((sessionResponse) => {
      return getAuthMethodsAndUser(serviceUrl, sessionResponse.session);
    });
  }

  if (!sessionWithData || !sessionWithData.factors || !sessionWithData.factors.user) {
    return (
      <Alert>
        <Translated i18nKey="unknownContext" namespace="error" />
      </Alert>
    );
  }

  const branding = await getBrandingSettings({
    serviceUrl,
    organization: sessionWithData.factors.user?.organizationId,
  });

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: sessionWithData.factors.user?.organizationId,
  });

  // check if user was verified recently
  const isUserVerified = await checkUserVerification(sessionWithData.factors.user?.id);

  if (!isUserVerified) {
    const params = new URLSearchParams({
      loginName: sessionWithData.factors.user.loginName as string,
      invite: "true",
      send: "true", // set this to true to request a new code immediately
    });

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || sessionWithData.factors.user.organizationId) {
      params.append("organization", organization ?? (sessionWithData.factors.user.organizationId as string));
    }

    redirect(`/verify?` + params);
  }

  const identityProviders = await getActiveIdentityProviders({
    serviceUrl,
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

  if (requestId) {
    params.set("requestId", requestId);
  }

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
          loginName={sessionWithData.factors?.user?.loginName}
          displayName={sessionWithData.factors?.user?.displayName}
          showDropdown
          searchParams={searchParams}
        ></UserAvatar>
      </div>

      <div className="w-full">
        {loginSettings && (
          <ChooseAuthenticatorToSetup
            authMethods={sessionWithData.authMethods}
            loginSettings={loginSettings}
            params={params}
          ></ChooseAuthenticatorToSetup>
        )}

        {loginSettings?.allowExternalIdp && !!identityProviders.length && (
          <>
            <div className="flex flex-col py-3">
              <p className="ztdl-p text-center">
                <Translated i18nKey="linkWithIDP" namespace="authenticator" />
              </p>
            </div>

            <SignInWithIdp
              showLabel={false}
              identityProviders={identityProviders}
              requestId={requestId}
              organization={sessionWithData.factors?.user?.organizationId}
              linkOnly={true} // tell the callback function to just link the IDP and not login, to get an error when user is already available
            ></SignInWithIdp>
          </>
        )}

        <div className="mt-8 flex w-full flex-row items-center">
          <BackButton />
          <span className="flex-grow"></span>
        </div>
      </div>
    </DynamicTheme>
  );
}
