"use client";

import { Alert } from "@/components/alert";
import { getDeviceAuthorizationRequest } from "@/lib/server/oidc";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

type Inputs = {
  userCode: string;
};

export function DeviceCodeForm({ userCode }: { userCode?: string }) {
  const t = useTranslations("verify");

  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      userCode: userCode || "",
    },
  });

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
            {...register("userCode", { required: "This field is required" })}
            label="Code"
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
            {loading && <Spinner className="h-5 w-5 mr-2" />}
            {t("verify.submit")}
          </Button>
        </div>
      </form>
    </>
  );
}
