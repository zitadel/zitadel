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
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getNextUrl } from "../client";
import { getSessionCookieByLoginName } from "../cookies";
import { createSessionAndUpdateCookie } from "./cookie";
import { checkMFAFactors } from "./password";

type VerifyUserByEmailCommand = {
  userId: string;
  code: string;
  isInvite: boolean;
  authRequestId?: string;
};

export async function sendVerification(command: VerifyUserByEmailCommand) {
  const verifyResponse = command.isInvite
    ? await verifyInviteCode(command.userId, command.code).catch(() => {
        return { error: "Could not verify invite" };
      })
    : await verifyEmail(command.userId, command.code).catch(() => {
        return { error: "Could not verify email" };
      });

  if (!verifyResponse) {
    return { error: "Could not verify user" };
  }

  const userResponse = await getUserByID(command.userId);

  if (!userResponse || !userResponse.user) {
    return { error: "Could not load user" };
  }

  const checks = create(ChecksSchema, {
    user: {
      search: {
        case: "loginName",
        value: userResponse.user.preferredLoginName,
      },
    },
  });

  const session = await createSessionAndUpdateCookie(
    checks,
    undefined,
    command.authRequestId,
  );

  const authMethodResponse = await listAuthenticationMethodTypes(
    command.userId,
  );

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
}

type resendVerifyEmailCommand = {
  userId: string;
  isInvite: boolean;
};

export async function resendVerification(command: resendVerifyEmailCommand) {
  return command.isInvite
    ? resendInviteCode(command.userId)
    : resendEmailCode(command.userId);
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
  if (!("loginName" in command || "userId" in command)) {
    return { error: "No userId, nor loginname provided" };
  }

  let session: Session | undefined;
  let user: User | undefined;

  const loginSettings = await getLoginSettings(command.organization);

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

    const userResponse = await getUserByID(session?.factors?.user?.id);

    if (!userResponse?.user) {
      return { error: "Could not load user" };
    }

    user = userResponse.user;
  } else if ("userId" in command) {
    const userResponse = await getUserByID(command.userId);

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

    // this is a fake error message to hide that the user does not even exist
    return { error: "Could not verify password" };
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

  const authMethodResponse = await listAuthenticationMethodTypes(user.userId);

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
  checkMFAFactors(
    session,
    loginSettings,
    authMethodResponse.authMethodTypes,
    command.organization,
    command.authRequestId,
  );

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
