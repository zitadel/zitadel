"use client";

import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { useState } from "react";
import { Alert } from "./alert";
import { SessionItem } from "./session-item";
import { Translated } from "./translated";

type Props = {
  sessions: Session[];
  avatarUrls?: Record<string, string | undefined>;
  requestId?: string;
};

export function SessionsList({ sessions, requestId, avatarUrls }: Props) {
  const [list, setList] = useState<Session[]>(sessions);
  return sessions ? (
    <div className="flex flex-col space-y-2">
      {list
        .filter((session) => session?.factors?.user?.loginName)
        // sort by change date descending
        .sort((a, b) => {
          const dateA = a.changeDate ? timestampDate(a.changeDate).getTime() : 0;
          const dateB = b.changeDate ? timestampDate(b.changeDate).getTime() : 0;
          return dateB - dateA;
        })
        // TODO: add sorting to move invalid sessions to the bottom
        .map((session, index) => {
          return (
            <SessionItem
              imageUrl={avatarUrls?.[session.factors?.user?.id ?? ""]}
              session={session}
              requestId={requestId}
              reload={() => {
                setList(list.filter((s) => s.id !== session.id));
              }}
              key={"session-" + index}
            />
          );
        })}
    </div>
  ) : (
    <Alert>
      <Translated i18nKey="noResults" namespace="accounts" />
    </Alert>
  );
}
