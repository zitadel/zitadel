"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";

type Inputs = {
  code: string;
};

type Props = {
  userId: string;
};

export default function VerifyEmailForm({ userId }: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitCode(values: Inputs) {
    setLoading(true);
    const res = await fetch("/verifyemail", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        code: values.code,
        userId,
      }),
    });

    if (!res.ok) {
      setLoading(false);
      throw new Error("Failed to verify email");
    }

    setLoading(false);
    return res.json();
  }

  function submitCodeAndContinue(value: Inputs): Promise<boolean | void> {
    console.log(value);
    return submitCode(value).then((resp: any) => {
      return router.push(`/accounts`);
    });
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      <div className="">
        <TextInput
          type="text"
          autoComplete="one-time-code"
          {...register("code", { required: "This field is required" })}
          label="Code"
          //   error={errors.username?.message as string}
        />
      </div>

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
          onClick={handleSubmit(submitCodeAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
