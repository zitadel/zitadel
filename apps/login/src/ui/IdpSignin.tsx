"use client";

import { useEffect, useState } from "react";
import { Spinner } from "./Spinner";
import Alert from "./Alert";
import { useRouter } from "next/navigation";

type Props = {
  userId: string;
  // organization: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  authRequestId?: string;
};

export default function IdpSignin(props: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const router = useRouter();

  async function createSessionForIdp() {
    setLoading(true);
    const res = await fetch("/api/session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        userId: props.userId,
        idpIntent: props.idpIntent,
        authRequestId: props.authRequestId,
        // organization: props.organization,
      }),
    });

    if (!res.ok) {
      const error = await res.json();
      throw error.details.details;
    }
    return res.json();
  }

  useEffect(() => {
    createSessionForIdp()
      .then((session) => {
        setLoading(false);
        if (props.authRequestId && session && session.sessionId) {
          return router.push(
            `/login?` +
              new URLSearchParams({
                sessionId: session.sessionId,
                authRequest: props.authRequestId,
              }),
          );
        } else {
          return router.push(
            `/signedin?` +
              new URLSearchParams(
                props.authRequestId
                  ? {
                      loginName: session.factors.user.loginName,
                      authRequestId: props.authRequestId,
                    }
                  : {
                      loginName: session.factors.user.loginName,
                    },
              ),
          );
        }
      })
      .catch((error) => {
        setLoading(false);
        setError(error.message);
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
