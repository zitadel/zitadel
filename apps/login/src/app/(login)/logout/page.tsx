import { DynamicTheme } from "@/components/dynamic-theme";
import { SessionsClearList } from "@/components/sessions-clear-list";
import { Translated } from "@/components/translated";
import { getAllSessionCookieIds } from "@/lib/cookies";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings, getDefaultOrg, listSessions, ServiceConfig } from "@/lib/zitadel";
import { verifyJwt } from "@zitadel/client/node";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("logout");
  return { title: t("title") };
}

async function loadSessions({ serviceConfig }: { serviceConfig: ServiceConfig }) {
  const cookieIds = await getAllSessionCookieIds();

  if (cookieIds && cookieIds.length) {
    const response = await listSessions({ serviceConfig, ids: cookieIds.filter((id) => !!id) as string[] });
    return response?.sessions ?? [];
  } else {
    console.info("No session cookie found.");
    return [];
  }
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const organization = searchParams?.organization;
  const logoutToken = searchParams?.logout_token;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let postLogoutRedirectUri, logoutHint;
  if (logoutToken) {
    try {
      const payload = await verifyJwt<{ post_logout_redirect_uri?: string; logoutHint?: string }>(
        logoutToken,
        `${serviceConfig.baseUrl}/oauth/v2/keys`,
        {
          instanceHost: serviceConfig.instanceHost,
          publicHost: serviceConfig.publicHost,
        },
      );
      console.log("logout token payload", payload);

      if (payload.post_logout_redirect_uri && typeof payload.post_logout_redirect_uri === "string") {
        postLogoutRedirectUri = payload.post_logout_redirect_uri;
      }
      if (payload.logout_hint && typeof payload.logout_hint === "string") {
        logoutHint = payload.logout_hint;
      }
    } catch (error) {
      console.error("Failed to verify logout token", error);
    }
  }

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg({ serviceConfig });
    if (org) {
      defaultOrganization = org.id;
    }
  }

  let sessions = await loadSessions({ serviceConfig });

  const branding = await getBrandingSettings({ serviceConfig, organization: organization ?? defaultOrganization });

  const params = new URLSearchParams();

  if (organization) {
    params.append("organization", organization);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="logout" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="logout" />
        </p>
      </div>

      <div className="w-full">
        <div className="flex w-full flex-col space-y-2">
          <SessionsClearList
            sessions={sessions}
            logoutHint={logoutHint}
            postLogoutRedirectUri={postLogoutRedirectUri}
            organization={organization ?? defaultOrganization}
          />
        </div>
      </div>
    </DynamicTheme>
  );
}
