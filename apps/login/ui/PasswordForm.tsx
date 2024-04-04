"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import { LoginSettings } from "@zitadel/server";

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
        password: values.password,
        authRequestId,
      }),
    });

    const response = await res.json();

    setLoading(false);
    if (!res.ok) {
      setError(response.details);
      return Promise.reject(response.details);
    }
    return response;
  }

  function submitPasswordAndContinue(value: Inputs): Promise<boolean | void> {
    return submitPassword(value).then((resp: any) => {
      if (
        resp.factors &&
        !resp.factors.passwordless && // if session was not verified with a passkey
        promptPasswordless && // if explicitly prompted due policy
        !isAlternative // escaped if password was used as an alternative method
      ) {
        const params = new URLSearchParams({
          loginName: resp.factors.user.loginName,
          promptPasswordless: "true",
        });

        if (organization) {
          params.append("organization", organization);
        }

        return router.push(`/passkey/add?` + params);
      } else {
        let continueWithMfa = undefined;
        if (
          loginSettings?.forceMfa &&
          loginSettings.secondFactors?.length >= 1 // TODO replace with user methods - if forceMFA is set and no user methods prompt to add method (/mfa/add)
        ) {
          if (loginSettings.secondFactors?.length === 1) {
            continueWithMfa = loginSettings.secondFactors[0];
          } else {
            // continueWithMfa = loginSettings.secondFactors[0];
            // render selection page for mfa (/mfa/select)
          }
        }
        // OIDC flows
        if (authRequestId && resp && resp.sessionId) {
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
                }
          );

          if (organization) {
            params.append("organization", organization);
          }

          return router.push(`/signedin?` + params);
        }
      }
    });
  }

  const { errors } = formState;

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
