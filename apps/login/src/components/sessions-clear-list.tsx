"use client";

import { clearSession } from "@/lib/server/session";
import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { useTranslations } from "next-intl";
import { redirect, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Alert, AlertType } from "./alert";
import { SessionClearItem } from "./session-clear-item";

type Props = {
  sessions: Session[];
  postLogoutRedirectUri?: string;
  loginHint?: string;
  organization?: string;
};

export function SessionsClearList({
  sessions,
  loginHint,
  postLogoutRedirectUri,
  organization,
}: Props) {
  const t = useTranslations("logout");
  const [list, setList] = useState<Session[]>(sessions);
  const router = useRouter();

  async function clearHintedSession() {
    // If a login hint is provided, we logout that specific session
    const sessionIdToBeCleared = sessions.find((session) => {
      return session.factors?.user?.loginName === loginHint;
    })?.id;

    if (sessionIdToBeCleared) {
      const clearSessionResponse = await clearSession({
        sessionId: sessionIdToBeCleared,
      });

      if (!clearSessionResponse) {
        console.error("Failed to clear session for login hint:", loginHint);
      }

      if (postLogoutRedirectUri) {
        return redirect(postLogoutRedirectUri);
      }

      const params = new URLSearchParams();

      if (organization) {
        params.set("organization", organization);
      }

      return router.push("/logout/success?" + params);
    } else {
      console.warn(`No session found for login hint: ${loginHint}`);
    }
  }

  useEffect(() => {
    clearHintedSession();
  }, []);

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
