"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import { AuthenticationMethodType, LoginSettings } from "@zitadel/server";

type Inputs = {
  loginName: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName: string | undefined;
};

export default function UsernameForm({ loginSettings, loginName }: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      loginName: loginName ? loginName : "",
    },
  });

  const router = useRouter();

  const [loading, setLoading] = useState<boolean>(false);

  async function submitLoginName(values: Inputs) {
    setLoading(true);
    const res = await fetch("/loginnames", {
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
      throw new Error("Failed to load authentication methods");
    }
    return res.json();
  }

  async function setLoginNameAndGetAuthMethods(values: Inputs) {
    return submitLoginName(values).then((response) => {
      console.log(response);
      if (response.authMethodTypes.length == 1) {
        const method = response.authMethodTypes[0];
        console.log(method);
        // switch (method) {
        //   case AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSWORD:
        //     return router.push(
        //       "/password?" +
        //         new URLSearchParams({ loginName: values.loginName })
        //     );
        //   case AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY:
        //     break;
        //   // return router.push(
        //   //   "/passkey/login?" +
        //   //     new URLSearchParams({ loginName: values.loginName })
        //   // );
        //   default:
        //     return router.push(
        //       "/password?" +
        //         new URLSearchParams({ loginName: values.loginName })
        //     );
        // }
      }
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
          onClick={handleSubmit(setLoginNameAndGetAuthMethods)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
