"use server";

import {
  getUserByID,
  resendEmailCode,
  resendInviteCode,
  verifyEmail,
  verifyInviteCode,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { createSessionAndUpdateCookie } from "./cookie";

type VerifyUserByEmailCommand = {
  userId: string;
  code: string;
  isInvite: boolean;
  authRequestId?: string;
};

export async function verifyUserAndCreateSession(
  command: VerifyUserByEmailCommand,
) {
  const verifyResponse = command.isInvite
    ? await verifyInviteCode(command.userId, command.code).catch((error) => {
        return { error: "Could not verify invite" };
      })
    : await verifyEmail(command.userId, command.code).catch((error) => {
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

  return {
    sessionId: session.id,
    factors: session.factors,
  };
}

type resendVerifyEmailCommand = {
  userId: string;
  isInvite: boolean;
};

export async function resendVerification(command: resendVerifyEmailCommand) {
  return command.isInvite
    ? resendEmailCode(command.userId)
    : resendInviteCode(command.userId);
}
