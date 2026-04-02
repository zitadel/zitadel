import { toJson } from "@zitadel/client"
import { ListAdministratorsResponseSchema } from "@zitadel/proto/zitadel/internal_permission/v2/internal_permission_service_pb"
import { listAdministrators } from "../../api/administrators"
import { AdministratorsClient } from "./administrators-client"

/**
 * Administrators page — server component that fetches administrators.
 */
export default async function AdministratorsPage() {
  let administrators: any[] = []
  let totalResult = 0
  let error: string | null = null

  try {
    const response = await listAdministrators({ pageSize: 10 })
    const json = toJson(ListAdministratorsResponseSchema, response)
    administrators = (json as any).administrators ?? []
    totalResult = parseInt((json as any).pagination?.totalResult ?? "0", 10)
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load administrators"
    console.error("Failed to load administrators:", e)
  }

  return <AdministratorsClient administrators={administrators} totalResult={totalResult} error={error} />
}
