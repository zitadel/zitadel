"use client";

import { getNextUrl } from "@/lib/client";
import { verifyTOTP } from "@/lib/server/verify";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { QRCodeSVG } from "qrcode.react";
import { useState } from "react";
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
export function TotpRegister({
  uri,
  secret,
  loginName,
  sessionId,
  requestId,
  organization,
  checkAfter,
  loginSettings,
}: Props) {
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();

  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
    defaultValues: {
      code: "",
    },
  });

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
          const url =
            requestId && sessionId
              ? await getNextUrl(
                  {
                    sessionId: sessionId,
                    requestId: requestId,
                    organization: organization,
                  },
                  loginSettings?.defaultRedirectUri,
                )
              : loginName
                ? await getNextUrl(
                    {
                      loginName: loginName,
                      organization: organization,
                    },
                    loginSettings?.defaultRedirectUri,
                  )
                : null;

          if (url) {
            return router.push(url);
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
    <div className="flex flex-col items-center ">
      {uri && (
        <>
          <QRCodeSVG
            className="rounded-md w-40 h-40 p-2 bg-white my-4"
            value={uri}
          />
          <div className="mb-4 w-96 flex text-sm my-2 border rounded-lg px-4 py-2 pr-2 border-divider-light dark:border-divider-dark">
            <Link href={uri} target="_blank" className="flex-1 overflow-x-auto">
              {uri}
            </Link>

            <CopyToClipboard value={uri}></CopyToClipboard>
          </div>
          <form className="w-full">
            <div className="">
              <TextInput
                type="text"
                {...register("code", { required: "This field is required" })}
                label="Code"
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
                {loading && <Spinner className="h-5 w-5 mr-2" />}
                <Translated i18nKey="set.submit" namespace="otp" />
              </Button>
            </div>
          </form>
        </>
      )}
    </div>
  );
}
