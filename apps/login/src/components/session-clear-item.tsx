"use client";

import { clearSession } from "@/lib/server/session";
import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import moment from "moment";
import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Avatar } from "./avatar";
import { isSessionValid } from "./session-item";
import { Translated } from "./translated";

export function SessionClearItem({
  session,
  reload,
}: {
  session: Session;
  reload: () => void;
}) {
  const currentLocale = useLocale();
  moment.locale(currentLocale === "zh" ? "zh-cn" : currentLocale);

  const [_loading, setLoading] = useState<boolean>(false);

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

  const [_error, setError] = useState<string | null>(null);

  // TODO: To we have to call this?
  useRouter();

  return (
    <button
      onClick={async () => {
        clearSessionId(session.id).then(() => {
          reload();
        });
      }}
      className="group flex flex-row items-center rounded-md border border-divider-light bg-background-light-400 px-4 py-2 transition-all hover:shadow-lg dark:bg-background-dark-400 dark:hover:bg-white/10"
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
        <span className="text-ellipsis text-xs opacity-80">
          {session.factors?.user?.loginName}
        </span>
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
              expired{" "}
              {session.expirationDate &&
                moment(timestampDate(session.expirationDate)).fromNow()}
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
  );
}
