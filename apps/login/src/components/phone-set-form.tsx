"use client";

import { Alert } from "@/components/alert";
import { handleServerActionResponse } from "@/lib/client-utils";
import { setPhoneAndContinue } from "@/lib/server/phone";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { AutoSubmitForm } from "./auto-submit-form";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  phone: string;
};

type Props = {
  userId: string;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
  checkAfter?: string;
};

function sanitizePhoneInput(value: string): string {
  const keepDigitsAndPlus = value.replace(/[^\d+]/g, "");
  const hasLeadingPlus = keepDigitsAndPlus.startsWith("+");
  const digits = keepDigitsAndPlus.replace(/\D/g, "");
  return hasLeadingPlus ? `+${digits}` : digits;
}

export function PhoneSetForm({ userId, loginName, sessionId, requestId, organization, checkAfter }: Props) {
  const router = useRouter();
  const t = useTranslations("otp");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onChange",
  });

  async function submit(values: Inputs) {
    setError("");
    setLoading(true);

    try {
      const response = await setPhoneAndContinue({
        userId,
        phone: values.phone,
        loginName,
        sessionId,
        requestId,
        organization,
        checkAfter,
      });

      handleServerActionResponse(response, router, setSamlData, setError);
    } catch {
      setError(t("set.errors.couldNotSavePhoneNumber"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      <form className="w-full">
        <div className="mt-4">
          <TextInput
            type="tel"
            autoComplete="tel"
            autoFocus
            inputMode="tel"
            autoCapitalize="none"
            autoCorrect="off"
            spellCheck={false}
            onInput={(event) => {
              const input = event.currentTarget;
              input.value = sanitizePhoneInput(input.value);
            }}
            {...register("phone", {
              required: t("set.required.phone"),
              pattern: {
                value: /^\+[1-9]\d{6,14}$/,
                message: t("set.validation.phoneInternationalFormat"),
              },
            })}
            label={t("set.labels.phone")}
            data-testid="phone-text-input"
            error={formState.errors.phone?.message}
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
            onClick={handleSubmit(submit)}
            data-testid="submit-button"
          >
            {loading && <Spinner className="mr-2 h-5 w-5" />}
            <Translated i18nKey="set.submit" namespace="otp" />
          </Button>
        </div>
      </form>
    </>
  );
}
