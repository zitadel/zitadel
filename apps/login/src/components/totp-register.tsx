"use client";

import { completeFlowOrGetUrl } from "@/lib/client";
import { verifyTOTP } from "@/lib/server/verify";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { QRCodeSVG } from "qrcode.react";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { Alert } from "./alert";
import { Button, ButtonVariants } from "./button";
import { CopyToClipboard } from "./copy-to-clipboard";
import { TextInput } from "./input";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Inputs = {
  code: string;
};

type Props = {
  uri: string;
  secret: string;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
  checkAfter?: boolean;
  loginSettings?: LoginSettings;
};
export function TotpRegister({ uri, loginName, sessionId, requestId, organization, checkAfter, loginSettings }: Props) {
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: "",
    },
  });

  const t = useTranslations("otp");

  async function continueWithCode(values: Inputs) {
    setLoading(true);
    return verifyTOTP(values.code, loginName, organization)
      .then(async () => {
        // if attribute is set, validate MFA after it is setup, otherwise proceed as usual (when mfa is enforced to login)
        if (checkAfter) {
          const params = new URLSearchParams({});

          if (loginName) {
            params.append("loginName", loginName);
          }
          if (requestId) {
            params.append("requestId", requestId);
          }
          if (organization) {
            params.append("organization", organization);
          }

          return router.push(`/otp/time-based?` + params);
        } else {
          if (requestId && sessionId) {
            const callbackResponse = await completeFlowOrGetUrl(
              {
                sessionId: sessionId,
                requestId: requestId,
                organization: organization,
              },
              loginSettings?.defaultRedirectUri,
            );

            if ("error" in callbackResponse) {
              setError(callbackResponse.error);
              return;
            }

            if ("redirect" in callbackResponse) {
              return router.push(callbackResponse.redirect);
            }
          } else if (loginName) {
            const callbackResponse = await completeFlowOrGetUrl(
              {
                loginName: loginName,
                organization: organization,
              },
              loginSettings?.defaultRedirectUri,
            );

            if ("error" in callbackResponse) {
              setError(callbackResponse.error);
              return;
            }

            if ("redirect" in callbackResponse) {
              return router.push(callbackResponse.redirect);
            }
          }
        }
      })
      .catch((e) => {
        setError(e.message);
        return;
      })
      .finally(() => {
        setLoading(false);
      });
  }

  return (
    <div className="flex flex-col items-center">
      {uri && (
        <>
          <QRCodeSVG className="my-4 h-40 w-40 rounded-md bg-white p-2" value={uri} />
          <div className="my-2 mb-4 flex w-96 rounded-lg border border-divider-light px-4 py-2 pr-2 text-sm dark:border-divider-dark">
            <Link href={uri} target="_blank" className="flex-1 overflow-x-auto">
              {uri}
            </Link>

            <CopyToClipboard value={uri}></CopyToClipboard>
          </div>
          <form className="w-full">
            <div className="">
              <TextInput
                type="text"
                {...register("code", { required: t("set.required.code") })}
                label={t("set.labels.code")}
                data-testid="code-text-input"
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
                onClick={handleSubmit(continueWithCode)}
                data-testid="submit-button"
              >
                {loading && <Spinner className="mr-2 h-5 w-5" />}
                <Translated i18nKey="set.submit" namespace="otp" />
              </Button>
            </div>
          </form>
        </>
      )}
    </div>
  );
}
