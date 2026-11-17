"use client";

import { sendLoginname } from "@/lib/actions/loginname";
import { resetPassword, sendPassword } from "@/lib/actions/password";
import { handleServerActionResponse } from "@/lib/client-utils";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { AutoSubmitForm } from "./auto-submit-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  loginName: string;
  password: string;
};

// The action in flight, so only the control that started it shows a spinner.
type Pending = "submit" | "passkey" | "reset";

type Props = {
  loginName: string | undefined;
  requestId: string | undefined;
  loginSettings: LoginSettings | undefined;
  organization?: string;
  defaultOrganization?: string;
  suffix?: string;
  hideSuffix?: boolean;
  submit: boolean;
  allowRegister: boolean;
};

export function LoginForm({
  loginName,
  requestId,
  organization,
  defaultOrganization,
  suffix,
  hideSuffix,
  loginSettings,
  submit,
  allowRegister,
}: Props) {
  const { register, handleSubmit, trigger, getValues, formState } = useForm<Inputs>({
    mode: "onChange",
    defaultValues: {
      loginName: loginName ? loginName : "",
      password: "",
    },
  });

  const t = useTranslations("loginname");
  const tPassword = useTranslations("password");

  const router = useRouter();

  const [pending, setPending] = useState<Pending | undefined>();
  const [error, setError] = useState<string>("");
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const loading = !!pending;

  // Discovery by login name only. It picks the passkey, the identity provider or the password screen for the user.
  const continueWithLoginName = useCallback(
    async (value: string, source: Pending) => {
      setError("");
      setPending(source);

      try {
        const res = await sendLoginname({
          loginName: value,
          organization,
          defaultOrganization,
          requestId,
          suffix,
        });

        handleServerActionResponse(res, router, setSamlData, setError);
        return res;
      } catch {
        setError(t("errors.internalError"));
      } finally {
        setPending(undefined);
      }
    },
    [organization, defaultOrganization, requestId, suffix, router, t],
  );

  useEffect(() => {
    if (submit && loginName) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      continueWithLoginName(loginName, "submit");
    }
  }, [submit, loginName, continueWithLoginName]);

  // With an org domain suffix only the local part is typed, the password actions need the full login name.
  const fullLoginName = (value: string) => (suffix ? `${value}@${suffix}` : value);

  async function submitPassword(values: Inputs) {
    setError("");
    setPending("submit");

    try {
      const response = await sendPassword({
        loginName: fullLoginName(values.loginName),
        organization,
        defaultOrganization,
        checks: create(ChecksSchema, {
          password: { password: values.password },
        }),
        requestId,
      });

      handleServerActionResponse(response, router, setSamlData, setError);
    } catch {
      setError(tPassword("verify.errors.couldNotVerifyPassword"));
    } finally {
      setPending(undefined);
    }
  }

  // The passkey and the reset link only need the login name, so the password stays out of the validation.
  async function validLoginName(): Promise<string | undefined> {
    const valid = await trigger("loginName", { shouldFocus: true });
    return valid ? getValues("loginName") : undefined;
  }

  async function continueWithPasskey() {
    const value = await validLoginName();
    if (value) {
      await continueWithLoginName(value, "passkey");
    }
  }

  async function resetPasswordAndContinue() {
    const value = await validLoginName();
    if (!value) {
      return;
    }

    setError("");
    setPending("reset");

    const resetLoginName = fullLoginName(value);

    const response = await resetPassword({
      loginName: resetLoginName,
      organization,
      defaultOrganization,
      requestId,
    })
      .catch(() => ({ error: tPassword("errors.couldNotSendResetLink") }))
      .finally(() => {
        setPending(undefined);
      });

    if (response && "error" in response) {
      setError(response.error as string);
      return;
    }

    const params = new URLSearchParams({ loginName: resetLoginName });

    if (organization) {
      params.append("organization", organization);
    }

    if (requestId) {
      params.append("requestId", requestId);
    }

    return router.push("/password/set?" + params);
  }

  function goToRegister() {
    const registerParams = new URLSearchParams();
    if (organization) {
      registerParams.append("organization", organization);
    }
    if (requestId) {
      registerParams.append("requestId", requestId);
    }

    router.push("/register?" + registerParams);
  }

  const emailLogin = !loginSettings?.disableLoginWithEmail;

  let inputLabel = t("labels.email");
  if (!emailLogin) {
    inputLabel = loginSettings?.disableLoginWithPhone ? t("labels.username") : t("labels.usernameOrPhoneNumber");
  }

  const { errors } = formState;

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <div className={`${error && "animate-shake transform-gpu"}`}>
          <TextInput
            // the suffix is rendered next to the field, so the typed value is not a complete address
            type={emailLogin && !suffix ? "email" : "text"}
            autoComplete="username"
            autoCapitalize="none"
            autoCorrect="off"
            spellCheck={false}
            autoFocus
            {...register("loginName", { required: t("required.loginName") })}
            label={inputLabel}
            error={errors.loginName?.message}
            data-testid="username-text-input"
            suffix={hideSuffix ? undefined : suffix}
          />
          <TextInput
            type="password"
            autoComplete="current-password"
            {...register("password", { required: tPassword("verify.required.password") })}
            label={tPassword("verify.labels.password")}
            data-testid="password-text-input"
          />
          {!loginSettings?.hidePasswordReset && (
            <button
              className="hover:text-primary-light-500 dark:hover:text-primary-dark-500 text-sm transition-all"
              onClick={() => resetPasswordAndContinue()}
              type="button"
              disabled={loading}
              data-testid="reset-button"
            >
              <Translated i18nKey="verify.resetPassword" namespace="password" />
            </button>
          )}
        </div>

        {error && (
          <div className="py-4" data-testid="error">
            <Alert>{error}</Alert>
          </div>
        )}

        <div className="mt-4 flex w-full flex-row items-center">
          <BackButton data-testid="back-button" />
          <span className="flex-grow"></span>
          <Button
            data-testid="submit-button"
            type="submit"
            className="self-end"
            variant={ButtonVariants.Primary}
            disabled={loading || !formState.isValid}
            onClick={handleSubmit(submitPassword)}
          >
            {pending === "submit" && <Spinner className="mr-2 h-5 w-5" />}
            <Translated i18nKey="submit" namespace="loginname" />
          </Button>
        </div>

        {loginSettings?.passkeysType === PasskeysType.ALLOWED && (
          <Button
            type="button"
            className="mt-4 w-full justify-center"
            variant={ButtonVariants.Secondary}
            disabled={loading}
            onClick={() => continueWithPasskey()}
            data-testid="passkey-button"
          >
            {pending === "passkey" && <Spinner className="mr-2 h-5 w-5" />}
            <Translated i18nKey="passkey" namespace="loginname" />
          </Button>
        )}

        {allowRegister && (
          <div className="border-divider-light dark:border-divider-dark mt-6 border-t pt-6">
            <Button
              type="button"
              className="text-primary-light-500 dark:text-primary-dark-500 border-primary-light-500 dark:border-primary-dark-500 w-full justify-center font-semibold"
              variant={ButtonVariants.Secondary}
              onClick={goToRegister}
              disabled={loading}
              data-testid="register-button"
            >
              <Translated i18nKey="register" namespace="loginname" />
            </Button>
          </div>
        )}
      </form>
    </>
  );
}
