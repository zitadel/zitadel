"use client";

import { Alert, AlertType } from "@/components/alert";
import { resendVerification, sendVerification } from "@/lib/server/verify";
import { handleServerActionResponse } from "@/lib/client";
import { UNKNOWN_USER_ID } from "@/lib/constants";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";
import { AutoSubmitForm } from "./auto-submit-form";

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
  submit: boolean;
};

export function VerifyForm({ userId, loginName, organization, requestId, code, isInvite, submit }: Props) {
  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onChange",
    defaultValues: {
      code: code ?? "",
    },
  });

  const t = useTranslations("verify");

  const [error, setError] = useState<string>("");
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const [loading, setLoading] = useState<boolean>(false);

  async function resendCode() {
    setError("");
    setLoading(true);

    // do not send code for dummy userid that is set to prevent user enumeration
    if (userId === UNKNOWN_USER_ID) {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      setLoading(false);
      return;
    }

    const response = await resendVerification({
      userId,
      isInvite: isInvite,
    })
      .catch(() => {
        setError(t("errors.couldNotResendEmail"));
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

  const processedCode = useRef<string | undefined>(undefined);

  const fcn = useCallback(
    async function submitCodeAndContinue(value: Inputs): Promise<boolean | void> {
      setError("");
      setLoading(true);

      try {
        const response = await sendVerification({
          code: value.code,
          userId,
          isInvite: isInvite,
          loginName: loginName,
          organization: organization,
          requestId: requestId,
        });

        handleServerActionResponse(response, router, setSamlData, setError);
      } catch {
        setError(t("errors.couldNotVerifyUser"));
      } finally {
        setLoading(false);
      }
    },
    [isInvite, userId],
  );

  useEffect(() => {
    if (submit && code && code !== processedCode.current) {
      processedCode.current = code;
      fcn({ code });
    }
  }, [submit, code, fcn]);

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <Alert type={AlertType.INFO}>
          <div className="flex flex-row">
            <span className="mr-auto flex-1 text-left">
              <Translated i18nKey="verify.noCodeReceived" namespace="verify" />
            </span>
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
            onClick={handleSubmit(fcn)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />}
            <Translated i18nKey="verify.submit" namespace="verify" />
          </Button>
        </div>
      </form>
    </>
  );
}
