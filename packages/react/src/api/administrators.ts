"use server";

import { create } from "@zitadel/client";
import {
  type ListAdministratorsResponse,
  ListAdministratorsRequestSchema,
} from "@zitadel/proto/zitadel/internal_permission/v2/internal_permission_service_pb";
import { getInternalPermissionService } from "./services";

/**
 * List administrators with optional filters.
 */
export async function listAdministrators(opts?: {
  pageSize?: number;
  offset?: number;
}): Promise<ListAdministratorsResponse> {
  const service = getInternalPermissionService();
  const request = create(ListAdministratorsRequestSchema, {
    pagination: {
      limit: opts?.pageSize ?? 10,
      offset: BigInt(opts?.offset ?? 0),
    },
  });
  return service.listAdministrators(request);
}
