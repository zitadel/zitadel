"use server";

import { resendEmailCode, verifyEmail } from "@/lib/zitadel";

type VerifyUserByEmailCommand = {
  userId: string;
  code: string;
};

export async function verifyUserByEmail(command: VerifyUserByEmailCommand) {
  return verifyEmail(command.userId, command.code);
}

type resendVerifyEmailCommand = {
  userId: string;
};

export async function resendVerifyEmail(command: resendVerifyEmailCommand) {
  return resendEmailCode(command.userId);
}
