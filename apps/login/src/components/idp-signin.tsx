"use client";

import { createNewSessionForIdp } from "@/lib/server/session";
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
  authRequestId?: string;
};

export function IdpSignin({
  userId,
  idpIntent: { idpIntentId, idpIntentToken },
  authRequestId,
}: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const router = useRouter();

  useEffect(() => {
    createNewSessionForIdp({
      userId,
      idpIntent: {
        idpIntentId,
        idpIntentToken,
      },
      authRequestId,
    })
      .catch((error) => {
        setError(error.message);
        return;
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return (
    <div className="flex items-center justify-center">
      {loading && <Spinner />}
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}
    </div>
  );
}
