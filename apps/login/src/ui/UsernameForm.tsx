"use client";

import { sendLoginname } from "@/lib/server/loginname";
import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useRouter } from "next/navigation";
import { ReactNode, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import Alert from "./Alert";
import BackButton from "./BackButton";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { Spinner } from "./Spinner";

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
  children?: ReactNode;
};

export default function UsernameForm({
  loginSettings,
  loginName,
  authRequestId,
  organization,
  submit,
  allowRegister,
  children,
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

    const res = await sendLoginname({
      loginName: values.loginName,
      organization,
      authRequestId,
    }).catch((error: Error) => {
      setError(error.message ?? "An internal error occurred");
      return Promise.reject(error ?? "An internal error occurred");
    });

    setLoading(false);

    return res;
  }

  async function setLoginNameAndGetAuthMethods(
    values: Inputs,
    organization?: string,
  ) {
    const response = await submitLoginName(values, organization);

    if (!response) {
      setError("An internal error occurred");
      return;
    }

    if (response?.authMethodTypes && response.authMethodTypes.length === 0) {
      setError(
        "User has no available authentication methods. Contact your administrator to setup authentication for the requested user.",
      );
      return;
    }

    if (response?.authMethodTypes.length == 1) {
      const method = response.authMethodTypes[0];
      switch (method) {
        case AuthenticationMethodType.PASSWORD: // user has only password as auth method
          const paramsPassword: any = {
            loginName: response?.factors?.user?.loginName,
          };

          // TODO: does this have to be checked in loginSettings.allowDomainDiscovery

          if (organization || response?.factors?.user?.organizationId) {
            paramsPassword.organization =
              organization ?? response?.factors?.user?.organizationId;
          }

          if (
            loginSettings?.passkeysType &&
            loginSettings?.passkeysType === PasskeysType.ALLOWED
          ) {
            paramsPassword.promptPasswordless = `true`;
          }

          if (authRequestId) {
            paramsPassword.authRequestId = authRequestId;
          }

          return router.push(
            "/password?" + new URLSearchParams(paramsPassword),
          );
        case AuthenticationMethodType.PASSKEY: // AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY
          const paramsPasskey: any = { loginName: values.loginName };
          if (authRequestId) {
            paramsPasskey.authRequestId = authRequestId;
          }

          if (organization || response?.factors?.user?.organizationId) {
            paramsPasskey.organization =
              organization ?? response?.factors?.user?.organizationId;
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

          if (organization || response?.factors?.user?.organizationId) {
            paramsPasskeyDefault.organization =
              organization ?? response?.factors?.user?.organizationId;
          }

          return router.push(
            "/password?" + new URLSearchParams(paramsPasskeyDefault),
          );
      }
    } else {
      // prefer passkey in favor of other methods
      if (response?.authMethodTypes.includes(2)) {
        const passkeyParams: any = {
          loginName: values.loginName,
          altPassword: `${response.authMethodTypes.includes(1)}`, // show alternative password option
        };

        if (authRequestId) {
          passkeyParams.authRequestId = authRequestId;
        }

        if (organization || response?.factors?.user?.organizationId) {
          passkeyParams.organization =
            organization ?? response?.factors?.user?.organizationId;
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

        if (organization || response?.factors?.user?.organizationId) {
          paramsPasswordDefault.organization =
            organization ?? response?.factors?.user?.organizationId;
        }

        return router.push(
          "/password?" + new URLSearchParams(paramsPasswordDefault),
        );
      }
    }
  }

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
        {allowRegister && (
          <button
            className="transition-all text-sm hover:text-primary-light-500 dark:hover:text-primary-dark-500"
            onClick={() => {
              const registerParams = new URLSearchParams();
              if (organization) {
                registerParams.append("organization", organization);
              }
              if (authRequestId) {
                registerParams.append("authRequestId", authRequestId);
              }

              router.push("/register?" + registerParams);
            }}
            type="button"
            disabled={loading}
          >
            Register new user
          </button>
        )}
      </div>

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="pt-6 pb-4">{children}</div>

      <div className="mt-4 flex w-full flex-row items-center">
        <BackButton />
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
