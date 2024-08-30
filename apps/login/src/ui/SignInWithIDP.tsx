"use client";

import { ReactNode, useState } from "react";
import {
  SignInWithGitlab,
  SignInWithAzureAD,
  SignInWithGoogle,
  SignInWithGithub,
} from "@zitadel/react";
import { useRouter } from "next/navigation";
import Alert from "./Alert";
import { IdentityProvider } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { idpTypeToSlug } from "@/lib/idp";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { startIDPFlow } from "@/lib/server/idp";

export interface SignInWithIDPProps {
  children?: ReactNode;
  host: string;
  identityProviders: IdentityProvider[];
  authRequestId?: string;
  organization?: string;
}

export function SignInWithIDP({
  host,
  identityProviders,
  authRequestId,
  organization,
}: SignInWithIDPProps) {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const router = useRouter();

  async function startFlow(idpId: string, provider: string) {
    setLoading(true);

    const params = new URLSearchParams();

    if (authRequestId) {
      params.set("authRequestId", authRequestId);
    }

    if (organization) {
      params.set("organization", organization);
    }

    const response = await startIDPFlow({
      idpId,
      successUrl:
        `${host}/idp/${provider}/success?` + new URLSearchParams(params),
      failureUrl:
        `${host}/idp/${provider}/failure?` + new URLSearchParams(params),
    }).catch((error: Error) => {
      setError(error.message ?? "Could not start IDP flow");
    });

    setLoading(false);

    return response;
  }

  return (
    <div className="flex flex-col w-full space-y-2 text-sm">
      {identityProviders &&
        identityProviders.map((idp, i) => {
          switch (idp.type) {
            case IdentityProviderType.GITHUB:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  onClick={() =>
                    startFlow(
                      idp.id,
                      idpTypeToSlug(IdentityProviderType.GITHUB),
                    ).then(({ authUrl }) => {
                      router.push(authUrl);
                    })
                  }
                ></SignInWithGithub>
              );
            case IdentityProviderType.GITHUB_ES:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  onClick={() => alert("TODO: unimplemented")}
                ></SignInWithGithub>
              );
            case IdentityProviderType.AZURE_AD:
              return (
                <SignInWithAzureAD
                  key={`idp-${i}`}
                  onClick={() =>
                    startFlow(
                      idp.id,
                      idpTypeToSlug(IdentityProviderType.AZURE_AD),
                    ).then(({ authUrl }) => {
                      router.push(authUrl);
                    })
                  }
                ></SignInWithAzureAD>
              );
            case IdentityProviderType.GOOGLE:
              return (
                <SignInWithGoogle
                  key={`idp-${i}`}
                  e2e="google"
                  name={idp.name}
                  onClick={() =>
                    startFlow(
                      idp.id,
                      idpTypeToSlug(IdentityProviderType.GOOGLE),
                    ).then(({ authUrl }) => {
                      router.push(authUrl);
                    })
                  }
                ></SignInWithGoogle>
              );
            case IdentityProviderType.GITLAB:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  onClick={() => alert("TODO: unimplemented")}
                ></SignInWithGitlab>
              );
            case IdentityProviderType.GITLAB_SELF_HOSTED:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  onClick={() => alert("TODO: unimplemented")}
                ></SignInWithGitlab>
              );
            default:
              return null;
          }
        })}
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}
    </div>
  );
}

SignInWithIDP.displayName = "SignInWithIDP";
