import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { getLoginPageMetadata, type LoginPageMetadataProps } from "@/lib/login-metadata";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings } from "@/lib/zitadel";
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(props: LoginPageMetadataProps): Promise<Metadata> {
  const t = await getTranslations("common");
  return getLoginPageMetadata(props, t("title"));
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const { organization } = searchParams;

  const branding = await getBrandingSettings({ serviceConfig, organization });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="success.title" namespace="logout" />
        </h1>
        <p className="ztdl-p mb-6 block">
          <Translated i18nKey="success.description" namespace="logout" />
        </p>
      </div>
      <div className="w-full"></div>
    </DynamicTheme>
  );
}
