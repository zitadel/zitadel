"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";

type Inputs = {
  loginName: string;
};

export default function UsernameForm() {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitUsername(values: Inputs) {
    setLoading(true);
    const res = await fetch("/session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName: values.loginName,
      }),
    });

    setLoading(false);
    if (!res.ok) {
      throw new Error("Failed to set user");
    }
    return res.json();
  }

  function submitUsernameAndContinue(value: Inputs): Promise<boolean | void> {
    return submitUsername(value).then(({ factors }) => {
      return router.push(
        `/password?` +
          new URLSearchParams({ loginName: `${factors.user.loginName}` })
      );
    });
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      <div className="">
        <TextInput
          type="text"
          autoComplete="username"
          {...register("loginName", { required: "This field is required" })}
          label="Loginname"
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
          onClick={handleSubmit(submitUsernameAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
