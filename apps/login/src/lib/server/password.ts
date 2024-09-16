"use server";

import { listUsers, passwordReset } from "@/lib/zitadel";

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
    return { error: "Could not find user" };
  }
  const userId = users.result[0].userId;

  return passwordReset(userId);
}
