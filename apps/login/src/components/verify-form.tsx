"use client";

import { Alert } from "@/components/alert";
import { resendVerification, verifyUser } from "@/lib/server/email";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { AuthenticatorMethods } from "./authenticator-methods";
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
  isInvite: boolean;
  verifyError?: string;
};

export function VerifyForm({
  userId,
  loginName,
  code,
  organization,
  authRequestId,
  sessionId,
  isInvite,
  verifyError,
}: Props) {
  const t = useTranslations("verify");
  const tError = useTranslations("error");

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ?? "",
    },
  });

  const [authMethods, setAuthMethods] = useState<
    AuthenticationMethodType[] | null
  >(null);

  const [error, setError] = useState<string>(verifyError || "");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  const params = new URLSearchParams({
    userId: userId,
  });

  if (isInvite) {
    params.append("initial", "true");
  }
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
    const params = new URLSearchParams({
      userId: userId,
      initial: "true", // defines that a code is not required and is therefore not shown in the UI
    });

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
    <>
      <h1>{t("verify.title")}</h1>
      <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

      <form className="w-full">
        <div className="">
          <TextInput
            type="text"
            autoComplete="one-time-code"
            {...register("code", { required: "This field is required" })}
            label="Code"
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
            {t("verify.resendCode")}
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
            {t("verify.submit")}
          </Button>
        </div>
      </form>
    </>
  ) : (
    <>
      <h1>{t("setup.title")}</h1>
      <p className="ztdl-p mb-6 block">{t("setup.description")}</p>

      <AuthenticatorMethods authMethods={authMethods} params={params} />
    </>
  );
}
