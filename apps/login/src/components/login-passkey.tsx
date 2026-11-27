"use client";

import { coerceToArrayBuffer, coerceToBase64Url } from "@/helpers/base64";
import { sendPasskey } from "@/lib/actions/passkeys";
import { updateOrCreateSession } from "@/lib/actions/session";
import { handleServerActionResponse } from "@/lib/client-utils";
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
import { RestartSignIn } from "./restart-sign-in";
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
  // the session of this step is gone, retrying cannot succeed
  const [sessionLost, setSessionLost] = useState<boolean>(false);
  // The page starts the passkey ceremony on mount, so Continue starts disabled.
  const [loading, setLoading] = useState<boolean>(true);
  const [samlData, setSamlData] = useState<{ url: string; fields: Record<string, string> } | null>(null);

  const t = useTranslations("passkey");
  const router = useRouter();

  const initialized = useRef(false);
  // At most one ceremony (challenge, navigator.credentials.get, verification) runs at a time.
  // A ref, not state: two taps can land before React commits loading=true.
  const ceremonyInFlight = useRef(false);
  // Aborts the pending navigator.credentials.get; also identifies the current attempt,
  // so results of an aborted attempt are ignored.
  const ceremonyController = useRef<AbortController | null>(null);

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
      runCeremony();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function runCeremony() {
    if (ceremonyInFlight.current) {
      return;
    }
    ceremonyInFlight.current = true;
    const controller = new AbortController();
    ceremonyController.current = controller;
    setError("");
    setSessionLost(false);
    setLoading(true);

    let navigating = false;
    try {
      const publicKey = await requestChallenge();
      if (!publicKey || controller.signal.aborted) {
        return;
      }

      const credential = await getCredential(publicKey, controller.signal);
      if (!credential || controller.signal.aborted) {
        return;
      }

      navigating = await submitLogin(credential);
    } finally {
      if (ceremonyController.current === controller) {
        ceremonyController.current = null;
      }
      // While the page navigates away after a successful verification, keep Continue
      // disabled and the guard set, so no new challenge replaces the verified one.
      if (!navigating) {
        ceremonyInFlight.current = false;
        setLoading(false);
      }
    }
  }

  function abortCeremony() {
    ceremonyController.current?.abort();
    ceremonyController.current = null;
  }

  async function requestChallenge(
    userVerificationRequirement: number = login
      ? UserVerificationRequirement.REQUIRED
      : UserVerificationRequirement.DISCOURAGED,
  ): Promise<any | undefined> {
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
      return undefined;
    });

    if (sessionResponse && "error" in sessionResponse && sessionResponse.error) {
      setSessionLost("sessionLost" in sessionResponse);
      setError(sessionResponse.error);
      return undefined;
    }

    const publicKey = sessionResponse?.challenges?.webAuthN?.publicKeyCredentialRequestOptions?.publicKey;

    if (!publicKey) {
      setError(t("verify.errors.couldNotRequestChallenge"));
      return undefined;
    }

    return publicKey;
  }

  async function getCredential(publicKey: any, signal: AbortSignal): Promise<JsonObject | undefined> {
    try {
      publicKey.challenge = coerceToArrayBuffer(publicKey.challenge, "publicKey.challenge");
      publicKey.allowCredentials.map((listItem: any) => {
        listItem.id = coerceToArrayBuffer(listItem.id, "publicKey.allowCredentials.id");
      });

      const assertedCredential: any = await navigator.credentials.get({
        publicKey,
        signal,
      });

      if (signal.aborted) {
        return undefined;
      }

      if (!assertedCredential) {
        setError(t("verify.errors.couldNotRetrievePasskey"));
        return undefined;
      }

      const authData = new Uint8Array(assertedCredential.response.authenticatorData);
      const clientDataJSON = new Uint8Array(assertedCredential.response.clientDataJSON);
      const rawId = new Uint8Array(assertedCredential.rawId);
      const sig = new Uint8Array(assertedCredential.response.signature);
      const userHandle = new Uint8Array(assertedCredential.response.userHandle);
      return {
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
    } catch (error: any) {
      if (signal.aborted) {
        return undefined;
      }
      // Handle passkey cancellation or errors
      if (error?.name === "NotAllowedError") {
        setError(t("verify.errors.verificationCancelled"));
      } else {
        setError(t("verify.errors.verificationFailed"));
      }
      console.error("Passkey verification error:", error);
      return undefined;
    }
  }

  // Returns true when the page is navigating away (redirect or SAML post).
  async function submitLogin(data: JsonObject): Promise<boolean> {
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

      setSessionLost(!!response && "sessionLost" in response);
      let failed = false;
      const handled = handleServerActionResponse(response, router, setSamlData, (message) => {
        failed = true;
        setError(message);
      });

      if (!handled) {
        if (!response) {
          setError(t("verify.errors.noResponseReceived"));
        } else {
          setError(t("verify.errors.noRedirectProvided"));
        }
        return false;
      }

      return !failed;
    } catch {
      setError(t("verify.errors.couldNotVerifyPasskey"));
      return false;
    }
  }

  return (
    <div className="w-full">
      {samlData && <AutoSubmitForm url={samlData.url} fields={samlData.fields} />}
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
          {sessionLost && <RestartSignIn loginName={loginName} organization={organization} requestId={requestId} />}
        </div>
      )}
      <div className="mt-8 flex w-full flex-row items-center">
        {altPassword ? (
          <Button
            type="button"
            variant={ButtonVariants.Secondary}
            onClick={() => {
              // Cancel the open passkey prompt before switching to the password step.
              abortCeremony();
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
          onClick={() => {
            runCeremony();
          }}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />} <Translated i18nKey="verify.submit" namespace="passkey" />
        </Button>
      </div>
    </div>
  );
}
