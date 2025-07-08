"use client";

import { Alert, AlertType } from "@/components/alert";
import { resendVerification, sendVerification } from "@/lib/server/verify";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
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
  isInvite: boolean;
  requestId?: string;
};

export function VerifyForm({
  userId,
  loginName,
  organization,
  requestId,
  code,
  isInvite,
}: Props) {
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

    if (response && "error" in response && response?.error) {
      setError(response.error);
      return;
    }

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
        loginName: loginName,
        organization: organization,
        requestId: requestId,
      })
        .catch(() => {
          setError("Could not verify user");
          return;
        })
        .finally(() => {
          setLoading(false);
        });

      if (response && "error" in response && response?.error) {
        setError(response.error);
        return;
      }

      if (response && "redirect" in response && response?.redirect) {
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
              <Translated i18nKey="verify.noCodeReceived" namespace="verify" />
            </span>
            <button
              aria-label="Resend Code"
              disabled={loading}
              type="button"
              className="ml-4 text-primary-light-500 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 hover:text-primary-light-400 cursor-pointer disabled:cursor-default disabled:text-gray-400 dark:disabled:text-gray-700"
              onClick={() => {
                resendCode();
              }}
              data-testid="resend-button"
            >
              <Translated i18nKey="verify.resendCode" namespace="verify" />
            </button>
          </div>
        </Alert>
        <div className="mt-4">
          <TextInput
            type="text"
            autoComplete="one-time-code"
            {...register("code", { required: "This field is required" })}
            label="Code"
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
            onClick={handleSubmit(fcn)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="h-5 w-5 mr-2" />}
            <Translated i18nKey="verify.submit" namespace="verify" />
          </Button>
        </div>
      </form>
    </>
  );
}
