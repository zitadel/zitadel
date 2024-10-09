"use client";

import { sendLoginname } from "@/lib/server/loginname";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { ReactNode, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs = {
  loginName: string;
};

type Props = {
  loginName: string | undefined;
  authRequestId: string | undefined;
  organization?: string;
  submit: boolean;
  allowRegister: boolean;
  children?: ReactNode;
};

export function UsernameForm({
  loginName,
  authRequestId,
  organization,
  submit,
  allowRegister,
  children,
}: Props) {
  const t = useTranslations("loginname");
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      loginName: loginName ? loginName : "",
    },
  });

  const router = useRouter();

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  async function submitLoginName(values: Inputs, organization?: string) {
    setLoading(true);

    const res = await sendLoginname({
      loginName: values.loginName,
      organization,
      authRequestId,
    }).catch(() => {
      setError("An internal error occurred");
    });

    setLoading(false);

    if (res?.error) {
      setError(res.error);
    }

    return res;
  }

  useEffect(() => {
    if (submit && loginName) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      submitLoginName({ loginName }, organization);
    }
  }, []);

  return (
    <form className="w-full">
      <div className="">
        <TextInput
          type="text"
          autoComplete="username"
          {...register("loginName", { required: "This field is required" })}
          label="Loginname"
        />
        {allowRegister && (
          <button
            className="transition-all text-sm hover:text-primary-light-500 dark:hover:text-primary-dark-500"
            onClick={() => {
              const registerParams = new URLSearchParams();
              if (organization) {
                registerParams.append("organization", organization);
              }
              if (authRequestId) {
                registerParams.append("authRequestId", authRequestId);
              }

              router.push("/register?" + registerParams);
            }}
            type="button"
            disabled={loading}
          >
            {t("register")}
          </button>
        )}
      </div>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="pt-6 pb-4">{children}</div>

      <div className="mt-4 flex w-full flex-row items-center">
        <BackButton />
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit((e) => submitLoginName(e, organization))}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
