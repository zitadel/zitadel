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
import { create } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { headers } from "next/headers";
import { getNextUrl } from "../client";
import { getSessionCookieByLoginName } from "../cookies";
import { getApiUrlOfHeaders } from "../service";
import { loadMostRecentSession } from "../session";
import { checkMFAFactors } from "../verify-helper";
import { createSessionAndUpdateCookie } from "./cookie";

export async function verifyTOTP(
  code: string,
  loginName?: string,
  organization?: string,
) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  return loadMostRecentSession({
    host,
    sessionParams: {
      loginName,
      organization,
    },
  }).then((session) => {
    if (session?.factors?.user?.id) {
      return verifyTOTPRegistration({
        host,
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
  authRequestId?: string;
};

export async function sendVerification(command: VerifyUserByEmailCommand) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  const verifyResponse = command.isInvite
    ? await verifyInviteCode({
        host,
        userId: command.userId,
        verificationCode: command.code,
      }).catch(() => {
        return { error: "Could not verify invite" };
      })
    : await verifyEmail({
        host,
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
      host,
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
      host,
      userId: session?.factors?.user?.id,
    });

    if (!userResponse?.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;
  } else {
    const userResponse = await getUserByID({ host, userId: command.userId });

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

    session = await createSessionAndUpdateCookie(
      checks,
      undefined,
      command.authRequestId,
    );
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
    host,
    organization: user.details?.resourceOwner,
  });

  const authMethodResponse = await listAuthenticationMethodTypes({
    host,
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

  // redirect to mfa factor if user has one, or redirect to set one up
  const mfaFactorCheck = checkMFAFactors(
    session,
    loginSettings,
    authMethodResponse.authMethodTypes,
    command.organization,
    command.authRequestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  // login user if no additional steps are required
  if (command.authRequestId && session.id) {
    const nextUrl = await getNextUrl(
      {
        sessionId: session.id,
        authRequestId: command.authRequestId,
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
  authRequestId?: string;
};

export async function resendVerification(command: resendVerifyEmailCommand) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host) {
    return { error: "No host found" };
  }

  return command.isInvite
    ? resendInviteCode({ host, userId: command.userId })
    : resendEmailCode({
        userId: command.userId,
        host,
        authRequestId: command.authRequestId,
      });
}

type sendEmailCommand = {
  host: string;
  userId: string;
  authRequestId?: string;
};

export async function sendEmailCode(command: sendEmailCommand) {
  return zitadelSendEmailCode({
    userId: command.userId,
    host: command.host,
    authRequestId: command.authRequestId,
  });
}

export type SendVerificationRedirectWithoutCheckCommand = {
  organization?: string;
  authRequestId?: string;
} & (
  | { userId: string; loginName?: never }
  | { userId?: never; loginName: string }
);

export async function sendVerificationRedirectWithoutCheck(
  command: SendVerificationRedirectWithoutCheckCommand,
) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);
  const host = instanceUrl;

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

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
      host,
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
      host,
      userId: session?.factors?.user?.id,
    });

    if (!userResponse?.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;
  } else if ("userId" in command) {
    const userResponse = await getUserByID({ host, userId: command.userId });

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

    session = await createSessionAndUpdateCookie(
      checks,
      undefined,
      command.authRequestId,
    );
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
    host,
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
    host,
    organization: user.details?.resourceOwner,
  });

  // redirect to mfa factor if user has one, or redirect to set one up
  const mfaFactorCheck = checkMFAFactors(
    session,
    loginSettings,
    authMethodResponse.authMethodTypes,
    command.organization,
    command.authRequestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  // login user if no additional steps are required
  if (command.authRequestId && session.id) {
    const nextUrl = await getNextUrl(
      {
        sessionId: session.id,
        authRequestId: command.authRequestId,
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
