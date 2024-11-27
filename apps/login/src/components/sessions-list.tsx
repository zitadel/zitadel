"use client";

import { timestampMs } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { useTranslations } from "next-intl";
import { useState } from "react";
import { Alert } from "./alert";
import { SessionItem } from "./session-item";

type Props = {
  sessions: Session[];
  authRequestId?: string;
};

function sortFc(a: Session, b: Session) {
  if (a.changeDate && b.changeDate) {
    return timestampMs(a.changeDate) - timestampMs(b.changeDate);
  } else {
    return 0;
  }
}

export function SessionsList({ sessions, authRequestId }: Props) {
  const t = useTranslations("accounts");
  const [list, setList] = useState<Session[]>(sessions);
  return sessions ? (
    <div className="flex flex-col space-y-2">
      {list
        .filter((session) => session?.factors?.user?.loginName)
        .sort(sortFc)
        .map((session, index) => {
          return (
            <SessionItem
              session={session}
              authRequestId={authRequestId}
              reload={() => {
                setList(list.filter((s) => s.id !== session.id));
              }}
              key={"session-" + index}
            />
          );
        })}
    </div>
  ) : (
    <Alert>{t("noResults")}</Alert>
  );
}
