"use client";

import { clearSession } from "@/lib/server/session";
import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import moment from "moment";
import { useLocale } from "next-intl";
import { useState } from "react";
import { Alert } from "./alert";
import { Avatar } from "./avatar";
import { isSessionValid } from "./session-item";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

export function SessionClearItem({ session, reload }: { session: Session; reload: () => void }) {
  const currentLocale = useLocale();
  moment.locale(currentLocale === "zh" ? "zh-cn" : currentLocale);

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  async function clearSessionId(id: string) {
    setLoading(true);
    setError(null);

    const response = await clearSession({ sessionId: id })
      .catch((err) => {
        setError(err.message);
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    return response;
  }

  const { valid, verifiedAt } = isSessionValid(session);

  return (
    <div>
      <button
        disabled={loading}
        onClick={async () => {
          clearSessionId(session.id).then(() => {
            reload();
          });
        }}
        className="group flex w-full flex-row items-center rounded-md border border-divider-light bg-background-light-400 px-4 py-2 transition-all hover:shadow-lg dark:bg-background-dark-400 dark:hover:bg-white/10"
      >
        <div className="pr-4">
          {loading ? (
            <Spinner className="h-8 w-8" />
          ) : (
            <Avatar
              size="small"
              loginName={session.factors?.user?.loginName as string}
              name={session.factors?.user?.displayName ?? ""}
            />
          )}
        </div>

        <div className="flex flex-col items-start overflow-hidden">
          <span>{session.factors?.user?.displayName}</span>
          <span className="text-ellipsis text-xs opacity-80">{session.factors?.user?.loginName}</span>
          {valid ? (
            <span className="text-ellipsis text-xs opacity-80">
              {verifiedAt && (
                <Translated
                  i18nKey="verifiedAt"
                  namespace="logout"
                  data={{ time: moment(timestampDate(verifiedAt)).fromNow() }}
                />
              )}
            </span>
          ) : (
            verifiedAt && (
              <span className="text-ellipsis text-xs opacity-80">
                expired {session.expirationDate && moment(timestampDate(session.expirationDate)).fromNow()}
              </span>
            )
          )}
        </div>

        <span className="flex-grow"></span>
        <div className="relative flex flex-row items-center">
          <div className="mr-6 flex hidden items-center justify-center rounded-full bg-[#ff0000]/10 px-2 py-[2px] text-xs text-warn-light-500 transition-all group-hover:block dark:bg-[#ff0000]/10 dark:text-warn-dark-500">
            <Translated i18nKey="clear" namespace="logout" />
          </div>

          {valid ? (
            <div className="absolute right-0 mx-2 h-2 w-2 transform rounded-full bg-green-500 transition-all"></div>
          ) : (
            <div className="absolute right-0 mx-2 h-2 w-2 transform rounded-full bg-red-500 transition-all"></div>
          )}
        </div>
      </button>

      {error && (
        <div className="mt-1">
          <Alert>{error}</Alert>
        </div>
      )}
    </div>
  );
}
