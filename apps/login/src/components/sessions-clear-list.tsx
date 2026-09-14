"use client";

import { clearSession } from "@/lib/server/session";
import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { useTranslations } from "next-intl";
import { redirect, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { Alert, AlertType } from "./alert";
import { SessionClearItem } from "./session-clear-item";
import { Translated } from "./translated";

type Props = {
  sessions: Session[];
  postLogoutRedirectUri?: string;
  logoutHint?: string;
  organization?: string;
};

export function SessionsClearList({ sessions, logoutHint, postLogoutRedirectUri, organization }: Props) {
  const [list, setList] = useState<Session[]>(sessions);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const t = useTranslations("error");

  const clearHintedSession = useCallback(async () => {
    console.log("Clearing session for login hint:", logoutHint);
    // If a login hint is provided, we logout that specific session
    const sessionIdToBeCleared = sessions.find((session) => {
      return session.factors?.user?.loginName === logoutHint;
    })?.id;

    if (sessionIdToBeCleared) {
      let clearError: string | undefined;
      try {
        const clearSessionResponse = await clearSession({ sessionId: sessionIdToBeCleared });
        if (clearSessionResponse && "error" in clearSessionResponse) {
          clearError = clearSessionResponse.error;
        }
      } catch (error) {
        console.error("Error clearing session:", error);
        clearError = t("couldNotClearSession");
      }

      // Do not tell the RP the logout completed when the session was kept:
      // show the error and leave the card so the user can retry.
      if (clearError) {
        setError(clearError);
        return;
      }

      if (postLogoutRedirectUri) {
        return redirect(postLogoutRedirectUri);
      }

      const params = new URLSearchParams();

      if (organization) {
        params.set("organization", organization);
      }

      return router.push("/logout/done?" + params);
    } else {
      console.warn(`No session found for login hint: ${logoutHint}`);
    }
  }, [logoutHint, sessions, postLogoutRedirectUri, organization, router, t]);

  useEffect(() => {
    if (logoutHint) {
      clearHintedSession();
    }
  }, [logoutHint, clearHintedSession]);

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
      {error && <Alert>{error}</Alert>}
      {list.length === 0 && (
        <Alert type={AlertType.INFO}>
          <Translated i18nKey="noResults" namespace="logout" />
        </Alert>
      )}
    </div>
  ) : (
    <Alert>
      <Translated i18nKey="noResults" namespace="logout" />
    </Alert>
  );
}
