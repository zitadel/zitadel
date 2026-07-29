import { DynamicTheme } from "@/components/dynamic-theme";
import { SessionsList } from "@/components/sessions-list";
import { Translated } from "@/components/translated";
import { getAllSessions } from "@/lib/cookies";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings, getDefaultOrg, listSessions, ServiceConfig } from "@/lib/zitadel";
import { UserPlusIcon } from "@heroicons/react/24/outline";
import { create } from "@zitadel/client";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { Session, SessionSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
// import { getLocale } from "next-intl/server";
import { headers } from "next/headers";
import Link from "next/link";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("accounts");
  return { title: t("title") };
}

async function loadSessions({ serviceConfig, organization }: { serviceConfig: ServiceConfig; organization?: string }) {
  const sessionCookies = await getAllSessions();

  if (!sessionCookies || !sessionCookies.length) {
    console.info("No session cookie found.");
    return [];
  }

  const ids = sessionCookies.map((s) => s.id).filter((id) => !!id) as string[];
  let liveSessions: Session[] = [];
  if (ids.length) {
    try {
      const response = await listSessions({ serviceConfig, ids });
      liveSessions = response?.sessions ?? [];
    } catch (error) {
      // listSessions can fail for stale/expired session IDs still in cookies,
      console.error("Failed to load sessions from API, falling back to cookie", error);
    }
  }

  // For cookie entries whose server-side session no longer exists
  // synthesize an invalid Session so the account stays selectable
  const liveIds = new Set(liveSessions.map((s) => s.id));
  const synthesized: Session[] = sessionCookies
    .filter((c) => !!c.id && !!c.loginName && !liveIds.has(c.id))
    .map((c) =>
      create(SessionSchema, {
        id: c.id,
        factors: {
          user: {
            loginName: c.loginName,
            displayName: c.loginName,
            organizationId: c.organization ?? "",
          },
        },
      }),
    );

  let sessions = [...liveSessions, ...synthesized];
  if (organization) {
    sessions = sessions.filter((s) => s.factors?.user?.organizationId === organization);
  }

  return sessions;
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const requestId = searchParams?.requestId;
  const organization = searchParams?.organization;
  const orgDomain = searchParams?.orgDomain;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg({ serviceConfig });
    if (org) {
      defaultOrganization = org.id;
    }
  }

  let sessions = await loadSessions({ serviceConfig, organization });

  const branding = await getBrandingSettings({ serviceConfig, organization: organization ?? defaultOrganization });

  const params = new URLSearchParams();

  if (requestId) {
    params.append("requestId", requestId);
  }

  if (organization) {
    params.append("organization", organization);
  }

  if (orgDomain) {
    params.append("orgDomain", orgDomain);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="accounts" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="accounts" />
        </p>
      </div>

      <div className="w-full">
        <div className="flex w-full flex-col space-y-2">
          <SessionsList sessions={sessions} requestId={requestId} />
          <Link href={`/loginname?` + params}>
            <div className="flex flex-row items-center rounded-md px-4 py-3 transition-all hover:bg-black/10 dark:hover:bg-white/10">
              <div className="mr-4 flex h-8 w-8 flex-row items-center justify-center rounded-full bg-black/5 dark:bg-white/5">
                <UserPlusIcon className="h-5 w-5" />
              </div>
              <span className="text-sm">
                <Translated i18nKey="addAnother" namespace="accounts" />
              </span>
            </div>
          </Link>
        </div>
      </div>
    </DynamicTheme>
  );
}
