"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";

type Inputs = {};

type Props = {
  sessionId: string;
};

export default function RegisterPasskey({ sessionId }: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitRegister() {
    // const link = await createPasskeyRegistrationLink(server, userId);
    // console.log(link);
    setError("");
    setLoading(true);
    const res = await fetch("/passkeys", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        sessionId,
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

  function submitRegisterAndContinue(value: Inputs): Promise<boolean | void> {
    return submitRegister().then((resp: any) => {
      return router.push(`/accounts`);
    });
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <Button
          type="button"
          variant={ButtonVariants.Secondary}
          onClick={() => router.back()}
        >
          back
        </Button>
        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitRegisterAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
