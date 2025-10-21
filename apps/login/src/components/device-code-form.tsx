"use client";

import { Alert } from "@/components/alert";
import { getDeviceAuthorizationRequest } from "@/lib/server/oidc";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  userCode: string;
};

export function DeviceCodeForm({ userCode }: { userCode?: string }) {
  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      userCode: userCode || "",
    },
  });

  const t = useTranslations("device");

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  async function submitCodeAndContinue(value: Inputs): Promise<boolean | void> {
    setLoading(true);

    const response = await getDeviceAuthorizationRequest(value.userCode)
      .catch(() => {
        setError("Could not continue the request");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (!response || !response.deviceAuthorizationRequest?.id) {
      setError("Could not continue the request");
      return;
    }

    return router.push(
      `/device/consent?` +
        new URLSearchParams({
          requestId: `device_${response.deviceAuthorizationRequest.id}`,
          user_code: value.userCode,
        }).toString(),
    );
  }

  return (
    <>
      <form className="w-full">
        <div className="mt-4">
          <TextInput
            type="text"
            autoComplete="one-time-code"
            {...register("userCode", { required: t("usercode.required.code") })}
            label={t("usercode.labels.code")}
            data-testid="code-text-input"
          />
        </div>

        {error && (
          <div className="py-4" data-testid="error">
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
            onClick={handleSubmit(submitCodeAndContinue)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />}{" "}
            <Translated i18nKey="usercode.submit" namespace="device" />
          </Button>
        </div>
      </form>
    </>
  );
}
