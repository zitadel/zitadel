"use server";

import {
  getLoginSettings,
  getUserByID,
  listAuthenticationMethodTypes,
  startIdentityProviderFlow,
  startLDAPIdentityProviderFlow,
} from "@/lib/zitadel";
import { headers } from "next/headers";
import { redirect } from "next/navigation";
import { completeFlowOrGetUrl } from "../client";
import { getServiceUrlFromHeaders } from "../service-url";
import { checkEmailVerification, checkMFAFactors } from "../verify-helper";
import { createSessionForIdpAndUpdateCookie } from "./cookie";
import { getOriginalHost } from "./host";

export type RedirectToIdpState = { error?: string | null } | undefined;

export async function redirectToIdp(prevState: RedirectToIdpState, formData: FormData): Promise<RedirectToIdpState> {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = await getOriginalHost();

  const params = new URLSearchParams();

  const linkOnly = formData.get("linkOnly") === "true";
  const requestId = formData.get("requestId") as string;
  const organization = formData.get("organization") as string;
  const idpId = formData.get("id") as string;
  const provider = formData.get("provider") as string;
  const postErrorRedirectUrl = formData.get("postErrorRedirectUrl") as string;

  if (linkOnly) params.set("link", "true");
  if (requestId) params.set("requestId", requestId);
  if (organization) params.set("organization", organization);
  if (postErrorRedirectUrl) params.set("postErrorRedirectUrl", postErrorRedirectUrl);

  // redirect to LDAP page where username and password is requested
  if (provider === "ldap") {
    params.set("idpId", idpId);
    redirect(`/idp/ldap?` + params.toString());
  }

  const response = await startIDPFlow({
    serviceUrl,
    host,
    idpId,
    successUrl: `/idp/${provider}/success?` + params.toString(),
    failureUrl: `/idp/${provider}/failure?` + params.toString(),
  });

  if (!response) {
    return { error: "Could not start IDP flow" };
  }

  if (response && "redirect" in response && response?.redirect) {
    redirect(response.redirect);
  }

  return { error: "Unexpected response from IDP flow" };
}

export type StartIDPFlowCommand = {
  serviceUrl: string;
  host: string;
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

async function startIDPFlow(command: StartIDPFlowCommand) {
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  const url = await startIdentityProviderFlow({
    serviceUrl: command.serviceUrl,
    idpId: command.idpId,
    urls: {
      successUrl: `${command.host.includes("localhost") ? "http://" : "https://"}${command.host}${basePath}${command.successUrl}`,
      failureUrl: `${command.host.includes("localhost") ? "http://" : "https://"}${command.host}${basePath}${command.failureUrl}`,
    },
  });

  if (!url) {
    return { error: "Could not start IDP flow" };
  }

  return { redirect: url };
}

export type CreateNewSessionCommand = {
  userId: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  loginName?: string;
  password?: string;
  organization?: string;
  requestId?: string;
};

export async function createNewSessionFromIdpIntent(command: CreateNewSessionCommand) {
  const _headers = await headers();

  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  if (!command.userId || !command.idpIntent) {
    throw new Error("No userId or loginName provided");
  }

  const userResponse = await getUserByID({
    serviceUrl,
    userId: command.userId,
  });

  if (!userResponse || !userResponse.user) {
    return { error: "User not found in the system" };
  }

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: userResponse.user.details?.resourceOwner,
  });

  const session = await createSessionForIdpAndUpdateCookie({
    userId: command.userId,
    idpIntent: command.idpIntent,
    requestId: command.requestId,
    lifetime: loginSettings?.externalLoginCheckLifetime,
  });

  if (!session || !session.factors?.user) {
    return { error: "Could not create session" };
  }

  const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  // check to see if user was verified
  const emailVerificationCheck = checkEmailVerification(session, humanUser, command.organization, command.requestId);

  if (emailVerificationCheck?.redirect) {
    return emailVerificationCheck;
  }

  // check if user has MFA methods
  let authMethods;
  if (session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({
      serviceUrl,
      userId: session.factors.user.id,
    });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  const mfaFactorCheck = await checkMFAFactors(
    serviceUrl,
    session,
    loginSettings,
    authMethods || [], // Pass empty array if no auth methods
    command.organization,
    command.requestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  return completeFlowOrGetUrl(
    command.requestId && session.id
      ? {
          sessionId: session.id,
          requestId: command.requestId,
          organization: session.factors.user.organizationId,
        }
      : {
          loginName: session.factors.user.loginName,
          organization: session.factors.user.organizationId,
        },
    loginSettings?.defaultRedirectUri,
  );
}

type createNewSessionForLDAPCommand = {
  username: string;
  password: string;
  idpId: string;
  link: boolean;
};

export async function createNewSessionForLDAP(command: createNewSessionForLDAPCommand) {
  const _headers = await headers();

  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  if (!command.username || !command.password) {
    return { error: "No username or password provided" };
  }

  const response = await startLDAPIdentityProviderFlow({
    serviceUrl,
    idpId: command.idpId,
    username: command.username,
    password: command.password,
  });

  if (!response || response.nextStep.case !== "idpIntent" || !response.nextStep.value) {
    return { error: "Could not start LDAP identity provider flow" };
  }

  const { userId, idpIntentId, idpIntentToken } = response.nextStep.value;

  const params = new URLSearchParams({
    userId,
    id: idpIntentId,
    token: idpIntentToken,
  });

  if (command.link) {
    params.set("link", "true");
  }

  return {
    redirect: `/idp/ldap/success?` + params.toString(),
  };
}
