import { getBrandingSettings, listSessions } from "@/lib/zitadel";
import { getAllSessionCookieIds } from "@/utils/cookies";
import { UserPlusIcon } from "@heroicons/react/24/outline";
import Link from "next/link";
import SessionsList from "@/ui/SessionsList";
import DynamicTheme from "@/ui/DynamicTheme";

async function loadSessions() {
  const ids = await getAllSessionCookieIds();

  if (ids && ids.length) {
    const response = await listSessions(
      ids.filter((id: string | undefined) => !!id),
    );
    return response?.sessions ?? [];
  } else {
    console.info("No session cookie found.");
    return [];
  }
}

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const authRequestId = searchParams?.authRequestId;
  const organization = searchParams?.organization;

  let sessions = await loadSessions();

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Accounts</h1>
        <p className="ztdl-p mb-6 block">Use your ZITADEL Account</p>

        <div className="flex flex-col w-full space-y-2">
          <SessionsList sessions={sessions} authRequestId={authRequestId} />
          <Link
            href={
              authRequestId
                ? `/loginname?` +
                  new URLSearchParams({
                    authRequestId,
                  })
                : "/loginname"
            }
          >
            <div className="flex flex-row items-center py-3 px-4 hover:bg-black/10 dark:hover:bg-white/10 rounded-md transition-all">
              <div className="w-8 h-8 mr-4 flex flex-row justify-center items-center rounded-full bg-black/5 dark:bg-white/5">
                <UserPlusIcon className="h-5 w-5" />
              </div>
              <span className="text-sm">Add another account</span>
            </div>
          </Link>
        </div>
      </div>
    </DynamicTheme>
  );
}
