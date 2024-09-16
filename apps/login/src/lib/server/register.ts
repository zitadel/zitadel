"use server";

import { addHumanUser } from "@/lib/zitadel";
import { createSessionForUserIdAndUpdateCookie } from "@/utils/session";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";

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

  return createSessionForUserIdAndUpdateCookie(
    human.userId,
    command.password,
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
