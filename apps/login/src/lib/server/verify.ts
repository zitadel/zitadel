"use server";

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
} from "@/lib/zitadel";
import crypto from "crypto";

import { create } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { cookies, headers } from "next/headers";
import { completeFlowOrGetUrl } from "../client";
import { getSessionCookieByLoginName } from "../cookies";
import { getOrSetFingerprintId } from "../fingerprint";
import { getServiceUrlFromHeaders } from "../service-url";
import { loadMostRecentSession } from "../session";
import { checkMFAFactors } from "../verify-helper";
import { createSessionAndUpdateCookie } from "./cookie";
import { getOriginalHostWithProtocol } from "./host";
import { getTranslations } from "next-intl/server";

export async function verifyTOTP(code: string, loginName?: string, organization?: string) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  return loadMostRecentSession({
    serviceUrl,
    sessionParams: {
      loginName,
      organization,
    },
  }).then((session) => {
    if (session?.factors?.user?.id) {
      return verifyTOTPRegistration({
        serviceUrl,
        code,
        userId: session.factors.user.id,
      });
    } else {
      throw Error("No user id found in session.");
    }
  });
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
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const verifyResponse = command.isInvite
    ? await verifyInviteCode({
        serviceUrl,
        userId: command.userId,
        verificationCode: command.code,
      }).catch((error) => {
        console.warn(error);
        return { error: t("errors.couldNotVerifyInvite") };
      })
    : await verifyEmail({
        serviceUrl,
        userId: command.userId,
        verificationCode: command.code,
      }).catch((error) => {
        console.warn(error);
        return { error: t("errors.couldNotVerifyEmail") };
      });

  if ("error" in verifyResponse) {
    return verifyResponse;
  }

  if (!verifyResponse) {
    return { error: t("errors.couldNotVerify") };
  }

  let session: Session | undefined;
  const userResponse = await getUserByID({
    serviceUrl,
    userId: command.userId,
  });

  if (!userResponse || !userResponse.user) {
    return { error: t("errors.couldNotLoadUser") };
  }

  const user = userResponse.user;

  const sessionCookie = await getSessionCookieByLoginName({
    loginName: "loginName" in command ? command.loginName : user.preferredLoginName,
    organization: command.organization,
  }).catch((error) => {
    console.warn("Ignored error:", error); // checked later
  });

  if (sessionCookie) {
    session = await getSession({
      serviceUrl,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    }).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  // load auth methods for user
  const authMethodResponse = await listAuthenticationMethodTypes({
    serviceUrl,
    userId: user.userId,
  });

  if (!authMethodResponse || !authMethodResponse.authMethodTypes) {
    return { error: t("errors.couldNotLoadAuthenticators") };
  }

  // if no authmethods are found on the user, redirect to set one up
  if (authMethodResponse && authMethodResponse.authMethodTypes && authMethodResponse.authMethodTypes.length == 0) {
    if (!sessionCookie) {
      const checks = create(ChecksSchema, {
        user: {
          search: {
            case: "loginName",
            value: userResponse.user.preferredLoginName,
          },
        },
      });

      session = await createSessionAndUpdateCookie({
        checks,
        requestId: command.requestId,
      });
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

    // set hash of userId and userAgentId to prevent attacks, checks are done for users with invalid sessions and invalid userAgentId
    const cookiesList = await cookies();
    const userAgentId = await getOrSetFingerprintId();

    const verificationCheck = crypto.createHash("sha256").update(`${user.userId}:${userAgentId}`).digest("hex");

    await cookiesList.set({
      name: "verificationCheck",
      value: verificationCheck,
      httpOnly: true,
      path: "/",
      maxAge: 300, // 5 minutes
    });

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

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: user.details?.resourceOwner,
  });

  // redirect to mfa factor if user has one, or redirect to set one up
  const mfaFactorCheck = await checkMFAFactors(
    serviceUrl,
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
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const hostWithProtocol = await getOriginalHostWithProtocol();

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  return command.isInvite
    ? createInviteCode({
        serviceUrl,
        userId: command.userId,
        urlTemplate:
          `${hostWithProtocol}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true` +
          (command.requestId ? `&requestId=${command.requestId}` : ""),
      }).catch((error) => {
        if (error.code === 9) {
          return { error: t("errors.userAlreadyVerified") };
        }
        return { error: t("errors.couldNotResendInvite") };
      })
    : zitadelSendEmailCode({
        userId: command.userId,
        serviceUrl,
        urlTemplate:
          `${hostWithProtocol}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
          (command.requestId ? `&requestId=${command.requestId}` : ""),
      });
}

type SendEmailCommand = {
  userId: string;
  urlTemplate: string;
};

export async function sendEmailCode(command: SendEmailCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  return zitadelSendEmailCode({
    serviceUrl,
    userId: command.userId,
    urlTemplate: command.urlTemplate,
  });
}

export async function sendInviteEmailCode(command: SendEmailCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  return createInviteCode({
    serviceUrl,
    userId: command.userId,
    urlTemplate: command.urlTemplate,
  });
}
