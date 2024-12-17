"use server";

import {
  getLoginSettings,
  getUserByID,
  listAuthenticationMethodTypes,
  listUsers,
  resendEmailCode,
  resendInviteCode,
  verifyEmail,
  verifyInviteCode,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getSessionCookieByLoginName } from "../cookies";
import {
  createSessionAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "./cookie";

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

type SendVerificationRedirectWithoutCheckCommand =
  | {
      loginName: string;
      organization?: string;
      authRequestId?: string;
    }
  | {
      userId: string;
      authRequestId?: string;
    };

export async function sendVerificationRedirectWithoutCheck(
  command: SendVerificationRedirectWithoutCheckCommand,
) {
  if (!("loginName" in command || "userId" in command)) {
    return { error: "No userId, nor loginname provided" };
  }

  let sessionCookie;
  let loginSettings: LoginSettings | undefined;
  let session;
  let user: User;

  if ("loginName" in command) {
    sessionCookie = await getSessionCookieByLoginName({
      loginName: command.loginName,
      organization: command.organization,
    }).catch((error) => {
      console.warn("Ignored error:", error);
    });
  } else if (command.userId) {
    const users = await listUsers({
      loginName: command.loginName,
      organizationId: command.organization,
    });

    if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
      user = users.result[0];

      const checks = create(ChecksSchema, {
        user: { search: { case: "userId", value: users.result[0].userId } },
        password: { password: command.checks.password?.password },
      });

      loginSettings = await getLoginSettings(command.organization);

      session = await createSessionAndUpdateCookie(
        checks,
        undefined,
        command.authRequestId,
        loginSettings?.passwordCheckLifetime,
      );
    }

    // this is a fake error message to hide that the user does not even exist
    return { error: "Could not verify password" };
  } else {
    session = await setSessionAndUpdateCookie(
      sessionCookie,
      command.checks,
      undefined,
      command.authRequestId,
      loginSettings?.passwordCheckLifetime,
    );

    if (!session?.factors?.user?.id) {
      return { error: "Could not create session for user" };
    }

    const userResponse = await getUserByID(session?.factors?.user?.id);

    if (!userResponse.user) {
      return { error: "Could not find user" };
    }

    user = userResponse.user;
  }

  if (!loginSettings) {
    loginSettings = await getLoginSettings(
      command.organization ?? session.factors?.user?.organizationId,
    );
  }

  if (!session?.factors?.user?.id || !sessionCookie) {
    return { error: "Could not create session for user" };
  }
  // const userResponse = await getUserByID(command.userId);

  // if (!userResponse || !userResponse.user) {
  //   return { error: "Could not load user" };
  // }

  // const checks = create(ChecksSchema, {
  //   user: {
  //     search: {
  //       case: "loginName",
  //       value: userResponse.user.preferredLoginName,
  //     },
  //   },
  // });

  // const session = await createSessionAndUpdateCookie(
  //   checks,
  //   undefined,
  //   command.authRequestId,
  // );

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
