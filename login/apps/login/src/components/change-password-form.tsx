"use client";

import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/helpers/validators";
import {
  checkSessionAndSetPassword,
  sendPassword,
} from "@/lib/server/password";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { PasswordComplexity } from "./password-complexity";
import { Spinner } from "./spinner";

type Inputs =
  | {
      password: string;
      confirmPassword: string;
    }
  | FieldValues;

type Props = {
  passwordComplexitySettings: PasswordComplexitySettings;
  sessionId: string;
  loginName: string;
  requestId?: string;
  organization?: string;
};

export function ChangePasswordForm({
  passwordComplexitySettings,
  sessionId,
  loginName,
  requestId,
  organization,
}: Props) {
  const t = useTranslations("password");
  const router = useRouter();

  const { register, handleSubmit, watch, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      password: "",
      comfirmPassword: "",
    },
  });

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  async function submitChange(values: Inputs) {
    setLoading(true);

    const changeResponse = checkSessionAndSetPassword({
      sessionId,
      password: values.password,
    })
      .catch(() => {
        setError("Could not change password");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (changeResponse && "error" in changeResponse && changeResponse.error) {
      setError(
        typeof changeResponse.error === "string"
          ? changeResponse.error
          : "Unknown error",
      );
      return;
    }

    if (!changeResponse) {
      setError("Could not change password");
      return;
    }

    await new Promise((resolve) => setTimeout(resolve, 1000)); // wait for a second, to prevent eventual consistency issues

    const passwordResponse = await sendPassword({
      loginName,
      organization,
      checks: create(ChecksSchema, {
        password: { password: values.password },
      }),
      requestId,
    })
      .catch(() => {
        setError("Could not verify password");
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
      <div className="pt-4 grid grid-cols-1 gap-4 mb-4">
        <div className="">
          <TextInput
            type="password"
            autoComplete="new-password"
            required
            {...register("password", {
              required: "You have to provide a new password!",
            })}
            label="New Password"
            error={errors.password?.message as string}
            data-testid="password-change-text-input"
          />
        </div>
        <div className="">
          <TextInput
            type="password"
            required
            autoComplete="new-password"
            {...register("confirmPassword", {
              required: "This field is required",
            })}
            label="Confirm Password"
            error={errors.confirmPassword?.message as string}
            data-testid="password-change-confirm-text-input"
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
          onClick={handleSubmit(submitChange)}
          data-testid="submit-button"
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("change.submit")}
        </Button>
      </div>
    </form>
  );
}
