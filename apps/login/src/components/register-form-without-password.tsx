"use client";

import { registerUser, RegisterUserResponse } from "@/lib/server/register";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import { Alert } from "./alert";
import {
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
  authRequestId?: string;
};

export function RegisterFormWithoutPassword({
  legal,
  email,
  firstname,
  lastname,
  organization,
  authRequestId,
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
  const [selected, setSelected] = useState(methods[0]);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitAndRegister(values: Inputs) {
    setLoading(true);
    const response = await registerUser({
      email: values.email,
      firstName: values.firstname,
      lastName: values.lastname,
      organization: organization,
    }).catch((error) => {
      setError("Could not register user");
      setLoading(false);
    });

    if (response && "error" in response) {
      setError(response.error);
    }

    setLoading(false);

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

    if (authRequestId) {
      registerParams.authRequestId = authRequestId;
    }

    if (withPassword) {
      return router.push(`/register?` + new URLSearchParams(registerParams));
    } else {
      const session = (await submitAndRegister(value)) as RegisterUserResponse;

      const params = new URLSearchParams({});
      if (session?.factors?.user?.loginName) {
        params.set("loginName", session.factors?.user?.loginName);
      }

      if (organization) {
        params.set("organization", organization);
      }

      if (authRequestId) {
        params.set("authRequestId", authRequestId);
      }

      return router.push(`/passkey/set?` + new URLSearchParams(params));
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
      </div>

      {legal && (
        <PrivacyPolicyCheckboxes
          legal={legal}
          onChange={setTosAndPolicyAccepted}
        />
      )}

      <p className="mt-4 ztdl-p mb-6 block text-text-light-secondary-500 dark:text-text-dark-secondary-500">
        {t("selectMethod")}
      </p>

      <div className="pb-4">
        <AuthenticationMethodRadio
          selected={selected}
          selectionChanged={setSelected}
        />
      </div>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center justify-between">
        <BackButton />
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid || !tosAndPolicyAccepted}
          onClick={handleSubmit((values) =>
            submitAndContinue(values, selected === methods[0] ? false : true),
          )}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("submit")}
        </Button>
      </div>
    </form>
  );
}
