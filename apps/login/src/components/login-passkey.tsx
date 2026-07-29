"use client";

import { coerceToArrayBuffer, coerceToBase64Url } from "@/helpers/base64";
import { handleServerActionResponse } from "@/lib/client-utils";
import { sendPasskey } from "@/lib/server/passkeys";
import { updateOrCreateSession } from "@/lib/server/session";
import { create, JsonObject } from "@zitadel/client";
import { RequestChallengesSchema, UserVerificationRequirement } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Alert } from "./alert";
import { AutoSubmitForm } from "./auto-submit-form";
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

export function LoginPasskey({ loginName, sessionId, requestId, altPassword, organization, login = true }: Props) {
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const t = useTranslations("passkey");
  const router = useRouter();

  const abortControllerRef = useRef<AbortController | null>(null);

  // Runs the full passkey ceremony: request a challenge, get the assertion from the
  // authenticator, then verify it. A single AbortController guards against overlapping
  // ceremonies — starting a new one aborts the WebAuthn request still in flight — so we
  // never sign an assertion against a challenge that a second, overlapping request has
  // already overwritten on the session. See https://github.com/zitadel/zitadel/issues/12495.
  async function startPasskeyLogin() {
    abortControllerRef.current?.abort();
    const abortController = new AbortController();
    abortControllerRef.current = abortController;

    setError("");
    setLoading(true);
    try {
      const response = await updateOrCreateSessionForChallenge();
      if (abortController.signal.aborted) {
        return;
      }

      const pK = response?.challenges?.webAuthN?.publicKeyCredentialRequestOptions?.publicKey;
      if (!pK) {
        setError(t("verify.errors.couldNotRequestChallenge"));
        return;
      }

      await submitLoginAndContinue(pK, abortController.signal);
    } catch (error) {
      if (!abortController.signal.aborted) {
        setError(error instanceof Error ? error.message : String(error));
      }
    } finally {
      // Only the most recent ceremony owns the shared loading state.
      if (abortControllerRef.current === abortController) {
        setLoading(false);
      }
    }
  }

  useEffect(() => {
    startPasskeyLogin();

    // Abort an in-flight WebAuthn request if the component unmounts (e.g. an RSC
    // navigation remounts the page) so it cannot resolve against a stale challenge.
    return () => {
      abortControllerRef.current?.abort();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function updateOrCreateSessionForChallenge(
    userVerificationRequirement: number = login
      ? UserVerificationRequirement.REQUIRED
      : UserVerificationRequirement.DISCOURAGED,
  ) {
    const sessionResponse = await updateOrCreateSession({
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
    }).catch((error) => {
      console.error(error);
      setError(t("verify.errors.couldNotRequestChallenge"));
      return;
    });

    if (sessionResponse && "error" in sessionResponse && sessionResponse.error) {
      setError(sessionResponse.error);
      return;
    }

    return sessionResponse;
  }

  async function submitLogin(data: JsonObject) {
    setLoading(true);
    try {
      const response = await sendPasskey({
        loginName,
        sessionId,
        organization,
        checks: {
          webAuthN: { credentialAssertionData: data },
        } as Checks,
        requestId,
      });

      const handled = handleServerActionResponse(response, router, setSamlData, setError);

      if (!handled) {
        if (!response) {
          setError(t("verify.errors.noResponseReceived"));
        } else {
          setError(t("verify.errors.noRedirectProvided"));
        }
      }
    } catch {
      setError(t("verify.errors.couldNotVerifyPasskey"));
    } finally {
      setLoading(false);
    }
  }

  async function submitLoginAndContinue(publicKey: any, signal?: AbortSignal): Promise<boolean | void> {
    publicKey.challenge = coerceToArrayBuffer(publicKey.challenge, "publicKey.challenge");
    publicKey.allowCredentials.map((listItem: any) => {
      listItem.id = coerceToArrayBuffer(listItem.id, "publicKey.allowCredentials.id");
    });

    return navigator.credentials
      .get({
        publicKey,
        signal,
      })
      .then((assertedCredential: any) => {
        if (!assertedCredential) {
          setError(t("verify.errors.couldNotRetrievePasskey"));
          return;
        }

        const authData = new Uint8Array(assertedCredential.response.authenticatorData);
        const clientDataJSON = new Uint8Array(assertedCredential.response.clientDataJSON);
        const rawId = new Uint8Array(assertedCredential.rawId);
        const sig = new Uint8Array(assertedCredential.response.signature);
        const userHandle = new Uint8Array(assertedCredential.response.userHandle);
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
      .catch((error) => {
        // A superseded or unmounted ceremony rejects with AbortError — that is
        // expected (we aborted it), not a verification failure to surface.
        if (error?.name === "AbortError") {
          return;
        }
        // Handle passkey cancellation or errors
        if (error?.name === "NotAllowedError") {
          setError(t("verify.errors.verificationCancelled"));
        } else {
          setError(t("verify.errors.verificationFailed"));
        }
        console.error("Passkey verification error:", error);
      });
  }

  return (
    <div className="w-full">
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
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
                "/password?" + params, // alt is set because password is requested as alternative auth method, so passkey prompt can be escaped
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
          onClick={() => startPasskeyLogin()}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />} <Translated i18nKey="verify.submit" namespace="passkey" />
        </Button>
      </div>
    </div>
  );
}
