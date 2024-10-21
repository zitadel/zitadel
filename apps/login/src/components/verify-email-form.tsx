"use client";

import { Alert } from "@/components/alert";
import { resendVerification, verifyUser } from "@/lib/server/email";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { PASSKEYS, PASSWORD } from "./auth-methods";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs = {
  code: string;
};

type Props = {
  userId: string;
  loginName: string;
  code: string;
  organization?: string;
  authRequestId?: string;
  sessionId?: string;
  loginSettings?: LoginSettings;
  isInvite: boolean;
};

export function VerifyEmailForm({
  userId,
  loginName,
  code,
  organization,
  authRequestId,
  sessionId,
  loginSettings,
  isInvite,
}: Props) {
  const t = useTranslations("verify");

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ?? "",
    },
  });

  const [authMethods, setAuthMethods] = useState<
    AuthenticationMethodType[] | null
  >(null);

  useEffect(() => {
    if (code && userId) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      // For programmatic verification, the /verifyemail API should be used.
      submitCodeAndContinue({ code });
    }
  }, []);

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  const params = new URLSearchParams({});

  if (loginName) {
    params.append("loginName", loginName);
  }
  if (sessionId) {
    params.append("sessionId", sessionId);
  }
  if (authRequestId) {
    params.append("authRequestId", authRequestId);
  }
  if (organization) {
    params.append("organization", organization);
  }

  async function resendCode() {
    setLoading(true);

    const response = await resendVerification({
      userId,
      isInvite: isInvite,
    }).catch(() => {
      setError("Could not resend email");
      setLoading(false);
      return;
    });

    setLoading(false);
    return response;
  }

  async function submitCodeAndContinue(value: Inputs): Promise<boolean | void> {
    setLoading(true);

    const verifyResponse = await verifyUser({
      code: value.code,
      userId,
      isInvite: isInvite,
    }).catch(() => {
      setError("Could not verify email");
      return;
    });

    setLoading(false);

    if (!verifyResponse) {
      setError("Could not verify email");
      return;
    }

    if (verifyResponse.authMethodTypes) {
      setAuthMethods(verifyResponse.authMethodTypes);
      return;
    }

    // if auth methods fall trough, we complete to login
    const params = new URLSearchParams({});

    if (organization) {
      params.set("organization", organization);
    }

    if (authRequestId && sessionId) {
      params.set("authRequest", authRequestId);
      params.set("sessionId", sessionId);
      return router.push(`/login?` + params);
    } else {
      return router.push(`/loginname?` + params);
    }
  }

  return !authMethods ? (
    <form className="w-full">
      <div className="">
        <TextInput
          type="text"
          autoComplete="one-time-code"
          {...register("code", { required: "This field is required" })}
          label="Code"
          //   error={errors.username?.message as string}
        />
      </div>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <Button
          type="button"
          onClick={() => resendCode()}
          variant={ButtonVariants.Secondary}
        >
          {t("resendCode")}
        </Button>
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitCodeAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("submit")}
        </Button>
      </div>
    </form>
  ) : (
    <div className="grid grid-cols-1 gap-5 w-full pt-4">
      {!authMethods.includes(AuthenticationMethodType.PASSWORD) &&
        PASSWORD(false, "/password/set?" + params)}
      {!authMethods.includes(AuthenticationMethodType.PASSKEY) &&
        PASSKEYS(false, "/passkeys/set?" + params)}
    </div>
  );
}
