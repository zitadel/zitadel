"use client";

import { createNewSessionFromIdpIntent } from "@/lib/server/idp";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
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

export function IdpSignin({
  userId,
  idpIntent: { idpIntentId, idpIntentToken },
  requestId,
}: Props) {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const router = useRouter();

  useEffect(() => {
    createNewSessionFromIdpIntent({
      userId,
      idpIntent: {
        idpIntentId,
        idpIntentToken,
      },
      requestId,
    })
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
