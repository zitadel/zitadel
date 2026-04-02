"use server"

import { create, toJson } from "@zitadel/client"
import {
  ListUsersRequestSchema,
  ListUsersResponseSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb"
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb"
import { createClient } from "@connectrpc/connect"
import { getTransport } from "./transport"

export interface UserSearchResult {
  userId: string
  username: string
  displayName: string
  email: string
  organizationId: string
  state: string
  type: string
}

/**
 * Search users via the v2 ListUsers RPC.
 * Uses an OR query across userName, displayName, and email with CONTAINS method.
 */
export async function searchUsers(
  query: string,
  limit: number = 10
): Promise<UserSearchResult[]> {
  const transport = getTransport()
  const client = createClient(UserService, transport)

  // TEXT_QUERY_METHOD_CONTAINS = 2
  const CONTAINS = 2

  const queries: any[] = []
  if (query.trim()) {
    queries.push({
      query: {
        case: "orQuery",
        value: {
          queries: [
            {
              query: {
                case: "userNameQuery",
                value: { userName: query, method: CONTAINS },
              },
            },
            {
              query: {
                case: "displayNameQuery",
                value: { displayName: query, method: CONTAINS },
              },
            },
            {
              query: {
                case: "emailQuery",
                value: { emailAddress: query, method: CONTAINS },
              },
            },
          ],
        },
      },
    })
  }

  const request = create(ListUsersRequestSchema, {
    query: { limit, asc: true },
    sortingColumn: 5, // USER_FIELD_NAME_DISPLAY_NAME
    queries,
  })

  const response = await client.listUsers(request)
  const json = toJson(ListUsersResponseSchema, response) as any
  const results = json.result ?? []

  return results.map((user: any) => ({
    userId: user.userId ?? "",
    username: user.username ?? "",
    displayName: user.human?.profile?.displayName ?? user.username ?? "",
    email: user.human?.email?.email ?? "",
    organizationId: user.details?.resourceOwner ?? "",
    state: user.state ?? "",
    type: user.type ?? "",
  }))
}
