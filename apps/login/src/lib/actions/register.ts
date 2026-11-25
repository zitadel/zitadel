"use server";

import { completeFlowOrGetUrl } from "@/lib/actions/flow";
import { answerRejection, isRejectedByUser } from "@/lib/server/api-error";
import { getBindingIdCookie } from "@/lib/server/fingerprint";
import { isIdpFlowProofValid } from "@/lib/server/idp-flow-proof";
import { createLogger } from "@/lib/server/logger";
import { getServiceConfig } from "@/lib/server/service-url";
import { createSessionAndUpdateCookie, createSessionForIdpAndUpdateCookie } from "@/lib/server/session-cookie";
import { setUserVerificationCheckCookie } from "@/lib/server/verification-check";
import { checkEmailVerification, checkMFAFactors } from "@/lib/server/verify-helper";
import {
  addHumanUser,
  addIDPLink,
  getLoginSettings,
  getUserByID,
  listAuthenticationMethodTypes,
  retrieveIDPIntent,
} from "@/lib/server/zitadel";
import { Code, ConnectError, Duration, create } from "@zitadel/client";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks, ChecksJson, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

const logger = createLogger("register");

const MAX_SESSION_RETRIES = 3;
const RETRY_DELAYS_MS = [500, 1000, 2000];

/**
 * After user creation, backend projections (users, login_names) may not be updated yet.
 * This helper retries createSessionAndUpdateCookie on NotFound errors with increasing delays.
 */
async function createSessionWithRetry(command: { checks: Checks; requestId: string | undefined; lifetime?: Duration }) {
  let lastError: unknown;
  for (let attempt = 0; attempt < MAX_SESSION_RETRIES; attempt++) {
    try {
      return await createSessionAndUpdateCookie(command);
    } catch (error) {
      lastError = error;
      const isNotFound = error instanceof ConnectError && error.code === Code.NotFound;
      const isLastAttempt = attempt + 1 >= MAX_SESSION_RETRIES;
      if (!isNotFound || isLastAttempt) {
        throw error;
      }
      const delay = RETRY_DELAYS_MS[attempt] ?? 2000;
      logger.warn(
        `Session creation failed with NotFound (attempt ${attempt + 1}/${MAX_SESSION_RETRIES}), retrying in ${delay}ms...`,
      );
      await new Promise((resolve) => setTimeout(resolve, delay));
    }
  }
  throw lastError;
}

/** Answers a rejected addHumanUser call: the user name is the email address, so AlreadyExists means it is taken. */
function answerAddHumanUserRejection(error: unknown, t: (key: string) => string) {
  const taken = isRejectedByUser(error) && error.code === Code.AlreadyExists;
  return answerRejection(error, taken ? t("errors.userAlreadyExists") : t("errors.couldNotCreateUser"));
}

type RegisterUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization: string;
  requestId?: string;
};

export type RegisterUserResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};
export async function registerUser(
  command: RegisterUserCommand,
): Promise<{ error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const t = await getTranslations("register");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const loginSettings = await getLoginSettings({ serviceConfig, organization: command.organization });

  if (!loginSettings) {
    return { error: t("errors.couldNotGetLoginSettings") };
  }

  if (!loginSettings.allowRegister) {
    return { error: t("errors.registerNotAllowed") };
  }

  if (command.password && !loginSettings.allowLocalAuthentication) {
    return { error: t("errors.localAuthenticationNotAllowed") };
  }

  let addResponse;
  try {
    addResponse = await addHumanUser({
      serviceConfig,
      email: command.email,
      firstName: command.firstName,
      lastName: command.lastName,
      password: command.password ? command.password : undefined,
      organization: command.organization,
    });
  } catch (error) {
    return answerAddHumanUserRejection(error, t);
  }

  let checkPayload: any = {
    user: { search: { case: "userId", value: addResponse.userId } },
  };

  if (command.password) {
    checkPayload = {
      ...checkPayload,
      password: { password: command.password },
    } as ChecksJson;
  }

  const checks = create(ChecksSchema, checkPayload);

  let result;
  try {
    result = await createSessionWithRetry({
      checks,
      requestId: command.requestId,
      lifetime: command.password ? loginSettings?.passwordCheckLifetime : undefined,
    });
  } catch (error) {
    return answerRejection(error, t("errors.couldNotCreateSession"), { context: { userId: addResponse.userId } });
  }

  const session = result?.session;

  if (!session || !session.factors?.user) {
    return { error: t("errors.couldNotCreateSession") };
  }

  if (!command.password) {
    const params = new URLSearchParams({
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    });

    if (command.requestId) {
      params.append("requestId", command.requestId);
    }

    // Set verification cookie for users registering with passkey (no password)
    // This allows them to proceed with passkey registration without additional verification
    await setUserVerificationCheckCookie(session.factors.user.id);

    return { redirect: "/passkey/set?" + params };
  } else {
    let userResponse;
    try {
      userResponse = await getUserByID({ serviceConfig, userId: session.factors.user.id });
    } catch (error) {
      return answerRejection(error, t("errors.userNotFound"));
    }

    if (!userResponse?.user) {
      return { error: t("errors.userNotFound") };
    }

    const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

    const emailVerificationCheck = await checkEmailVerification(
      session,
      humanUser,
      session.factors.user.organizationId,
      command.requestId,
    );

    if (emailVerificationCheck?.redirect) {
      return emailVerificationCheck;
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
}

type RegisterUserAndLinkToIDPommand = {
  email: string;
  firstName: string;
  lastName: string;
  organization: string;
  requestId?: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  /** only used when the IDP did not provide a user name */
  idpUserName: string;
  flowFingerprint?: string;
};

export type registerUserAndLinkToIDPResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};
export async function registerUserAndLinkToIDP(
  command: RegisterUserAndLinkToIDPommand,
): Promise<{ error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const t = await getTranslations("register");

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const loginSettings = await getLoginSettings({ serviceConfig, organization: command.organization });

  if (!loginSettings) {
    return { error: t("errors.couldNotGetLoginSettings") };
  }

  if (!loginSettings.allowRegister) {
    return { error: t("errors.registerNotAllowed") };
  }

  // the intent must belong to an IdP flow started in this browser (login CSRF)
  const bindingIdCookie = await getBindingIdCookie();
  if (!isIdpFlowProofValid(command.flowFingerprint, bindingIdCookie?.value)) {
    logger.warn("IDP flow fingerprint missing or mismatch on registration");
    return { error: t("errors.sessionMismatch") };
  }

  // The external identity to link is read from the intent, which proves the user authenticated at
  // the IDP. The browser only passes the intent id and token, so it can't link another identity.
  const intent = await retrieveIDPIntent({
    serviceConfig,
    id: command.idpIntent.idpIntentId,
    token: command.idpIntent.idpIntentToken,
  }).catch((error) => {
    logger.warn("registerUserAndLinkToIDP: could not retrieve IDP intent", { error });
    return undefined;
  });

  const idpInformation = intent?.idpInformation;

  if (!idpInformation?.idpId || !idpInformation.userId) {
    return { error: t("errors.couldNotLinkIDP") };
  }

  let addUserResponse;
  try {
    addUserResponse = await addHumanUser({
      serviceConfig,
      email: command.email,
      firstName: command.firstName,
      lastName: command.lastName,
      organization: command.organization,
    });
  } catch (error) {
    return answerAddHumanUserRejection(error, t);
  }

  let idpLink;
  try {
    idpLink = await addIDPLink({
      serviceConfig,
      idp: {
        id: idpInformation.idpId,
        userId: idpInformation.userId,
        userName: idpInformation.userName || command.idpUserName,
      },
      userId: addUserResponse.userId,
    });
  } catch (error) {
    // e.g. the external identity is already linked to another user
    return answerRejection(error, t("errors.couldNotLinkIDP"), { context: { userId: addUserResponse.userId } });
  }

  if (!idpLink) {
    return { error: t("errors.couldNotLinkIDP") };
  }

  let session;
  try {
    session = await createSessionForIdpAndUpdateCookie({
      requestId: command.requestId,
      userId: addUserResponse.userId, // the user we just created
      idpIntent: command.idpIntent,
      lifetime: loginSettings?.externalLoginCheckLifetime,
    });
  } catch (error) {
    // e.g. the intent expired or was already used
    return answerRejection(error, t("errors.couldNotCreateSession"), { context: { userId: addUserResponse.userId } });
  }

  if (!session || !session.factors?.user) {
    return { error: t("errors.couldNotCreateSession") };
  }

  // const userResponse = await getUserByID({
  //   serviceConfig.baseUrl,
  //   userId: session?.factors?.user?.id,
  // });

  // if (!userResponse.user) {
  //   return { error: "User not found in the system" };
  // }

  // const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  // check to see if user was verified
  // const emailVerificationCheck = checkEmailVerification(session, humanUser, command.organization, command.requestId);

  // if (emailVerificationCheck?.redirect) {
  //   return emailVerificationCheck;
  // }

  // check if user has MFA methods
  let authMethods;
  if (session.factors?.user?.id) {
    let response;
    try {
      response = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });
    } catch (error) {
      return answerRejection(error, t("errors.couldNotRegisterUser"));
    }
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  // Always check MFA factors, even if no auth methods are configured
  // This ensures that force MFA settings are respected
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
