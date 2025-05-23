"use client";

import {
  sendVerificationRedirectWithoutCheck,
  SendVerificationRedirectWithoutCheckCommand,
} from "@/lib/server/verify";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";

export function VerifyRedirectButton({
  userId,
  loginName,
  requestId,
  authMethods,
  organization,
}: {
  userId?: string;
  loginName?: string;
  requestId: string;
  authMethods: AuthenticationMethodType[] | null;
  organization?: string;
}) {
  const t = useTranslations("verify");
  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();

  async function submitAndContinue(): Promise<boolean | void> {
    setLoading(true);

    let command = {
      organization,
      requestId,
    } as SendVerificationRedirectWithoutCheckCommand;

    if (userId) {
      command = {
        ...command,
        userId,
      } as SendVerificationRedirectWithoutCheckCommand;
    } else if (loginName) {
      command = {
        ...command,
        loginName,
      } as SendVerificationRedirectWithoutCheckCommand;
    }

    const response = await sendVerificationRedirectWithoutCheck(command)
      .catch(() => {
        setError("Could not verify");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    if (response && "redirect" in response && response.redirect) {
      router.push(response.redirect);
      return true;
    }
  }

  return (
    <>
      <Alert type={AlertType.INFO}>{t("success")}</Alert>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton />
        <span className="flex-grow"></span>
        {authMethods?.length === 0 && (
          <Button
            onClick={() => submitAndContinue()}
            type="submit"
            className="self-end"
            variant={ButtonVariants.Primary}
          >
            {loading && <Spinner className="h-5 w-5 mr-2" />}
            {t("setupAuthenticator")}
          </Button>
        )}
      </div>
    </>
  );
}
