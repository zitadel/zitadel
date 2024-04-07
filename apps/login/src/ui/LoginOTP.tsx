"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { coerceToArrayBuffer, coerceToBase64Url } from "@/utils/base64";
import { Button, ButtonVariants } from "./Button";
import Alert, { AlertType } from "./Alert";
import { Spinner } from "./Spinner";
import { useForm } from "react-hook-form";
import { TextInput } from "./Input";
import { Checks } from "@zitadel/proto/zitadel/session/v2beta/session_service_pb";
import { PlainMessage } from "@zitadel/client2";
import { Challenges } from "@zitadel/proto/zitadel/session/v2beta/challenge_pb";

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

export default function LoginOTP({
  loginName,
  sessionId,
  authRequestId,
  organization,
  method,
  code,
}: Props) {
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
    const challenges: PlainMessage<Challenges> = {};

    if (method === "email") {
      challenges.otpEmail = "";
    }

    if (method === "sms") {
      challenges.otpSms = "";
    }
    setLoading(true);
    const res = await fetch("/api/session", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        sessionId,
        organization,
        challenges,
        authRequestId,
      }),
    });

    setLoading(false);
    if (!res.ok) {
      const error = await res.json();
      throw error.details.details;
    }
    return res.json();
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

    const checks: PlainMessage<Checks> = {};
    if (method === "sms") {
      checks.otpSms = { code: values.code };
    }
    if (method === "email") {
      checks.otpEmail = { code: values.code };
    }
    if (method === "time-based") {
      checks.totp = { code: values.code };
    }

    const res = await fetch("/api/session", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        sessionId,
        organization,
        checks,
        authRequestId,
      }),
    });

    setLoading(false);
    if (!res.ok) {
      const response = await res.json();

      setError(response.details.details ?? "An internal error occurred");
      return Promise.reject(
        response.details.details ?? "An internal error occurred",
      );
    }
    return res.json();
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

        return router.push(`/login?` + params);
      } else {
        const params = new URLSearchParams(
          authRequestId
            ? {
                loginName: response.factors.user.loginName,
                authRequestId,
              }
            : {
                loginName: response.factors.user.loginName,
              },
        );

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/signedin?` + params);
      }
    });
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      {["email", "sms"].includes(method) && (
        <Alert type={AlertType.INFO}>
          <div className="flex flex-row">
            <span className="flex-1 mr-auto text-left">
              Did not get the Code?
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
              Resend
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
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit((e) => setCodeAndContinue(e, organization))}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
