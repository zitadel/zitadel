"use client";

import { handleServerActionResponse } from "@/lib/client-utils";
import { sendRecoveryCode } from "@/lib/server/recovery-code";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { AutoSubmitForm } from "./auto-submit-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

// either loginName or sessionId must be provided
type Props = {
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
};

type Inputs = {
  code: string;
};

export function LoginRecoveryCode({ loginName, sessionId, requestId, organization }: Props) {
  const t = useTranslations("recoveryCode");

  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onChange",
  });

  async function submitCodeAndContinue(values: Inputs) {
    setLoading(true);
    setError("");

    try {
      const response = await sendRecoveryCode({
        loginName,
        sessionId,
        organization,
        requestId,
        code: values.code,
      });

      if (!handleServerActionResponse(response, router, setSamlData, setError)) {
        setError(t("verify.errors.couldNotVerifyCode"));
      }
    } catch {
      setError(t("verify.errors.couldNotVerifyCode"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <div className="mt-4">
          <TextInput
            type="text"
            autoFocus
            {...register("code", { required: t("verify.required.code") })}
            label={t("verify.labels.code")}
            autoComplete="off"
            autoCapitalize="off"
            spellCheck={false}
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
            onClick={handleSubmit(submitCodeAndContinue)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />} <Translated i18nKey="verify.submit" namespace="recoveryCode" />
          </Button>
        </div>
      </form>
    </>
  );
}
