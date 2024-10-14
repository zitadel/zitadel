import { DynamicTheme } from "@/components/dynamic-theme";
import { SelfServiceMenu } from "@/components/self-service-menu";
import { UserAvatar } from "@/components/user-avatar";
import { getMostRecentCookieWithLoginname } from "@/lib/cookies";
import { createCallback, getBrandingSettings, getSession } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  CreateCallbackRequestSchema,
  SessionSchema,
} from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { redirect } from "next/navigation";

async function loadSession(loginName: string, authRequestId?: string) {
  const recent = await getMostRecentCookieWithLoginname({ loginName });

  if (authRequestId) {
    return createCallback(
      create(CreateCallbackRequestSchema, {
        authRequestId,
        callbackKind: {
          case: "session",
          value: create(SessionSchema, {
            sessionId: recent.id,
            sessionToken: recent.token,
          }),
        },
      }),
    ).then(({ callbackUrl }) => {
      return redirect(callbackUrl);
    });
  }
  return getSession({ sessionId: recent.id, sessionToken: recent.token }).then(
    (response) => {
      if (response?.session) {
        return response.session;
      }
    },
  );
}

export default async function Page({ searchParams }: { searchParams: any }) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "signedin" });

  const { loginName, authRequestId, organization } = searchParams;
  const sessionFactors = await loadSession(loginName, authRequestId);

  const branding = await getBrandingSettings(organization);

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
      </div>
    </DynamicTheme>
  );
}
