"use client";

import { sendLoginname } from "@/lib/server/loginname";
import { clearSession, continueWithSession } from "@/lib/server/session";
import { XCircleIcon } from "@heroicons/react/24/outline";
import { Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import moment from "moment";
import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Avatar } from "./avatar";

export function isSessionValid(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const validIDP = session?.factors?.intent?.verifiedAt;

  const stillValid = session.expirationDate
    ? timestampDate(session.expirationDate) > new Date()
    : true;

  const verifiedAt = validPassword || validPasskey || validIDP;
  const valid = !!((validPassword || validPasskey || validIDP) && stillValid);

  return { valid, verifiedAt };
}

export function SessionItem({
  session,
  reload,
  requestId,
}: {
  session: Session;
  reload: () => void;
  requestId?: string;
}) {
  const currentLocale = useLocale();
  moment.locale(currentLocale === "zh" ? "zh-cn" : currentLocale);

  const [loading, setLoading] = useState<boolean>(false);

  async function clearSessionId(id: string) {
    setLoading(true);
    const response = await clearSession({
      sessionId: id,
    })
      .catch((error) => {
        setError(error.message);
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    return response;
  }

  const { valid, verifiedAt } = isSessionValid(session);

  const [error, setError] = useState<string | null>(null);

  const router = useRouter();

  return (
    <button
      onClick={async () => {
        if (valid && session?.factors?.user) {
          const resp = await continueWithSession({
            ...session,
            requestId: requestId,
          });

          if (resp?.redirect) {
            return router.push(resp.redirect);
          }
        } else if (session.factors?.user) {
          setLoading(true);
          const res = await sendLoginname({
            loginName: session.factors?.user?.loginName,
            organization: session.factors.user.organizationId,
            requestId: requestId,
          })
            .catch(() => {
              setError("An internal error occurred");
              return;
            })
            .finally(() => {
              setLoading(false);
            });

          if (res && "redirect" in res && res.redirect) {
            return router.push(res.redirect);
          }

          if (res && "error" in res && res.error) {
            setError(res.error);
            return;
          }
        }
      }}
      className="group flex flex-row items-center bg-background-light-400 dark:bg-background-dark-400  border border-divider-light hover:shadow-lg dark:hover:bg-white/10 py-2 px-4 rounded-md transition-all"
    >
      <div className="pr-4">
        <Avatar
          size="small"
          loginName={session.factors?.user?.loginName as string}
          name={session.factors?.user?.displayName ?? ""}
        />
      </div>

      <div className="flex flex-col items-start overflow-hidden">
        <span className="">{session.factors?.user?.displayName}</span>
        <span className="text-xs opacity-80 text-ellipsis">
          {session.factors?.user?.loginName}
        </span>
        {valid ? (
          <span className="text-xs opacity-80 text-ellipsis">
            {verifiedAt && moment(timestampDate(verifiedAt)).fromNow()}
          </span>
        ) : (
          verifiedAt && (
            <span className="text-xs opacity-80 text-ellipsis">
              expired{" "}
              {session.expirationDate &&
                moment(timestampDate(session.expirationDate)).fromNow()}
            </span>
          )
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
            event.stopPropagation();
            clearSessionId(session.id).then(() => {
              reload();
            });
          }}
        />
      </div>
    </button>
  );
}
