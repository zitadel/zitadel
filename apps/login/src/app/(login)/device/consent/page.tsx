import { ConsentScreen } from "@/components/consent";
import { DynamicTheme } from "@/components/dynamic-theme";
import { getServiceUrlFromHeaders } from "@/lib/service";
import {
  getBrandingSettings,
  getDefaultOrg,
  getDeviceAuthorizationRequest,
} from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "device" });

  const userCode = searchParams?.user_code;
  const requestId = searchParams?.requestId;
  const organization = searchParams?.organization;

  if (!userCode || !requestId) {
    return <div>{t("error.no_user_code")}</div>;
  }

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const { deviceAuthorizationRequest } = await getDeviceAuthorizationRequest({
    serviceUrl,
    userCode,
  });

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
    <DynamicTheme
      branding={branding}
      appName={deviceAuthorizationRequest?.appName}
    >
      <div className="flex flex-col items-center space-y-4">
        {!userCode && (
          <>
            <h1>{t("usercode.title")}</h1>
            <p className="ztdl-p">{t("usercode.description")}</p>
            <ConsentScreen scope={deviceAuthorizationRequest?.scope} />
          </>
        )}
      </div>
    </DynamicTheme>
  );
}
