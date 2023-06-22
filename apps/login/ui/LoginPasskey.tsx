"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import { RegisterPasskeyResponse } from "@zitadel/server";
import { coerceToArrayBuffer, coerceToBase64Url } from "#/utils/base64";
type Inputs = {};

type Props = {
  sessionId: string;
};

export default function LoginPasskey({ sessionId }: Props) {
  const { login, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitLogin(
    passkeyId: string,
    passkeyName: string,
    publicKeyCredential: any,
    sessionId: string
  ) {
    setLoading(true);
    const res = await fetch("/passkeys/verify", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        passkeyId,
        passkeyName,
        publicKeyCredential,
        sessionId,
      }),
    });

    const response = await res.json();

    setLoading(false);
    if (!res.ok) {
      setError(response.details);
      return Promise.reject(response.details);
    }
    return response;
  }

  function submitLoginAndContinue(value: Inputs): Promise<boolean | void> {
    return submitLogin().then((resp: LoginPasskeyResponse) => {
      const passkeyId = resp.passkeyId;

      if (
        resp.publicKeyCredentialCreationOptions &&
        resp.publicKeyCredentialCreationOptions.publicKey
      ) {
        resp.publicKeyCredentialCreationOptions.publicKey.challenge =
          coerceToArrayBuffer(
            resp.publicKeyCredentialCreationOptions.publicKey.challenge,
            "challenge"
          );
        resp.publicKeyCredentialCreationOptions.publicKey.user.id =
          coerceToArrayBuffer(
            resp.publicKeyCredentialCreationOptions.publicKey.user.id,
            "challenge"
          );
        if (
          resp.publicKeyCredentialCreationOptions.publicKey.excludeCredentials
        ) {
          resp.publicKeyCredentialCreationOptions.publicKey.excludeCredentials.map(
            (cred: any) => {
              cred.id = coerceToArrayBuffer(
                cred.id as string,
                "excludeCredentials.id"
              );
              return cred;
            }
          );
        }

        navigator.credentials
          .get({
            publicKey: resp.publicKeyCredentialCreationOptions,
          })
          .then((assertedCredential: any) => {
            if (assertedCredential) {
              let authData = new Uint8Array(
                assertedCredential.response.authenticatorData
              );
              let clientDataJSON = new Uint8Array(
                assertedCredential.response.clientDataJSON
              );
              let rawId = new Uint8Array(assertedCredential.rawId);
              let sig = new Uint8Array(assertedCredential.response.signature);
              let userHandle = new Uint8Array(
                assertedCredential.response.userHandle
              );

              let data = JSON.stringify({
                id: assertedCredential.id,
                rawId: coerceToBase64Url(rawId, "rawId"),
                type: assertedCredential.type,
                response: {
                  authenticatorData: coerceToBase64Url(authData, "authData"),
                  clientDataJSON: coerceToBase64Url(
                    clientDataJSON,
                    "clientDataJSON"
                  ),
                  signature: coerceToBase64Url(sig, "sig"),
                  userHandle: coerceToBase64Url(userHandle, "userHandle"),
                },
              });

              return submitVerify(passkeyId, "", data, sessionId);
            } else {
              setLoading(false);
              setError("An error on retrieving passkey");
              return null;
            }
          })
          .catch((error) => {
            console.error(error);
            setLoading(false);
            //   setError(error);

            return null;
          });
      }
      //   return router.push(`/accounts`);
    });
  }

  const { errors } = formState;

  return (
    <form className="w-full">
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        {isPrompt ? (
          <Button
            type="button"
            variant={ButtonVariants.Secondary}
            onClick={() => router.push("/accounts")}
          >
            skip
          </Button>
        ) : (
          <Button
            type="button"
            variant={ButtonVariants.Secondary}
            onClick={() => router.back()}
          >
            back
          </Button>
        )}

        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading || !formState.isValid}
          onClick={handleSubmit(submitRegisterAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
