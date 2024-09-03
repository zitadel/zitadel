"use server";

import { resendEmailCode, verifyEmail } from "@/lib/zitadel";

type VerifyUserByEmailCommand = {
  userId: string;
  code: string;
};

export async function verifyUserByEmail(command: VerifyUserByEmailCommand) {
  const { userId, code } = command;
  return verifyEmail(userId, code);
}

type resendVerifyEmailCommand = {
  userId: string;
};

export async function resendVerifyEmail(command: resendVerifyEmailCommand) {
  const { userId } = command;

  // replace with resend Mail method once its implemented
  return resendEmailCode(userId);
}
