import { Alert, AlertType } from "@/components/alert";
import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getMostRecentCookieWithLoginname, getSessionCookieById } from "@/lib/cookies";
import { completeDeviceAuthorization } from "@/lib/server/device";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getLoginSettings, getSession } from "@/lib/zitadel";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import Link from "next/link";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("signedin");
  return { title: t("title", { user: "" }) };
}

async function loadSessionById(serviceUrl: string, sessionId: string, organization?: string) {
  const recent = await getSessionCookieById({ sessionId, organization });
  return getSession({
    serviceUrl,
    sessionId: recent.id,
    sessionToken: recent.token,
  }).then((response) => {
    if (response?.session) {
      return response.session;
    }
  });
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const { loginName, requestId, organization, sessionId } = searchParams;

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  // complete device authorization flow if device requestId is present
  if (requestId && requestId.startsWith("device_")) {
    const cookie = sessionId
      ? await getSessionCookieById({ sessionId, organization })
      : await getMostRecentCookieWithLoginname({
          loginName: loginName,
          organization: organization,
        });

    await completeDeviceAuthorization(requestId.replace("device_", ""), {
      sessionId: cookie.id,
      sessionToken: cookie.token,
    }).catch((err) => {
      return (
        <DynamicTheme branding={branding}>
          <div className="flex flex-col space-y-4">
            <h1>
              <Translated i18nKey="error.title" namespace="signedin" />
            </h1>
            <p className="ztdl-p mb-6 block">
              <Translated i18nKey="error.description" namespace="signedin" />
            </p>
            <Alert>{err.message}</Alert>
          </div>
          <div className="w-full"></div>
        </DynamicTheme>
      );
    });
  }

  const sessionFactors = sessionId
    ? await loadSessionById(serviceUrl, sessionId, organization)
    : await loadMostRecentSession({
        serviceUrl,
        sessionParams: { loginName, organization },
      });

  let loginSettings;
  if (!requestId) {
    loginSettings = await getLoginSettings({
      serviceUrl,
      organization,
    });
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="signedin" data={{ user: sessionFactors?.factors?.user?.displayName }} />
        </h1>
        <p className="ztdl-p mb-6 block">
          <Translated i18nKey="description" namespace="signedin" />
        </p>

        <UserAvatar
          loginName={loginName ?? sessionFactors?.factors?.user?.loginName}
          displayName={sessionFactors?.factors?.user?.displayName}
          showDropdown={!(requestId && requestId.startsWith("device_"))}
          searchParams={searchParams}
        />
      </div>

      <div className="w-full">
        {requestId && requestId.startsWith("device_") && (
          <Alert type={AlertType.INFO}>
            You can now close this window and return to the device where you started the authorization process to continue.
          </Alert>
        )}

        {loginSettings?.defaultRedirectUri && (
          <div className="mt-8 flex w-full flex-row items-center">
            <span className="flex-grow"></span>

            <Link href={loginSettings?.defaultRedirectUri}>
              <Button type="submit" className="self-end" variant={ButtonVariants.Primary}>
                <Translated i18nKey="continue" namespace="signedin" />
              </Button>
            </Link>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
