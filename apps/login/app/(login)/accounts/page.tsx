import { listSessions, server } from "#/lib/zitadel";
import Alert from "#/ui/Alert";
import { Avatar } from "#/ui/Avatar";
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
          sessions.map((session: any, index: number) => {
            const validPassword = session.factors.password?.verifiedAt;
            console.log(session);
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
                className="group flex flex-row items-center hover:bg-black/10 dark:hover:bg-white/10 py-3 px-4 rounded-md"
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
      </div>
    </div>
  );
}
