"use client";

import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/helpers/validators";
import {
  changePassword,
  resetPassword,
  sendPassword,
} from "@/lib/server/password";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { FieldValues, useForm } from "react-hook-form";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { PasswordComplexity } from "./password-complexity";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs =
  | {
      code: string;
      password: string;
      confirmPassword: string;
    }
  | FieldValues;

type Props = {
  code?: string;
  passwordComplexitySettings: PasswordComplexitySettings;
  loginName: string;
  userId: string;
  organization?: string;
  requestId?: string;
  codeRequired: boolean;
};

export function SetPasswordForm({
  passwordComplexitySettings,
  organization,
  requestId,
  loginName,
  userId,
  code,
  codeRequired,
}: Props) {
  const { register, handleSubmit, watch, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ?? "",
    },
  });

  const t = useTranslations("password");

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function resendCode() {
    setError("");
    setLoading(true);

    const response = await resetPassword({
      loginName,
      organization,
      requestId,
    })
      .catch(() => {
        setError(t("set.errors.couldNotResetPassword"));
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response) {
      setError(response.error);
      return;
    }
  }

  async function submitPassword(values: Inputs) {
    setLoading(true);
    let payload: { userId: string; password: string; code?: string } = {
      userId: userId,
      password: values.password,
    };

    // this is not required for initial password setup
    if (codeRequired) {
      payload = { ...payload, code: values.code };
    }

    const changeResponse = await changePassword(payload)
      .catch(() => {
        setError(t("set.errors.couldNotSetPassword"));
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (changeResponse && "error" in changeResponse) {
      setError(changeResponse.error);
      return;
    }

    if (!changeResponse) {
      setError(t("set.errors.couldNotSetPassword"));
      return;
    }

    const params = new URLSearchParams({});

    if (loginName) {
      params.append("loginName", loginName);
    }
    if (organization) {
      params.append("organization", organization);
    }

    await new Promise((resolve) => setTimeout(resolve, 2000)); // Wait for a second to avoid eventual consistency issues with an initial password being set

    const passwordResponse = await sendPassword({
      loginName,
      organization,
      checks: create(ChecksSchema, {
        password: { password: values.password },
      }),
      requestId,
    })
      .catch(() => {
        setError(t("set.errors.couldNotVerifyPassword"));
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (
      passwordResponse &&
      "error" in passwordResponse &&
      passwordResponse.error
    ) {
      setError(passwordResponse.error);
      return;
    }

    if (
      passwordResponse &&
      "redirect" in passwordResponse &&
      passwordResponse.redirect
    ) {
      return router.push(passwordResponse.redirect);
    }

    return;
  }

  const { errors } = formState;

  const watchPassword = watch("password", "");
  const watchConfirmPassword = watch("confirmPassword", "");

  const hasMinLength =
    passwordComplexitySettings &&
    watchPassword?.length >= passwordComplexitySettings.minLength;
  const hasSymbol = symbolValidator(watchPassword);
  const hasNumber = numberValidator(watchPassword);
  const hasUppercase = upperCaseValidator(watchPassword);
  const hasLowercase = lowerCaseValidator(watchPassword);

  const policyIsValid =
    passwordComplexitySettings &&
    (passwordComplexitySettings.requiresLowercase ? hasLowercase : true) &&
    (passwordComplexitySettings.requiresNumber ? hasNumber : true) &&
    (passwordComplexitySettings.requiresUppercase ? hasUppercase : true) &&
    (passwordComplexitySettings.requiresSymbol ? hasSymbol : true) &&
    hasMinLength;

  return (
    <form className="w-full">
      <div className="mb-4 grid grid-cols-1 gap-4 pt-4">
        {codeRequired && (
          <Alert type={AlertType.INFO}>
            <div className="flex flex-row">
              <span className="mr-auto flex-1 text-left">
                <Translated i18nKey="set.noCodeReceived" namespace="password" />
              </span>
              <button
                aria-label="Resend OTP Code"
                disabled={loading}
                type="button"
                className="ml-4 cursor-pointer text-primary-light-500 hover:text-primary-light-400 disabled:cursor-default disabled:text-gray-400 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 dark:disabled:text-gray-700"
                onClick={() => {
                  resendCode();
                }}
                data-testid="resend-button"
              >
                <Translated i18nKey="set.resend" namespace="password" />
              </button>
            </div>
          </Alert>
        )}
        {codeRequired && (
          <div>
            <TextInput
              type="text"
              required
              {...register("code", {
                required: t("set.required.code"),
              })}
              label={t("set.labels.code")}
              autoComplete="one-time-code"
              error={errors.code?.message as string}
              data-testid="code-text-input"
            />
          </div>
        )}
        <div>
          <TextInput
            type="password"
            autoComplete="new-password"
            required
            {...register("password", {
              required: t("set.required.newPassword"),
            })}
            label={t("set.labels.newPassword")}
            error={errors.password?.message as string}
            data-testid="password-set-text-input"
          />
        </div>
        <div>
          <TextInput
            type="password"
            required
            autoComplete="new-password"
            {...register("confirmPassword", {
              required: t("set.required.confirmPassword"),
            })}
            label={t("set.labels.confirmPassword")}
            error={errors.confirmPassword?.message as string}
            data-testid="password-set-confirm-text-input"
          />
        </div>
      </div>

      {passwordComplexitySettings && (
        <PasswordComplexity
          passwordComplexitySettings={passwordComplexitySettings}
          password={watchPassword}
          equals={!!watchPassword && watchPassword === watchConfirmPassword}
        />
      )}

      {error && <Alert>{error}</Alert>}

      <div className="mt-8 flex w-full flex-row items-center justify-between">
        <BackButton data-testid="back-button" />
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={
            loading ||
            !policyIsValid ||
            !formState.isValid ||
            watchPassword !== watchConfirmPassword
          }
          onClick={handleSubmit(submitPassword)}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}{" "}
          <Translated i18nKey="set.submit" namespace="password" />
        </Button>
      </div>
    </form>
  );
}
