"use server";

import {
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  resendEmailCode,
  resendInviteCode,
  verifyEmail,
  verifyInviteCode,
  verifyTOTPRegistration,
  sendEmailCode as zitadelSendEmailCode,
} from "@/lib/zitadel";
import crypto from "crypto";

import { create } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { cookies, headers } from "next/headers";
import { getNextUrl } from "../client";
import { getSessionCookieByLoginName } from "../cookies";
import { getOrSetFingerprintId } from "../fingerprint";
import { getServiceUrlFromHeaders } from "../service-url";
import { loadMostRecentSession } from "../session";
import { checkMFAFactors } from "../verify-helper";
import { createSessionAndUpdateCookie } from "./cookie";

export async function verifyTOTP(
  code: string,
  loginName?: string,
  organization?: string,
) {
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
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const verifyResponse = command.isInvite
    ? await verifyInviteCode({
        serviceUrl,
        userId: command.userId,
        verificationCode: command.code,
      }).catch(() => {
        return { error: "Could not verify invite" };
      })
    : await verifyEmail({
        serviceUrl,
        userId: command.userId,
        verificationCode: command.code,
      }).catch(() => {
        return { error: "Could not verify email" };
      });

  if ("error" in verifyResponse) {
    return verifyResponse;
  }

  if (!verifyResponse) {
    return { error: "Could not verify" };
  }

  let session: Session | undefined;
  let user: User | undefined;

  if ("loginName" in command) {
    const sessionCookie = await getSessionCookieByLoginName({
      loginName: command.loginName,
      organization: command.organization,
    }).catch((error) => {
      console.warn("Ignored error:", error);
    });

    if (!sessionCookie) {
      return { error: "Could not load session cookie" };
    }

    session = await getSession({
      serviceUrl,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    }).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });

    if (!session?.factors?.user?.id) {
      return { error: "Could not create session for user" };
    }

    const userResponse = await getUserByID({
      serviceUrl,
      userId: session?.factors?.user?.id,
    });

    if (!userResponse?.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;
  } else {
    const userResponse = await getUserByID({
      serviceUrl,
      userId: command.userId,
    });

    if (!userResponse || !userResponse.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;

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

  if (!session?.factors?.user?.id) {
    return { error: "Could not create session for user" };
  }

  if (!session?.factors?.user?.id) {
    return { error: "Could not create session for user" };
  }

  if (!user) {
    return { error: "Could not load user" };
  }

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: user.details?.resourceOwner,
  });

  const authMethodResponse = await listAuthenticationMethodTypes({
    serviceUrl,
    userId: user.userId,
  });

  if (!authMethodResponse || !authMethodResponse.authMethodTypes) {
    return { error: "Could not load possible authenticators" };
  }

  // if no authmethods are found on the user, redirect to set one up
  if (
    authMethodResponse &&
    authMethodResponse.authMethodTypes &&
    authMethodResponse.authMethodTypes.length == 0
  ) {
    const params = new URLSearchParams({
      sessionId: session.id,
    });

    if (session.factors?.user?.loginName) {
      params.set("loginName", session.factors?.user?.loginName);
    }

    // set hash of userId and userAgentId to prevent attacks, checks are done for users with invalid sessions and invalid userAgentId
    const cookiesList = await cookies();
    const userAgentId = await getOrSetFingerprintId();

    const verificationCheck = crypto
      .createHash("sha256")
      .update(`${user.userId}:${userAgentId}`)
      .digest("hex");

    await cookiesList.set({
      name: "verificationCheck",
      value: verificationCheck,
      httpOnly: true,
      path: "/",
      maxAge: 300, // 5 minutes
    });

    return { redirect: `/authenticator/set?${params}` };
  }

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
    const nextUrl = await getNextUrl(
      {
        sessionId: session.id,
        requestId: command.requestId,
        organization:
          command.organization ?? session.factors?.user?.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );

    return { redirect: nextUrl };
  }

  const url = await getNextUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors?.user?.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );

  return { redirect: url };
}

type resendVerifyEmailCommand = {
  userId: string;
  isInvite: boolean;
  requestId?: string;
};

export async function resendVerification(command: resendVerifyEmailCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = _headers.get("host");

  if (!host) {
    return { error: "No host found" };
  }

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  return command.isInvite
    ? resendInviteCode({ serviceUrl, userId: command.userId })
    : resendEmailCode({
        userId: command.userId,
        serviceUrl,
        urlTemplate:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/password/set?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
          (command.requestId ? `&requestId=${command.requestId}` : ""),
      });
}

type sendEmailCommand = {
  serviceUrl: string;

  userId: string;
  urlTemplate: string;
};

export async function sendEmailCode(command: sendEmailCommand) {
  return zitadelSendEmailCode({
    serviceUrl: command.serviceUrl,
    userId: command.userId,
    urlTemplate: command.urlTemplate,
  });
}

export type SendVerificationRedirectWithoutCheckCommand = {
  organization?: string;
  requestId?: string;
} & (
  | { userId: string; loginName?: never }
  | { userId?: never; loginName: string }
);

export async function sendVerificationRedirectWithoutCheck(
  command: SendVerificationRedirectWithoutCheckCommand,
) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  if (!("loginName" in command || "userId" in command)) {
    return { error: "No userId, nor loginname provided" };
  }

  let session: Session | undefined;
  let user: User | undefined;

  if ("loginName" in command) {
    const sessionCookie = await getSessionCookieByLoginName({
      loginName: command.loginName,
      organization: command.organization,
    }).catch((error) => {
      console.warn("Ignored error:", error);
    });

    if (!sessionCookie) {
      return { error: "Could not load session cookie" };
    }

    session = await getSession({
      serviceUrl,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    }).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });

    if (!session?.factors?.user?.id) {
      return { error: "Could not create session for user" };
    }

    const userResponse = await getUserByID({
      serviceUrl,
      userId: session?.factors?.user?.id,
    });

    if (!userResponse?.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;
  } else if ("userId" in command) {
    const userResponse = await getUserByID({
      serviceUrl,
      userId: command.userId,
    });

    if (!userResponse?.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;

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

  if (!session?.factors?.user?.id) {
    return { error: "Could not create session for user" };
  }

  if (!session?.factors?.user?.id) {
    return { error: "Could not create session for user" };
  }

  if (!user) {
    return { error: "Could not load user" };
  }

  const authMethodResponse = await listAuthenticationMethodTypes({
    serviceUrl,
    userId: user.userId,
  });

  if (!authMethodResponse || !authMethodResponse.authMethodTypes) {
    return { error: "Could not load possible authenticators" };
  }

  // if no authmethods are found on the user, redirect to set one up
  if (
    authMethodResponse &&
    authMethodResponse.authMethodTypes &&
    authMethodResponse.authMethodTypes.length == 0
  ) {
    const params = new URLSearchParams({
      sessionId: session.id,
    });

    if (session.factors?.user?.loginName) {
      params.set("loginName", session.factors?.user?.loginName);
    }
    return { redirect: `/authenticator/set?${params}` };
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
    const nextUrl = await getNextUrl(
      {
        sessionId: session.id,
        requestId: command.requestId,
        organization:
          command.organization ?? session.factors?.user?.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );

    return { redirect: nextUrl };
  }

  const url = await getNextUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors?.user?.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );

  return { redirect: url };
}
