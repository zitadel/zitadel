"use client";

import { sendLoginname } from "@/lib/server/loginname";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { ReactNode, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  loginName: string;
};

type Props = {
  loginName: string | undefined;
  requestId: string | undefined;
  loginSettings: LoginSettings | undefined;
  organization?: string;
  suffix?: string;
  submit: boolean;
  allowRegister: boolean;
  children?: ReactNode;
};

export function UsernameForm({
  loginName,
  requestId,
  organization,
  suffix,
  loginSettings,
  submit,
  allowRegister,
  children,
}: Props) {
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
      requestId,
      suffix,
    })
      .catch(() => {
        setError("An internal error occurred");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (res && "redirect" in res && res.redirect) {
      return router.push(res.redirect);
    }

    if (res && "error" in res && res.error) {
      setError(res.error);
      return;
    }

    return res;
  }

  useEffect(() => {
    if (submit && loginName) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      submitLoginName({ loginName }, organization);
    }
  }, []);

  let inputLabel = "Loginname";
  if (
    loginSettings?.disableLoginWithEmail &&
    loginSettings?.disableLoginWithPhone
  ) {
    inputLabel = "Username";
  } else if (loginSettings?.disableLoginWithEmail) {
    inputLabel = "Username or phone number";
  } else if (loginSettings?.disableLoginWithPhone) {
    inputLabel = "Username or email";
  }

  return (
    <form className="w-full">
      <div className="">
        <TextInput
          type="text"
          autoComplete="username"
          {...register("loginName", { required: "This field is required" })}
          label={inputLabel}
          data-testid="username-text-input"
          suffix={suffix}
        />
        {allowRegister && (
          <button
            className="transition-all text-sm hover:text-primary-light-500 dark:hover:text-primary-dark-500"
            onClick={() => {
              const registerParams = new URLSearchParams();
              if (organization) {
                registerParams.append("organization", organization);
              }
              if (requestId) {
                registerParams.append("requestId", requestId);
              }

              router.push("/register?" + registerParams);
            }}
            type="button"
            disabled={loading}
            data-testid="register-button"
          >
            <Translated i18nKey="register" namespace="loginname" />
          </button>
        )}
      </div>

      {error && (
        <div className="py-4" data-testid="error">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="pt-6 pb-4">{children}</div>

      <div className="mt-4 flex w-full flex-row items-center">
        <BackButton data-testid="back-button" />
        <span className="flex-grow"></span>
        <Button
          data-testid="submit-button"
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit((e) => submitLoginName(e, organization))}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          <Translated i18nKey="submit" namespace="loginname" />
        </Button>
      </div>
    </form>
  );
}
