import { toJson } from "@zitadel/client"
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { listOrganizations } from "../../../api/organizations"
import { AddUserForm } from "./add-user-form"

/**
 * Add User page — server component that fetches orgs for the org selector.
 */
export default async function AddUserPage() {
  let organizations: any[] = []

  try {
    const response = await listOrganizations({ pageSize: 100 })
    const json = toJson(ListOrganizationsResponseSchema, response)
    organizations = (json as any).result ?? []
  } catch (e) {
    console.error("Failed to load organizations for user creation:", e)
  }

  return <AddUserForm organizations={organizations} />
}
