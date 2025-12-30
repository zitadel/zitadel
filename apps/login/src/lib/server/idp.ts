"use server";

import {
  getLoginSettings,
  getUserByID,
  listAuthenticationMethodTypes,
  startIdentityProviderFlow,
  startLDAPIdentityProviderFlow,
  ServiceConfig,
} from "@/lib/zitadel";
import crypto from "crypto";
import { headers } from "next/headers";
import { redirect } from "next/navigation";
import { completeFlowOrGetUrl } from "../client";
import { getSessionCookieById } from "../cookies";
import { getOrSetFingerprintId } from "../fingerprint";
import { getServiceConfig } from "../service-url";
import { checkEmailVerification, checkMFAFactors } from "../verify-helper";
import { createSessionForIdpAndUpdateCookie } from "./cookie";
import { getPublicHost } from "./host";

export type RedirectToIdpState = { error?: string | null } | undefined;

export async function redirectToIdp(prevState: RedirectToIdpState, formData: FormData): Promise<RedirectToIdpState> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  const params = new URLSearchParams();

  const sessionId = formData.get("sessionId") as string;
  const requestId = formData.get("requestId") as string;
  const organization = formData.get("organization") as string;
  const idpId = formData.get("id") as string;
  const provider = formData.get("provider") as string;
  const postErrorRedirectUrl = formData.get("postErrorRedirectUrl") as string;

  if (sessionId) {
    try {
      // Validate that the requestor owns the session they are trying to link
      await getSessionCookieById({ sessionId });

      // Get fingerprint (ensure it exists)
      const fingerprintId = await getOrSetFingerprintId();

      // Create hash to verify intent upon return
      const linkFingerprint = crypto
        .createHash("sha256")
        .update(sessionId + fingerprintId)
        .digest("hex");

      params.set("linkToSessionId", sessionId);
      params.set("linkFingerprint", linkFingerprint);
    } catch {
      return { error: "Invalid session for linking" };
    }
  }
  if (requestId) params.set("requestId", requestId);
  if (organization) params.set("organization", organization);
  if (postErrorRedirectUrl) params.set("postErrorRedirectUrl", postErrorRedirectUrl);

  // redirect to LDAP page where username and password is requested
  if (provider === "ldap") {
    params.set("idpId", idpId);
    redirect(`/idp/ldap?` + params.toString());
  }

  const response = await startIDPFlow({
    serviceConfig,
    host,
    idpId,
    successUrl: `/idp/${provider}/process?` + params.toString(),
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
  serviceConfig: ServiceConfig;
  host: string;
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

async function startIDPFlow(command: StartIDPFlowCommand) {
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  const url = await startIdentityProviderFlow({
    serviceConfig: command.serviceConfig,
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

  const { serviceConfig } = getServiceConfig(_headers);

  if (!command.userId || !command.idpIntent) {
    throw new Error("No userId or loginName provided");
  }

  const userResponse = await getUserByID({ serviceConfig, userId: command.userId });

  if (!userResponse || !userResponse.user) {
    return { error: "User not found in the system" };
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: userResponse.user.details?.resourceOwner });

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
    const response = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  const mfaFactorCheck = await checkMFAFactors(
    serviceConfig,
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

  const { serviceConfig } = getServiceConfig(_headers);

  if (!command.username || !command.password) {
    return { error: "No username or password provided" };
  }

  const response = await startLDAPIdentityProviderFlow({
    serviceConfig,
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
