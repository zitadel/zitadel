"use client";

import Link from "next/link";
import { useState } from "react";
import { Avatar } from "./Avatar";
import moment from "moment";
import { XCircleIcon } from "@heroicons/react/24/outline";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { timestampDate } from "@zitadel/client";
import { deleteSession } from "@/lib/zitadel";
import { cleanupSession } from "@/lib/server/session";

export default function SessionItem({
  session,
  reload,
  authRequestId,
}: {
  session: Session;
  reload: () => void;
  authRequestId?: string;
}) {
  const [loading, setLoading] = useState<boolean>(false);

  async function clearSession(id: string) {
    setLoading(true);
    const response = await cleanupSession({
      sessionId: id,
    }).catch((error) => {
      setError(error.message);
    });

    setLoading(false);
    return response;
  }

  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const stillValid = session.expirationDate
    ? timestampDate(session.expirationDate) > new Date()
    : true;

  const validDate = validPassword || validPasskey;
  const validUser = (validPassword || validPasskey) && stillValid;

  const [error, setError] = useState<string | null>(null);

  return (
    <Link
      prefetch={false}
      href={
        validUser && authRequestId
          ? `/login?` +
            new URLSearchParams({
              // loginName: session.factors?.user?.loginName as string,
              sessionId: session.id,
              authRequest: authRequestId,
            })
          : !validUser
            ? `/loginname?` +
              new URLSearchParams(
                authRequestId
                  ? {
                      loginName: session.factors?.user?.loginName as string,
                      submit: "true",
                      authRequestId,
                    }
                  : {
                      loginName: session.factors?.user?.loginName as string,
                      submit: "true",
                    },
              )
            : "/signedin?" +
              new URLSearchParams(
                authRequestId
                  ? {
                      loginName: session.factors?.user?.loginName as string,
                      authRequestId,
                    }
                  : {
                      loginName: session.factors?.user?.loginName as string,
                    },
              )
      }
      className="group flex flex-row items-center bg-background-light-400 dark:bg-background-dark-400  border border-divider-light hover:shadow-lg dark:hover:bg-white/10 py-2 px-4 rounded-md transition-all"
    >
      <div className="pr-4">
        <Avatar
          size="small"
          loginName={session.factors?.user?.loginName as string}
          name={session.factors?.user?.displayName ?? ""}
        />
      </div>

      <div className="flex flex-col">
        <span className="">{session.factors?.user?.displayName}</span>
        <span className="text-xs opacity-80">
          {session.factors?.user?.loginName}
        </span>
        {validUser && (
          <span className="text-xs opacity-80">
            {validDate && moment(timestampDate(validDate)).fromNow()}
          </span>
        )}
      </div>

      <span className="flex-grow"></span>
      <div className="relative flex flex-row items-center">
        {validUser ? (
          <div className="absolute h-2 w-2 bg-green-500 rounded-full mx-2 transform right-0 group-hover:right-6 transition-all"></div>
        ) : (
          <div className="absolute h-2 w-2 bg-red-500 rounded-full mx-2 transform right-0 group-hover:right-6 transition-all"></div>
        )}

        <XCircleIcon
          className="hidden group-hover:block h-5 w-5 transition-all opacity-50 hover:opacity-100"
          onClick={(event) => {
            event.preventDefault();
            clearSession(session.id).then(() => {
              reload();
            });
          }}
        />
      </div>
    </Link>
  );
}
