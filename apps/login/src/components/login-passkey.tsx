"use client";

import { coerceToArrayBuffer, coerceToBase64Url } from "@/helpers/base64";
import { sendPasskey } from "@/lib/server/passkeys";
import { updateSession } from "@/lib/server/session";
import { create, JsonObject } from "@zitadel/client";
import {
  RequestChallengesSchema,
  UserVerificationRequirement,
} from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

// either loginName or sessionId must be provided
type Props = {
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  altPassword: boolean;
  login?: boolean;
  organization?: string;
};

export function LoginPasskey({
  loginName,
  sessionId,
  requestId,
  altPassword,
  organization,
  login = true,
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
          const pK =
            response?.challenges?.webAuthN?.publicKeyCredentialRequestOptions
              ?.publicKey;

          if (!pK) {
            setError("Could not request passkey challenge");
            setLoading(false);
            return;
          }

          return submitLoginAndContinue(pK)
            .catch((error) => {
              setError(error);
              return;
            })
            .finally(() => {
              setLoading(false);
            });
        })
        .catch((error) => {
          setError(error);
          return;
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, []);

  async function updateSessionForChallenge(
    userVerificationRequirement: number = login
      ? UserVerificationRequirement.REQUIRED
      : UserVerificationRequirement.DISCOURAGED,
  ) {
    setError("");
    setLoading(true);
    const session = await updateSession({
      loginName,
      sessionId,
      organization,
      challenges: create(RequestChallengesSchema, {
        webAuthN: {
          domain: "",
          userVerificationRequirement,
        },
      }),
      requestId,
    })
      .catch(() => {
        setError("Could not request passkey challenge");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (session && "error" in session && session.error) {
      setError(session.error);
      return;
    }

    return session;
  }

  async function submitLogin(data: JsonObject) {
    setLoading(true);
    const response = await sendPasskey({
      loginName,
      sessionId,
      organization,
      checks: {
        webAuthN: { credentialAssertionData: data },
      } as Checks,
      requestId,
    })
      .catch(() => {
        setError("Could not verify passkey");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response.error) {
      setError(response.error);
      return;
    }

    if (response && "redirect" in response && response.redirect) {
      return router.push(response.redirect);
    }
  }

  async function submitLoginAndContinue(
    publicKey: any,
  ): Promise<boolean | void> {
    publicKey.challenge = coerceToArrayBuffer(
      publicKey.challenge,
      "publicKey.challenge",
    );
    publicKey.allowCredentials.map((listItem: any) => {
      listItem.id = coerceToArrayBuffer(
        listItem.id,
        "publicKey.allowCredentials.id",
      );
    });

    navigator.credentials
      .get({
        publicKey,
      })
      .then((assertedCredential: any) => {
        if (!assertedCredential) {
          setError("An error on retrieving passkey");
          return;
        }

        const authData = new Uint8Array(
          assertedCredential.response.authenticatorData,
        );
        const clientDataJSON = new Uint8Array(
          assertedCredential.response.clientDataJSON,
        );
        const rawId = new Uint8Array(assertedCredential.rawId);
        const sig = new Uint8Array(assertedCredential.response.signature);
        const userHandle = new Uint8Array(
          assertedCredential.response.userHandle,
        );
        const data = {
          id: assertedCredential.id,
          rawId: coerceToBase64Url(rawId, "rawId"),
          type: assertedCredential.type,
          response: {
            authenticatorData: coerceToBase64Url(authData, "authData"),
            clientDataJSON: coerceToBase64Url(clientDataJSON, "clientDataJSON"),
            signature: coerceToBase64Url(sig, "sig"),
            userHandle: coerceToBase64Url(userHandle, "userHandle"),
          },
        };

        return submitLogin(data);
      })
      .finally(() => {
        setLoading(false);
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
              const params = new URLSearchParams();

              if (loginName) {
                params.append("loginName", loginName);
              }

              if (sessionId) {
                params.append("sessionId", sessionId);
              }

              if (requestId) {
                params.append("requestId", requestId);
              }

              if (organization) {
                params.append("organization", organization);
              }

              return router.push(
                "/password?" + params, // alt is set because password is requested as alternative auth method, so passwordless prompt can be escaped
              );
            }}
            data-testid="password-button"
          >
            <Translated i18nKey="verify.usePassword" namespace="passkey" />
          </Button>
        ) : (
          <BackButton />
        )}

        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading}
          onClick={async () => {
            const response = await updateSessionForChallenge().finally(() => {
              setLoading(false);
            });

            const pK =
              response?.challenges?.webAuthN?.publicKeyCredentialRequestOptions
                ?.publicKey;

            if (!pK) {
              setError("Could not request passkey challenge");
              return;
            }

            setLoading(true);

            return submitLoginAndContinue(pK)
              .catch((error) => {
                setError(error);
                return;
              })
              .finally(() => {
                setLoading(false);
              });
          }}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />}{" "}
          <Translated i18nKey="verify.submit" namespace="passkey" />
        </Button>
      </div>
    </div>
  );
}
