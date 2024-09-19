"use server";

import { createSessionAndUpdateCookie } from "@/lib/server/cookie";
import { addHumanUser } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";

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
  const human = await addHumanUser({
    email: command.email,
    firstName: command.firstName,
    lastName: command.lastName,
    password: command.password ? command.password : undefined,
    organization: command.organization,
  });

  if (!human) {
    return { error: "Could not create user" };
  }

  const checks = create(ChecksSchema, {
    user: { search: { case: "userId", value: human.userId } },
    password: { password: command.password },
  });

  return createSessionAndUpdateCookie(
    checks,
    undefined,
    command.authRequestId,
  ).then((session) => {
    return {
      userId: human.userId,
      sessionId: session.id,
      factors: session.factors,
    };
  });
}
