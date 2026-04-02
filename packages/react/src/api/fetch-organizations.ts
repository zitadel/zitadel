"use server";

import { toJson } from "@zitadel/client";
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { listOrganizations } from "./organizations";

/**
 * Fetch organizations as JSON-safe data with pagination info.
 */
export async function fetchOrganizationsPage(pageSize: number = 10, offset: number = 0) {
  const response = await listOrganizations({ pageSize, offset });
  const json = toJson(ListOrganizationsResponseSchema, response) as any;
  return {
    organizations: json.result ?? [],
    totalResult: parseInt(json.details?.totalResult ?? "0", 10),
  };
}
