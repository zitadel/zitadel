"use client";

import { coerceToArrayBuffer, coerceToBase64Url } from "@/helpers/base64";
import { sendPasskey } from "@/lib/actions/passkeys";
import { updateOrCreateSession } from "@/lib/actions/session";
import { handleServerActionResponse, isExternalUrl } from "@/lib/client-utils";
import { useSettledAction } from "@/lib/settled-action";
import { create, JsonObject } from "@zitadel/client";
import { RequestChallengesSchema, UserVerificationRequirement } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { RestartSignIn } from "./restart-sign-in";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

// How long Continue stays disabled after an external redirect that may not unload the page.
const EXTERNAL_REDIRECT_RELEASE_MS = 5000;

// Shows an error of the running ceremony; sessionLost offers to restart the sign-in.
type ReportError = (message: string, sessionLost?: boolean) => void;

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

  const t = useTranslations("passkey");
  const router = useRouter();
  const runAction = useSettledAction();

  // The running ceremony (challenge, navigator.credentials.get, verification), or null.
  // At most one runs at a time; a ceremony whose controller is no longer current drops
  // its results. A ref, not state: two taps can land before React commits loading=true.
  const ceremony = useRef<AbortController | null>(null);
  const releaseTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  useEffect(() => {
    runCeremony();

    // Safari can restore the page from the back/forward cache after it navigated away,
    // for example when going back from an app's callback; make Continue usable again.
    const onPageShow = (event: PageTransitionEvent) => {
      if (event.persisted) {
        releaseCeremony();
      }
    };
    window.addEventListener("pageshow", onPageShow);

    // Leaving the page any way (Use Password, Back, navigation) aborts the ceremony.
    // Under strict mode the first run is aborted here before it reaches the prompt.
    return () => {
      window.removeEventListener("pageshow", onPageShow);
      abortCeremony();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function runCeremony() {
    if (ceremony.current) {
      return;
    }
    const controller = new AbortController();
    ceremony.current = controller;
    clearTimeout(releaseTimer.current);
    setError("");
    setSessionLost(false);
    setLoading(true);

    const isCurrent = () => ceremony.current === controller;
    // Ignored once this ceremony was aborted or replaced.
    const report: ReportError = (message, sessionLost = false) => {
      if (isCurrent()) {
        setError(message);
        setSessionLost(sessionLost);
      }
    };

    void (async () => {
      let navigating = false;
      try {
        const publicKey = await requestChallenge(report);
        if (!publicKey || !isCurrent()) {
          return;
        }

        const credential = await getCredential(publicKey, controller.signal, report);
        if (!credential || !isCurrent()) {
          return;
        }

        navigating = await submitLogin(credential, isCurrent, report);
      } finally {
        // While the page navigates away after a successful verification, keep Continue
        // disabled and the ceremony current, so no new challenge replaces the verified one.
        if (isCurrent() && !navigating) {
          ceremony.current = null;
          setLoading(false);
        }
      }
    })();
  }

  function abortCeremony() {
    clearTimeout(releaseTimer.current);
    ceremony.current?.abort();
    ceremony.current = null;
  }

  function releaseCeremony() {
    clearTimeout(releaseTimer.current);
    ceremony.current = null;
    setLoading(false);
  }

  async function requestChallenge(report: ReportError): Promise<any> {
    // Called directly, without waiting for the page update: no navigation follows, and the
    // passkey prompt that does must open while the user's activation (Continue) still holds.
    const sessionResponse = await updateOrCreateSession({
      loginName,
      sessionId,
      organization,
      challenges: create(RequestChallengesSchema, {
        webAuthN: {
          domain: "",
          userVerificationRequirement: login
            ? UserVerificationRequirement.REQUIRED
            : UserVerificationRequirement.DISCOURAGED,
        },
      }),
      requestId,
    }).catch((error) => {
      console.error(error);
      return undefined;
    });

    if (sessionResponse && "error" in sessionResponse && sessionResponse.error) {
      report(sessionResponse.error, "sessionLost" in sessionResponse);
      return undefined;
    }

    const publicKey = sessionResponse?.challenges?.webAuthN?.publicKeyCredentialRequestOptions?.publicKey;

    if (!publicKey) {
      report(t("verify.errors.couldNotRequestChallenge"));
      return undefined;
    }

    return publicKey;
  }

  async function getCredential(publicKey: any, signal: AbortSignal, report: ReportError): Promise<JsonObject | undefined> {
    try {
      publicKey.challenge = coerceToArrayBuffer(publicKey.challenge, "publicKey.challenge");
      publicKey.allowCredentials.map((listItem: any) => {
        listItem.id = coerceToArrayBuffer(listItem.id, "publicKey.allowCredentials.id");
      });

      const assertedCredential: any = await navigator.credentials.get({
        publicKey,
        signal,
      });

      if (!assertedCredential) {
        report(t("verify.errors.couldNotRetrievePasskey"));
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
      // Handle passkey cancellation or errors
      if (error?.name === "NotAllowedError") {
        report(t("verify.errors.verificationCancelled"));
      } else {
        report(t("verify.errors.verificationFailed"));
      }
      console.error("Passkey verification error:", error);
      return undefined;
    }
  }

  // Returns true when the page is navigating away (redirect or SAML post).
  async function submitLogin(data: JsonObject, isCurrent: () => boolean, report: ReportError): Promise<boolean> {
    try {
      const response = await runAction(() =>
        sendPasskey({
          loginName,
          sessionId,
          organization,
          checks: {
            webAuthN: { credentialAssertionData: data },
          } as Checks,
          requestId,
        }),
      );

      // The user left for another step while the assertion was being verified.
      if (!isCurrent()) {
        return false;
      }

      let failed = false;
      const handled = handleServerActionResponse(response, router, (message) => {
        failed = true;
        report(message, !!response && "sessionLost" in response);
      });

      if (!handled) {
        report(t(response ? "verify.errors.noRedirectProvided" : "verify.errors.noResponseReceived"));
        return false;
      }

      if (failed) {
        return false;
      }

      // A custom-scheme or other external redirect (a native app's callback) may open
      // another app and leave this page in place, so make Continue usable again later.
      if (response && "redirect" in response && response.redirect && isExternalUrl(response.redirect)) {
        releaseTimer.current = setTimeout(releaseCeremony, EXTERNAL_REDIRECT_RELEASE_MS);
      }

      return true;
    } catch {
      report(t("verify.errors.couldNotVerifyPasskey"));
      return false;
    }
  }

  return (
    <div className="w-full">
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
              // Close the open passkey prompt now, not only once the next page unmounts this one.
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
