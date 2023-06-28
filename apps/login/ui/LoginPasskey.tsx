"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import { Challenges_Passkey } from "@zitadel/server";
import { coerceToArrayBuffer, coerceToBase64Url } from "#/utils/base64";

type Props = {
  challenge: Challenges_Passkey;
};

export default function LoginPasskey({ challenge }: Props) {
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

  function submitLoginAndContinue(): Promise<boolean | void> {
      navigator.credentials
        .get({
          publicKey: challenge.publicKeyCredentialRequestOptions,
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

            return submitLogin(passkeyId, "", data, sessionId);
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
  }

  return (
    <form className="w-full">
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        
          <Button
            type="button"
            variant={ButtonVariants.Secondary}
            onClick={() => router.back()}
          >
            back
          </Button>

        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading}
          onClick={handleSubmit(submitLoginAndContinue)}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </form>
  );
}
