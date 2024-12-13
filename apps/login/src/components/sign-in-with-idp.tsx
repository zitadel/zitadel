"use client";

import { idpTypeToSlug } from "@/lib/idp";
import { startIDPFlow } from "@/lib/server/idp";
import {
  IdentityProvider,
  IdentityProviderType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { useRouter } from "next/navigation";
import { ReactNode, useCallback, useState } from "react";
import { Alert } from "./alert";
import { SignInWithIdentityProviderProps } from "./idps/base-button";
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
}: Readonly<SignInWithIDPProps>) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const startFlow = useCallback(
    async (idpId: string, provider: string) => {
      setLoading(true);
      const params = new URLSearchParams();
      if (linkOnly) params.set("link", "true");
      if (authRequestId) params.set("authRequestId", authRequestId);
      if (organization) params.set("organization", organization);

      try {
        const response = await startIDPFlow({
          idpId,
          successUrl: `/idp/${provider}/success?` + params.toString(),
          failureUrl: `/idp/${provider}/failure?` + params.toString(),
        });

        if (response && "error" in response && response?.error) {
          setError(response.error);
          return;
        }

        if (response && "redirect" in response && response?.redirect) {
          return router.push(response.redirect);
        }
      } catch {
        setError("Could not start IDP flow");
      } finally {
        setLoading(false);
      }
    },
    [authRequestId, organization, linkOnly, router],
  );

  const renderIDPButton = (idp: IdentityProvider) => {
    const { id, name, type } = idp;
    const onClick = () => startFlow(id, idpTypeToSlug(type));
    /* - TODO: Implement after https://github.com/zitadel/zitadel/issues/8981  */

    //   .filter((idp) =>
    //     linkOnly ? idp.config?.options.isLinkingAllowed : true,
    //   )
    const components: Partial<
      Record<
        IdentityProviderType,
        (props: SignInWithIdentityProviderProps) => ReactNode
      >
    > = {
      [IdentityProviderType.APPLE]: SignInWithApple,
      [IdentityProviderType.OAUTH]: SignInWithGeneric,
      [IdentityProviderType.OIDC]: SignInWithGeneric,
      [IdentityProviderType.GITHUB]: SignInWithGithub,
      [IdentityProviderType.GITHUB_ES]: SignInWithGithub,
      [IdentityProviderType.AZURE_AD]: SignInWithAzureAd,
      [IdentityProviderType.GOOGLE]: (props) => (
        <SignInWithGoogle {...props} e2e="google" />
      ),
      [IdentityProviderType.GITLAB]: SignInWithGitlab,
      [IdentityProviderType.GITLAB_SELF_HOSTED]: SignInWithGitlab,
    };

    const Component = components[type];
    return Component ? (
      <Component key={id} name={name} onClick={onClick} />
    ) : null;
  };

  return (
    <div className="flex flex-col w-full space-y-2 text-sm">
      {identityProviders?.map(renderIDPButton)}
      {error && (
        <div className="py-4">
          <Alert>{error}</Alert>
        </div>
      )}
    </div>
  );
}

SignInWithIdp.displayName = "SignInWithIDP";
