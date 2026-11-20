"use client";

import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from "@/helpers/validators";
import { registerUser } from "@/lib/actions/register";
import { handleServerActionResponse } from "@/lib/client-utils";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { AutoSubmitForm } from "./auto-submit-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { PasswordComplexity } from "./password-complexity";
import { PrivacyPolicyCheckboxes } from "./privacy-policy-checkboxes";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  firstname: string;
  lastname: string;
  email: string;
  password: string;
  confirmPassword: string;
};

type Props = {
  legal: LegalAndSupportSettings;
  passwordComplexitySettings: PasswordComplexitySettings;
  firstname?: string;
  lastname?: string;
  email?: string;
  organization: string;
  requestId?: string;
};

export function RegisterForm({
  legal,
  passwordComplexitySettings,
  email,
  firstname,
  lastname,
  organization,
  requestId,
}: Props) {
  const { register, handleSubmit, watch, formState } = useForm<Inputs>({
    mode: "onChange",
    defaultValues: {
      email: email ?? "",
      firstname: firstname ?? "",
      lastname: lastname ?? "",
      password: "",
      confirmPassword: "",
    },
  });

  const t = useTranslations("register");

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);
  const [tosAndPolicyAccepted, setTosAndPolicyAccepted] = useState(false);

  const router = useRouter();

  async function submitRegister(values: Inputs) {
    setError("");
    setLoading(true);
    try {
      const response = await registerUser({
        email: values.email,
        firstName: values.firstname,
        lastName: values.lastname,
        password: values.password,
        organization,
        requestId,
      });

      handleServerActionResponse(response, router, setSamlData, setError);
    } catch {
      setError(t("errors.couldNotRegisterUser"));
    } finally {
      setLoading(false);
    }
  }

  const { errors } = formState;

  const watchPassword = watch("password", "");
  const watchConfirmPassword = watch("confirmPassword", "");

  const policyIsValid =
    watchPassword.length >= passwordComplexitySettings.minLength &&
    (!passwordComplexitySettings.requiresLowercase || lowerCaseValidator(watchPassword)) &&
    (!passwordComplexitySettings.requiresNumber || numberValidator(watchPassword)) &&
    (!passwordComplexitySettings.requiresUppercase || upperCaseValidator(watchPassword)) &&
    (!passwordComplexitySettings.requiresSymbol || symbolValidator(watchPassword));

  const passwordsMatch = !!watchPassword && watchPassword === watchConfirmPassword;

  const isLegalAcceptanceRequired = !!(legal?.tosLink || legal?.privacyPolicyLink);
  const canSubmit =
    formState.isValid && policyIsValid && passwordsMatch && (!isLegalAcceptanceRequired || tosAndPolicyAccepted);

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <div className="mb-4 grid grid-cols-2 gap-4">
          <div className="">
            <TextInput
              type="text"
              autoComplete="given-name"
              autoFocus
              required
              {...register("firstname", { required: t("required.firstname") })}
              label={t("labels.firstname")}
              error={errors.firstname?.message}
              data-testid="firstname-text-input"
            />
          </div>
          <div className="">
            <TextInput
              type="text"
              autoComplete="family-name"
              required
              {...register("lastname", { required: t("required.lastname") })}
              label={t("labels.lastname")}
              error={errors.lastname?.message}
              data-testid="lastname-text-input"
            />
          </div>
          <div className="col-span-2">
            <TextInput
              type="email"
              autoComplete="email"
              required
              {...register("email", { required: t("required.email") })}
              label={t("labels.email")}
              error={errors.email?.message}
              data-testid="email-text-input"
            />
          </div>
          <div className="col-span-2">
            <TextInput
              type="password"
              autoComplete="new-password"
              required
              {...register("password", { required: t("password.required.password") })}
              label={t("password.labels.password")}
              error={errors.password?.message}
              data-testid="password-text-input"
            />
          </div>
          <div className="col-span-2">
            <TextInput
              type="password"
              autoComplete="new-password"
              required
              {...register("confirmPassword", { required: t("password.required.confirmPassword") })}
              label={t("password.labels.confirmPassword")}
              error={errors.confirmPassword?.message}
              data-testid="password-confirm-text-input"
            />
          </div>
        </div>

        <PasswordComplexity
          passwordComplexitySettings={passwordComplexitySettings}
          password={watchPassword}
          equals={passwordsMatch}
        />

        {isLegalAcceptanceRequired && <PrivacyPolicyCheckboxes legal={legal} onChange={setTosAndPolicyAccepted} />}

        {error && (
          <div className="py-4" data-testid="error">
            <Alert>{error}</Alert>
          </div>
        )}

        <div className="mt-8 flex w-full flex-row items-center justify-between">
          <BackButton data-testid="back-button" />
          <Button
            type="submit"
            variant={ButtonVariants.Primary}
            disabled={loading || !canSubmit}
            onClick={handleSubmit(submitRegister)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />}
            <Translated i18nKey="submit" namespace="register" />
          </Button>
        </div>
      </form>
    </>
  );
}
