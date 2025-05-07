import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterPasskey } from "@/components/register-passkey";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "passkey" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { loginName, prompt, organization, requestId, userId } = searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const session = await loadMostRecentSession({
    serviceUrl,
    sessionParams: {
      loginName,
      organization,
    },
  });

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("set.title")}</h1>

        {session && (
          <UserAvatar
            loginName={loginName ?? session.factors?.user?.loginName}
            displayName={session.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}
        <p className="ztdl-p mb-6 block">{t("set.description")}</p>

        <Alert type={AlertType.INFO}>
          <span>
            {t("set.info.description")}
            <a
              className="text-primary-light-500 dark:text-primary-dark-500 hover:text-primary-light-300 hover:dark:text-primary-dark-300"
              target="_blank"
              href="https://zitadel.com/docs/guides/manage/user/reg-create-user#with-passwordless"
            >
              {t("set.info.link")}
            </a>
          </span>
        </Alert>

        {!session && (
          <div className="py-4">
            <Alert>{tError("unknownContext")}</Alert>
          </div>
        )}

        {session?.id && (
          <RegisterPasskey
            sessionId={session.id}
            isPrompt={!!prompt}
            organization={organization}
            requestId={requestId}
          />
        )}
      </div>
    </DynamicTheme>
  );
}
