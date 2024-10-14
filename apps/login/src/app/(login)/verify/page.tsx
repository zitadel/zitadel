import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { VerifyEmailForm } from "@/components/verify-email-form";
import { getBrandingSettings, getLoginSettings } from "@/lib/zitadel";
import { ExclamationTriangleIcon } from "@heroicons/react/24/outline";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({ searchParams }: { searchParams: any }) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });

  const {
    userId,
    loginName,
    sessionId,
    code,
    submit,
    organization,
    authRequestId,
  } = searchParams;

  const branding = await getBrandingSettings(organization);

  const loginSettings = await getLoginSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>
        <p className="ztdl-p mb-6 block">{t("description")}</p>

        {!userId && (
          <div className="py-4">
            <Alert>{t("error:unknownContext")}</Alert>
          </div>
        )}

        {userId ? (
          <VerifyEmailForm
            userId={userId}
            loginName={loginName}
            code={code}
            submit={submit === "true"}
            organization={organization}
            authRequestId={authRequestId}
            sessionId={sessionId}
            loginSettings={loginSettings}
          />
        ) : (
          <div className="w-full flex flex-row items-center justify-center border border-yellow-600/40 dark:border-yellow-500/20 bg-yellow-200/30 text-yellow-600 dark:bg-yellow-700/20 dark:text-yellow-200 rounded-md py-2 scroll-px-40">
            <ExclamationTriangleIcon className="h-5 w-5 mr-2" />
            <span className="text-center text-sm">{t("userIdMissing")}</span>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
