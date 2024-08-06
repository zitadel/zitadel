"use client";

import { useEffect, useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

type Inputs = {
  loginName: string;
};

type Props = {
  loginSettings: LoginSettings | undefined;
  loginName: string | undefined;
  authRequestId: string | undefined;
  organization?: string;
  submit: boolean;
  allowRegister: boolean;
};

export default function UsernameForm({
  loginSettings,
  loginName,
  authRequestId,
  organization,
  submit,
  allowRegister,
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
      if (response.nextStep) {
        return router.push(response.nextStep);
      } else if (response.authMethodTypes.length == 1) {
        const method = response.authMethodTypes[0];
        switch (method) {
          case 1: // user has only password as auth method
            const paramsPassword: any = {
              loginName: response.factors.user.loginName,
            };

            // TODO: does this have to be checked in loginSettings.allowDomainDiscovery

            if (organization || response.factors.user.organizationId) {
              paramsPassword.organization =
                organization ?? response.factors.user.organizationId;
            }

            if (
              loginSettings?.passkeysType &&
              (loginSettings?.passkeysType === PasskeysType.ALLOWED ||
                (loginSettings.passkeysType as string) ===
                  "PASSKEYS_TYPE_ALLOWED")
            ) {
              paramsPassword.promptPasswordless = `true`;
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

            if (organization || response.factors.user.organizationId) {
              paramsPasskey.organization =
                organization ?? response.factors.user.organizationId;
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

            if (organization || response.factors.user.organizationId) {
              paramsPasskeyDefault.organization =
                organization ?? response.factors.user.organizationId;
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

          if (organization || response.factors.user.organizationId) {
            passkeyParams.organization =
              organization ?? response.factors.user.organizationId;
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

          if (organization || response.factors.user.organizationId) {
            paramsPasswordDefault.organization =
              organization ?? response.factors.user.organizationId;
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
        {allowRegister && (
          <Button
            type="button"
            className="self-end"
            variant={ButtonVariants.Secondary}
            onClick={() => router.push("/register")}
          >
            register
          </Button>
        )}
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
