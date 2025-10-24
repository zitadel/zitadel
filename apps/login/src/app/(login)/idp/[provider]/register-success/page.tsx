import { DynamicTheme } from "@/components/dynamic-theme";
import { IdpSignin } from "@/components/idp-signin";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getBrandingSettings } from "@/lib/zitadel";
import { headers } from "next/headers";

/**
 * Register success page - shown when user is auto-created via IDP
 */
export default async function RegisterSuccessPage(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const searchParams = await props.searchParams;
  const { id, userId, requestId, organization } = searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  if (!userId || !id) {
    throw new Error("Missing required parameters");
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="registerSuccess.title" namespace="idp" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="registerSuccess.description" namespace="idp" />
        </p>
      </div>

      <div className="w-full">
        <IdpSignin userId={userId} idpIntent={{ idpIntentId: id, idpIntentToken: "processed" }} requestId={requestId} />
      </div>
    </DynamicTheme>
  );
}
