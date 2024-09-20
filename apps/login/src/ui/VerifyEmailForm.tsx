"use client";

import { resendVerifyEmail, verifyUserByEmail } from "@/lib/server/email";
import Alert from "@/ui/Alert";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { Spinner } from "./Spinner";

type Inputs = {
  code: string;
};

type Props = {
  userId: string;
  loginName: string;
  code: string;
  submit: boolean;
  organization?: string;
  authRequestId?: string;
  sessionId?: string;
  loginSettings?: LoginSettings;
};

export default function VerifyEmailForm({
  userId,
  loginName,
  code,
  submit,
  organization,
  authRequestId,
  sessionId,
  loginSettings,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ?? "",
    },
  });

  useEffect(() => {
    if (submit && code && userId) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      // For programmatic verification, the /verifyemail API should be used.
      submitCodeAndContinue({ code });
    }
  }, []);

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function resendCode() {
    setLoading(true);
    const response = await resendVerifyEmail({
      userId,
    }).catch(() => {
      setError("Could not resend email");
    });

    setLoading(false);
    return response;
  }

  async function submitCodeAndContinue(value: Inputs): Promise<boolean | void> {
    setLoading(true);
    const verifyResponse = await verifyUserByEmail({
      code: value.code,
      userId,
    }).catch(() => {
      setError("Could not verify email");
    });

    setLoading(false);

    if (!verifyResponse) {
      setError("Could not verify email");
      return;
    }

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

  return (
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
          resend code
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
          continue
        </Button>
      </div>
    </form>
  );
}
