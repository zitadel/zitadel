"use client";

import { addOTP } from "@/lib/server/otp";
import { useTranslations } from "next-intl";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Props = {
  method: "email" | "sms";
  sessionId: string;
  urlToContinue: string;
  checkAfter?: boolean;
};

export function OTPAddForm({ method, sessionId, urlToContinue, checkAfter }: Props) {
  const router = useRouter();
  const t = useTranslations("otp");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [added, setAdded] = useState<boolean>(false);

  async function submit() {
    setLoading(true);
    setError("");

    const response = await addOTP({ method, sessionId })
      .catch(() => ({ error: t(method === "email" ? "errors.couldNotAddEmail" : "errors.couldNotAddSms") }))
      .finally(() => setLoading(false));

    if ("error" in response) {
      setError(response.error);
      return;
    }

    // if attribute is set, validate MFA after it is set up
    if (checkAfter) {
      router.push(urlToContinue);
      return;
    }

    setAdded(true);
  }

  if (added) {
    return (
      <div className="flex w-full flex-col space-y-4">
        <Alert type={AlertType.INFO}>
          <Translated i18nKey={method === "email" ? "set.emailAdded" : "set.smsAdded"} namespace="otp" />
        </Alert>
        <div className="mt-8 flex w-full flex-row items-center">
          <span className="flex-grow"></span>
          <Link href={urlToContinue}>
            <Button type="button" className="self-end" variant={ButtonVariants.Primary} data-testid="submit-button">
              <Translated i18nKey="set.submit" namespace="otp" />
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="flex w-full flex-col space-y-4">
      {error && (
        <div className="py-4" data-testid="error">
          <Alert>{error}</Alert>
        </div>
      )}
      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton />
        <span className="flex-grow"></span>
        <Button
          type="button"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading}
          onClick={submit}
          data-testid="add-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}
          <Translated i18nKey="set.add" namespace="otp" />
        </Button>
      </div>
    </div>
  );
}
