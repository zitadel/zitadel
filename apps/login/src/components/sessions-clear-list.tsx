"use client";

import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Alert, AlertType } from "./alert";
import { SessionClearItem } from "./session-clear-item";

type Props = {
  sessions: Session[];
  postLogoutRedirectUri?: string;
};

export function SessionsClearList({ sessions, postLogoutRedirectUri }: Props) {
  const t = useTranslations("logout");
  const [list, setList] = useState<Session[]>(sessions);
  const router = useRouter();

  return sessions ? (
    <div className="flex flex-col space-y-2">
      {list
        .filter((session) => session?.factors?.user?.loginName)
        // sort by change date descending
        .sort((a, b) => {
          const dateA = a.changeDate
            ? timestampDate(a.changeDate).getTime()
            : 0;
          const dateB = b.changeDate
            ? timestampDate(b.changeDate).getTime()
            : 0;
          return dateB - dateA;
        })
        // TODO: add sorting to move invalid sessions to the bottom
        .map((session, index) => {
          return (
            <SessionClearItem
              session={session}
              reload={() => {
                setList(list.filter((s) => s.id !== session.id));
                if (postLogoutRedirectUri) {
                  router.push(postLogoutRedirectUri);
                }
              }}
              key={"session-" + index}
            />
          );
        })}
      {list.length === 0 && (
        <Alert type={AlertType.INFO}>{t("noResults")}</Alert>
      )}
    </div>
  ) : (
    <Alert>{t("noResults")}</Alert>
  );
}
