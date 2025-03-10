import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SelfServiceMenu } from "@/components/self-service-menu";
import { UserAvatar } from "@/components/user-avatar";
import { getMostRecentCookieWithLoginname } from "@/lib/cookies";
import { getServiceUrlFromHeaders } from "@/lib/service";
import {
  createCallback,
  createResponse,
  getBrandingSettings,
  getLoginSettings,
  getSession,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  CreateCallbackRequestSchema,
  SessionSchema,
} from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { CreateResponseRequestSchema } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import Link from "next/link";
import { redirect } from "next/navigation";

async function loadSession(
  serviceUrl: string,

  loginName: string,
  requestId?: string,
) {
  const recent = await getMostRecentCookieWithLoginname({ loginName });

  if (requestId && requestId.startsWith("oidc_")) {
    return createCallback({
      serviceUrl,
      req: create(CreateCallbackRequestSchema, {
        authRequestId: requestId,
        callbackKind: {
          case: "session",
          value: create(SessionSchema, {
            sessionId: recent.id,
            sessionToken: recent.token,
          }),
        },
      }),
    }).then(({ callbackUrl }) => {
      return redirect(callbackUrl);
    });
  } else if (requestId && requestId.startsWith("saml_")) {
    return createResponse({
      serviceUrl,
      req: create(CreateResponseRequestSchema, {
        samlRequestId: requestId.replace("saml_", ""),
        responseKind: {
          case: "session",
          value: {
            sessionId: recent.id,
            sessionToken: recent.token,
          },
        },
      }),
    }).then(({ url }) => {
      return redirect(url);
    });
  }

  return getSession({
    serviceUrl,
    sessionId: recent.id,
    sessionToken: recent.token,
  }).then((response) => {
    if (response?.session) {
      return response.session;
    }
  });
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "signedin" });

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const { loginName, requestId, organization } = searchParams;
  const sessionFactors = await loadSession(serviceUrl, loginName, requestId);

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  let loginSettings;
  if (!requestId) {
    loginSettings = await getLoginSettings({
      serviceUrl,
      organization,
    });
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          {t("title", { user: sessionFactors?.factors?.user?.displayName })}
        </h1>
        <p className="ztdl-p mb-6 block">{t("description")}</p>

        <UserAvatar
          loginName={loginName ?? sessionFactors?.factors?.user?.loginName}
          displayName={sessionFactors?.factors?.user?.displayName}
          showDropdown
          searchParams={searchParams}
        />

        {sessionFactors?.id && (
          <SelfServiceMenu sessionId={sessionFactors?.id} />
        )}

        {loginSettings?.defaultRedirectUri && (
          <div className="mt-8 flex w-full flex-row items-center">
            <span className="flex-grow"></span>

            <Link href={loginSettings?.defaultRedirectUri}>
              <Button
                type="submit"
                className="self-end"
                variant={ButtonVariants.Primary}
              >
                {t("continue")}
              </Button>
            </Link>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
