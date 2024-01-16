"use client";

import {
  LegalAndSupportSettings,
  PasswordComplexitySettings,
} from "@zitadel/server";
import PasswordComplexity from "./PasswordComplexity";
import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { PrivacyPolicyCheckboxes } from "./PrivacyPolicyCheckboxes";
import { FieldValues, useForm } from "react-hook-form";
import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "#/utils/validators";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
      password: string;
      confirmPassword: string;
    }
  | FieldValues;

type Props = {
  legal: LegalAndSupportSettings;
  passwordComplexitySettings: PasswordComplexitySettings;
};

export default function RegisterForm({
  legal,
  passwordComplexitySettings,
}: Props) {
  const { register, handleSubmit, watch, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitRegister(values: Inputs) {
    setLoading(true);
    const res = await fetch("/api/registeruser", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: values.email,
        password: values.password,
        firstName: values.firstname,
        lastName: values.lastname,
      }),
    });

    if (!res.ok) {
      setLoading(false);
      throw new Error("Failed to register user");
    }

    setLoading(false);
    return res.json();
  }

  function submitAndLink(value: Inputs): Promise<boolean | void> {
    return submitRegister(value).then((resp: any) => {
      return router.push(`/verify?userID=${resp.userId}`);
    });
  }

  const { errors } = formState;

  const watchPassword = watch("password", "");
  const watchConfirmPassword = watch("confirmPassword", "");

  const [tosAndPolicyAccepted, setTosAndPolicyAccepted] = useState(false);

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
      <div className="grid grid-cols-2 gap-4 mb-4">
        <div className="">
          <TextInput
            type="firstname"
            autoComplete="firstname"
            required
            {...register("firstname", { required: "This field is required" })}
            label="First name"
            error={errors.firstname?.message as string}
          />
        </div>
        <div className="">
          <TextInput
            type="lastname"
            autoComplete="lastname"
            required
            {...register("lastname", { required: "This field is required" })}
            label="Last name"
            error={errors.lastname?.message as string}
          />
        </div>
        <div className="col-span-2">
          <TextInput
            type="email"
            autoComplete="email"
            required
            {...register("email", { required: "This field is required" })}
            label="E-mail"
            error={errors.email?.message as string}
          />
        </div>
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

      {legal && (
        <PrivacyPolicyCheckboxes
          legal={legal}
          onChange={setTosAndPolicyAccepted}
        />
      )}

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
            !tosAndPolicyAccepted ||
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
