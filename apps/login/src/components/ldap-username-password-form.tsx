"use client";

import { createNewSessionForLDAP } from "@/lib/server/idp";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  loginName: string;
  password: string;
};

type Props = {
  idpId: string;
  link: boolean;
};

export function LDAPUsernamePasswordForm({ idpId, link }: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const t = useTranslations("ldap");

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitUsernamePassword(values: Inputs) {
    setError("");
    setLoading(true);

    const response = await createNewSessionForLDAP({
      idpId: idpId,
      username: values.loginName,
      password: values.password,
      link: link,
    })
      .catch(() => {
        setError("Could not start LDAP flow");
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
  }

  return (
    <form className="w-full space-y-4">
      <TextInput
        type="text"
        autoComplete="username"
        {...register("loginName", { required: t("required.username") })}
        label={t("labels.username")}
        data-testid="username-text-input"
      />

      <div className={`${error && "transform-gpu animate-shake"}`}>
        <TextInput
          type="password"
          autoComplete="password"
          {...register("password", { required: t("required.password") })}
          label={t("labels.password")}
          data-testid="password-text-input"
        />
      </div>

      {error && (
        <div className="py-4" data-testid="error">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton data-testid="back-button" />
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitUsernamePassword)}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}
          <Translated i18nKey="submit" namespace="ldap" />
        </Button>
      </div>
    </form>
  );
}
