"use client";

import { processIDPCallback } from "@/lib/actions/idp-intent";
import { submitSamlForm } from "@/lib/client-utils";
import { useSettledAction } from "@/lib/settled-action";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Alert } from "./alert";
import { Spinner } from "./spinner";

type Props = {
  provider: string;
  id: string;
  token: string;
  requestId?: string;
  organization?: string;
  link?: string;
  sessionId?: string;
  flowFingerprint?: string;
  postErrorRedirectUrl?: string;
};

/**
 * Client component that handles IDP callback processing.
 * Must be client-side to allow cookie modifications via server actions.
 */
export function IdpProcessHandler({
  provider,
  id,
  token,
  requestId,
  organization,
  link,
  sessionId,
  flowFingerprint,
  postErrorRedirectUrl,
}: Props) {
  const t = useTranslations("idp");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const executedRef = useRef(false);
  const router = useRouter();
  const runAction = useSettledAction();

  useEffect(() => {
    // Prevent double execution in React Strict Mode
    if (executedRef.current) {
      return;
    }

    executedRef.current = true;

    console.log("[IDP Process Handler] Starting IDP callback processing from client");

    runAction(() =>
      processIDPCallback({
        provider,
        id,
        token,
        requestId,
        organization,
        sessionId,
        flowFingerprint,
        postErrorRedirectUrl,
      }),
    )
      .then((result) => {
        if (result.error) {
          console.error("[IDP Process Handler] Error:", result.error);
          setError(result.error);
          setLoading(false);
          return;
        }

        if (result.redirect) {
          console.log("[IDP Process Handler] Redirecting to:", result.redirect);
          router.push(result.redirect);
          return;
        }

        if (result.samlData) {
          console.log("[IDP Process Handler] Received samlData, posting it to the service provider");
          submitSamlForm(result.samlData);
          setLoading(false);
          return;
        }

        setError(t("processing.noRedirect"));
        setLoading(false);
      })
      .catch((err) => {
        console.error("[IDP Process Handler] Unexpected error:", err);
        setError(err instanceof Error ? err.message : t("processing.unexpectedError"));
        setLoading(false);
      });
  }, [
    provider,
    id,
    token,
    requestId,
    organization,
    link,
    sessionId,
    flowFingerprint,
    postErrorRedirectUrl,
    router,
    runAction,
    t,
  ]);

  return (
    <div className="flex items-center justify-center">
      {loading && (
        <div className="flex flex-col items-center space-y-4">
          <Spinner className="h-8 w-8" />
          <p className="text-sm text-gray-700 dark:text-gray-300">{t("processing.message")}</p>
        </div>
      )}
      {error && (
        <div className="max-w-md py-4">
          <Alert>{error}</Alert>
        </div>
      )}
    </div>
  );
}
