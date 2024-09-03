"use server";

import { addHumanUser } from "@/lib/zitadel";
import {
  createSessionAndUpdateCookie,
  createSessionForUserIdAndUpdateCookie,
} from "@/utils/session";

type RegisterUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization?: string;
  authRequestId?: string;
};
export async function registerUser(command: RegisterUserCommand) {
  const { email, password, firstName, lastName, organization, authRequestId } =
    command;

  const human = await addHumanUser({
    email: email,
    firstName,
    lastName,
    password: password ? password : undefined,
    organization,
  });
  if (!human) {
    throw Error("Could not create user");
  }

  return createSessionForUserIdAndUpdateCookie(
    human.userId,
    password,
    undefined,
    authRequestId,
  ).then((session) => {
    return {
      userId: human.userId,
      sessionId: session.id,
      factors: session.factors,
    };
  });
}
