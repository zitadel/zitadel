"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { coerceToArrayBuffer, coerceToBase64Url } from "@/utils/base64";
import { Button, ButtonVariants } from "./Button";
import Alert from "./Alert";
import { Spinner } from "./Spinner";
import { Checks } from "@zitadel/proto/zitadel/session/v2beta/session_service_pb";

// either loginName or sessionId must be provided
type Props = {
  loginName?: string;
  sessionId?: string;
  authRequestId?: string;
  altPassword: boolean;
  login?: boolean;
  organization?: string;
};

export default function LoginPasskey({
  loginName,
  sessionId,
  authRequestId,
  altPassword,
  organization,
  login = true,
}: Props) {
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  const initialized = useRef(false);

  // TODO: move this to server side
  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
      setLoading(true);
      updateSessionForChallenge()
        .then((response) => {
          const pK =
            response.challenges.webAuthN.publicKeyCredentialRequestOptions
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

  async function updateSessionForChallenge(
    userVerificationRequirement: number = login ? 1 : 3,
  ) {
    setLoading(true);
    const res = await fetch("/api/session", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        sessionId,
        organization,
        challenges: {
          webAuthN: {
            domain: "",
            // USER_VERIFICATION_REQUIREMENT_UNSPECIFIED = 0;
            // USER_VERIFICATION_REQUIREMENT_REQUIRED = 1; - passkey login
            // USER_VERIFICATION_REQUIREMENT_PREFERRED = 2;
            // USER_VERIFICATION_REQUIREMENT_DISCOURAGED = 3; - mfa
            userVerificationRequirement: userVerificationRequirement,
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
        sessionId,
        organization,
        checks: {
          webAuthN: { credentialAssertionData: data },
        } as Checks,
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
        if (assertedCredential) {
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
              clientDataJSON: coerceToBase64Url(
                clientDataJSON,
                "clientDataJSON",
              ),
              signature: coerceToBase64Url(sig, "sig"),
              userHandle: coerceToBase64Url(userHandle, "userHandle"),
            },
          };
          return submitLogin(data).then((resp) => {
            if (authRequestId && resp && resp.sessionId) {
              return router.push(
                `/login?` +
                  new URLSearchParams({
                    sessionId: resp.sessionId,
                    authRequest: authRequestId,
                  }),
              );
            } else {
              return router.push(
                `/signedin?` +
                  new URLSearchParams(
                    authRequestId
                      ? {
                          loginName: resp.factors.user.loginName,
                          authRequestId,
                        }
                      : {
                          loginName: resp.factors.user.loginName,
                        },
                  ),
              );
            }
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
              const params: any = { alt: "true" };

              if (loginName) {
                params.loginName = loginName;
              }

              if (sessionId) {
                params.sessionId = sessionId;
              }

              if (authRequestId) {
                params.authRequestId = authRequestId;
              }

              if (organization) {
                params.organization = organization;
              }

              return router.push(
                "/password?" + new URLSearchParams(params), // alt is set because password is requested as alternative auth method, so passwordless prompt can be escaped
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
