"use client";

import { registerUser } from "@/lib/server/register";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import { Alert } from "./alert";
import {
  AuthenticationMethod,
  AuthenticationMethodRadio,
  methods,
} from "./authentication-method-radio";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { PrivacyPolicyCheckboxes } from "./privacy-policy-checkboxes";
import { Spinner } from "./spinner";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
    }
  | FieldValues;

type Props = {
  legal: LegalAndSupportSettings;
  firstname?: string;
  lastname?: string;
  email?: string;
  organization?: string;
  requestId?: string;
  loginSettings?: LoginSettings;
};

export function RegisterForm({
  legal,
  email,
  firstname,
  lastname,
  organization,
  requestId,
  loginSettings,
}: Props) {
  const t = useTranslations("register");

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      email: email ?? "",
      firstName: firstname ?? "",
      lastname: lastname ?? "",
    },
  });

  const [loading, setLoading] = useState<boolean>(false);
  const [selected, setSelected] = useState<AuthenticationMethod>(methods[0]);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitAndRegister(values: Inputs) {
    setLoading(true);
    const response = await registerUser({
      email: values.email,
      firstName: values.firstname,
      lastName: values.lastname,
      organization: organization,
      requestId: requestId,
    })
      .catch(() => {
        setError("Could not register user");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    if (response && "redirect" in response && response.redirect) {
      return router.push(response.redirect);
    }

    return response;
  }

  async function submitAndContinue(
    value: Inputs,
    withPassword: boolean = false,
  ) {
    const registerParams: any = value;

    if (organization) {
      registerParams.organization = organization;
    }

    if (requestId) {
      registerParams.requestId = requestId;
    }

    // redirect user to /register/password if password is chosen
    if (withPassword) {
      return router.push(
        `/register/password?` + new URLSearchParams(registerParams),
      );
    } else {
      return submitAndRegister(value);
    }
  }

  const { errors } = formState;

  const [tosAndPolicyAccepted, setTosAndPolicyAccepted] = useState(false);
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
            data-testid="firstname-text-input"
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
            data-testid="lastname-text-input"
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
            data-testid="email-text-input"
          />
        </div>
      </div>
      {legal && (
        <PrivacyPolicyCheckboxes
          legal={legal}
          onChange={setTosAndPolicyAccepted}
        />
      )}
      <p className="mt-4 ztdl-p mb-6 block text-left">{t("selectMethod")}</p>
      {/* show chooser if both methods are allowed */}
      {loginSettings &&
        loginSettings.allowUsernamePassword &&
        loginSettings.passkeysType == PasskeysType.ALLOWED && (
          <div className="pb-4">
            <AuthenticationMethodRadio
              selected={selected}
              selectionChanged={setSelected}
            />
          </div>
        )}

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}
      <div className="mt-8 flex w-full flex-row items-center justify-between">
        <BackButton data-testid="back-button" />
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid || !tosAndPolicyAccepted}
          onClick={handleSubmit((values) => {
            const usePasswordToContinue: boolean =
              loginSettings?.allowUsernamePassword &&
              loginSettings?.passkeysType == PasskeysType.ALLOWED
                ? !!!(selected === methods[0]) // choose selection if both available
                : !!loginSettings?.allowUsernamePassword; // if password is chosen
            // set password as default if only password is allowed
            return submitAndContinue(values, usePasswordToContinue);
          })}
          data-testid="submit-button"
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("submit")}
        </Button>
      </div>
    </form>
  );
}
