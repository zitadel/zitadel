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
        clearSessionId(session.id).then(() => {
          reload();
        });
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
            {verifiedAt && (
              <Translated
                i18nKey="verfiedAt"
                namespace="logout"
                data={{ time: moment(timestampDate(verifiedAt)).fromNow() }}
              />
            )}
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
        <div className="mr-6 px-2 py-[2px] text-xs hidden group-hover:block transition-all text-warn-light-500 dark:text-warn-dark-500 bg-[#ff0000]/10 dark:bg-[#ff0000]/10 rounded-full flex items-center justify-center">
          <Translated i18nKey="clear" namespace="logout" />
        </div>

        {valid ? (
          <div className="absolute h-2 w-2 bg-green-500 rounded-full mx-2 transform right-0 transition-all"></div>
        ) : (
          <div className="absolute h-2 w-2 bg-red-500 rounded-full mx-2 transform right-0 transition-all"></div>
        )}
      </div>
    </button>
  );
}
