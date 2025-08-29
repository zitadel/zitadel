"use client";

import { getNextUrl } from "@/lib/client";
import { updateSession } from "@/lib/server/session";
import { create } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { useTranslations } from "next-intl";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

// either loginName or sessionId must be provided
type Props = {
  host: string | null;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
  method: string;
  code?: string;
  loginSettings?: LoginSettings;
};

type Inputs = {
  code: string;
};

export function LoginOTP({
  host,
  loginName,
  sessionId,
  requestId,
  organization,
  method,
  code,
  loginSettings,
}: Props) {
  const t = useTranslations("otp");

  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  const initialized = useRef(false);

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ? code : "",
    },
  });

  useEffect(() => {
    if (!initialized.current && ["email", "sms"].includes(method) && !code) {
      initialized.current = true;
      setLoading(true);
      updateSessionForOTPChallenge()
        .catch((error) => {
          setError(error);
          return;
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, []);

  async function updateSessionForOTPChallenge() {
    let challenges;

    const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

    if (method === "email") {
      challenges = create(RequestChallengesSchema, {
        otpEmail: {
          deliveryType: {
            case: "sendCode",
            value: host
              ? {
                  urlTemplate:
                    `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/otp/${method}?code={{.Code}}&userId={{.UserID}}&sessionId={{.SessionID}}` +
                    (requestId ? `&requestId=${requestId}` : ""),
                }
              : {},
          },
        },
      });
    }

    if (method === "sms") {
      challenges = create(RequestChallengesSchema, {
        otpSms: {},
      });
    }

    setLoading(true);
    const response = await updateSession({
      loginName,
      sessionId,
      organization,
      challenges,
      requestId,
    })
      .catch(() => {
        setError("Could not request OTP challenge");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    return response;
  }

  async function submitCode(values: Inputs, organization?: string) {
    setLoading(true);

    let body: any = {
      code: values.code,
      method,
    };

    if (organization) {
      body.organization = organization;
    }

    if (requestId) {
      body.requestId = requestId;
    }

    let checks;

    if (method === "sms") {
      checks = create(ChecksSchema, {
        otpSms: { code: values.code },
      });
    }
    if (method === "email") {
      checks = create(ChecksSchema, {
        otpEmail: { code: values.code },
      });
    }
    if (method === "time-based") {
      checks = create(ChecksSchema, {
        totp: { code: values.code },
      });
    }

    const response = await updateSession({
      loginName,
      sessionId,
      organization,
      checks,
      requestId,
    })
      .catch(() => {
        setError("Could not verify OTP code");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    return response;
  }

  function setCodeAndContinue(values: Inputs, organization?: string) {
    return submitCode(values, organization).then(async (response) => {
      if (response && "sessionId" in response) {
        setLoading(true);
        // Wait for 2 seconds to avoid eventual consistency issues with an OTP code being verified in the /login endpoint
        await new Promise((resolve) => setTimeout(resolve, 2000));

        const url =
          requestId && response.sessionId
            ? await getNextUrl(
                {
                  sessionId: response.sessionId,
                  requestId: requestId,
                  organization: response.factors?.user?.organizationId,
                },
                loginSettings?.defaultRedirectUri,
              )
            : response.factors?.user
              ? await getNextUrl(
                  {
                    loginName: response.factors.user.loginName,
                    organization: response.factors?.user?.organizationId,
                  },
                  loginSettings?.defaultRedirectUri,
                )
              : null;

        setLoading(false);
        if (url) {
          router.push(url);
        }
      }
    });
  }

  return (
    <form className="w-full">
      {["email", "sms"].includes(method) && (
        <Alert type={AlertType.INFO}>
          <div className="flex flex-row">
            <span className="mr-auto flex-1 text-left">
              <Translated i18nKey="verify.noCodeReceived" namespace="otp" />
            </span>
            <button
              aria-label="Resend OTP Code"
              disabled={loading}
              type="button"
              className="ml-4 cursor-pointer text-primary-light-500 hover:text-primary-light-400 disabled:cursor-default disabled:text-gray-400 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 dark:disabled:text-gray-700"
              onClick={() => {
                setLoading(true);
                updateSessionForOTPChallenge()
                  .catch((error) => {
                    setError(error);
                    return;
                  })
                  .finally(() => {
                    setLoading(false);
                  });
              }}
              data-testid="resend-button"
            >
              <Translated i18nKey="verify.resendCode" namespace="otp" />
            </button>
          </div>
        </Alert>
      )}
      <div className="mt-4">
        <TextInput
          type="text"
          {...register("code", { required: t("verify.required.code") })}
          label={t("verify.labels.code")}
          autoComplete="one-time-code"
          data-testid="code-text-input"
        />
      </div>

      {error && (
        <div className="py-4" data-testid="error">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton data-testid="back-button" />
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit((e) => {
            setCodeAndContinue(e, organization);
          })}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}{" "}
          <Translated i18nKey="verify.submit" namespace="otp" />
        </Button>
      </div>
    </form>
  );
}
