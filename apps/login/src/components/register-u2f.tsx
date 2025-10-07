"use client";

import { coerceToArrayBuffer, coerceToBase64Url } from "@/helpers/base64";
import { completeFlowOrGetUrl } from "@/lib/client";
import { addU2F, verifyU2F } from "@/lib/server/u2f";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { RegisterU2FResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Alert } from "./alert";
import { BackButton } from "./back-button";
import { Button, ButtonVariants } from "./button";
import { Spinner } from "./spinner";
import { Translated } from "./translated";

type Props = {
  loginName?: string;
  sessionId: string;
  requestId?: string;
  organization?: string;
  checkAfter: boolean;
  loginSettings?: LoginSettings;
};

export function RegisterU2f({ loginName, sessionId, organization, requestId, checkAfter, loginSettings }: Props) {
  const [error, setError] = useState<string>("");

  const [loading, setLoading] = useState<boolean>(false);

  const router = useRouter();

  async function submitVerify(u2fId: string, passkeyName: string, publicKeyCredential: any, sessionId: string) {
    setError("");
    setLoading(true);
    const response = await verifyU2F({
      u2fId,
      passkeyName,
      publicKeyCredential,
      sessionId,
    })
      .catch(() => {
        setError("An error on verifying passkey occurred");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response?.error) {
      setError(response?.error);
      return;
    }

    return response;
  }

  async function submitRegisterAndContinue(): Promise<boolean | void | null> {
    setError("");
    setLoading(true);
    const response = await addU2F({
      sessionId,
    })
      .catch(() => {
        setError("An error on registering passkey");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response?.error) {
      setError(response?.error);
      return;
    }

    if (!response || !("u2fId" in response)) {
      setError("An error on registering passkey");
      return;
    }

    const u2fResponse = response as unknown as RegisterU2FResponse;

    const u2fId = u2fResponse.u2fId;
    const options: CredentialCreationOptions =
      (u2fResponse?.publicKeyCredentialCreationOptions as CredentialCreationOptions) ?? {};

    if (options.publicKey) {
      options.publicKey.challenge = coerceToArrayBuffer(options.publicKey.challenge, "challenge");
      options.publicKey.user.id = coerceToArrayBuffer(options.publicKey.user.id, "userid");
      if (options.publicKey.excludeCredentials) {
        options.publicKey.excludeCredentials.map((cred: any) => {
          cred.id = coerceToArrayBuffer(cred.id as string, "excludeCredentials.id");
          return cred;
        });
      }

      const resp = await navigator.credentials.create(options);

      if (
        !resp ||
        !(resp as any).response.attestationObject ||
        !(resp as any).response.clientDataJSON ||
        !(resp as any).rawId
      ) {
        setError("An error on registering passkey");
        return;
      }

      const attestationObject = (resp as any).response.attestationObject;
      const clientDataJSON = (resp as any).response.clientDataJSON;
      const rawId = (resp as any).rawId;

      const data = {
        id: resp.id,
        rawId: coerceToBase64Url(rawId, "rawId"),
        type: resp.type,
        response: {
          attestationObject: coerceToBase64Url(attestationObject, "attestationObject"),
          clientDataJSON: coerceToBase64Url(clientDataJSON, "clientDataJSON"),
        },
      };

      const submitResponse = await submitVerify(u2fId, "", data, sessionId);

      if (!submitResponse) {
        setError("An error on verifying passkey");
        return;
      }

      if (checkAfter) {
        const paramsToContinue = new URLSearchParams({});

        if (sessionId) {
          paramsToContinue.append("sessionId", sessionId);
        }
        if (loginName) {
          paramsToContinue.append("loginName", loginName);
        }
        if (organization) {
          paramsToContinue.append("organization", organization);
        }
        if (requestId) {
          paramsToContinue.append("requestId", requestId);
        }

        return router.push(`/u2f?` + paramsToContinue);
      } else {
        if (requestId && sessionId) {
          const callbackResponse = await completeFlowOrGetUrl(
            {
              sessionId: sessionId,
              requestId: requestId,
              organization: organization,
            },
            loginSettings?.defaultRedirectUri,
          );

          if ("error" in callbackResponse) {
            setError(callbackResponse.error);
            return;
          }

          if ("redirect" in callbackResponse) {
            return router.push(callbackResponse.redirect);
          }
        } else if (loginName) {
          const callbackResponse = await completeFlowOrGetUrl(
            {
              loginName: loginName,
              organization: organization,
            },
            loginSettings?.defaultRedirectUri,
          );

          if ("error" in callbackResponse) {
            setError(callbackResponse.error);
            return;
          }

          if ("redirect" in callbackResponse) {
            return router.push(callbackResponse.redirect);
          }
        }
      }
    }
  }

  return (
    <form className="w-full">
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}

      <div className="mt-8 flex w-full flex-row items-center">
        <BackButton data-testid="back-button" />

        <span className="flex-grow"></span>
        <Button
          type="submit"
          className="self-end"
          variant={ButtonVariants.Primary}
          disabled={loading}
          onClick={submitRegisterAndContinue}
          data-testid="submit-button"
        >
          {loading && <Spinner className="mr-2 h-5 w-5" />} <Translated i18nKey="set.submit" namespace="u2f" />
        </Button>
      </div>
    </form>
  );
}
