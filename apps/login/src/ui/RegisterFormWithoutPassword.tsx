"use client";

import { registerUser } from "@/lib/server/register";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import Alert from "./Alert";
import AuthenticationMethodRadio, {
  methods,
} from "./AuthenticationMethodRadio";
import BackButton from "./BackButton";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { PrivacyPolicyCheckboxes } from "./PrivacyPolicyCheckboxes";
import { Spinner } from "./Spinner";

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

export default function RegisterFormWithoutPassword({
  legal,
  email,
  firstname,
  lastname,
  organization,
  authRequestId,
}: Props) {
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
      setError(error.message ?? "Could not register user");
      setLoading(false);
    });

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
      const session = await submitAndRegister(value).catch((error) => {
        setError(error.message ?? "Could not register user");
      });

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

      return router.push(`/passkey/add?` + new URLSearchParams(params));
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
          continue
        </Button>
      </div>
    </form>
  );
}
