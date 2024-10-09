import { DynamicTheme } from "@/components/dynamic-theme";
import { getBrandingSettings } from "@/lib/zitadel";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { getLocale, getTranslations } from "next-intl/server";

// This configuration shows the given name in the respective IDP button as fallback
const PROVIDER_NAME_MAPPING: {
  [provider: string]: string;
} = {
  [IdentityProviderType.GOOGLE]: "Google",
  [IdentityProviderType.GITHUB]: "GitHub",
  [IdentityProviderType.AZURE_AD]: "Microsoft",
};

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: { provider: string };
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  const { organization } = searchParams;

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("loginError.title")}</h1>
        <div>{t("loginError.description")}</div>
      </div>
    </DynamicTheme>
  );
}
