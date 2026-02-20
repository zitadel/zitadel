"use client";

import { registerUserAndLinkToIDP } from "@/lib/server/register";
import { handleServerActionResponse } from "@/lib/client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { FieldValues, useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";
import { AutoSubmitForm } from "./auto-submit-form";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
      username?: string;
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
  idpUserName?: string;
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
    mode: "onChange",
    defaultValues: {
      email: defaultValues?.email ?? "",
      firstname: defaultValues?.firstname ?? "",
      lastname: defaultValues?.lastname ?? "",
    },
  });

  const t = useTranslations("register");
  const router = useRouter();

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  async function submitAndRegister(values: Inputs) {
    setLoading(true);
    try {
      const response = await registerUserAndLinkToIDP({
        idpId: idpId,
        idpUserName: idpUserName ? idpUserName : values.username,
        idpUserId: idpUserId,
        email: values.email,
        firstName: values.firstname,
        lastName: values.lastname,
        organization: organization,
        requestId: requestId,
        idpIntent: idpIntent,
      });

      handleServerActionResponse(response, router, setSamlData, setError);
    } catch {
      setError("Could not register user");
    } finally {
      setLoading(false);
    }
  }

  const { errors } = formState;

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <div className="mb-4 grid grid-cols-1 gap-4">
          {!idpUserName && (
            <div className="">
              <TextInput
                type="text"
                autoComplete="username"
                autoFocus
                required
                {...register("username", { required: "Username is required" })}
                label="Username"
                error={errors.username?.message as string}
                data-testid="username-text-input"
              />
            </div>
          )}
          <div className="grid grid-cols-2 gap-4">
            <div className="">
              <TextInput
                type="firstname"
                autoComplete="firstname"
                autoFocus={!!idpUserName}
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
          </div>
          <div className="">
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
            {loading && <Spinner className="mr-2 h-5 w-5" />} <Translated i18nKey="submit" namespace="register" />
          </Button>
        </div>
      </form>
    </>
  );
}
