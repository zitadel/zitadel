"use client";

import { Alert, AlertType } from "@/components/alert";
import { handleServerActionResponse } from "@/lib/client-utils";
import { resendPhoneVerification, verifyPhoneAndContinue } from "@/lib/server/phone";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useRedirectLoading } from "@/lib/use-redirect-loading";
import { useCallback, useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { AutoSubmitForm } from "./auto-submit-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  code: string;
};

type Props = {
  userId: string;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
  checkAfter?: string;
  send?: boolean;
};

export function VerifyPhoneForm({ userId, loginName, sessionId, requestId, organization, checkAfter, send }: Props) {
  const router = useRouter();
  const t = useTranslations("otp");
  const initialized = useRef(false);

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onChange",
    defaultValues: {
      code: "",
    },
  });

  const [error, setError] = useState<string>("");
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);
  const { loading, setLoading, startRedirectLoading } = useRedirectLoading();

  const resendCode = useCallback(async () => {
    setError("");
    setLoading(true);

    try {
      const response = await resendPhoneVerification({ userId });
      if (response && "error" in response && response.error) {
        setError(response.error);
      }
    } catch {
      setError("Could not resend SMS code");
    } finally {
      setLoading(false);
    }
  }, [userId]);

  useEffect(() => {
    if (!initialized.current && send) {
      initialized.current = true;
      resendCode();
    }
  }, [send, resendCode]);

  async function submit(values: Inputs) {
    setError("");
    setLoading(true);

    try {
      const response = await verifyPhoneAndContinue({
        userId,
        code: values.code,
        loginName,
        sessionId,
        requestId,
        organization,
        checkAfter,
      });

      handleServerActionResponse(response, router, setSamlData, setError, undefined, startRedirectLoading);
    } catch {
      setError("Could not verify phone number");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <Alert type={AlertType.INFO}>
          <div className="flex flex-row">
            <span className="mr-auto flex-1 text-left">
              <Translated i18nKey="verify.noCodeReceived" namespace="otp" />
            </span>
            <button
              aria-label="Resend SMS Code"
              disabled={loading}
              type="button"
              className="ml-4 cursor-pointer text-primary-light-500 hover:text-primary-light-400 disabled:cursor-default disabled:text-gray-400 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 dark:disabled:text-gray-700"
              onClick={resendCode}
              data-testid="resend-button"
            >
              <Translated i18nKey="verify.resendCode" namespace="otp" />
            </button>
          </div>
        </Alert>
        <div className="mt-4">
          <TextInput
            type="text"
            autoComplete="one-time-code"
            autoFocus
            {...register("code", { required: t("verify.required.code") })}
            label={t("verify.labels.code")}
            data-testid="code-text-input"
          />
        </div>

        {error && (
          <div className="py-4" data-testid="error">
            <Alert>{error}</Alert>
          </div>
        )}

        <div className="mt-8 flex w-full flex-row items-center">
          <BackButton />
          <span className="flex-grow"></span>
          <Button
            type="submit"
            className="self-end"
            variant={ButtonVariants.Primary}
            disabled={loading || !formState.isValid}
            onClick={handleSubmit(submit)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />}
            <Translated i18nKey="verify.submit" namespace="otp" />
          </Button>
        </div>
      </form>
    </>
  );
}
