"use client";

import { coerceToArrayBuffer, coerceToBase64Url } from "@/helpers/base64";
import { finishFlow } from "@/lib/login";
import { addU2F, verifyU2F } from "@/lib/server/u2f";
import { RegisterU2FResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";

type Props = {
  loginName?: string;
  sessionId: string;
  authRequestId?: string;
  organization?: string;
  checkAfter: boolean;
};

export function RegisterU2f({
  loginName,
  sessionId,
  organization,
  authRequestId,
  checkAfter,
}: Props) {
  const t = useTranslations("u2f");

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitVerify(
    u2fId: string,
    passkeyName: string,
    publicKeyCredential: any,
    sessionId: string,
  ) {
    setError("");
    setLoading(true);
    const response = await verifyU2F({
      u2fId,
      passkeyName,
      publicKeyCredential,
      sessionId,
    })
      .catch(() => {
        setError("An error on verifying passkey occurred");
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response?.error) {
      setError(response?.error);
    }

    return response;
  }

  async function submitRegisterAndContinue(): Promise<boolean | void | null> {
    setError("");
    setLoading(true);
    const response = await addU2F({
      sessionId,
    })
      .catch(() => {
        setError("An error on registering passkey");
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response?.error) {
      setError(response?.error);
    }

    if (!response || !("u2fId" in response)) {
      setError("An error on registering passkey");
      return;
    }

    const u2fResponse = response as unknown as RegisterU2FResponse;

    const u2fId = u2fResponse.u2fId;
    const options: CredentialCreationOptions =
      (u2fResponse?.publicKeyCredentialCreationOptions as CredentialCreationOptions) ??
      {};

    if (options.publicKey) {
      options.publicKey.challenge = coerceToArrayBuffer(
        options.publicKey.challenge,
        "challenge",
      );
      options.publicKey.user.id = coerceToArrayBuffer(
        options.publicKey.user.id,
        "userid",
      );
      if (options.publicKey.excludeCredentials) {
        options.publicKey.excludeCredentials.map((cred: any) => {
          cred.id = coerceToArrayBuffer(
            cred.id as string,
            "excludeCredentials.id",
          );
          return cred;
        });
      }

      const resp = await navigator.credentials.create(options);

      if (
        !resp ||
        !(resp as any).response.attestationObject ||
        !(resp as any).response.clientDataJSON ||
        !(resp as any).rawId
      ) {
        setError("An error on registering passkey");
        return;
      }

      const attestationObject = (resp as any).response.attestationObject;
      const clientDataJSON = (resp as any).response.clientDataJSON;
      const rawId = (resp as any).rawId;

      const data = {
        id: resp.id,
        rawId: coerceToBase64Url(rawId, "rawId"),
        type: resp.type,
        response: {
          attestationObject: coerceToBase64Url(
            attestationObject,
            "attestationObject",
          ),
          clientDataJSON: coerceToBase64Url(clientDataJSON, "clientDataJSON"),
        },
      };

      const submitResponse = await submitVerify(u2fId, "", data, sessionId);

      if (!submitResponse) {
        setError("An error on verifying passkey");
        return;
      }

      if (checkAfter) {
        const paramsToContinue = new URLSearchParams({});

        if (sessionId) {
          paramsToContinue.append("sessionId", sessionId);
        }
        if (loginName) {
          paramsToContinue.append("loginName", loginName);
        }
        if (organization) {
          paramsToContinue.append("organization", organization);
        }
        if (authRequestId) {
          paramsToContinue.append("authRequestId", authRequestId);
        }

        return router.push(`/u2f?` + paramsToContinue);
      } else {
        return authRequestId && sessionId
          ? finishFlow({
              sessionId: sessionId,
              authRequestId: authRequestId,
              organization: organization,
            })
          : loginName
            ? finishFlow({
                loginName: loginName,
                organization: organization,
              })
            : null;
      }
    }
  }

  return (
    <form className="w-full">
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton />

        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading}
          onClick={submitRegisterAndContinue}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          {t("set.submit")}
        </Button>
      </div>
    </form>
  );
}
