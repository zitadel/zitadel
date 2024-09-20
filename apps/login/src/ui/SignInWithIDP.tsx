"use client";

import { idpTypeToSlug } from "@/lib/idp";
import { startIDPFlow } from "@/lib/server/idp";
import {
  IdentityProvider,
  IdentityProviderType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { ReactNode, useState } from "react";
import Alert from "./Alert";
import { SignInWithApple } from "./idps/SignInWithApple";
import { SignInWithAzureAD } from "./idps/SignInWithAzureAD";
import { SignInWithGeneric } from "./idps/SignInWithGeneric";
import { SignInWithGithub } from "./idps/SignInWithGithub";
import { SignInWithGitlab } from "./idps/SignInWithGitlab";
import { SignInWithGoogle } from "./idps/SignInWithGoogle";

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
    }).catch(() => {
      setError("Could not start IDP flow");
    });

    setLoading(false);

    return response;
  }

  async function navigateToAuthUrl(id: string, type: IdentityProviderType) {
    const startFlowResponse = await startFlow(id, idpTypeToSlug(type));
    if (
      startFlowResponse &&
      startFlowResponse.nextStep.case === "authUrl" &&
      startFlowResponse?.nextStep.value
    ) {
      router.push(startFlowResponse.nextStep.value);
    }
  }

  return (
    <div className="flex flex-col w-full space-y-2 text-sm">
      {identityProviders &&
        identityProviders.map((idp, i) => {
          switch (idp.type) {
            case IdentityProviderType.APPLE:
              return (
                <SignInWithApple
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.APPLE)
                  }
                ></SignInWithApple>
              );
            case IdentityProviderType.OAUTH:
              return (
                <SignInWithGeneric
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.OAUTH)
                  }
                ></SignInWithGeneric>
              );
            case IdentityProviderType.OIDC:
              return (
                <SignInWithGeneric
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.OIDC)
                  }
                ></SignInWithGeneric>
              );
            case IdentityProviderType.GITHUB:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.GITHUB)
                  }
                ></SignInWithGithub>
              );
            case IdentityProviderType.GITHUB_ES:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.GITHUB_ES)
                  }
                ></SignInWithGithub>
              );
            case IdentityProviderType.AZURE_AD:
              return (
                <SignInWithAzureAD
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.AZURE_AD)
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
                    navigateToAuthUrl(idp.id, IdentityProviderType.GOOGLE)
                  }
                ></SignInWithGoogle>
              );
            case IdentityProviderType.GITLAB:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(idp.id, IdentityProviderType.GITLAB)
                  }
                ></SignInWithGitlab>
              );
            case IdentityProviderType.GITLAB_SELF_HOSTED:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    navigateToAuthUrl(
                      idp.id,
                      IdentityProviderType.GITLAB_SELF_HOSTED,
                    )
                  }
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
