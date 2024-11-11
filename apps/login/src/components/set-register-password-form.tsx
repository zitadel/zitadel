"use client";

import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/helpers/validators";
import { registerUser } from "@/lib/server/register";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { useTranslations } from "next-intl";
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
  email: string;
  firstname: string;
  lastname: string;
  organization?: string;
  authRequestId?: string;
};

export function SetRegisterPasswordForm({
  passwordComplexitySettings,
  email,
  firstname,
  lastname,
  organization,
  authRequestId,
}: Props) {
  const t = useTranslations("register");

  const { register, handleSubmit, watch, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      email: email ?? "",
      firstname: firstname ?? "",
      lastname: lastname ?? "",
    },
  });

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  async function submitRegister(values: Inputs) {
    setLoading(true);
    const response = await registerUser({
      email: email,
      firstName: firstname,
      lastName: lastname,
      organization: organization,
      authRequestId: authRequestId,
      password: values.password,
    })
      .catch(() => {
        setError("Could not register user");
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response) {
      setError(response.error);
      return;
    }
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
              required: "You have to provide a password!",
            })}
            label="Password"
            error={errors.password?.message as string}
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
        <BackButton />
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={
            loading ||
            !policyIsValid ||
            !formState.isValid ||
            watchPassword !== watchConfirmPassword
          }
          onClick={handleSubmit(submitRegister)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("password.submit")}
        </Button>
      </div>
    </form>
  );
}
