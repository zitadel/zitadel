"use client";

import { registerUser, RegisterUserResponse } from "@/lib/server/register";
import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/utils/validators";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import Alert from "./Alert";
import BackButton from "./BackButton";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import PasswordComplexity from "./PasswordComplexity";
import { Spinner } from "./Spinner";

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

export default function SetPasswordForm({
  passwordComplexitySettings,
  email,
  firstname,
  lastname,
  organization,
  authRequestId,
}: Props) {
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

  const router = useRouter();

  async function submitRegister(values: Inputs) {
    setLoading(true);
    const response = await registerUser({
      email: email,
      firstName: firstname,
      lastName: lastname,
      organization: organization,
      authRequestId: authRequestId,
      password: values.password,
    }).catch(() => {
      setError("Could not register user");
    });

    if (response && "error" in response) {
      setError(response.error);
    }

    setLoading(false);

    if (!response) {
      setError("Could not register user");
      return;
    }

    const userResponse = response as RegisterUserResponse;

    const params = new URLSearchParams({ userId: userResponse.userId });

    if (userResponse.factors?.user?.loginName) {
      params.append("loginName", userResponse.factors.user.loginName);
    }
    if (organization) {
      params.append("organization", organization);
    }
    if (userResponse && userResponse.sessionId) {
      params.append("sessionId", userResponse.sessionId);
    }

    // skip verification for now as it is an app based flow
    // return router.push(`/verify?` + params);

    // check for mfa force to continue with mfa setup

    if (authRequestId && userResponse.sessionId) {
      if (authRequestId) {
        params.append("authRequest", authRequestId);
      }
      return router.push(`/login?` + params);
    } else {
      if (authRequestId) {
        params.append("authRequestId", authRequestId);
      }
      return router.push(`/signedin?` + params);
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
          continue
        </Button>
      </div>
    </form>
  );
}
