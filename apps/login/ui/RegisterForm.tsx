"use client";

import { PasswordComplexityPolicy, PrivacyPolicy } from "@zitadel/server";
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
  privacyPolicy: PrivacyPolicy;
  passwordComplexityPolicy: PasswordComplexityPolicy;
};

export default function RegisterForm({
  privacyPolicy,
  passwordComplexityPolicy,
}: Props) {
  const { register, handleSubmit, watch, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const router = useRouter();

  async function submitRegister(values: Inputs) {
    const res = await fetch("/registeruser", {
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
      throw new Error("Failed to register user");
    }

    return res.json();
  }

  function submitAndLink(value: Inputs): Promise<boolean | void> {
    return submitRegister(value).then((resp: any) => {
      return router.push(`/register/success?userid=${resp.userId}`);
    });
  }

  const { errors } = formState;

  const watchPassword = watch("password", "");
  const watchConfirmPassword = watch("confirmPassword", "");

  const [tosAndPolicyAccepted, setTosAndPolicyAccepted] = useState(false);

  const hasMinLength =
    passwordComplexityPolicy &&
    watchPassword?.length >= passwordComplexityPolicy.minLength;
  const hasSymbol = symbolValidator(watchPassword);
  const hasNumber = numberValidator(watchPassword);
  const hasUppercase = upperCaseValidator(watchPassword);
  const hasLowercase = lowerCaseValidator(watchPassword);

  const policyIsValid =
    passwordComplexityPolicy &&
    (passwordComplexityPolicy.hasLowercase ? hasLowercase : true) &&
    (passwordComplexityPolicy.hasNumber ? hasNumber : true) &&
    (passwordComplexityPolicy.hasUppercase ? hasUppercase : true) &&
    (passwordComplexityPolicy.hasSymbol ? hasSymbol : true) &&
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

      {passwordComplexityPolicy && (
        <PasswordComplexity
          passwordComplexityPolicy={passwordComplexityPolicy}
          password={watchPassword}
          equals={!!watchPassword && watchPassword === watchConfirmPassword}
        />
      )}

      {privacyPolicy && (
        <PrivacyPolicyCheckboxes
          privacyPolicy={privacyPolicy}
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
            !policyIsValid ||
            !formState.isValid ||
            !tosAndPolicyAccepted ||
            watchPassword !== watchConfirmPassword
          }
          onClick={handleSubmit(submitAndLink)}
        >
          continue
        </Button>
      </div>
    </form>
  );
}
