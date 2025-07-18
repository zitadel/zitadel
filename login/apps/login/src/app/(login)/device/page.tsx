import { DeviceCodeForm } from "@/components/device-code-form";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getBrandingSettings } from "@/lib/zitadel";
import { getEffectiveOrganizationId } from "@/lib/organization";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;

  const userCode = searchParams?.user_code;
  const organization = searchParams?.organization;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const effectiveOrganization = await getEffectiveOrganizationId({
    serviceUrl,
    organization,
  });

  const branding = await getBrandingSettings({
    serviceUrl,
    organization: effectiveOrganization,
  });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="usercode.title" namespace="device" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="usercode.description" namespace="device" />
        </p>
        <DeviceCodeForm userCode={userCode}></DeviceCodeForm>
      </div>
    </DynamicTheme>
  );
}
