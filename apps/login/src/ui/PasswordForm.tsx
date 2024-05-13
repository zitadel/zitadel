"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import {
  LoginSettings,
  AuthFactor,
  Checks,
  AuthenticationMethodType,
} from "@zitadel/server";

type Inputs = {
  password: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName?: string;
  organization?: string;
  authRequestId?: string;
  isAlternative?: boolean; // whether password was requested as alternative auth method
  promptPasswordless?: boolean;
};

export default function PasswordForm({
  loginSettings,
  loginName,
  organization,
  authRequestId,
  promptPasswordless,
  isAlternative,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitPassword(values: Inputs) {
    setError("");
    setLoading(true);

    const res = await fetch("/api/session", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        organization,
        checks: {
          password: { password: values.password },
        } as Checks,
        authRequestId,
      }),
    });

    const response = await res.json();

    setLoading(false);
    if (!res.ok) {
      console.log(response.details.details);
      setError(response.details?.details ?? "Could not verify password");
      return Promise.reject(response.details);
    }
    return response;
  }

  function submitPasswordAndContinue(value: Inputs): Promise<boolean | void> {
    return submitPassword(value).then((resp) => {
      // if user has mfa -> /otp/[method] or /u2f
      // if mfa is forced and user has no mfa -> /mfa/set
      // if no passwordless -> /passkey/add

      // exclude password
      const availableSecondFactors = resp.authMethods?.filter(
        (m: AuthenticationMethodType) => m !== 1,
      );
      if (availableSecondFactors.length == 1) {
        const params = new URLSearchParams({
          loginName: resp.factors.user.loginName,
        });

        if (authRequestId) {
          params.append("authRequestId", authRequestId);
        }

        if (organization) {
          params.append("organization", organization);
        }

        const factor = availableSecondFactors[0];
        if (factor === 4) {
          return router.push(`/otp/time-based?` + params);
        } else if (factor === 6) {
          return router.push(`/otp/sms?` + params);
        } else if (factor === 7) {
          return router.push(`/otp/email?` + params);
        } else if (factor === 5) {
          return router.push(`/u2f?` + params);
        }
      } else if (availableSecondFactors.length >= 1) {
        const params = new URLSearchParams({
          loginName: resp.factors.user.loginName,
        });

        if (authRequestId) {
          params.append("authRequestId", authRequestId);
        }

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/mfa?` + params);
      } else if (loginSettings?.forceMfa && !availableSecondFactors.length) {
        const params = new URLSearchParams({
          loginName: resp.factors.user.loginName,
          checkAfter: "true", // this defines if the check is directly made after the setup
        });

        if (authRequestId) {
          params.append("authRequestId", authRequestId);
        }

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/mfa/set?` + params);
      } else if (
        resp.factors &&
        !resp.factors.passwordless && // if session was not verified with a passkey
        promptPasswordless && // if explicitly prompted due policy
        !isAlternative // escaped if password was used as an alternative method
      ) {
        const params = new URLSearchParams({
          loginName: resp.factors.user.loginName,
          promptPasswordless: "true",
        });

        if (authRequestId) {
          params.append("authRequestId", authRequestId);
        }

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/passkey/add?` + params);
      } else if (authRequestId && resp && resp.sessionId) {
        const params = new URLSearchParams({
          sessionId: resp.sessionId,
          authRequest: authRequestId,
        });

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/login?` + params);
      } else {
        // without OIDC flow
        const params = new URLSearchParams(
          authRequestId
            ? {
                loginName: resp.factors.user.loginName,
                authRequestId,
              }
            : {
                loginName: resp.factors.user.loginName,
              },
        );

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/signedin?` + params);
      }
    });
  }

  return (
    <form className="w-full">
      <div className={`${error && "transform-gpu animate-shake"}`}>
        <TextInput
          type="password"
          autoComplete="password"
          {...register("password", { required: "This field is required" })}
          label="Password"
          //   error={errors.username?.message as string}
        />

        {loginName && (
          <input type="hidden" name="loginName" value={loginName} />
        )}
      </div>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        {/* <Button type="button" variant={ButtonVariants.Secondary}>
          back
        </Button> */}
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitPasswordAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
