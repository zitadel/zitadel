"use server";

import {
  createSessionAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "@/lib/server/cookie";
import {
  listAuthenticationMethodTypes,
  listUsers,
  passwordReset,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  Checks,
  ChecksSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getSessionCookieByLoginName } from "../cookies";

type ResetPasswordCommand = {
  loginName: string;
  organization?: string;
};

export async function resetPassword(command: ResetPasswordCommand) {
  const users = await listUsers({
    loginName: command.loginName,
    organizationId: command.organization,
  });

  if (
    !users.details ||
    users.details.totalResult !== BigInt(1) ||
    !users.result[0].userId
  ) {
    return { error: "Could not send Password Reset Link" };
  }
  const userId = users.result[0].userId;

  return passwordReset(userId);
}

export type UpdateSessionCommand = {
  loginName: string;
  organization?: string;
  checks: Checks;
  authRequestId?: string;
};

export async function sendPassword(command: UpdateSessionCommand) {
  let sessionCookie = await getSessionCookieByLoginName({
    loginName: command.loginName,
    organization: command.organization,
  }).catch((error) => {
    console.warn("Ignored error:", error);
  });

  let session;
  if (!sessionCookie) {
    const users = await listUsers({
      loginName: command.loginName,
      organizationId: command.organization,
    });

    if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
      const checks = create(ChecksSchema, {
        user: { search: { case: "userId", value: users.result[0].userId } },
        password: { password: command.checks.password?.password },
      });

      session = await createSessionAndUpdateCookie(
        checks,
        undefined,
        command.authRequestId,
      );
    }

    // this is a fake error message to hide that the user does not even exist
    return { error: "Could not verify password!" };
  } else {
    session = await setSessionAndUpdateCookie(
      sessionCookie,
      command.checks,
      undefined,
      command.authRequestId,
    );
  }

  if (!session?.factors?.user?.id || !sessionCookie) {
    return { error: "Could not create session for user" };
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (command.checks && command.checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes(
      session.factors.user.id,
    );
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  return {
    sessionId: session.id,
    factors: session.factors,
    challenges: session.challenges,
    authMethods,
  };
}
