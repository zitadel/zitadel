"use server";

import { createSessionAndUpdateCookie } from "@/lib/server/cookie";
import { addHumanUser, getLoginSettings, getUserByID } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import {
  ChecksJson,
  ChecksSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getNextUrl } from "../client";
import { checkEmailVerification } from "../verify-helper";

type RegisterUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization?: string;
  authRequestId?: string;
};

export type RegisterUserResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};
export async function registerUser(command: RegisterUserCommand) {
  const addResponse = await addHumanUser({
    email: command.email,
    firstName: command.firstName,
    lastName: command.lastName,
    password: command.password ? command.password : undefined,
    organization: command.organization,
  });

  if (!addResponse) {
    return { error: "Could not create user" };
  }

  const loginSettings = await getLoginSettings(command.organization);

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

  const session = await createSessionAndUpdateCookie(
    checks,
    undefined,
    command.authRequestId,
    command.password ? loginSettings?.passwordCheckLifetime : undefined,
  );

  if (!session || !session.factors?.user) {
    return { error: "Could not create session" };
  }

  if (!command.password) {
    const params = new URLSearchParams({
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    });

    if (command.authRequestId) {
      params.append("authRequestId", command.authRequestId);
    }

    return { redirect: "/passkey/set?" + params };
  } else {
    const userResponse = await getUserByID(session?.factors?.user?.id);

    if (!userResponse.user) {
      return { error: "Could not find user" };
    }

    const humanUser =
      userResponse.user.type.case === "human"
        ? userResponse.user.type.value
        : undefined;

    const emailVerificationCheck = checkEmailVerification(
      session,
      humanUser,
      session.factors.user.organizationId,
      command.authRequestId,
    );

    if (emailVerificationCheck?.redirect) {
      return emailVerificationCheck;
    }

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

    return { redirect: url };
  }
}
