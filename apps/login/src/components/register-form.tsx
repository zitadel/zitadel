"use client";

import { registerUser } from "@/lib/server/register";
import { LegalAndSupportSettings } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { LoginSettings, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { FieldValues, useForm } from "react-hook-form";
import { Alert, AlertType } from "./alert";
import { AuthenticationMethod, AuthenticationMethodRadio, methods } from "./authentication-method-radio";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { TextInput } from "./input";
import { PrivacyPolicyCheckboxes } from "./privacy-policy-checkboxes";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs =
  | {
      firstname: string;
      lastname: string;
      email: string;
    }
  | FieldValues;

type Props = {
  legal: LegalAndSupportSettings;
  firstname?: string;
  lastname?: string;
  email?: string;
  organization: string;
  requestId?: string;
  loginSettings?: LoginSettings;
  idpCount: number;
};

export function RegisterForm({
  legal,
  email,
  firstname,
  lastname,
  organization,
  requestId,
  loginSettings,
  idpCount = 0,
}: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      email: email ?? "",
      firstName: firstname ?? "",
      lastname: lastname ?? "",
    },
  });

  const t = useTranslations("register");

  const [loading, setLoading] = useState<boolean>(false);
  const [selected, setSelected] = useState<AuthenticationMethod>(methods[0]);
  const [error, setError] = useState<string>("");

  const router = useRouter();

  async function submitAndRegister(values: Inputs) {
    setLoading(true);
    const response = await registerUser({
      email: values.email,
      firstName: values.firstname,
      lastName: values.lastname,
      organization: organization,
      requestId: requestId,
    })
      .catch(() => {
        setError("Could not register user");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    if (response && "redirect" in response && response.redirect) {
      return router.push(response.redirect);
    }

    return response;
  }

  async function submitAndContinue(value: Inputs, withPassword: boolean = false) {
    const registerParams: any = value;

    if (organization) {
      registerParams.organization = organization;
    }

    if (requestId) {
      registerParams.requestId = requestId;
    }

    // redirect user to /register/password if password is chosen
    if (withPassword) {
      return router.push(`/register/password?` + new URLSearchParams(registerParams));
    } else {
      return submitAndRegister(value);
    }
  }

  const { errors } = formState;

  const [tosAndPolicyAccepted, setTosAndPolicyAccepted] = useState(false);

  // Check if legal acceptance is required
  const isLegalAcceptanceRequired = !!(legal?.tosLink || legal?.privacyPolicyLink);
  const canSubmit = formState.isValid && (!isLegalAcceptanceRequired || tosAndPolicyAccepted);

  return (
    <form className="w-full">
      <div className="mb-4 grid grid-cols-2 gap-4">
        <div className="">
          <TextInput
            type="firstname"
            autoComplete="firstname"
            required
            {...register("firstname", { required: t("required.firstname") })}
            label={t("labels.firstname")}
            error={errors.firstname?.message as string}
            data-testid="firstname-text-input"
          />
        </div>
        <div className="">
          <TextInput
            type="lastname"
            autoComplete="lastname"
            required
            {...register("lastname", { required: t("required.lastname") })}
            label={t("labels.lastname")}
            error={errors.lastname?.message as string}
            data-testid="lastname-text-input"
          />
        </div>
        <div className="col-span-2">
          <TextInput
            type="email"
            autoComplete="email"
            required
            {...register("email", { required: t("required.email") })}
            label={t("labels.email")}
            error={errors.email?.message as string}
            data-testid="email-text-input"
          />
        </div>
      </div>
      {(legal?.tosLink || legal?.privacyPolicyLink) && (
        <PrivacyPolicyCheckboxes legal={legal} onChange={setTosAndPolicyAccepted} />
      )}
      {/* show chooser if both methods are allowed */}
      {loginSettings && loginSettings.allowUsernamePassword && loginSettings.passkeysType == PasskeysType.ALLOWED && (
        <>
          <p className="ztdl-p mb-6 mt-4 block text-left">
            <Translated i18nKey="selectMethod" namespace="register" />
          </p>

          <div className="pb-4">
            <AuthenticationMethodRadio selected={selected} selectionChanged={setSelected} />
          </div>
        </>
      )}
      {!loginSettings?.allowUsernamePassword &&
        loginSettings?.passkeysType !== PasskeysType.ALLOWED &&
        (!loginSettings?.allowExternalIdp || !idpCount) && (
          <div className="py-4">
            <Alert type={AlertType.INFO}>
              <Translated i18nKey="noMethodAvailableWarning" namespace="register" />
            </Alert>
          </div>
        )}

      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center justify-between">
        <BackButton data-testid="back-button" />
        <Button
          type="submit"
          variant={ButtonVariants.Primary}
          disabled={loading || !canSubmit}
          onClick={handleSubmit((values) => {
            const usePasswordToContinue: boolean =
              loginSettings?.allowUsernamePassword && loginSettings?.passkeysType == PasskeysType.ALLOWED
                ? !(selected === methods[0]) // choose selection if both available
                : !!loginSettings?.allowUsernamePassword; // if password is chosen
            // set password as default if only password is allowed
            return submitAndContinue(values, usePasswordToContinue);
          })}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}
          <Translated i18nKey="submit" namespace="register" />
        </Button>
      </div>
    </form>
  );
}
