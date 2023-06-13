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

export default function RegisterPasskey({ sessionId }: Props) {
  const { register, handleSubmit, formState } = useForm<Inputs>({
    mode: "onBlur",
  });

  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitRegister() {
    // const link = await createPasskeyRegistrationLink(server, userId);
    // console.log(link);
    setError("");
    setLoading(true);
    const res = await fetch("/passkeys", {
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

  function submitRegisterAndContinue(value: Inputs): Promise<boolean | void> {
    return submitRegister().then((resp: RegisterPasskeyResponse) => {
      console.log(resp.publicKeyCredentialCreationOptions?.publicKey);
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
          .create(resp.publicKeyCredentialCreationOptions)
          .then((resp) => {
            console.log(resp);
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

              const data = JSON.stringify({
                id: resp.id,
                rawId: coerceToBase64Url(rawId, "rawId"),
                type: resp.type,
                response: {
                  attestationObject: coerceToBase64Url(
                    attestationObject,
                    "attestationObject"
                  ),
                  clientDataJSON: coerceToBase64Url(
                    clientDataJSON,
                    "clientDataJSON"
                  ),
                },
              });

              const base64 = btoa(data);

              return base64;
              // if (this.type === U2FComponentDestination.MFA) {
              //   this.service
              //     .verifyMyMultiFactorU2F(base64, this.name)
              //     .then(() => {
              //       this.translate
              //         .get("USER.MFA.U2F_SUCCESS")
              //         .pipe(take(1))
              //         .subscribe((msg) => {
              //           this.toast.showInfo(msg);
              //         });
              //       this.dialogRef.close(true);
              //       this.loading = false;
              //     })
              //     .catch((error) => {
              //       this.loading = false;
              //       this.toast.showError(error);
              //     });
              // } else if (this.type === U2FComponentDestination.PASSWORDLESS) {
              //   this.service
              //     .verifyMyPasswordless(base64, this.name)
              //     .then(() => {
              //       this.translate
              //         .get("USER.PASSWORDLESS.U2F_SUCCESS")
              //         .pipe(take(1))
              //         .subscribe((msg) => {
              //           this.toast.showInfo(msg);
              //         });
              //       this.dialogRef.close(true);
              //       this.loading = false;
              //     })
              //     .catch((error) => {
              //       this.loading = false;
              //       this.toast.showError(error);
              //     });
              // }
            } else {
              setLoading(false);
              setError("An error on registering passkey");
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
