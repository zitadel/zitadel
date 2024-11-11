"use client";

import { inviteUser } from "@/lib/server/invite";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
    }
  | FieldValues;

type Props = {
  firstname?: string;
  lastname?: string;
  email?: string;
  organization?: string;
};

export function InviteForm({
  email,
  firstname,
  lastname,
  organization,
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
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitAndContinue(values: Inputs) {
    setLoading(true);
    const response = await inviteUser({
      email: values.email,
      firstName: values.firstname,
      lastName: values.lastname,
      organization: organization,
    })
      .catch(() => {
        setError("Could not create invitation Code");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && typeof response === "object" && "error" in response) {
      setError(response.error);
      return;
    }

    if (!response) {
      setError("Could not create invitation Code");
      return;
    }

    const params = new URLSearchParams({});

    if (response) {
      params.append("userId", response);
    }

    return router.push(`/invite/success?` + params);
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      <div className="grid grid-cols-2 gap-4 mb-4">
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
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("submit")}
        </Button>
      </div>
    </form>
  );
}
