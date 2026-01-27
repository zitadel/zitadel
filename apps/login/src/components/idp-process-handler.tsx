"use client";

import { processIDPCallback } from "@/lib/server/idp-intent";
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
  linkFingerprint?: string;
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
  linkFingerprint,
  postErrorRedirectUrl,
}: Props) {
  const t = useTranslations("idp");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const executedRef = useRef(false);
  const router = useRouter();

  useEffect(() => {
    // Prevent double execution in React Strict Mode
    if (executedRef.current) {
      return;
    }

    executedRef.current = true;

    console.log("[IDP Process Handler] Starting IDP callback processing from client");

    processIDPCallback({
      provider,
      id,
      token,
      requestId,
      organization,
      sessionId,
      linkFingerprint,
      postErrorRedirectUrl,
    })
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

        setError(t("processing.noRedirect"));
        setLoading(false);
      })
      .catch((err) => {
        console.error("[IDP Process Handler] Unexpected error:", err);
        setError(err instanceof Error ? err.message : t("processing.unexpectedError"));
        setLoading(false);
      });
  }, [provider, id, token, requestId, organization, link, sessionId, postErrorRedirectUrl, router]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      {loading && (
        <div className="flex flex-col items-center space-y-4">
          <Spinner className="h-8 w-8" />
          <p className="text-sm text-gray-600">{t("processing.message")}</p>
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
