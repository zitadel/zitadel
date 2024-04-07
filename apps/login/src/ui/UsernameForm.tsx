"use client";

import { useEffect, useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2beta/login_settings_pb";

type Inputs = {
  loginName: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName: string | undefined;
  authRequestId: string | undefined;
  organization?: string;
  submit: boolean;
};

export default function UsernameForm({
  loginSettings,
  loginName,
  authRequestId,
  organization,
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

  async function submitLoginName(values: Inputs, organization?: string) {
    setLoading(true);

    let body: any = {
      loginName: values.loginName,
    };

    if (organization) {
      body.organization = organization;
    }

    if (authRequestId) {
      body.authRequestId = authRequestId;
    }

    const res = await fetch("/api/loginname", {
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

  function setLoginNameAndGetAuthMethods(
    values: Inputs,
    organization?: string,
  ) {
    return submitLoginName(values, organization).then((response) => {
      if (response.authMethodTypes.length == 1) {
        const method = response.authMethodTypes[0];
        switch (method) {
          case 1: // user has only password as auth method
            const paramsPassword: any = {
              loginName: response.factors.user.loginName,
            };

            if (organization) {
              paramsPassword.organization = organization;
            }

            if (loginSettings?.passkeysType === 1) {
              paramsPassword.promptPasswordless = `true`; // PasskeysType.PASSKEYS_TYPE_ALLOWED,
            }

            if (authRequestId) {
              paramsPassword.authRequestId = authRequestId;
            }

            return router.push(
              "/password?" + new URLSearchParams(paramsPassword),
            );
          case 2: // AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY
            const paramsPasskey: any = { loginName: values.loginName };
            if (authRequestId) {
              paramsPasskey.authRequestId = authRequestId;
            }
            if (organization) {
              paramsPasskey.organization = organization;
            }

            return router.push(
              "/passkey/login?" + new URLSearchParams(paramsPasskey),
            );
          default:
            const paramsPasskeyDefault: any = { loginName: values.loginName };

            if (loginSettings?.passkeysType === 1) {
              paramsPasskeyDefault.promptPasswordless = `true`; // PasskeysType.PASSKEYS_TYPE_ALLOWED,
            }

            if (authRequestId) {
              paramsPasskeyDefault.authRequestId = authRequestId;
            }
            if (organization) {
              paramsPasskeyDefault.organization = organization;
            }

            return router.push(
              "/password?" + new URLSearchParams(paramsPasskeyDefault),
            );
        }
      } else if (
        response.authMethodTypes &&
        response.authMethodTypes.length === 0
      ) {
        setError(
          "User has no available authentication methods. Contact your administrator to setup authentication for the requested user.",
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

          if (organization) {
            passkeyParams.organization = organization;
          }

          return router.push(
            "/passkey/login?" + new URLSearchParams(passkeyParams),
          );
        } else {
          // user has no passkey setup and login settings allow passkeys
          const paramsPasswordDefault: any = { loginName: values.loginName };

          if (loginSettings?.passkeysType === 1) {
            paramsPasswordDefault.promptPasswordless = `true`; // PasskeysType.PASSKEYS_TYPE_ALLOWED,
          }

          if (authRequestId) {
            paramsPasswordDefault.authRequestId = authRequestId;
          }

          if (organization) {
            paramsPasswordDefault.organization = organization;
          }

          return router.push(
            "/password?" + new URLSearchParams(paramsPasswordDefault),
          );
        }
      }
    });
  }

  const { errors } = formState;

  useEffect(() => {
    if (submit && loginName) {
      // When we navigate to this page, we always want to be redirected if submit is true and the parameters are valid.
      setLoginNameAndGetAuthMethods({ loginName }, organization);
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
          onClick={handleSubmit((e) =>
            setLoginNameAndGetAuthMethods(e, organization),
          )}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
