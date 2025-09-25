"use client";

import { CreateNewSessionCommand, createNewSessionFromIdpIntent } from "@/lib/server/idp";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Alert } from "./alert";
import { Spinner } from "./spinner";

type Props = {
  userId: string;
  // organization: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  requestId?: string;
};

export function IdpSignin({ userId, idpIntent: { idpIntentId, idpIntentToken }, requestId }: Props) {
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
    let request: CreateNewSessionCommand = {
      userId,
      idpIntent: {
        idpIntentId,
        idpIntentToken,
      },
    };

    if (requestId) {
      request = { ...request, requestId: requestId };
    }

    createNewSessionFromIdpIntent(request)
      .then((response) => {
        if (response && "error" in response && response?.error) {
          setError(response?.error);
          return;
        }

        if (response && "redirect" in response && response?.redirect) {
          return router.push(response.redirect);
        }
      })
      .catch(() => {
        setError("An internal error occurred");
        return;
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return (
    <div className="flex items-center justify-center py-4">
      {loading && <Spinner className="h-5 w-5" />}
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}
    </div>
  );
}
