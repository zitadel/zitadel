"use client";

import { sendLoginname } from "@/lib/server/loginname";
import { clearSession, continueWithSession } from "@/lib/server/session";
import { XCircleIcon } from "@heroicons/react/24/outline";
import * as Tooltip from "@radix-ui/react-tooltip";
import { Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import moment from "moment";
import { useLocale } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Avatar } from "./avatar";
import { Translated } from "./translated";

export function isSessionValid(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const validIDP = session?.factors?.intent?.verifiedAt;

  const stillValid = session.expirationDate
    ? timestampDate(session.expirationDate) > new Date()
    : true;

  const verifiedAt = validPassword || validPasskey || validIDP;
  const valid = !!((validPassword || validPasskey || validIDP) && stillValid);

  return { valid, verifiedAt };
}

export function SessionItem({
  session,
  reload,
  requestId,
}: {
  session: Session;
  reload: () => void;
  requestId?: string;
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

  const router = useRouter();

  return (
    <Tooltip.Root delayDuration={300}>
      <Tooltip.Trigger asChild>
        <button
          onClick={async () => {
            if (valid && session?.factors?.user) {
              await continueWithSession({
                ...session,
                requestId: requestId,
              });
            } else if (session.factors?.user) {
              setLoading(true);
              const res = await sendLoginname({
                loginName: session.factors?.user?.loginName,
                organization: session.factors.user.organizationId,
                requestId: requestId,
              })
                .catch(() => {
                  setError("An internal error occurred");
                  return;
                })
                .finally(() => {
                  setLoading(false);
                });

              if (res && "redirect" in res && res.redirect) {
                return router.push(res.redirect);
              }

              if (res && "error" in res && res.error) {
                setError(res.error);
                return;
              }
            }
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
                <Translated i18nKey="verified" namespace="accounts" />{" "}
                {verifiedAt && moment(timestampDate(verifiedAt)).fromNow()}
              </span>
            ) : (
              verifiedAt && (
                <span className="text-ellipsis text-xs opacity-80">
                  <Translated i18nKey="expired" namespace="accounts" />{" "}
                  {session.expirationDate &&
                    moment(timestampDate(session.expirationDate)).fromNow()}
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
              className="h-5 w-5 opacity-50 transition-all hover:opacity-100 group-hover:block sm:hidden"
              onClick={(event) => {
                event.preventDefault();
                event.stopPropagation();
                clearSessionId(session.id).then(() => {
                  reload();
                });
              }}
            />
          </div>
        </button>
      </Tooltip.Trigger>
      {valid && session.expirationDate && (
        <Tooltip.Portal>
          <Tooltip.Content
            className="z-50 select-none rounded-md border bg-background-light-500 px-3 py-2 text-xs text-black shadow-xl dark:border-white/20 dark:bg-background-dark-500 dark:text-white"
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
