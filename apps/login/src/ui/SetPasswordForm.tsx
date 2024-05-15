"use client";

import { PasswordComplexitySettings } from "@zitadel/server";
import PasswordComplexity from "./PasswordComplexity";
import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { FieldValues, useForm } from "react-hook-form";
import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/utils/validators";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";

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
  });

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitRegister(values: Inputs) {
    setLoading(true);
    const res = await fetch("/api/registeruser", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: email,
        firstName: firstname,
        lastName: lastname,
        organization: organization,
        authRequestId: authRequestId,
        password: values.password,
      }),
    });
    setLoading(false);
    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.details);
    }
    return res.json();
  }

  function submitAndLink(value: Inputs): Promise<boolean | void> {
    return submitRegister(value)
      .then((registerResponse) => {
        setError("");

        setLoading(false);
        const params: any = { userId: registerResponse.userId };

        if (authRequestId) {
          params.authRequestId = authRequestId;
        }
        if (organization) {
          params.organization = organization;
        }
        if (registerResponse && registerResponse.sessionId) {
          params.sessionId = registerResponse.sessionId;
        }

        return router.push(`/verify?` + new URLSearchParams(params));
      })
      .catch((errorDetails: Error) => {
        setLoading(false);
        setError(errorDetails.message);
      });
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
        <Button type="button" variant={ButtonVariants.Secondary}>
          back
        </Button>
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={
            loading ||
            !policyIsValid ||
            !formState.isValid ||
            watchPassword !== watchConfirmPassword
          }
          onClick={handleSubmit(submitAndLink)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
