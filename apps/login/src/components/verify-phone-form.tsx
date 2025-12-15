"use client";

import { Alert } from "@/components/alert";
import { resendPhoneVerification, sendPhoneVerification } from "@/lib/server/verify";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
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
  organization?: string;
  code?: string;
  requestId?: string;
};

export function VerifyPhoneForm({
  userId,
  loginName,
  organization,
  requestId,
  code,
}: Props) {
  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onChange",
    defaultValues: {
      code: code ?? "",
    },
  });

  const t = useTranslations("verify");

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  async function resendCode() {
    setError("");
    setLoading(true);

    const response = await resendPhoneVerification({
      userId,
    }).catch(() => {
      setError(t("errors.couldNotResendPhone"));
      setLoading(false);
      return;
    });

    if (response && "error" in response && response?.error) {
      setError(response.error);
      setLoading(false);
      return;
    }

    setLoading(false);
    return response;
  }

  const fcn = useCallback(
    async function submitCodeAndContinue(
      value: Inputs,
    ): Promise<boolean | void> {
      setLoading(true);

      const response = await sendPhoneVerification({
        code: value.code,
        userId,
        loginName: loginName,
        organization: organization,
        requestId: requestId,
      }).catch(() => {
        setError(t("errors.couldNotVerifyUser"));
        setLoading(false);
        return;
      });

      if (response && "error" in response && response?.error) {
        setError(response.error);
        setLoading(false);
        return;
      }

      if (response && "redirect" in response && response?.redirect) {
        // Keep loading state true during redirect
        return router.push(response?.redirect);
      }

      setLoading(false);
    },
    [userId, loginName, organization, requestId, t, router],
  );

  useEffect(() => {
    if (code) {
      fcn({ code });
    }
  }, [code, fcn]);

  return (
    <>
      <form className="w-full">
        <div className="mt-4">
          <TextInput
            type="text"
            autoComplete="one-time-code"
            {...register("code", { required: t("verify.required.code") })}
            label={t("verify.labels.code")}
            data-testid="code-text-input"
          />
        </div>
        <div className="w-full">  
          <div className="flex flex-row justify-end pt-1">
            <button
              aria-label="Resend Code"
              disabled={loading}
              type="button"
              className="ml-4 cursor-pointer text-primary-light-500 hover:text-primary-light-400 disabled:cursor-default disabled:text-gray-400 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 dark:disabled:text-gray-700"
              onClick={() => {
                resendCode();
              }}
              data-testid="resend-button"
            >
              <Translated i18nKey="verify.resendCode" namespace="verify" />
            </button>
          </div>
        </div>

        {error && (
          <div className="py-4" data-testid="error">
            <Alert>{error}</Alert>
          </div>
        )}

        <div className="mt-8 flex w-full flex-col items-center gap-2">
          <Button
            type="submit"
            className="self-end w-full"
            variant={ButtonVariants.Primary}
            disabled={loading || !formState.isValid}
            onClick={handleSubmit(fcn)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />} <Translated i18nKey="verify.submit" namespace="otp" />
          </Button>
          <BackButton data-testid="back-button" />
        </div>
      </form>
    </>
  );
}
