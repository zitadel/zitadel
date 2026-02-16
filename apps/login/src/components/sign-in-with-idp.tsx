"use client";

import { idpTypeToSlug } from "@/lib/idp";
import { redirectToIdp } from "@/lib/server/idp";
import { IdentityProvider, IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { ReactNode, useActionState } from "react";
import { Alert } from "./alert";
import { SignInWithIdentityProviderProps } from "./idps/base-button";
import { SignInWithApple } from "./idps/sign-in-with-apple";
import { SignInWithAzureAd } from "./idps/sign-in-with-azure-ad";
import { SignInWithGeneric } from "./idps/sign-in-with-generic";
import { SignInWithGithub } from "./idps/sign-in-with-github";
import { SignInWithGitlab } from "./idps/sign-in-with-gitlab";
import { SignInWithGoogle } from "./idps/sign-in-with-google";
import { Translated } from "./translated";
import { AutoSubmitForm } from "./auto-submit-form";
import { trackEvent, MixpanelEvents } from "@/lib/mixpanel";

export interface SignInWithIDPProps {
  children?: ReactNode;
  identityProviders: IdentityProvider[];
  requestId?: string;
  organization?: string;
  sessionId?: string;
  postErrorRedirectUrl?: string;
  showLabel?: boolean;
}

export function SignInWithIdp({
  identityProviders,
  requestId,
  organization,
  sessionId,
  postErrorRedirectUrl,
  showLabel = true,
}: Readonly<SignInWithIDPProps>) {
  const [state, action, _isPending] = useActionState(redirectToIdp, {});

  const renderIDPButton = (idp: IdentityProvider, index: number) => {
    const { id, name, type } = idp;

    const components: Partial<Record<IdentityProviderType, (props: SignInWithIdentityProviderProps) => ReactNode>> = {
      [IdentityProviderType.APPLE]: SignInWithApple,
      [IdentityProviderType.OAUTH]: SignInWithGeneric,
      [IdentityProviderType.OIDC]: SignInWithGeneric,
      [IdentityProviderType.GITHUB]: SignInWithGithub,
      [IdentityProviderType.GITHUB_ES]: SignInWithGithub,
      [IdentityProviderType.AZURE_AD]: SignInWithAzureAd,
      [IdentityProviderType.GOOGLE]: (props) => <SignInWithGoogle {...props} e2e="google" />,
      [IdentityProviderType.GITLAB]: SignInWithGitlab,
      [IdentityProviderType.GITLAB_SELF_HOSTED]: SignInWithGitlab,
      [IdentityProviderType.SAML]: SignInWithGeneric,
      [IdentityProviderType.LDAP]: SignInWithGeneric,
      [IdentityProviderType.JWT]: SignInWithGeneric,
    };

    const Component = components[type];
    return Component ? (
      <form action={action} className="flex" key={`idp-${index}`} onSubmit={() => trackEvent(MixpanelEvents.idp_button_clicked, { idp_name: name, idp_type: String(type) })}>
        <input type="hidden" name="id" value={id} />
        <input type="hidden" name="provider" value={idpTypeToSlug(type)} />
        <input type="hidden" name="requestId" value={requestId} />
        <input type="hidden" name="organization" value={organization} />
        {sessionId && <input type="hidden" name="sessionId" value={sessionId} />}
        {postErrorRedirectUrl && <input type="hidden" name="postErrorRedirectUrl" value={postErrorRedirectUrl} />}
        <Component key={id} name={name} />
      </form>
    ) : null;
  };

  return (
    <div className="flex w-full flex-col space-y-2 text-sm">
      {state?.samlData && <AutoSubmitForm url={state.samlData.url} fields={state.samlData.fields} />}
      {showLabel && (
        <p className="ztdl-p text-center">
          <Translated i18nKey="orSignInWith" namespace="idp" />
        </p>
      )}
      {!!identityProviders?.length && identityProviders?.map(renderIDPButton)}
      {state?.error && (
        <div className="py-4">
          <Alert>{state?.error}</Alert>
        </div>
      )}
    </div>
  );
}

SignInWithIdp.displayName = "SignInWithIDP";
