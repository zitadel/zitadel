"use client";
import { ReactNode, useState } from "react";

// import { IdentityProviderType } from "@zitadel/server";
// import { IdentityProvider } from "@zitadel/client";

import {
  SignInWithGitlab,
  SignInWithAzureAD,
  SignInWithGoogle,
  SignInWithGithub,
} from "@zitadel/react";
import { useRouter } from "next/navigation";

export interface SignInWithIDPProps {
  children?: ReactNode;
  identityProviders: any[];
}

export function SignInWithIDP({ identityProviders }: SignInWithIDPProps) {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const router = useRouter();

  async function startFlow(idp: any) {
    console.log("start flow");
    const host = "http://localhost:3000";
    setLoading(true);

    const res = await fetch("/api/idp/start", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        idpId: idp.id,
        successUrl: `${host}`,
        failureUrl: `${host}`,
      }),
    });

    const response = await res.json();

    setLoading(false);
    if (!res.ok) {
      setError(response.details);
      return Promise.reject(response.details);
    }
    return response;
  }

  return (
    <div className="flex flex-col w-full space-y-2 text-sm">
      {identityProviders &&
        identityProviders.map((idp, i) => {
          switch (idp.type) {
            case 6: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGithub>
              );
            case 7: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB_ES:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGithub>
              );
            case 5: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_AZURE_AD:
              return (
                <SignInWithAzureAD
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithAzureAD>
              );
            case 10: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GOOGLE:
              return (
                <SignInWithGoogle
                  key={`idp-${i}`}
                  name={idp.name}
                  onClick={() =>
                    startFlow(idp).then(({ authUrl }) => {
                      console.log("done");
                      router.push(authUrl);
                    })
                  }
                ></SignInWithGoogle>
              );
            case 8: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGitlab>
              );
            case 9: //IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGitlab>
              );
            default:
              return <div>{idp.name}</div>;
          }
        })}
    </div>
  );
}

SignInWithIDP.displayName = "SignInWithIDP";
