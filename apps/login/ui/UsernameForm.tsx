"use client";

import { useEffect, useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import { LoginSettings } from "@zitadel/server";

type Inputs = {
  loginName: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName: string | undefined;
  authRequestId: string | undefined;
  submit: boolean;
};

export default function UsernameForm({
  loginSettings,
  loginName,
  authRequestId,
  submit,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      loginName: loginName ? loginName : "",
    },
  });

  const router = useRouter();

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  async function submitLoginName(values: Inputs) {
    setLoading(true);

    const body = {
      loginName: values.loginName,
    };

    const res = await fetch("/api/loginname", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(authRequestId ? { ...body, authRequestId } : body),
    });

    setLoading(false);
    if (!res.ok) {
      throw new Error("Failed to load authentication methods");
    }
    return res.json();
  }

  async function setLoginNameAndGetAuthMethods(values: Inputs) {
    return submitLoginName(values).then((response) => {
      if (response.authMethodTypes.length == 1) {
        const method = response.authMethodTypes[0];
        switch (method) {
          case 1: //AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSWORD:
            const paramsPassword: any = { loginName: values.loginName };

            if (loginSettings?.passkeysType === 1) {
              paramsPassword.promptPasswordless = `true`; // PasskeysType.PASSKEYS_TYPE_ALLOWED,
            }

            if (authRequestId) {
              paramsPassword.authRequestId = authRequestId;
            }

            return router.push(
              "/password?" + new URLSearchParams(paramsPassword)
            );
          case 2: // AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY
            const paramsPasskey: any = { loginName: values.loginName };
            if (authRequestId) {
              paramsPasskey.authRequestId = authRequestId;
            }

            return router.push(
              "/passkey/login?" + new URLSearchParams(paramsPasskey)
            );
          default:
            const paramsPasskeyDefault: any = { loginName: values.loginName };

            if (loginSettings?.passkeysType === 1) {
              paramsPasskeyDefault.promptPasswordless = `true`; // PasskeysType.PASSKEYS_TYPE_ALLOWED,
            }

            if (authRequestId) {
              paramsPasskeyDefault.authRequestId = authRequestId;
            }
            return router.push(
              "/password?" + new URLSearchParams(paramsPasskeyDefault)
            );
        }
      } else if (
        response.authMethodTypes &&
        response.authMethodTypes.length === 0
      ) {
        setError(
          "User has no available authentication methods. Contact your administrator to setup authentication for the requested user."
        );
      } else {
        // prefer passkey in favor of other methods
        if (response.authMethodTypes.includes(2)) {
          const passkeyParams: any = {
            loginName: values.loginName,
            altPassword: `${response.authMethodTypes.includes(1)}`, // show alternative password option
          };

          if (authRequestId) {
            passkeyParams.authRequestId = authRequestId;
          }

          return router.push(
            "/passkey/login?" + new URLSearchParams(passkeyParams)
          );
        }
      }
    });
  }

  const { errors } = formState;

  useEffect(() => {
    if (submit && loginName) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      setLoginNameAndGetAuthMethods({ loginName });
    }
  }, []);

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
