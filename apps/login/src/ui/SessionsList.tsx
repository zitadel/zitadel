"use client";

import SessionItem from "./SessionItem";
import Alert from "./Alert";
import { useEffect, useState } from "react";
import { Session } from "@zitadel/proto/zitadel/session/v2beta/session_pb";

type Props = {
  sessions: Session[];
  authRequestId?: string;
};

export default function SessionsList({ sessions, authRequestId }: Props) {
  const [list, setList] = useState<Session[]>(sessions);
  return sessions ? (
    <div className="flex flex-col space-y-2">
      {list
        .filter((session) => session?.factors?.user?.loginName)
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
    <Alert>No Sessions available!</Alert>
  );
}
