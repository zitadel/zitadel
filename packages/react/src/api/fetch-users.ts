"use server";

import { toJson } from "@zitadel/client";
import { ListUsersResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { listUsers, listUsersByOrganization } from "./users";
import { listOrganizations } from "./organizations";

export interface FetchUsersResult {
  users: any[];
  organizations: any[];
  totalResult: number;
}

/**
 * Fetch users, optionally filtered by organization ID.
 * Returns JSON-safe data for client component consumption.
 */
export async function fetchUsers(
  orgId?: string | null,
  pageSize: number = 10,
  offset: number = 0,
): Promise<FetchUsersResult> {
  const [usersResponse, orgsResponse] = await Promise.all([
    orgId
      ? listUsersByOrganization(orgId, { pageSize, offset })
      : listUsers({ pageSize, offset }),
    listOrganizations({ pageSize: 100 }),
  ]);

  const usersJson = toJson(ListUsersResponseSchema, usersResponse) as any;
  const orgsJson = toJson(ListOrganizationsResponseSchema, orgsResponse) as any;

  return {
    users: usersJson.result ?? [],
    organizations: orgsJson.result ?? [],
    totalResult: parseInt(usersJson.details?.totalResult ?? "0", 10),
  };
}
