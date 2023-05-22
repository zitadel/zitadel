"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";

type Inputs = {
  password: string;
};

type Props = {
  loginName: string;
};

export default function PasswordForm({ loginName }: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitPassword(values: Inputs) {
    setLoading(true);
    const res = await fetch("/session", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        password: values.password,
      }),
    });

    const response = await res.json();

    if (!res.ok) {
      setLoading(false);
      console.log(response);
      setError(response.details);
      return Promise.reject(response.details);
    } else {
      setLoading(false);
      return response;
    }
  }

  function submitPasswordAndContinue(value: Inputs): Promise<boolean | void> {
    return submitPassword(value).then((resp: any) => {
      return router.push(`/accounts`);
    });
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      <div className="">
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
