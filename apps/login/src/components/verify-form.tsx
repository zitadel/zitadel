"use client";

import { Alert, AlertType } from "@/components/alert";
import { resendVerification, sendVerification } from "@/lib/server/email";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs = {
  code: string;
};

type Props = {
  userId: string;
  code?: string;
  isInvite: boolean;
  params: URLSearchParams;
};

export function VerifyForm({ userId, code, isInvite, params }: Props) {
  const t = useTranslations("verify");

  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ?? "",
    },
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  async function resendCode() {
    setError("");
    setLoading(true);

    const response = await resendVerification({
      userId,
      isInvite: isInvite,
    })
      .catch(() => {
        setError("Could not resend email");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    return response;
  }

  const fcn = useCallback(
    async function submitCodeAndContinue(
      value: Inputs,
    ): Promise<boolean | void> {
      setLoading(true);

      const response = await sendVerification({
        code: value.code,
        userId,
        isInvite: isInvite,
      })
        .catch(() => {
          setError("Could not verify user");
          return;
        })
        .finally(() => {
          setLoading(false);
        });

      if (response?.error) {
        setError(response.error);
        return;
      }

      if (response?.redirect) {
        return router.push(response?.redirect);
      }
    },
    [isInvite, userId],
  );

  useEffect(() => {
    if (code) {
      fcn({ code });
    }
  }, [code, fcn]);

  return (
    <>
      <form className="w-full">
        <Alert type={AlertType.INFO}>
          <div className="flex flex-row">
            <span className="flex-1 mr-auto text-left">
              {t("verify.noCodeReceived")}
            </span>
            <button
              aria-label="Resend OTP Code"
              disabled={loading}
              type="button"
              className="ml-4 text-primary-light-500 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 hover:text-primary-light-400 cursor-pointer disabled:cursor-default disabled:text-gray-400 dark:disabled:text-gray-700"
              onClick={() => {
                resendCode();
              }}
              data-testid="resend-button"
            >
              {t("verify.resendCode")}
            </button>
          </div>
        </Alert>
        <div className="mt-4">
          <TextInput
            type="text"
            autoComplete="one-time-code"
            {...register("code", { required: "This field is required" })}
            label="Code"
          />
        </div>

        {error && (
          <div className="py-4">
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
            onClick={handleSubmit(fcn)}
          >
            {loading && <Spinner className="h-5 w-5 mr-2" />}
            {t("verify.submit")}
          </Button>
        </div>
      </form>
    </>
  );
}
