import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings } from "@/lib/zitadel";
import { headers } from "next/headers";

/**
 * Linking failed page - shown when IDP linking fails
 */
export default async function LinkingFailedPage(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const searchParams = await props.searchParams;
  const { organization, error } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="idp" />
        </h1>
        <p className="ztdl-p text-center">
          <Translated i18nKey="errors.linkingFailed" namespace="idp" />
        </p>
        {error && <p className="text-sm text-red-600 dark:text-red-400">{error}</p>}
      </div>
    </DynamicTheme>
  );
}
