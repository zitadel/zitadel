"use server";

import { addHumanUser } from "@/lib/zitadel";
import { createSessionForUserIdAndUpdateCookie } from "@/utils/session";

type RegisterUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization?: string;
  authRequestId?: string;
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
    throw Error("Could not create user");
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
