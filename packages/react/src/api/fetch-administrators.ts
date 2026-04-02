"use server";

import { toJson } from "@zitadel/client";
import { ListAdministratorsResponseSchema } from "@zitadel/proto/zitadel/internal_permission/v2/internal_permission_service_pb";
import { listAdministrators } from "./administrators";

export interface FetchAdministratorsResult {
  administrators: any[];
  totalResult: number;
}

/**
 * Fetch administrators and return JSON-safe data for client component consumption.
 */
export async function fetchAdministrators(
  pageSize: number = 10,
  offset: number = 0,
): Promise<FetchAdministratorsResult> {
  const response = await listAdministrators({ pageSize, offset });
  const json = toJson(ListAdministratorsResponseSchema, response) as any;

  return {
    administrators: json.administrators ?? [],
    totalResult: parseInt(json.pagination?.totalResult ?? "0", 10),
  };
}
