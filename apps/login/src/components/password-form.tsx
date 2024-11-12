"use client";

import { resetPassword, sendPassword } from "@/lib/server/password";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useTranslations } from "next-intl";
import { redirect, useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs = {
  password: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName: string;
  organization?: string;
  authRequestId?: string;
  isAlternative?: boolean; // whether password was requested as alternative auth method
  promptPasswordless?: boolean;
};

export function PasswordForm({
  loginSettings,
  loginName,
  organization,
  authRequestId,
  promptPasswordless,
  isAlternative,
}: Props) {
  const t = useTranslations("password");

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [info, setInfo] = useState<string>("");
  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitPassword(values: Inputs) {
    setError("");
    setLoading(true);

    const response = await sendPassword({
      loginName,
      organization,
      checks: create(ChecksSchema, {
        password: { password: values.password },
      }),
      authRequestId,
      forceMfa: loginSettings?.forceMfa,
    })
      .catch(() => {
        setError("Could not verify password");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    if (response && response.nextStep) {
      return redirect(response.nextStep);
    }

    return;
  }

  async function resetPasswordAndContinue() {
    setError("");
    setInfo("");
    setLoading(true);

    const response = await resetPassword({
      loginName,
      organization,
    })
      .catch(() => {
        setError("Could not reset password");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response) {
      setError(response.error);
      return;
    }

    setInfo("Password was reset. Please check your email.");

    const params = new URLSearchParams({
      loginName: loginName,
    });

    if (organization) {
      params.append("organization", organization);
    }

    if (authRequestId) {
      params.append("authRequestId", authRequestId);
    }

    return router.push("/password/set?" + params);
  }

  return (
    <form className="w-full">
      <div className={`${error && "transform-gpu animate-shake"}`}>
        <TextInput
          type="password"
          autoComplete="password"
          {...register("password", { required: "This field is required" })}
          label="Password"
        />
        {!loginSettings?.hidePasswordReset && (
          <button
            className="transition-all text-sm hover:text-primary-light-500 dark:hover:text-primary-dark-500"
            onClick={() => resetPasswordAndContinue()}
            type="button"
            disabled={loading}
          >
            {t("verify.resetPassword")}
          </button>
        )}

        {loginName && (
          <input
            type="hidden"
            name="loginName"
            autoComplete="username"
            value={loginName}
          />
        )}
      </div>

      {info && (
        <div className="py-4">
          <Alert type={AlertType.INFO}>{info}</Alert>
        </div>
      )}

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
          onClick={handleSubmit(submitPassword)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("verify.submit")}
        </Button>
      </div>
    </form>
  );
}
