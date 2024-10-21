"use server";

import {
  listAuthenticationMethodTypes,
  resendEmailCode,
  resendInviteCode,
  verifyEmail,
  verifyInviteCode,
} from "@/lib/zitadel";

type VerifyUserByEmailCommand = {
  userId: string;
  code: string;
  isInvite: boolean;
};

export async function verifyUser(command: VerifyUserByEmailCommand) {
  const verifyResponse = command.isInvite
    ? await verifyInviteCode(command.userId, command.code).catch((error) => {
        console.log(error.code);
        return { error: "Could not verify invite" };
      })
    : await verifyEmail(command.userId, command.code).catch((error) => {
        console.log(error.code);
        return { error: "Could not verify email" };
      });

  if (!verifyResponse) {
    return { error: "Could not verify user" };
  }

  const authMethodResponse = await listAuthenticationMethodTypes(
    command.userId,
  );

  if (!authMethodResponse || !authMethodResponse.authMethodTypes) {
    return { error: "Could not load possible authenticators" };
  }

  return { authMethodTypes: authMethodResponse.authMethodTypes };
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
