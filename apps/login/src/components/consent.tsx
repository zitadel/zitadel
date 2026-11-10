"use client";

import { approveDeviceAuthorization, denyDeviceAuthorization } from "@/lib/server/device";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Alert } from "./alert";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

export function ConsentScreen({
  scope,
  userCode,
  sessionId,
  organization,
  appName,
}: {
  scope?: string[];
  userCode: string;
  sessionId: string;
  organization?: string;
  appName?: string;
}) {
  const t = useTranslations();
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const router = useRouter();

  async function run(action: () => Promise<unknown>, onSuccess: (response: { redirect?: string }) => void) {
    setLoading(true);
    const response = (await action()
      .catch(() => ({ error: "Could not complete the request" }))
      .finally(() => {
        setLoading(false);
      })) as { error?: string; redirect?: string };

    if (response?.error) {
      setError(response.error);
      return;
    }

    onSuccess(response);
  }

  const scopes = scope?.filter((s) => !!s);

  return (
    <div className="flex w-full flex-col items-center space-y-4 pt-4">
      <ul className="w-full list-disc space-y-2">
        {scopes?.length === 0 && (
          <span className="border-divider-light bg-background-light-400 dark:bg-background-dark-400 flex w-full flex-row items-center rounded-md border px-4 py-2 text-sm transition-all">
            <Translated i18nKey="scope.openid" namespace="device" />
          </span>
        )}
        {scopes?.map((s) => {
          const translationKey = `device.scope.${s}`;
          const description = t(translationKey);

          // Check if the key itself is returned and provide a fallback
          const resolvedDescription = description === translationKey ? "" : description;

          return (
            <li
              key={s}
              className="border-divider-light bg-background-light-400 dark:bg-background-dark-400 flex w-full flex-row items-center rounded-md border px-4 py-2 text-sm transition-all"
            >
              <span>{resolvedDescription}</span>
            </li>
          );
        })}
      </ul>

      {/* RFC 8628 5.4: show the user code so the user can compare it with the one on the device */}
      <p className="ztdl-p w-full text-left text-sm">
        <Translated i18nKey="request.usercode" namespace="device" />{" "}
        <span className="font-mono font-bold" data-testid="user-code">
          {userCode}
        </span>
      </p>

      <p className="ztdl-p text-left text-xs">
        <Translated i18nKey="request.disclaimer" namespace="device" data={{ appName: appName }} />
      </p>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-4 flex w-full flex-row items-center">
        <Button
          onClick={() =>
            run(
              () => denyDeviceAuthorization({ userCode, sessionId, organization }),
              (response) => router.push(response.redirect ?? "/device"),
            )
          }
          variant={ButtonVariants.Secondary}
          data-testid="deny-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}
          <Translated i18nKey="request.deny" namespace="device" />
        </Button>
        <span className="flex-grow"></span>

        <Button
          onClick={() =>
            run(
              () => approveDeviceAuthorization({ userCode, sessionId, organization }),
              (response) => response.redirect && router.push(response.redirect),
            )
          }
          disabled={loading}
          data-testid="submit-button"
          className="self-end"
          variant={ButtonVariants.Primary}
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}
          <Translated i18nKey="request.submit" namespace="device" />
        </Button>
      </div>
    </div>
  );
}
