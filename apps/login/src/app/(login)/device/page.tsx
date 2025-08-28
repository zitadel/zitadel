import { DeviceCodeForm } from "@/components/device-code-form";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getBrandingSettings, getDefaultOrg } from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("device");
  return { title: t("usercode.title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const userCode = searchParams?.user_code;
  const organization = searchParams?.organization;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg({
      serviceUrl,
    });
    if (org) {
      defaultOrganization = org.id;
    }
  }

  const branding = await getBrandingSettings({
    serviceUrl,
    organization: organization ?? defaultOrganization,
  });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="usercode.title" namespace="device" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="usercode.description" namespace="device" />
        </p>
      </div>

      <div className="w-full">
        <DeviceCodeForm userCode={userCode}></DeviceCodeForm>
      </div>
    </DynamicTheme>
  );
}
