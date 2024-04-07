"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { PrivacyPolicyCheckboxes } from "./PrivacyPolicyCheckboxes";
import { FieldValues, useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import AuthenticationMethodRadio, {
  methods,
} from "./AuthenticationMethodRadio";
import Alert from "./Alert";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2beta/legal_settings_pb";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
    }
  | FieldValues;

type Props = {
  legal: LegalAndSupportSettings;
  organization?: string;
  authRequestId?: string;
};

export default function RegisterFormWithoutPassword({
  legal,
  organization,
  authRequestId,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [loading, setLoading] = useState<boolean>(false);
  const [selected, setSelected] = useState(methods[0]);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitAndRegister(values: Inputs) {
    setLoading(true);
    const res = await fetch("/api/registeruser", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: values.email,
        firstName: values.firstname,
        lastName: values.lastname,
        organization: organization,
      }),
    });
    setLoading(false);
    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.details);
    }
    return res.json();
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

    return withPassword
      ? router.push(`/register?` + new URLSearchParams(registerParams))
      : submitAndRegister(value)
          .then((session) => {
            setError("");

            const params: any = { loginName: session.factors.user.loginName };

            if (organization) {
              params.organization = organization;
            }

            if (authRequestId) {
              params.authRequestId = authRequestId;
            }

            return router.push(`/passkey/add?` + new URLSearchParams(params));
          })
          .catch((errorDetails: Error) => {
            setLoading(false);
            setError(errorDetails.message);
          });
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
        Select the method you would like to authenticate
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
        <Button
          type="button"
          variant={ButtonVariants.Secondary}
          onClick={() => router.back()}
        >
          back
        </Button>
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid || !tosAndPolicyAccepted}
          onClick={handleSubmit((values) =>
            submitAndContinue(values, selected === methods[0] ? false : true),
          )}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
