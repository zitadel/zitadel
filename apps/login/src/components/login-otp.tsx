"use client";

import { updateSession } from "@/lib/server/session";
import { create } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertType } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { Spinner } from "./spinner";

// either loginName or sessionId must be provided
type Props = {
  loginName?: string;
  sessionId?: string;
  authRequestId?: string;
  organization?: string;
  method: string;
  code?: string;
};

type Inputs = {
  code: string;
};

export function LoginOTP({
  loginName,
  sessionId,
  authRequestId,
  organization,
  method,
  code,
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
    if (!initialized.current && ["email", "sms"].includes(method)) {
      initialized.current = true;
      setLoading(true);
      updateSessionForOTPChallenge()
        .then((response) => {
          setLoading(false);
        })
        .catch((error) => {
          setError(error);
          setLoading(false);
        });
    }
  }, []);

  async function updateSessionForOTPChallenge() {
    let challenges;

    if (method === "email") {
      challenges = create(RequestChallengesSchema, {
        otpEmail: { deliveryType: { case: "sendCode", value: {} } },
      });
    }

    if (method === "sms") {
      challenges = create(RequestChallengesSchema, {
        otpSms: { returnCode: true },
      });
    }

    setLoading(true);
    const response = await updateSession({
      loginName,
      sessionId,
      organization,
      challenges,
      authRequestId,
    }).catch((error) => {
      setError(error.message ?? "Could not request OTP challenge");
      setLoading(false);
    });

    setLoading(false);

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

    if (authRequestId) {
      body.authRequestId = authRequestId;
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
      authRequestId,
    }).catch((error) => {
      setError(error.message ?? "Could not verify OTP code");
      setLoading(false);
    });

    setLoading(false);

    return response;
  }

  function setCodeAndContinue(values: Inputs, organization?: string) {
    return submitCode(values, organization).then((response) => {
      if (authRequestId && response && response.sessionId) {
        const params = new URLSearchParams({
          sessionId: response.sessionId,
          authRequest: authRequestId,
        });

        if (organization) {
          params.append("organization", organization);
        }

        if (authRequestId) {
          params.append("authRequest", authRequestId);
        }

        if (sessionId) {
          params.append("sessionId", sessionId);
        }

        return router.push(`/login?` + params);
      } else {
        const params = new URLSearchParams();
        if (response?.factors?.user?.loginName) {
          params.append("loginName", response.factors.user.loginName);
        }
        if (authRequestId) {
          params.append("authRequestId", authRequestId);
        }

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/signedin?` + params);
      }
    });
  }

  return (
    <form className="w-full">
      {["email", "sms"].includes(method) && (
        <Alert type={AlertType.INFO}>
          <div className="flex flex-row">
            <span className="flex-1 mr-auto text-left">
              {t("noCodeReceived")}
            </span>
            <button
              aria-label="Resend OTP Code"
              disabled={loading}
              className="ml-4 text-primary-light-500 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 hover:text-primary-light-400 cursor-pointer disabled:cursor-default disabled:text-gray-400 dark:disabled:text-gray-700"
              onClick={() => {
                setLoading(true);
                updateSessionForOTPChallenge()
                  .then((response) => {
                    setLoading(false);
                  })
                  .catch((error) => {
                    setError(error);
                    setLoading(false);
                  });
              }}
            >
              {t("resendCode")}
            </button>
          </div>
        </Alert>
      )}
      <div className="mt-4">
        <TextInput
          type="text"
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
        <BackButton />
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit((e) => {
            setCodeAndContinue(e, organization);
          })}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("submit")}
        </Button>
      </div>
    </form>
  );
}
