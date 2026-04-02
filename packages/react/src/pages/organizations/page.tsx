import { toJson } from "@zitadel/client"
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { listOrganizations } from "../../api/organizations"
import { OrganizationsClient } from "./organizations-client"

/**
 * Organizations list page — server component.
 */
export default async function OrganizationsPage() {
  let organizations: any[] = []
  let totalResult = 0
  let error: string | null = null

  try {
    const response = await listOrganizations({ pageSize: 10 })
    const json = toJson(ListOrganizationsResponseSchema, response)
    organizations = (json as any).result ?? []
    totalResult = parseInt((json as any).details?.totalResult ?? "0", 10)
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load organizations"
    console.error("Failed to load organizations:", e)
  }

  return <OrganizationsClient organizations={organizations} totalResult={totalResult} error={error} />
}
