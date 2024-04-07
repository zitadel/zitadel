"use client";

import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import { coerceToArrayBuffer, coerceToBase64Url } from "@/utils/base64";
import { RegisterU2FResponse } from "@zitadel/proto/zitadel/user/v2beta/user_service_pb";
type Inputs = {};

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
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitRegister() {
    setError("");
    setLoading(true);
    const res = await fetch("/api/u2f", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
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

  async function submitVerify(
    u2fId: string,
    passkeyName: string,
    publicKeyCredential: any,
    sessionId: string,
  ) {
    setLoading(true);
    const res = await fetch("/api/u2f/verify", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        u2fId,
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

  function submitRegisterAndContinue(value: Inputs): Promise<boolean | void> {
    return submitRegister().then((resp: RegisterU2FResponse) => {
      const u2fId = resp.u2fId;
      const options: CredentialCreationOptions =
        (resp.publicKeyCredentialCreationOptions?.toJson() as CredentialCreationOptions) ??
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

        navigator.credentials
          .create(options)
          .then((resp) => {
            if (
              resp &&
              (resp as any).response.attestationObject &&
              (resp as any).response.clientDataJSON &&
              (resp as any).rawId
            ) {
              const attestationObject = (resp as any).response
                .attestationObject;
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
                  clientDataJSON: coerceToBase64Url(
                    clientDataJSON,
                    "clientDataJSON",
                  ),
                },
              };
              return submitVerify(u2fId, "", data, sessionId).then(() => {
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
              });
            } else {
              setLoading(false);
              setError("An error on registering passkey");
              return null;
            }
          })
          .catch((error) => {
            console.error(error);
            setLoading(false);
            setError(error);

            return null;
          });
      }
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
