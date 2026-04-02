"use server";

import { create } from "@zitadel/client";
import {
  LockUserRequestSchema,
  UnlockUserRequestSchema,
  DeactivateUserRequestSchema,
  ReactivateUserRequestSchema,
  DeleteUserRequestSchema,
  PasswordResetRequestSchema,
  UpdateUserRequestSchema,
  UpdateUserRequest_HumanSchema,
  UpdateUserRequest_Human_ProfileSchema,
  UserService,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { SetHumanEmailSchema } from "@zitadel/proto/zitadel/user/v2/email_pb";
import { SetHumanPhoneSchema } from "@zitadel/proto/zitadel/user/v2/phone_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

function getUserClient() {
  return createClient(UserService, getTransport());
}

export async function lockUser(userId: string) {
  const client = getUserClient();
  const req = create(LockUserRequestSchema, { userId });
  await client.lockUser(req);
}

export async function unlockUser(userId: string) {
  const client = getUserClient();
  const req = create(UnlockUserRequestSchema, { userId });
  await client.unlockUser(req);
}

export async function deactivateUser(userId: string) {
  const client = getUserClient();
  const req = create(DeactivateUserRequestSchema, { userId });
  await client.deactivateUser(req);
}

export async function reactivateUser(userId: string) {
  const client = getUserClient();
  const req = create(ReactivateUserRequestSchema, { userId });
  await client.reactivateUser(req);
}

export async function deleteUser(userId: string) {
  const client = getUserClient();
  const req = create(DeleteUserRequestSchema, { userId });
  await client.deleteUser(req);
}

export async function resetPassword(userId: string) {
  const client = getUserClient();
  const req = create(PasswordResetRequestSchema, { userId });
  return client.passwordReset(req);
}

export interface UpdateUserData {
  username?: string;
  profile?: {
    givenName?: string;
    familyName?: string;
    nickName?: string;
    displayName?: string;
    preferredLanguage?: string;
  };
  email?: string;
  phone?: string;
}

/**
 * Update a human user's profile, email, phone, and/or username.
 */
export async function updateUser(userId: string, data: UpdateUserData) {
  const client = getUserClient();

  // Build profile sub-message if any profile fields provided
  let human = undefined;
  if (data.profile || data.email !== undefined || data.phone !== undefined) {
    const profileMsg = data.profile
      ? create(UpdateUserRequest_Human_ProfileSchema, {
          ...(data.profile.givenName !== undefined && { givenName: data.profile.givenName }),
          ...(data.profile.familyName !== undefined && { familyName: data.profile.familyName }),
          ...(data.profile.nickName !== undefined && { nickName: data.profile.nickName }),
          ...(data.profile.displayName !== undefined && { displayName: data.profile.displayName }),
          ...(data.profile.preferredLanguage !== undefined && { preferredLanguage: data.profile.preferredLanguage }),
        })
      : undefined;

    const emailMsg = data.email !== undefined
      ? create(SetHumanEmailSchema, { email: data.email })
      : undefined;

    const phoneMsg = data.phone !== undefined
      ? create(SetHumanPhoneSchema, { phone: data.phone || undefined })
      : undefined;

    human = create(UpdateUserRequest_HumanSchema, {
      ...(profileMsg && { profile: profileMsg }),
      ...(emailMsg && { email: emailMsg }),
      ...(phoneMsg && { phone: phoneMsg }),
    });
  }

  const req = create(UpdateUserRequestSchema, {
    userId,
    ...(data.username !== undefined && { username: data.username }),
    ...(human && {
      userType: {
        case: "human" as const,
        value: human,
      },
    }),
  });

  return client.updateUser(req);
}

