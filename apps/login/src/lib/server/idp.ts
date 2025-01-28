"use server";

import {
  getLoginSettings,
  getUserByID,
  startIdentityProviderFlow,
} from "@/lib/zitadel";
import { headers } from "next/headers";
import { getNextUrl } from "../client";
import { getApiUrlOfHeaders } from "../service";
import { checkEmailVerification } from "../verify-helper";
import { createSessionForIdpAndUpdateCookie } from "./cookie";

export type StartIDPFlowCommand = {
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

export async function startIDPFlow(command: StartIDPFlowCommand) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host) {
    return { error: "Could not get host" };
  }

  return startIdentityProviderFlow({
    host,
    idpId: command.idpId,
    urls: {
      successUrl: `${host.includes("localhost") ? "http://" : "https://"}${host}${command.successUrl}`,
      failureUrl: `${host.includes("localhost") ? "http://" : "https://"}${host}${command.failureUrl}`,
    },
  }).then((response) => {
    if (
      response &&
      response.nextStep.case === "authUrl" &&
      response?.nextStep.value
    ) {
      return { redirect: response.nextStep.value };
    }
  });
}

type CreateNewSessionCommand = {
  userId: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  loginName?: string;
  password?: string;
  organization?: string;
  authRequestId?: string;
};

export async function createNewSessionFromIdpIntent(
  command: CreateNewSessionCommand,
) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host) {
    return { error: "Could not get domain" };
  }

  if (!command.userId || !command.idpIntent) {
    throw new Error("No userId or loginName provided");
  }

  const userResponse = await getUserByID({ host, userId: command.userId });

  if (!userResponse || !userResponse.user) {
    return { error: "User not found in the system" };
  }

  const loginSettings = await getLoginSettings({
    host,
    organization: userResponse.user.details?.resourceOwner,
  });

  const session = await createSessionForIdpAndUpdateCookie(
    command.userId,
    command.idpIntent,
    command.authRequestId,
    loginSettings?.externalLoginCheckLifetime,
  );

  if (!session || !session.factors?.user) {
    return { error: "Could not create session" };
  }

  const humanUser =
    userResponse.user.type.case === "human"
      ? userResponse.user.type.value
      : undefined;

  // check to see if user was verified
  const emailVerificationCheck = checkEmailVerification(
    session,
    humanUser,
    command.organization,
    command.authRequestId,
  );

  if (emailVerificationCheck?.redirect) {
    return emailVerificationCheck;
  }

  // TODO: check if user has MFA methods
  // const mfaFactorCheck = checkMFAFactors(session, loginSettings, authMethods, organization, authRequestId);
  // if (mfaFactorCheck?.redirect) {
  //   return mfaFactorCheck;
  // }

  const url = await getNextUrl(
    command.authRequestId && session.id
      ? {
          sessionId: session.id,
          authRequestId: command.authRequestId,
          organization: session.factors.user.organizationId,
        }
      : {
          loginName: session.factors.user.loginName,
          organization: session.factors.user.organizationId,
        },
    loginSettings?.defaultRedirectUri,
  );

  if (url) {
    return { redirect: url };
  }
}
