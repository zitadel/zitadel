"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListProjectsRequestSchema,
  ListProjectsResponseSchema,
  GetProjectRequestSchema,
  GetProjectResponseSchema,
  DeleteProjectRequestSchema,
  ProjectService,
} from "@zitadel/proto/zitadel/project/v2/project_service_pb";
import {
  PaginationRequestSchema,
} from "@zitadel/proto/zitadel/filter/v2/filter_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

function getProjectService() {
  return createClient(ProjectService, getTransport());
}

/**
 * List projects with pagination.
 */
export async function listProjects(opts?: { pageSize?: number; offset?: number }) {
  const service = getProjectService();
  const pagination = create(PaginationRequestSchema, {
    limit: opts?.pageSize ?? 10,
    offset: BigInt(opts?.offset ?? 0),
  });
  const request = create(ListProjectsRequestSchema, {
    pagination,
  });
  return service.listProjects(request);
}

/**
 * Get a single project by ID.
 */
export async function getProject(projectId: string) {
  const service = getProjectService();
  const request = create(GetProjectRequestSchema, { projectId });
  return service.getProject(request);
}

/**
 * Delete a project by ID.
 */
export async function deleteProject(projectId: string) {
  const service = getProjectService();
  const request = create(DeleteProjectRequestSchema, { projectId });
  return service.deleteProject(request);
}

/**
 * Fetch projects as JSON-safe data for client consumption.
 */
export async function fetchProjects(pageSize: number = 10, offset: number = 0) {
  const response = await listProjects({ pageSize, offset });
  const json = toJson(ListProjectsResponseSchema, response) as any;
  return {
    projects: json.projects ?? [],
    totalResult: parseInt(json.pagination?.totalResult ?? "0", 10),
  };
}

/**
 * Fetch a single project as JSON-safe data.
 */
export async function fetchProject(projectId: string) {
  const response = await getProject(projectId);
  const json = toJson(GetProjectResponseSchema, response) as any;
  return json.project ?? null;
}
