import { Session } from "#/../../packages/zitadel-server/dist";
import { listSessions, server } from "#/lib/zitadel";
import Alert from "#/ui/Alert";
import { Avatar } from "#/ui/Avatar";
import { getAllSessionIds } from "#/utils/cookies";
import { UserPlusIcon, XCircleIcon } from "@heroicons/react/24/outline";
import moment from "moment";
import Link from "next/link";

async function loadSessions(): Promise<Session[]> {
  const ids = await getAllSessionIds();

  if (ids && ids.length) {
    const response = await listSessions(
      server,
      ids.filter((id: string | undefined) => !!id)
    );
    return response?.sessions ?? [];
  } else {
    return [];
  }
}

export default async function Page() {
  const sessions = await loadSessions();

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Accounts</h1>
      <p className="ztdl-p mb-6 block">Use your ZITADEL Account</p>

      <div className="flex flex-col w-full space-y-2">
        {sessions ? (
          sessions.map((session: any, index: number) => {
            const validPassword = session.factors.password?.verifiedAt;
            return (
              <Link
                key={"session-" + index}
                href={
                  validPassword
                    ? `/signedin?` +
                      new URLSearchParams({
                        loginName: session.factors.user.loginName,
                      })
                    : `/password?` +
                      new URLSearchParams({
                        loginName: session.factors.user.loginName,
                      })
                }
                className="group flex flex-row items-center bg-background-light-400 dark:bg-background-dark-400  border border-divider-light hover:shadow-lg dark:hover:bg-white/10 py-2 px-4 rounded-md transition-all"
              >
                <div className="pr-4">
                  <Avatar
                    size="small"
                    loginName={session.factors.user.loginName}
                    name={session.factors.user.displayName}
                  />
                </div>

                <div className="flex flex-col">
                  <span className="">{session.factors.user.displayName}</span>
                  <span className="text-xs opacity-80">
                    {session.factors.user.loginName}
                  </span>
                  {validPassword && (
                    <span className="text-xs opacity-80">
                      {moment(new Date(validPassword)).fromNow()}
                    </span>
                  )}
                </div>

                <span className="flex-grow"></span>
                <div className="relative flex flex-row items-center">
                  {validPassword ? (
                    <div className="absolute h-2 w-2 bg-green-500 rounded-full mx-2 transform right-0 group-hover:right-6 transition-all"></div>
                  ) : (
                    <div className="absolute h-2 w-2 bg-red-500 rounded-full mx-2 transform right-0 group-hover:right-6 transition-all"></div>
                  )}

                  <XCircleIcon className="hidden group-hover:block h-5 w-5 transition-all opacity-50 hover:opacity-100" />
                </div>
              </Link>
            );
          })
        ) : (
          <Alert>No Sessions available!</Alert>
        )}
        <Link href="/username">
          <div className="flex flex-row items-center py-3 px-4 hover:bg-black/10 dark:hover:bg-white/10 rounded-md transition-all">
            <div className="w-8 h-8 mr-4 flex flex-row justify-center items-center rounded-full bg-black/5 dark:bg-white/5">
              <UserPlusIcon className="h-5 w-5" />
            </div>
            <span className="text-sm">Add another account</span>
          </div>
        </Link>
      </div>
    </div>
  );
}
