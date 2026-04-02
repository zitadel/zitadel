"use server";

import { create } from "@zitadel/client";
import {
  type ListUsersResponse,
  type GetUserByIDResponse,
  type AddHumanUserResponse,
  ListUsersRequestSchema,
  GetUserByIDRequestSchema,
  AddHumanUserRequestSchema,
  DeleteUserRequestSchema,
  UpdateUserRequestSchema,
  type UpdateUserRequest,
  type AddHumanUserRequest,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import {
  SearchQuerySchema as UserSearchQuerySchema,
  type SearchQuery as UserSearchQuery,
} from "@zitadel/proto/zitadel/user/v2/query_pb";
import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { getUserService } from "./services";

/**
 * List users with optional search queries and pagination.
 */
export async function listUsers(opts?: {
  queries?: UserSearchQuery[];
  pageSize?: number;
  offset?: number;
  sortingColumn?: number;
  asc?: boolean;
}): Promise<ListUsersResponse> {
  const userService = getUserService();
  const request = create(ListUsersRequestSchema, {
    query: {
      limit: opts?.pageSize ?? 10,
      offset: BigInt(opts?.offset ?? 0),
      asc: opts?.asc ?? true,
    },
    queries: opts?.queries ?? [],
  });
  return userService.listUsers(request);
}

/**
 * Get a single user by ID.
 */
export async function getUserById(
  userId: string
): Promise<GetUserByIDResponse> {
  const userService = getUserService();
  const request = create(GetUserByIDRequestSchema, { userId });
  return userService.getUserByID(request);
}

/**
 * Search users by email address.
 */
export async function searchUsersByEmail(
  email: string
): Promise<ListUsersResponse> {
  const query = create(UserSearchQuerySchema, {
    query: {
      case: "emailQuery",
      value: {
        emailAddress: email,
        method: TextQueryMethod.CONTAINS,
      },
    },
  });
  return listUsers({ queries: [query] });
}

/**
 * Search users by login name.
 */
export async function searchUsersByLoginName(
  loginName: string
): Promise<ListUsersResponse> {
  const query = create(UserSearchQuerySchema, {
    query: {
      case: "loginNameQuery",
      value: {
        loginName,
        method: TextQueryMethod.CONTAINS,
      },
    },
  });
  return listUsers({ queries: [query] });
}

/**
 * List users in a specific organization.
 */
export async function listUsersByOrganization(
  organizationId: string,
  opts?: { pageSize?: number; offset?: number }
): Promise<ListUsersResponse> {
  const query = create(UserSearchQuerySchema, {
    query: {
      case: "organizationIdQuery",
      value: { organizationId },
    },
  });
  return listUsers({ queries: [query], ...opts });
}

/**
 * Create a new human user.
 */
export async function addHumanUser(
  request: AddHumanUserRequest
): Promise<AddHumanUserResponse> {
  const userService = getUserService();
  const req = create(AddHumanUserRequestSchema, request);
  return userService.addHumanUser(req);
}

/**
 * Update a user.
 */
export async function updateUser(request: UpdateUserRequest) {
  const userService = getUserService();
  const req = create(UpdateUserRequestSchema, request);
  return userService.updateUser(req);
}

/**
 * Delete a user by ID.
 */
export async function deleteUser(userId: string) {
  const userService = getUserService();
  const req = create(DeleteUserRequestSchema, { userId });
  return userService.deleteUser(req);
}
