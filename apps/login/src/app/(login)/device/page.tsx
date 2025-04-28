import { DeviceCodeForm } from "@/components/device-code-form";
import { DynamicTheme } from "@/components/dynamic-theme";
import { getServiceUrlFromHeaders } from "@/lib/service";
import {
  getBrandingSettings,
  getDefaultOrg,
  getDeviceAuthorizationRequest,
} from "@/lib/zitadel";
import { DeviceAuthorizationRequest } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "device" });

  const loginName = searchParams?.loginName;
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

  let deviceAuthRequest: DeviceAuthorizationRequest | null = null;
  if (userCode) {
    const deviceAuthorizationRequestResponse =
      await getDeviceAuthorizationRequest({
        serviceUrl,
        userCode,
      });

    if (deviceAuthorizationRequestResponse.deviceAuthorizationRequest) {
      deviceAuthRequest =
        deviceAuthorizationRequestResponse.deviceAuthorizationRequest;
    }
  }

  const branding = await getBrandingSettings({
    serviceUrl,
    organization: organization ?? defaultOrganization,
  });

  return (
    <DynamicTheme branding={branding} appName={deviceAuthRequest?.appName}>
      <div className="flex flex-col items-center space-y-4">
        {!userCode && (
          <>
            <h1>{t("usercode.title")}</h1>
            <p className="ztdl-p">{t("usercode.description")}</p>
            <DeviceCodeForm
            // loginSettings={contextLoginSettings}
            ></DeviceCodeForm>
          </>
        )}

        {deviceAuthRequest && (
          <div>
            <h1>
              {deviceAuthRequest.appName}
              <br />
              {t("request.title")}
            </h1>
            <p className="ztdl-p text-left text-xs mt-4">
              {t("request.description")}
            </p>
            {/* {JSON.stringify(deviceAuthRequest)} */}
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
