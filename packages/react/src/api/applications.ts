"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListApplicationsRequestSchema,
  ListApplicationsResponseSchema,
  GetApplicationRequestSchema,
  GetApplicationResponseSchema,
  DeleteApplicationRequestSchema,
  ApplicationService,
} from "@zitadel/proto/zitadel/application/v2/application_service_pb";
import {
  ApplicationSearchFilterSchema,
  ProjectIDFilterSchema,
} from "@zitadel/proto/zitadel/application/v2/application_pb";
import {
  PaginationRequestSchema,
} from "@zitadel/proto/zitadel/filter/v2/filter_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

function getAppService() {
  return createClient(ApplicationService, getTransport());
}

/**
 * List applications for a project using the project_id filter.
 */
export async function listApplications(projectId: string, opts?: { pageSize?: number; offset?: number }) {
  const service = getAppService();
  const pagination = create(PaginationRequestSchema, {
    limit: opts?.pageSize ?? 10,
    offset: BigInt(opts?.offset ?? 0),
  });
  const projectFilter = create(ApplicationSearchFilterSchema, {
    filter: {
      case: "projectIdFilter",
      value: create(ProjectIDFilterSchema, { projectId }),
    },
  });
  const request = create(ListApplicationsRequestSchema, {
    pagination,
    filters: [projectFilter],
  });
  return service.listApplications(request);
}

/**
 * Get a single application by ID.
 */
export async function getApplication(projectId: string, appId: string) {
  const service = getAppService();
  const request = create(GetApplicationRequestSchema, { applicationId: appId });
  return service.getApplication(request);
}

/**
 * Delete an application.
 */
export async function deleteApplication(projectId: string, appId: string) {
  const service = getAppService();
  const request = create(DeleteApplicationRequestSchema, {
    applicationId: appId,
    projectId,
  });
  return service.deleteApplication(request);
}

/**
 * Fetch applications as JSON-safe data for client consumption.
 */
export async function fetchApplications(projectId: string, pageSize: number = 10, offset: number = 0) {
  const response = await listApplications(projectId, { pageSize, offset });
  const json = toJson(ListApplicationsResponseSchema, response) as any;
  return {
    applications: json.applications ?? [],
    totalResult: parseInt(json.pagination?.totalResult ?? "0", 10),
  };
}

/**
 * Fetch a single application as JSON-safe data.
 */
export async function fetchApplication(projectId: string, appId: string) {
  const response = await getApplication(projectId, appId);
  const json = toJson(GetApplicationResponseSchema, response) as any;
  return json.application ?? null;
}
