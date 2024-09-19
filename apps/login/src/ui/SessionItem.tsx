"use client";

import { cleanupSession } from "@/lib/server/session";
import { XCircleIcon } from "@heroicons/react/24/outline";
import { Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import moment from "moment";
import Link from "next/link";
import { useState } from "react";
import { Avatar } from "./Avatar";

export function isSessionValid(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const stillValid = session.expirationDate
    ? timestampDate(session.expirationDate) > new Date()
    : true;

  const verifiedAt = validPassword || validPasskey;
  const valid = !!((validPassword || validPasskey) && stillValid);

  return { valid, verifiedAt };
}

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

  const { valid, verifiedAt } = isSessionValid(session);

  const [error, setError] = useState<string | null>(null);

  return (
    <Link
      prefetch={false}
      href={
        valid && authRequestId
          ? `/login?` +
            new URLSearchParams({
              // loginName: session.factors?.user?.loginName as string,
              sessionId: session.id,
              authRequest: authRequestId,
            })
          : !valid
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

      <div className="flex flex-col overflow-hidden">
        <span className="">{session.factors?.user?.displayName}</span>
        <span className="text-xs opacity-80 text-ellipsis">
          {session.factors?.user?.loginName}
        </span>
        {valid && (
          <span className="text-xs opacity-80 text-ellipsis">
            {verifiedAt && moment(timestampDate(verifiedAt)).fromNow()}
          </span>
        )}
      </div>

      <span className="flex-grow"></span>
      <div className="relative flex flex-row items-center">
        {valid ? (
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
