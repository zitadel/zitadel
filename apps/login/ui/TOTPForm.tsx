"use client";

import { useEffect, useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";

type Inputs = {
  code: string;
};

type Props = {
  loginName: string | undefined;
  sessionId: string | undefined;
  code: string | undefined;
  authRequestId?: string;
  organization?: string;
  submit: boolean;
};

export default function TOTPForm({
  loginName,
  code,
  authRequestId,
  organization,
  submit,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: code ? code : "",
    },
  });

  const router = useRouter();

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  async function submitCode(values: Inputs, organization?: string) {
    setLoading(true);

    let body: any = {
      code: values.code,
    };

    if (organization) {
      body.organization = organization;
    }

    if (authRequestId) {
      body.authRequestId = authRequestId;
    }

    const res = await fetch("/api/totp/verify", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
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

  useEffect(() => {
    if (submit && code) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      setCodeAndContinue({ code }, organization);
    }
  }, []);

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
