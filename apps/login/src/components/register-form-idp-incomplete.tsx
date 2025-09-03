"use client";

import { registerUserAndLinkToIDP } from "@/lib/server/register";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { FieldValues, useForm } from "react-hook-form";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
    }
  | FieldValues;

type Props = {
  organization: string;
  requestId?: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  defaultValues?: {
    firstname?: string;
    lastname?: string;
    email?: string;
  };
  idpUserId: string;
  idpId: string;
  idpUserName: string;
};

export function RegisterFormIDPIncomplete({
  organization,
  requestId,
  idpIntent,
  defaultValues,
  idpUserId,
  idpId,
  idpUserName,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      email: defaultValues?.email ?? "",
      firstname: defaultValues?.firstname ?? "",
      lastname: defaultValues?.lastname ?? "",
    },
  });

  const t = useTranslations("register");

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitAndRegister(values: Inputs) {
    setLoading(true);
    const response = await registerUserAndLinkToIDP({
      idpId: idpId,
      idpUserName: idpUserName,
      idpUserId: idpUserId,
      email: values.email,
      firstName: values.firstname,
      lastName: values.lastname,
      organization: organization,
      requestId: requestId,
      idpIntent: idpIntent,
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

  const { errors } = formState;

  return (
    <form className="w-full">
      <div className="mb-4 grid grid-cols-2 gap-4">
        <div className="">
          <TextInput
            type="firstname"
            autoComplete="firstname"
            required
            {...register("firstname", { required: t("required.firstname") })}
            label={t("labels.firstname")}
            error={errors.firstname?.message as string}
            data-testid="firstname-text-input"
          />
        </div>
        <div className="">
          <TextInput
            type="lastname"
            autoComplete="lastname"
            required
            {...register("lastname", { required: t("required.lastname") })}
            label={t("labels.lastname")}
            error={errors.lastname?.message as string}
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
            error={errors.email?.message as string}
            data-testid="email-text-input"
          />
        </div>
      </div>

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
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitAndRegister)}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}{" "}
          <Translated i18nKey="submit" namespace="register" />
        </Button>
      </div>
    </form>
  );
}
