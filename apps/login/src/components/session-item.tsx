"use client";

import { handleServerActionResponse } from "@/lib/client-utils";
import { sendLoginname } from "@/lib/server/loginname";
import { clearSession, continueWithSession, ContinueWithSessionCommand } from "@/lib/server/session";
import { XCircleIcon } from "@heroicons/react/24/outline";
import * as Tooltip from "@radix-ui/react-tooltip";
import { Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import moment from "moment";
import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { AutoSubmitForm } from "./auto-submit-form";
import { Avatar } from "./avatar";
import { Translated } from "./translated";

export function isSessionValid(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const validIDP = session?.factors?.intent?.verifiedAt;

  const stillValid = session.expirationDate ? timestampDate(session.expirationDate) > new Date() : true;

  const verifiedAt = validPassword || validPasskey || validIDP;
  const valid = !!((validPassword || validPasskey || validIDP) && stillValid);

  return { valid, verifiedAt };
}

export function SessionItem({ session, reload, requestId }: { session: Session; reload: () => void; requestId?: string }) {
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
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const [_error, setError] = useState<string | null>(null);

  const router = useRouter();

  return (
    <Tooltip.Root delayDuration={300}>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <Tooltip.Trigger asChild>
        <button
          onClick={async () => {
            if (valid && session?.factors?.user) {
              const sessionPayload: ContinueWithSessionCommand = session;
              if (requestId) {
                sessionPayload.requestId = requestId;
              }

              const callbackResponse = await continueWithSession(sessionPayload);

              handleServerActionResponse(callbackResponse, router, setSamlData, (e) => setError(e));
            } else if (session.factors?.user) {
              setLoading(true);
              try {
                const res = await sendLoginname({
                  loginName: session.factors?.user?.loginName,
                  organization: session.factors.user.organizationId,
                  requestId: requestId,
                });

                handleServerActionResponse(res, router, setSamlData, (e) => setError(e));
              } catch {
                setError("An internal error occurred");
              } finally {
                setLoading(false);
              }
            }
          }}
          className="group border-divider-light bg-background-light-400 dark:bg-background-dark-400 flex flex-row items-center rounded-md border px-4 py-2 transition-all hover:shadow-lg dark:hover:bg-white/10"
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
            <span className="text-xs text-ellipsis opacity-80">{session.factors?.user?.loginName}</span>
            {valid ? (
              <span className="text-xs text-ellipsis opacity-80">
                <Translated i18nKey="verified" namespace="accounts" />{" "}
                {verifiedAt && moment(timestampDate(verifiedAt)).fromNow()}
              </span>
            ) : (
              verifiedAt && (
                <span className="text-xs text-ellipsis opacity-80">
                  <Translated i18nKey="expired" namespace="accounts" />{" "}
                  {session.expirationDate && moment(timestampDate(session.expirationDate)).fromNow()}
                </span>
              )
            )}
          </div>

          <span className="flex-grow"></span>
          <div className="relative flex flex-row items-center">
            {valid ? (
              <div className="absolute right-6 mx-2 h-2 w-2 transform rounded-full bg-green-500 transition-all group-hover:right-6 sm:right-0"></div>
            ) : (
              <div className="absolute right-6 mx-2 h-2 w-2 transform rounded-full bg-red-500 transition-all group-hover:right-6 sm:right-0"></div>
            )}

            <XCircleIcon
              className="h-5 w-5 opacity-50 transition-all group-hover:block hover:opacity-100 sm:hidden"
              onClick={async (event: React.MouseEvent) => {
                event.preventDefault();
                event.stopPropagation();
                await clearSessionId(session.id);
                reload();
              }}
            />
          </div>
        </button>
      </Tooltip.Trigger>
      {valid && session.expirationDate && (
        <Tooltip.Portal>
          <Tooltip.Content
            className="bg-background-light-500 dark:bg-background-dark-500 z-50 rounded-md border px-3 py-2 text-xs text-black shadow-xl select-none dark:border-white/20 dark:text-white"
            sideOffset={5}
          >
            Expires {moment(timestampDate(session.expirationDate)).fromNow()}
            <Tooltip.Arrow className="fill-white dark:fill-white/20" />
          </Tooltip.Content>
        </Tooltip.Portal>
      )}
    </Tooltip.Root>
  );
}
