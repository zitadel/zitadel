"use client";

import { addU2F, verifyU2F } from "@/lib/server/u2f";
import { coerceToArrayBuffer, coerceToBase64Url } from "@/utils/base64";
import { RegisterU2FResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useRouter } from "next/navigation";
import { useState } from "react";
import Alert from "./Alert";
import BackButton from "./BackButton";
import { Button, ButtonVariants } from "./Button";
import { Spinner } from "./Spinner";

type Props = {
  sessionId: string;
  authRequestId?: string;
  organization?: string;
};

export default function RegisterU2F({
  sessionId,
  organization,
  authRequestId,
}: Props) {
  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitVerify(
    u2fId: string,
    passkeyName: string,
    publicKeyCredential: any,
    sessionId: string,
  ) {
    setLoading(true);
    const response = await verifyU2F({
      u2fId,
      passkeyName,
      publicKeyCredential,
      sessionId,
    }).catch((error: Error) => {
      console.error(error);
      setLoading(false);
      setError("An error on verifying passkey occurred");
    });

    setLoading(false);

    return response;
  }

  async function submitRegisterAndContinue(): Promise<boolean | void> {
    setError("");
    setLoading(true);
    const response = await addU2F({
      sessionId,
    }).catch((error: Error) => {
      console.error(error);
      setLoading(false);
      setError("An error on registering passkey");
    });

    if (response && "error" in response && response?.error) {
      setError(response?.error);
    }

    if (!response || !("u2fId" in response)) {
      setLoading(false);
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
        setLoading(false);
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
        setLoading(false);
        setError("An error on verifying passkey");
        return;
      }

      const params = new URLSearchParams();

      if (organization) {
        params.set("organization", organization);
      }

      if (authRequestId) {
        params.set("authRequestId", authRequestId);
        params.set("sessionId", sessionId);
        // params.set("altPassword", ${false}); // without setting altPassword this does not allow password
        // params.set("loginName", resp.loginName);

        router.push("/u2f?" + params);
      } else {
        router.push("/accounts?" + params);
      }
    }

    setLoading(false);
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
          continue
        </Button>
      </div>
    </form>
  );
}
