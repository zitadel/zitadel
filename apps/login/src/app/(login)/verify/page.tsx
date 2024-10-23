import { DynamicTheme } from "@/components/dynamic-theme";
import { VerifyEmailForm } from "@/components/verify-email-form";
import { getBrandingSettings, getLoginSettings } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({ searchParams }: { searchParams: any }) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const {
    userId,
    loginName,
    sessionId,
    code,
    organization,
    authRequestId,
    invite,
  } = searchParams;

  const branding = await getBrandingSettings(organization);

  const loginSettings = await getLoginSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <VerifyEmailForm
          userId={userId}
          loginName={loginName}
          code={code}
          organization={organization}
          authRequestId={authRequestId}
          sessionId={sessionId}
          loginSettings={loginSettings}
          isInvite={invite === "true"}
        />
      </div>
    </DynamicTheme>
  );
}
