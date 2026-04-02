"use server";

import { create, toJson } from "@zitadel/client";
import { ListUsersResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import {
  ListProjectsRequestSchema,
  ListProjectsResponseSchema,
  ProjectService,
} from "@zitadel/proto/zitadel/project/v2/project_service_pb";
import {
  ListApplicationsRequestSchema,
  ListApplicationsResponseSchema,
  ApplicationService,
} from "@zitadel/proto/zitadel/application/v2/application_service_pb";
import { PaginationRequestSchema } from "@zitadel/proto/zitadel/filter/v2/filter_pb";
import { createClient } from "@connectrpc/connect";
import { listUsers, listUsersByOrganization } from "./users";
import { listOrganizations } from "./organizations";
import { getTransport } from "./transport";

export interface NavCounts {
  users: number;
  organizations: number;
  projects: number;
  applications: number;
}

/**
 * Lightweight fetch for sidebar nav counts.
 * Makes 4 parallel pageSize=1 calls to get totalResult counts.
 * Scoped to orgId when provided.
 */
export async function fetchNavCounts(orgId?: string | null): Promise<NavCounts> {
  try {
    const transport = getTransport();
    const projectService = createClient(ProjectService, transport);
    const appService = createClient(ApplicationService, transport);

    const [usersResp, orgsResp, projectsResp, appsResp] = await Promise.all([
      // Users — scoped to org if selected
      orgId
        ? listUsersByOrganization(orgId, { pageSize: 1 })
        : listUsers({ pageSize: 1 }),
      // Organizations — always global
      listOrganizations({ pageSize: 1 }),
      // Projects
      (async () => {
        const pagination = create(PaginationRequestSchema, { limit: 1 });
        const request = create(ListProjectsRequestSchema, { pagination });
        return projectService.listProjects(request);
      })(),
      // Applications
      (async () => {
        const pagination = create(PaginationRequestSchema, { limit: 1 });
        const request = create(ListApplicationsRequestSchema, { pagination });
        return appService.listApplications(request);
      })(),
    ]);

    const usersJson = toJson(ListUsersResponseSchema, usersResp) as any;
    const orgsJson = toJson(ListOrganizationsResponseSchema, orgsResp) as any;
    const projectsJson = toJson(ListProjectsResponseSchema, projectsResp) as any;
    const appsJson = toJson(ListApplicationsResponseSchema, appsResp) as any;

    return {
      users: Number(usersJson.details?.totalResult ?? 0),
      organizations: Number(orgsJson.details?.totalResult ?? 0),
      projects: Number(projectsJson.pagination?.totalResult ?? 0),
      applications: Number(appsJson.pagination?.totalResult ?? 0),
    };
  } catch (e) {
    console.error("Failed to fetch nav counts:", e);
    return { users: 0, organizations: 0, projects: 0, applications: 0 };
  }
}
