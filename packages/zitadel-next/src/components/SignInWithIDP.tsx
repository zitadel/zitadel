"use server";
import { ReactNode } from "react";

import {
  ZitadelServer,
  settings,
  GetActiveIdentityProvidersResponse,
  IdentityProvider,
  IdentityProviderType,
} from "@zitadel/server";
// import {
//   SignInWithGitlab,
//   SignInWithAzureAD,
//   SignInWithGoogle,
//   SignInWithGithub,
// } from "@zitadel/react";

export interface SignInWithIDPProps {
  children?: ReactNode;
  server: ZitadelServer;
  orgId?: string;
}

function getIdentityProviders(
  server: ZitadelServer,
  orgId?: string
): Promise<IdentityProvider[] | undefined> {
  const settingsService = settings.getSettings(server);
  console.log("req");
  return settingsService
    .getActiveIdentityProviders(
      orgId ? { ctx: { orgId } } : { ctx: { instance: true } },
      {}
    )
    .then((resp: GetActiveIdentityProvidersResponse) => {
      return resp.identityProviders;
    });
}

export function SignInWithIDP(props: SignInWithIDPProps) {
  console.log(props.server);
  // const identityProviders = await getIdentityProviders(
  //   props.server,
  //   props.orgId
  // );

  // console.log(identityProviders);

  return (
    <div className="ztdl-next-flex ztdl-next-flex-col ztdl-next-w-full ztdl-next-space-y-2 ztdl-next-text-sm">
      {/* {identityProviders &&
        identityProviders.map((idp, i) => {
          switch (idp.type) {
            case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGithub>
              );
            case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB_ES:
              return (
                <SignInWithGithub
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGithub>
              );
            case IdentityProviderType.IDENTITY_PROVIDER_TYPE_AZURE_AD:
              return (
                <SignInWithAzureAD
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithAzureAD>
              );
            case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GOOGLE:
              return (
                <SignInWithGoogle
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGoogle>
              );
            case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGitlab>
              );
            case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED:
              return (
                <SignInWithGitlab
                  key={`idp-${i}`}
                  name={idp.name}
                ></SignInWithGitlab>
              );
            default:
              return <div>{idp.name}</div>;
          }
        })} */}
      {props.children}
    </div>
  );
}

SignInWithIDP.displayName = "SignInWithIDP";
