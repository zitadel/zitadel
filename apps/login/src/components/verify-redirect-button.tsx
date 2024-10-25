"use client";

import { sendVerificationRedirectWithoutCheck } from "@/lib/server/email";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { useState } from "react";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";

export function VerifyRedirectButton({
  userId,
  authRequestId,
  authMethods,
}: {
  userId: string;
  authRequestId: string;
  authMethods: AuthenticationMethodType[] | null;
}) {
  const t = useTranslations("verify");
  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  async function submitAndContinue(): Promise<boolean | void> {
    setLoading(true);

    await sendVerificationRedirectWithoutCheck({
      userId,
      authRequestId,
    }).catch((error) => {
      setError("Could not verify user");
      setLoading(false);
      return;
    });

    setLoading(false);
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
