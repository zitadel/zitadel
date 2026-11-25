import "server-only";

import { completeFlowOrGetUrl } from "@/lib/actions/flow";
import { getServiceConfig } from "@/lib/server/service-url";
import { createSessionForIdpAndUpdateCookie } from "@/lib/server/session-cookie";
import { checkEmailVerification, checkMFAFactors } from "@/lib/server/verify-helper";
import { getLoginSettings, getUserByID, listAuthenticationMethodTypes } from "@/lib/server/zitadel";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

// Not a server action: the user id comes from the caller, so only the IDP callback handling,
// which resolved the user from the retrieved intent, may call this.

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

  const t = await getTranslations("idp");

  // A call ZITADEL rejects throws on purpose: processIDPCallback turns it into the failure page,
  // which offers the user's other sign-in methods.
  const userResponse = await getUserByID({ serviceConfig, userId: command.userId });

  if (!userResponse || !userResponse.user) {
    return { error: t("errors.userNotFound") };
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: userResponse.user.details?.resourceOwner });

  const session = await createSessionForIdpAndUpdateCookie({
    userId: command.userId,
    idpIntent: command.idpIntent,
    requestId: command.requestId,
    lifetime: loginSettings?.externalLoginCheckLifetime,
  });

  if (!session || !session.factors?.user) {
    return { error: t("errors.couldNotCreateSession") };
  }

  const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  // check to see if user was verified
  const emailVerificationCheck = await checkEmailVerification(session, humanUser, command.organization, command.requestId);

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
