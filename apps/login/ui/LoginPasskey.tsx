"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { ChallengeKind, Challenges_Passkey } from "@zitadel/server";
import { coerceToArrayBuffer, coerceToBase64Url } from "#/utils/base64";
import { Button, ButtonVariants } from "./Button";
import Alert from "./Alert";
import { Spinner } from "./Spinner";

type Props = {
  loginName: string;
  challenge: Challenges_Passkey;
};

export default function LoginPasskey({ loginName, challenge }: Props) {
  const [error, setError] = useState<string>("");
  const [publicKey, setPublicKey] = useState();
  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  const initialized = useRef(false);

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
      updateSessionForChallenge()
        .then((response) => {
          const pK =
            response.challenges.passkey.publicKeyCredentialRequestOptions
              .publicKey;
          if (pK) {
            setPublicKey(pK);
          } else {
            setError("Could not request passkey challenge");
          }
        })
        .catch((error) => {
          setError(error);
        });
    }
  }, []);

  async function updateSessionForChallenge() {
    setLoading(true);
    const res = await fetch("/api/session", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        challenges: [1], // request passkey challenge
      }),
    });

    setLoading(false);
    if (!res.ok) {
      const error = await res.json();
      throw error.details.details;
    }
    return res.json();
  }

  async function submitLogin(data: any) {
    setLoading(true);
    const res = await fetch("/api/session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        passkey: data,
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

  async function submitLoginAndContinue(): Promise<boolean | void> {
    console.log("login", publicKey);
    if (publicKey) {
      console.log(publicKey);
      (publicKey as any).challenge = coerceToArrayBuffer(
        (publicKey as any).challenge,
        "publicKey.challenge"
      );
      (publicKey as any).allowCredentials.map((listItem: any) => {
        listItem.id = coerceToArrayBuffer(
          listItem.id,
          "publicKey.allowCredentials.id"
        );
      });
      console.log(publicKey);
      navigator.credentials
        .get({
          publicKey,
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
            console.log(data);
            return submitLogin(data);
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
  }

  return (
    <div className="w-full">
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
          disabled={loading || !publicKey}
          onClick={() => submitLoginAndContinue()}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </div>
  );
}
