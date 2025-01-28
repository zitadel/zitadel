import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SelfServiceMenu } from "@/components/self-service-menu";
import { UserAvatar } from "@/components/user-avatar";
import { getMostRecentCookieWithLoginname } from "@/lib/cookies";
import { getApiUrlOfHeaders } from "@/lib/service";
import {
  createCallback,
  getBrandingSettings,
  getLoginSettings,
  getSession,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  CreateCallbackRequestSchema,
  SessionSchema,
} from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import Link from "next/link";
import { redirect } from "next/navigation";

async function loadSession(
  host: string,
  loginName: string,
  authRequestId?: string,
) {
  const recent = await getMostRecentCookieWithLoginname({ loginName });

  if (authRequestId) {
    return createCallback({
      host,
      req: create(CreateCallbackRequestSchema, {
        authRequestId,
        callbackKind: {
          case: "session",
          value: create(SessionSchema, {
            sessionId: recent.id,
            sessionToken: recent.token,
          }),
        },
      }),
    }).then(({ callbackUrl }) => {
      return redirect(callbackUrl);
    });
  }
  return getSession({
    host,
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
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "signedin" });

  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  const { loginName, authRequestId, organization } = searchParams;
  const sessionFactors = await loadSession(host, loginName, authRequestId);

  const branding = await getBrandingSettings({ host, organization });

  let loginSettings;
  if (!authRequestId) {
    loginSettings = await getLoginSettings({ host, organization });
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          {t("title", { user: sessionFactors?.factors?.user?.displayName })}
        </h1>
        <p className="ztdl-p mb-6 block">{t("description")}</p>

        <UserAvatar
          loginName={loginName ?? sessionFactors?.factors?.user?.loginName}
          displayName={sessionFactors?.factors?.user?.displayName}
          showDropdown
          searchParams={searchParams}
        />

        {sessionFactors?.id && (
          <SelfServiceMenu sessionId={sessionFactors?.id} />
        )}

        {loginSettings?.defaultRedirectUri && (
          <div className="mt-8 flex w-full flex-row items-center">
            <span className="flex-grow"></span>

            <Link href={loginSettings?.defaultRedirectUri}>
              <Button
                type="submit"
                className="self-end"
                variant={ButtonVariants.Primary}
              >
                {t("continue")}
              </Button>
            </Link>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
