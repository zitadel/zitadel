"use server";

import { answerRejection } from "@/lib/server/api-error";
import { isClassifiedError } from "@/lib/server/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/server/logger";
import {
  createInviteCode,
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  verifyEmail,
  verifyInviteCode,
  verifyTOTPRegistration,
  sendEmailCode as zitadelSendEmailCode,
} from "@/lib/server/zitadel";

import { completeFlowOrGetUrl } from "@/lib/actions/flow";
import { getSessionCookieByLoginName } from "@/lib/server/cookies";
import { getSetupAuthorizationError } from "@/lib/server/enrollment-guard";
import { getPublicHostWithProtocol } from "@/lib/server/host";
import { checkMFAFactors } from "@/lib/server/mfa-helper";
import { getServiceConfig } from "@/lib/server/service-url";
import { loadMostRecentSession } from "@/lib/server/session";
import { createSessionAndUpdateCookie } from "@/lib/server/session-cookie";
import { setUserVerificationCheckCookie } from "@/lib/server/verification-check";
import { buildVerificationUrlTemplate } from "@/lib/server/verification-email";
import { Code, create } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

const logger = createLogger("verify");

export async function verifyTOTP(code: string, loginName?: string, organization?: string) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const t = await getTranslations("verify");

  const session = await loadMostRecentSession({
    serviceConfig,
    sessionParams: {
      loginName,
      organization,
    },
  });

  const userId = session?.factors?.user?.id;
  if (!session || !userId) {
    return { error: t("errors.couldNotVerify") };
  }

  // Enrollment must be authorized: a bare identify-only session must not be able to
  // activate a TOTP factor on the account (GHSA-45f2-5q3r-xgg6).
  const enrollmentError = await getSetupAuthorizationError({
    serviceConfig,
    session,
    userId,
    method: "totp",
  });
  if (enrollmentError) {
    return { error: enrollmentError };
  }

  try {
    return await verifyTOTPRegistration({ serviceConfig, code, userId });
  } catch (error) {
    // A wrong code (InvalidArgument) is the routine failure here: answer it instead of failing the
    // action with a 500. Server faults keep throwing.
    const wrongCode = isClassifiedError(error) && error.code === Code.InvalidArgument;
    return answerRejection(error, t(wrongCode ? "errors.invalidCode" : "errors.couldNotVerify"));
  }
}

type VerifyUserByEmailCommand = {
  userId: string;
  loginName?: string; // to determine already existing session
  organization?: string;
  code: string;
  isInvite: boolean;
  requestId?: string;
};

export async function sendVerification(command: VerifyUserByEmailCommand) {
  const t = await getTranslations("verify");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const verifyResponse = command.isInvite
    ? await verifyInviteCode({ serviceConfig, userId: command.userId, verificationCode: command.code }).catch((error) => {
        logger.warn("Could not verify invite:", { error });
        return { error: t("errors.couldNotVerifyInvite") };
      })
    : await verifyEmail({ serviceConfig, userId: command.userId, verificationCode: command.code }).catch((error) => {
        logger.warn("Could not verify email:", { error });
        return { error: t("errors.couldNotVerifyEmail") };
      });

  if ("error" in verifyResponse) {
    return verifyResponse;
  }

  if (!verifyResponse) {
    return { error: t("errors.couldNotVerify") };
  }

  let session: Session | undefined;
  const userResponse = await getUserByID({ serviceConfig, userId: command.userId });

  if (!userResponse || !userResponse.user) {
    return { error: t("errors.couldNotLoadUser") };
  }

  const user = userResponse.user;

  const sessionCookie = await getSessionCookieByLoginName({
    loginName: command.loginName ?? user.preferredLoginName,
    organization: command.organization,
  });

  if (sessionCookie) {
    session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token })
      .then((response) => {
        if (response?.session) {
          return response.session;
        }
      })
      .catch((error) => {
        // user session is not found, so we create a new one
        logger.warn("[verify] user session is not found, so we create a new one", { error });
        return undefined;
      });
  }

  // load auth methods for user
  const authMethodResponse = await listAuthenticationMethodTypes({ serviceConfig, userId: user.userId });

  if (!authMethodResponse || !authMethodResponse.authMethodTypes) {
    return { error: t("errors.couldNotLoadAuthenticators") };
  }

  const hasPrimaryMethod =
    authMethodResponse?.authMethodTypes?.some(
      (m: AuthenticationMethodType) =>
        m === AuthenticationMethodType.PASSWORD ||
        m === AuthenticationMethodType.PASSKEY ||
        m === AuthenticationMethodType.IDP,
    ) ?? false;

  // if no primary auth methods are found on the user, redirect to set one up
  if (!hasPrimaryMethod) {
    if (!session) {
      const checks = create(ChecksSchema, {
        user: {
          search: {
            case: "loginName",
            value: userResponse.user.preferredLoginName,
          },
        },
      });

      const result = await createSessionAndUpdateCookie({
        checks,
        requestId: command.requestId,
      });
      session = result.session;
    }

    if (!session) {
      return { error: t("errors.couldNotCreateSession") };
    }

    const params = new URLSearchParams({
      sessionId: session.id,
    });

    if (session.factors?.user?.loginName) {
      params.set("loginName", session.factors?.user?.loginName);
    }

    if (command.requestId) {
      params.set("requestId", command.requestId);
    }

    // remember the successful verification for this browser (signed with a server-only secret),
    // checks are done for users with invalid sessions and invalid userAgentId
    await setUserVerificationCheckCookie(user.userId);

    return { redirect: `/authenticator/set?${params}` };
  }

  // if no session found only show success page,
  // if user is invited, recreate invite flow to not depend on session
  if (!session?.factors?.user?.id) {
    const verifySuccessParams = new URLSearchParams({});

    if (command.userId) {
      verifySuccessParams.set("userId", command.userId);
    }

    if (("loginName" in command && command.loginName) || user.preferredLoginName) {
      verifySuccessParams.set(
        "loginName",
        "loginName" in command && command.loginName ? command.loginName : user.preferredLoginName,
      );
    }
    if (command.requestId) {
      verifySuccessParams.set("requestId", command.requestId);
    }
    if (command.organization) {
      verifySuccessParams.set("organization", command.organization);
    }

    return { redirect: `/verify/success?${verifySuccessParams}` };
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: user.details?.resourceOwner });

  // redirect to mfa factor if user has one, or redirect to set one up
  const mfaFactorCheck = await checkMFAFactors(
    serviceConfig,
    session,
    loginSettings,
    authMethodResponse.authMethodTypes,
    command.organization,
    command.requestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  // login user if no additional steps are required
  if (command.requestId && session.id) {
    return completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: command.requestId,
        organization: command.organization ?? session.factors?.user?.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  // Regular flow - return URL for client-side navigation
  return completeFlowOrGetUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors?.user?.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );
}

type resendVerifyEmailCommand = {
  userId: string;
  isInvite: boolean;
  requestId?: string;
};

export async function resendVerification(command: resendVerifyEmailCommand) {
  const t = await getTranslations("verify");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const hostWithProtocol = await getPublicHostWithProtocol(_headers);

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";
  const urlTemplate = buildVerificationUrlTemplate(hostWithProtocol, basePath, command.isInvite, command.requestId);

  // FailedPrecondition: the user is already verified, so there is nothing to resend.
  // Other rejections get a message too; server faults keep throwing.
  const alreadyVerified = (error: unknown) => isClassifiedError(error) && error.code === Code.FailedPrecondition;

  if (command.isInvite) {
    return createInviteCode({ serviceConfig, userId: command.userId, urlTemplate }).catch((error: unknown) =>
      answerRejection(error, t(alreadyVerified(error) ? "errors.userAlreadyVerified" : "errors.couldNotResendInvite")),
    );
  }

  return zitadelSendEmailCode({ serviceConfig, userId: command.userId, urlTemplate }).catch((error: unknown) =>
    answerRejection(error, t(alreadyVerified(error) ? "errors.userAlreadyVerified" : "errors.couldNotResendEmail")),
  );
}
