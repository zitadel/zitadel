import { ConsentScreen } from "@/components/consent";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getBrandingSettings, getDefaultOrg, getDeviceAuthorizationRequest } from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { headers } from "next/headers";

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const userCode = searchParams?.user_code;
  const requestId = searchParams?.requestId;
  const organization = searchParams?.organization;

  if (!userCode || !requestId) {
    return (
      <div>
        <Translated i18nKey="noUserCode" namespace="error" />
      </div>
    );
  }

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const { deviceAuthorizationRequest } = await getDeviceAuthorizationRequest({
    serviceUrl,
    userCode,
  });

  if (!deviceAuthorizationRequest) {
    return (
      <div>
        <Translated i18nKey="noDeviceRequest" namespace="error" />
      </div>
    );
  }

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

  const params = new URLSearchParams();

  if (requestId) {
    params.append("requestId", requestId);
  }

  if (organization) {
    params.append("organization", organization);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="request.title" namespace="device" data={{ appName: deviceAuthorizationRequest?.appName }} />
        </h1>

        <p className="ztdl-p">
          <Translated
            i18nKey="request.description"
            namespace="device"
            data={{ appName: deviceAuthorizationRequest?.appName }}
          />
        </p>
      </div>

      <div className="w-full">
        <ConsentScreen
          deviceAuthorizationRequestId={deviceAuthorizationRequest?.id}
          scope={deviceAuthorizationRequest.scope}
          appName={deviceAuthorizationRequest?.appName}
          nextUrl={`/loginname?` + params}
        />
      </div>
    </DynamicTheme>
  );
}
