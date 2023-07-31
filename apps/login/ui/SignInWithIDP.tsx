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
import { ProviderSlug } from "#/lib/demos";

export interface SignInWithIDPProps {
  children?: ReactNode;
  instanceUrl: string;
  identityProviders: any[];
  startIDPFlowPath?: (idpId: string) => string;
}

const START_IDP_FLOW_PATH = (idpId: string) =>
  `/v2alpha/users/idps/${idpId}/start`;

export function SignInWithIDP({
  instanceUrl,
  identityProviders,
  startIDPFlowPath = START_IDP_FLOW_PATH,
}: SignInWithIDPProps) {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const router = useRouter();

  async function startFlow(idpId: string, provider: ProviderSlug) {
    const host = process.env.VERCEL_URL ?? "http://localhost:3000";
    setLoading(true);

    // const path = startIDPFlowPath(idpId);
    // const res = await fetch(`${instanceUrl}${path}`, {
    const res = await fetch("/api/idp/start", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        idpId,
        successUrl: `${host}/register/idp/${provider}/success`,
        failureUrl: `${host}/register/idp/${provider}/failure`,
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
                  onClick={() =>
                    startFlow(idp.id, ProviderSlug.GITHUB).then(
                      ({ authUrl }) => {
                        console.log("done");
                        router.push(authUrl);
                      }
                    )
                  }
                ></SignInWithGithub>
              );
            case 7: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB_ES:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  // onClick={() =>
                  //   startFlow(idp, ProviderSlug.GITHUB).then(({ authUrl }) => {
                  //     console.log("done");
                  //     router.push(authUrl);
                  //   })
                  // }
                ></SignInWithGithub>
              );
            case 5: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_AZURE_AD:
              return <SignInWithAzureAD key={`idp-${i}`}></SignInWithAzureAD>;
            case 10: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GOOGLE:
              return (
                <SignInWithGoogle
                  key={`idp-${i}`}
                  onClick={() =>
                    startFlow(idp.id, ProviderSlug.GOOGLE).then(
                      ({ authUrl }) => {
                        console.log("done");
                        router.push(authUrl);
                      }
                    )
                  }
                ></SignInWithGoogle>
              );
            case 8: // IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB:
              return <SignInWithGitlab key={`idp-${i}`}></SignInWithGitlab>;
            case 9: //IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED:
              return <SignInWithGitlab key={`idp-${i}`}></SignInWithGitlab>;
            default:
              return <div>{idp.name}</div>;
          }
        })}
    </div>
  );
}

SignInWithIDP.displayName = "SignInWithIDP";
