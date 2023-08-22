"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { coerceToArrayBuffer, coerceToBase64Url } from "#/utils/base64";
import { Button, ButtonVariants } from "./Button";
import Alert from "./Alert";
import { Spinner } from "./Spinner";

type Props = {
  loginName: string;
  authRequestId?: string;
  altPassword: boolean;
};

export default function LoginPasskey({
  loginName,
  authRequestId,
  altPassword,
}: Props) {
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  const initialized = useRef(false);

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
      setLoading(true);
      updateSessionForChallenge()
        .then((response) => {
          console.log(response);
          const pK =
            response.challenges.passkey.publicKeyCredentialRequestOptions
              .publicKey;
          if (pK) {
            submitLoginAndContinue(pK)
              .then(() => {
                setLoading(false);
              })
              .catch((error) => {
                setError(error);
                setLoading(false);
              });
          } else {
            setError("Could not request passkey challenge");
            setLoading(false);
          }
        })
        .catch((error) => {
          setError(error);
          setLoading(false);
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
        challenges: {
          webAuthN: {
            domain: "",
            userVerificationRequirement: 2,
          },
        },
        authRequestId,
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
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        passkey: data,
        authRequestId,
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

  async function submitLoginAndContinue(
    publicKey: any
  ): Promise<boolean | void> {
    publicKey.challenge = coerceToArrayBuffer(
      publicKey.challenge,
      "publicKey.challenge"
    );
    publicKey.allowCredentials.map((listItem: any) => {
      listItem.id = coerceToArrayBuffer(
        listItem.id,
        "publicKey.allowCredentials.id"
      );
    });

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
          return submitLogin(data).then(() => {
            return router.push(`/accounts`);
          });
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

  return (
    <div className="w-full">
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}
      <div className="mt-8 flex w-full flex-row items-center">
        {altPassword ? (
          <Button
            type="button"
            variant={ButtonVariants.Secondary}
            onClick={() => {
              const params = { loginName, alt: "true" };

              return router.push(
                "/password?" +
                  new URLSearchParams(
                    authRequestId ? { ...params, authRequestId } : params
                  ) // alt is set because password is requested as alternative auth method, so passwordless prompt can be escaped
              );
            }}
          >
            use password
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
          disabled={loading}
          onClick={() => updateSessionForChallenge()}
        >
          {loading && <Spinner className="h-5 w-5 mr-2" />}
          continue
        </Button>
      </div>
    </div>
  );
}
