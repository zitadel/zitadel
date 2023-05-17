import { listSessions, server } from "#/lib/zitadel";
import { Avatar, AvatarSize } from "#/ui/Avatar";
import { getAllSessionIds } from "#/utils/cookies";
import {
  ExclamationTriangleIcon,
  XCircleIcon,
} from "@heroicons/react/24/outline";
import moment from "moment";
import Link from "next/link";

async function loadSessions() {
  const ids = await getAllSessionIds().catch((error) => {
    console.log("err", error);
  });

  if (ids && ids.length) {
    return listSessions(
      server,
      ids.filter((id: string | undefined) => !!id)
    ).then((sessions) => {
      console.log("ss", sessions.sessions);
      return sessions;
    });
  } else {
    return [];
  }
}

export default async function Page() {
  const { sessions } = await loadSessions();

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Accounts</h1>
      <p className="ztdl-p mb-6 block">Use your ZITADEL Account</p>

      <div className="flex flex-col w-full space-y-1">
        {sessions ? (
          sessions.map((session: any) => {
            return (
              <Link
                href={
                  `/password?` + new URLSearchParams({ session: session.id })
                }
                className="group flex flex-row items-center hover:bg-black/10 dark:hover:bg-white/10 py-3 px-4 rounded-md"
              >
                <div className="pr-4">
                  <Avatar
                    size={AvatarSize.SMALL}
                    loginName={session.factors.user.loginName}
                    name={session.factors.user.displayName}
                  />
                </div>

                <div className="flex flex-col">
                  <span className="">{session.factors.user.displayName}</span>
                  <span className="text-xs opacity-80">
                    {session.factors.user.loginName}
                  </span>
                  {session.factors.password?.verifiedAt && (
                    <span className="text-xs opacity-80">
                      {moment(
                        new Date(session.factors.password.verifiedAt)
                      ).fromNow()}
                    </span>
                  )}
                </div>

                <span className="flex-grow"></span>
                <div className="relative flex flex-row items-center">
                  {session.factors.password?.verifiedAt ? (
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
          <div className="flex flex-row items-center justify-center border border-yellow-600/40 dark:border-yellow-500/20 bg-yellow-200/30 text-yellow-600 dark:bg-yellow-700/20 dark:text-yellow-200 rounded-md py-2 scroll-px-40">
            <ExclamationTriangleIcon className="h-5 w-5 mr-2" />
            <span className="text-center text-sm">No Sessions available!</span>
          </div>
        )}
      </div>
    </div>
  );
}
