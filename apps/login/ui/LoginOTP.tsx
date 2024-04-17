"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { coerceToArrayBuffer, coerceToBase64Url } from "#/utils/base64";
import { Button, ButtonVariants } from "./Button";
import Alert from "./Alert";
import { Spinner } from "./Spinner";
import { Checks } from "@zitadel/server";
import { useForm } from "react-hook-form";
import { TextInput } from "./Input";
import { Challenges } from "@zitadel/server";

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
    const challenges: Challenges = {};

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

  // async function submitLogin(inputs: Inputs) {
  //   setLoading(true);

  //   const checks: Checks = {};

  //   if (method === "email") {
  //     checks.otpEmail = {
  //       code: inputs.code,
  //     };
  //   }

  //   if (method === "sms") {
  //     checks.otpSms = {
  //       code: inputs.code,
  //     };
  //   }

  //   const res = await fetch("/api/session", {
  //     method: "PUT",
  //     headers: {
  //       "Content-Type": "application/json",
  //     },
  //     body: JSON.stringify({
  //       loginName,
  //       sessionId,
  //       organization,
  //       authRequestId,
  //       checks,
  //     }),
  //   });

  //   const response = await res.json();

  //   setLoading(false);
  //   if (!res.ok) {
  //     setError(response.details);
  //     return Promise.reject(response.details);
  //   }
  //   return response;
  // }

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

    const checks: Checks = {};
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
      method: "POST",
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

      setError(response.message ?? "An internal error occurred");
      return Promise.reject(response.message ?? "An internal error occurred");
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
              }
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
      <div className="">
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
