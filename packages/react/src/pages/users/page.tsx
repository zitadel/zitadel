import { toJson } from "@zitadel/client"
import { ListUsersResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb"
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { listUsers } from "../../api/users"
import { listOrganizations } from "../../api/organizations"
import { UsersClient } from "./users-client"

/**
 * Users list page — server component that fetches users and orgs.
 * Orgs are needed for the Add User sheet's org selector.
 */
export default async function UsersPage() {
  let users: any[] = []
  let organizations: any[] = []
  let totalResult = 0
  let error: string | null = null

  try {
    const [usersResponse, orgsResponse] = await Promise.all([
      listUsers({ pageSize: 10 }),
      listOrganizations({ pageSize: 100 }),
    ])
    const usersJson = toJson(ListUsersResponseSchema, usersResponse)
    const orgsJson = toJson(ListOrganizationsResponseSchema, orgsResponse)
    users = (usersJson as any).result ?? []
    organizations = (orgsJson as any).result ?? []
    totalResult = parseInt((usersJson as any).details?.totalResult ?? "0", 10)
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load users"
    console.error("Failed to load users:", e)
  }

  return <UsersClient users={users} organizations={organizations} totalResult={totalResult} error={error} />
}
