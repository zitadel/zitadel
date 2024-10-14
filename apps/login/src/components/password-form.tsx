"use client";

import { resetPassword, sendPassword } from "@/lib/server/password";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs = {
  password: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName: string;
  organization?: string;
  authRequestId?: string;
  isAlternative?: boolean; // whether password was requested as alternative auth method
  promptPasswordless?: boolean;
};

export function PasswordForm({
  loginSettings,
  loginName,
  organization,
  authRequestId,
  promptPasswordless,
  isAlternative,
}: Props) {
  const t = useTranslations("password");

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [info, setInfo] = useState<string>("");
  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitPassword(values: Inputs) {
    setError("");
    setLoading(true);

    const response = await sendPassword({
      loginName,
      organization,
      checks: create(ChecksSchema, {
        password: { password: values.password },
      }),
      authRequestId,
    }).catch(() => {
      setError("Could not verify password");
    });

    if (response && "error" in response && response.error) {
      setError(response.error);
    }

    setLoading(false);

    return response;
  }

  async function resetPasswordAndContinue() {
    setError("");
    setInfo("");
    setLoading(true);

    const response = await resetPassword({
      loginName,
      organization,
    }).catch(() => {
      setError("Could not reset password");
    });

    if (response && "error" in response) {
      setError(response.error);
    } else {
      setInfo("Password was reset. Please check your email.");
    }

    setLoading(false);

    return response;
  }

  async function submitPasswordAndContinue(
    value: Inputs,
  ): Promise<boolean | void> {
    const submitted = await submitPassword(value);
    setInfo("");
    // if user has mfa -> /otp/[method] or /u2f
    // if mfa is forced and user has no mfa -> /mfa/set
    // if no passwordless -> /passkey/set

    // exclude password and passwordless
    if (
      !submitted ||
      !submitted.authMethods ||
      !submitted.factors?.user?.loginName
    ) {
      return;
    }

    const availableSecondFactors = submitted?.authMethods?.filter(
      (m: AuthenticationMethodType) =>
        m !== AuthenticationMethodType.PASSWORD &&
        m !== AuthenticationMethodType.PASSKEY,
    );

    if (availableSecondFactors?.length == 1) {
      const params = new URLSearchParams({
        loginName: submitted.factors?.user.loginName,
      });

      if (authRequestId) {
        params.append("authRequestId", authRequestId);
      }

      if (organization) {
        params.append("organization", organization);
      }

      const factor = availableSecondFactors[0];
      // if passwordless is other method, but user selected password as alternative, perform a login
      if (factor === AuthenticationMethodType.TOTP) {
        return router.push(`/otp/time-based?` + params);
      } else if (factor === AuthenticationMethodType.OTP_SMS) {
        return router.push(`/otp/sms?` + params);
      } else if (factor === AuthenticationMethodType.OTP_EMAIL) {
        return router.push(`/otp/email?` + params);
      } else if (factor === AuthenticationMethodType.U2F) {
        return router.push(`/u2f?` + params);
      }
    } else if (availableSecondFactors?.length >= 1) {
      const params = new URLSearchParams({
        loginName: submitted.factors.user.loginName,
      });

      if (authRequestId) {
        params.append("authRequestId", authRequestId);
      }

      if (organization) {
        params.append("organization", organization);
      }

      return router.push(`/mfa?` + params);
    } else if (loginSettings?.forceMfa && !availableSecondFactors.length) {
      const params = new URLSearchParams({
        loginName: submitted.factors.user.loginName,
        force: "true", // this defines if the mfa is forced in the settings
        checkAfter: "true", // this defines if the check is directly made after the setup
      });

      if (authRequestId) {
        params.append("authRequestId", authRequestId);
      }

      if (organization) {
        params.append("organization", organization);
      }

      // TODO: provide a way to setup passkeys on mfa page?
      return router.push(`/mfa/set?` + params);
    } else if (
      submitted.factors &&
      !submitted.factors.webAuthN && // if session was not verified with a passkey
      promptPasswordless && // if explicitly prompted due policy
      !isAlternative // escaped if password was used as an alternative method
    ) {
      const params = new URLSearchParams({
        loginName: submitted.factors.user.loginName,
        prompt: "true",
      });

      if (authRequestId) {
        params.append("authRequestId", authRequestId);
      }

      if (organization) {
        params.append("organization", organization);
      }

      return router.push(`/passkey/set?` + params);
    } else if (authRequestId && submitted.sessionId) {
      const params = new URLSearchParams({
        sessionId: submitted.sessionId,
        authRequest: authRequestId,
      });

      if (organization) {
        params.append("organization", organization);
      }

      return router.push(`/login?` + params);
    }

    // without OIDC flow
    const params = new URLSearchParams(
      authRequestId
        ? {
            loginName: submitted.factors.user.loginName,
            authRequestId,
          }
        : {
            loginName: submitted.factors.user.loginName,
          },
    );

    if (organization) {
      params.append("organization", organization);
    }

    return router.push(`/signedin?` + params);
  }

  return (
    <form className="w-full">
      <div className={`${error && "transform-gpu animate-shake"}`}>
        <TextInput
          type="password"
          autoComplete="password"
          {...register("password", { required: "This field is required" })}
          label="Password"
        />
        {!loginSettings?.hidePasswordReset && (
          <button
            className="transition-all text-sm hover:text-primary-light-500 dark:hover:text-primary-dark-500"
            onClick={() => resetPasswordAndContinue()}
            type="button"
            disabled={loading}
          >
            {t("resetPassword")}
          </button>
        )}

        {loginName && (
          <input
            type="hidden"
            name="loginName"
            autoComplete="username"
            value={loginName}
          />
        )}
      </div>

      {info && (
        <div className="py-4">
          <Alert type={AlertType.INFO}>{info}</Alert>
        </div>
      )}

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton />
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitPasswordAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("submit")}
        </Button>
      </div>
    </form>
  );
}
