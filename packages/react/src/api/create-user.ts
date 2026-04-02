"use server";

import { create, toJson } from "@zitadel/client";
import {
  CreateUserRequestSchema,
  CreateUserResponseSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

export interface CreateUserInput {
  organizationId: string;
  username?: string;
  givenName: string;
  familyName: string;
  displayName?: string;
  email: string;
  isEmailVerified?: boolean;
  password?: string;
  requirePasswordChange?: boolean;
}

/**
 * Create a new human user via the v2 CreateUser RPC.
 * Returns the created user's ID or throws an error.
 */
export async function createUser(input: CreateUserInput): Promise<{ userId: string }> {
  const transport = getTransport();
  const client = createClient(UserService, transport);

  const request = create(CreateUserRequestSchema, {
    organizationId: input.organizationId,
    username: input.username || undefined,
    userType: {
      case: "human",
      value: {
        profile: {
          givenName: input.givenName,
          familyName: input.familyName,
          displayName: input.displayName || `${input.givenName} ${input.familyName}`,
        },
        email: {
          email: input.email,
          verification: input.isEmailVerified
            ? { case: "isVerified", value: true }
            : { case: "sendCode", value: {} },
        },
        passwordType: input.password
          ? {
              case: "password",
              value: {
                password: input.password,
                changeRequired: input.requirePasswordChange ?? true,
              },
            }
          : undefined,
      },
    },
  });

  const response = await client.createUser(request);
  const json = toJson(CreateUserResponseSchema, response) as any;
  return { userId: json.userId ?? "" };
}
