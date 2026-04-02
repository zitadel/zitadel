"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListApplicationsRequestSchema,
  ListApplicationsResponseSchema,
  ApplicationService,
} from "@zitadel/proto/zitadel/application/v2/application_service_pb";
import { PaginationRequestSchema } from "@zitadel/proto/zitadel/filter/v2/filter_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

/**
 * Fetch all applications across all projects using the v2 ApplicationService.
 * Uses a single API call without a project filter to list all apps.
 */
export async function fetchAllApplications(pageSize: number = 10, offset: number = 0) {
  const service = createClient(ApplicationService, getTransport());
  const pagination = create(PaginationRequestSchema, {
    limit: pageSize,
    offset: BigInt(offset),
  });
  const request = create(ListApplicationsRequestSchema, { pagination });
  const response = await service.listApplications(request);
  const json = toJson(ListApplicationsResponseSchema, response) as any;

  return {
    applications: json.applications ?? [],
    totalResult: parseInt(json.pagination?.totalResult ?? "0", 10),
  };
}
