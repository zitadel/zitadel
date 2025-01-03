"use server";

import { addHumanUser, createInviteCode } from "@/lib/zitadel";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { headers } from "next/headers";

type InviteUserCommand = {
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

export async function inviteUser(command: InviteUserCommand) {
  const host = (await headers()).get("host");

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

  const codeResponse = await createInviteCode(human.userId, host);

  if (!codeResponse || !human) {
    return { error: "Could not create invite code" };
  }

  return human.userId;
}
