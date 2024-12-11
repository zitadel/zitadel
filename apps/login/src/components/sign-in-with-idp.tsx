"use client";

import { idpTypeToSlug } from "@/lib/idp";
import { startIDPFlow } from "@/lib/server/idp";
import {
  IdentityProvider,
  IdentityProviderType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { ReactNode, useState } from "react";
import { Alert } from "./alert";
import { SignInWithApple } from "./idps/sign-in-with-apple";
import { SignInWithAzureAd } from "./idps/sign-in-with-azure-ad";
import { SignInWithGeneric } from "./idps/sign-in-with-generic";
import { SignInWithGithub } from "./idps/sign-in-with-github";
import { SignInWithGitlab } from "./idps/sign-in-with-gitlab";
import { SignInWithGoogle } from "./idps/sign-in-with-google";

export interface SignInWithIDPProps {
  children?: ReactNode;
  identityProviders: IdentityProvider[];
  authRequestId?: string;
  organization?: string;
  linkOnly?: boolean;
}

export function SignInWithIdp({
  identityProviders,
  authRequestId,
  organization,
  linkOnly,
}: SignInWithIDPProps) {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const router = useRouter();

  async function startFlow(idpId: string, provider: string) {
    setLoading(true);

    const params = new URLSearchParams();

    if (linkOnly) {
      params.set("link", "true");
    }

    if (authRequestId) {
      params.set("authRequestId", authRequestId);
    }

    if (organization) {
      params.set("organization", organization);
    }

    const response = await startIDPFlow({
      idpId,
      successUrl: `/idp/${provider}/success?` + new URLSearchParams(params),
      failureUrl: `/idp/${provider}/failure?` + new URLSearchParams(params),
    })
      .catch(() => {
        setError("Could not start IDP flow");
        return;
      })
      .finally(() => {
        setLoading(false);
      });

    if (response && "error" in response && response?.error) {
      setError(response.error);
      return;
    }

    if (response && "redirect" in response && response?.redirect) {
      return router.push(response.redirect);
    }
  }

  return (
    <div className="flex flex-col w-full space-y-2 text-sm">
      {identityProviders &&
        identityProviders
          /* - TODO: Implement after https://github.com/zitadel/zitadel/issues/8981  */

          //   .filter((idp) =>
          //     linkOnly ? idp.config?.options.isLinkingAllowed : true,
          //   )
          .map((idp, i) => {
            switch (idp.type) {
              case IdentityProviderType.APPLE:
                return (
                  <SignInWithApple
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.APPLE),
                      )
                    }
                  ></SignInWithApple>
                );
              case IdentityProviderType.OAUTH:
                return (
                  <SignInWithGeneric
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.OAUTH),
                      )
                    }
                  ></SignInWithGeneric>
                );
              case IdentityProviderType.OIDC:
                return (
                  <SignInWithGeneric
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.OIDC),
                      )
                    }
                  ></SignInWithGeneric>
                );
              case IdentityProviderType.GITHUB:
                return (
                  <SignInWithGithub
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.GITHUB),
                      )
                    }
                  ></SignInWithGithub>
                );
              case IdentityProviderType.GITHUB_ES:
                return (
                  <SignInWithGithub
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.GITHUB_ES),
                      )
                    }
                  ></SignInWithGithub>
                );
              case IdentityProviderType.AZURE_AD:
                return (
                  <SignInWithAzureAd
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.AZURE_AD),
                      )
                    }
                  ></SignInWithAzureAd>
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
                      )
                    }
                  ></SignInWithGoogle>
                );
              case IdentityProviderType.GITLAB:
                return (
                  <SignInWithGitlab
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.GITLAB),
                      )
                    }
                  ></SignInWithGitlab>
                );
              case IdentityProviderType.GITLAB_SELF_HOSTED:
                return (
                  <SignInWithGitlab
                    key={`idp-${i}`}
                    name={idp.name}
                    onClick={() =>
                      startFlow(
                        idp.id,
                        idpTypeToSlug(IdentityProviderType.GITLAB_SELF_HOSTED),
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

SignInWithIdp.displayName = "SignInWithIDP";
