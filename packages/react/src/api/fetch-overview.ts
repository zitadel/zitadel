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
  ProjectSearchFilterSchema,
  ProjectOrganizationIDFilterSchema,
} from "@zitadel/proto/zitadel/project/v2/query_pb";
import {
  ListApplicationsRequestSchema,
  ListApplicationsResponseSchema,
  ApplicationService,
} from "@zitadel/proto/zitadel/application/v2/application_service_pb";
import {
  ApplicationSearchFilterSchema,
  ProjectIDFilterSchema,
} from "@zitadel/proto/zitadel/application/v2/application_pb";
import { PaginationRequestSchema } from "@zitadel/proto/zitadel/filter/v2/filter_pb";
import { createClient } from "@connectrpc/connect";
import { listUsers, listUsersByOrganization } from "./users";
import { listOrganizations } from "./organizations";
import { fetchAllSessions } from "./all-sessions";
import { getTransport } from "./transport";
import { getInstanceService } from "./services";

export interface OverviewStats {
  userCount: number;
  orgCount: number;
  projectCount: number;
  appCount: number;
  recentSessions: any[];
  version: string | null;
}

/**
 * Get unique project IDs for an org. The v2 ListProjects API returns duplicate entries
 * for projects that have grants (one entry per grant + the owner entry).
 */
async function getOrgProjectIds(orgId: string): Promise<string[]> {
  const service = createClient(ProjectService, getTransport());
  const pagination = create(PaginationRequestSchema, { limit: 100 });
  const orgFilter = create(ProjectSearchFilterSchema, {
    filter: {
      case: "organizationIdFilter",
      value: create(ProjectOrganizationIDFilterSchema, {
        organizationId: orgId,
        type: 1, // OWNED
      }),
    },
  });
  const request = create(ListProjectsRequestSchema, {
    pagination,
    filters: [orgFilter],
  });
  const response = await service.listProjects(request);
  const json = toJson(ListProjectsResponseSchema, response) as any;
  const projects = json.projects ?? [];

  // Deduplicate by projectId (projects with grants appear multiple times)
  const uniqueIds = Array.from(new Set<string>(projects.map((p: any) => String(p.projectId))));
  return uniqueIds;
}

/**
 * Count projects, optionally filtered by organization.
 */
async function countProjects(orgId?: string | null): Promise<number> {
  if (orgId) {
    return (await getOrgProjectIds(orgId)).length;
  }
  const service = createClient(ProjectService, getTransport());
  const pagination = create(PaginationRequestSchema, { limit: 1 });
  const request = create(ListProjectsRequestSchema, { pagination });
  const response = await service.listProjects(request);
  const json = toJson(ListProjectsResponseSchema, response) as any;
  return Number(json.pagination?.totalResult ?? json.projects?.length ?? 0);
}

/**
 * Count applications, optionally filtered by org's projects.
 */
async function countApplications(orgId?: string | null): Promise<number> {
  const service = createClient(ApplicationService, getTransport());

  if (!orgId) {
    // No org filter — count all apps
    const pagination = create(PaginationRequestSchema, { limit: 1 });
    const request = create(ListApplicationsRequestSchema, { pagination });
    const response = await service.listApplications(request);
    const json = toJson(ListApplicationsResponseSchema, response) as any;
    return Number(json.pagination?.totalResult ?? json.applications?.length ?? 0);
  }

  // Get unique project IDs for this org
  const projectIds = await getOrgProjectIds(orgId);
  if (projectIds.length === 0) return 0;

  // Count apps for each unique project in parallel
  const counts = await Promise.all(
    projectIds.map(async (projectId: string) => {
      const pagination = create(PaginationRequestSchema, { limit: 1 });
      const filter = create(ApplicationSearchFilterSchema, {
        filter: {
          case: "projectIdFilter",
          value: create(ProjectIDFilterSchema, { projectId }),
        },
      });
      const request = create(ListApplicationsRequestSchema, {
        pagination,
        filters: [filter],
      });
      const response = await service.listApplications(request);
      const json = toJson(ListApplicationsResponseSchema, response) as any;
      return Number(json.pagination?.totalResult ?? json.applications?.length ?? 0);
    })
  );

  return counts.reduce((sum, c) => sum + c, 0);
}

/**
 * Fetch overview stats, optionally scoped to an organization.
 */
export async function fetchOverviewStats(orgId?: string | null): Promise<{
  stats: OverviewStats;
  error: string | null;
}> {
  try {
    const [usersResp, orgsResp, projectCount, appCount, sessions] = await Promise.all([
      orgId
        ? listUsersByOrganization(orgId, { pageSize: 1 })
        : listUsers({ pageSize: 1 }),
      listOrganizations({ pageSize: 1 }),
      countProjects(orgId),
      countApplications(orgId),
      fetchAllSessions(5, 0),
    ]);

    const usersJson = toJson(ListUsersResponseSchema, usersResp) as any;
    const orgsJson = toJson(ListOrganizationsResponseSchema, orgsResp) as any;

    // Fetch instance version separately — may fail if user lacks permissions
    let version: string | null = null;
    try {
      const instanceService = getInstanceService();
      const instanceResp = await instanceService.getInstance({});
      version = instanceResp.instance?.version ?? null;
    } catch {
      // Instance API not accessible (e.g., cloud instance without system permissions)
    }

    return {
      stats: {
        userCount: Number(usersJson.details?.totalResult ?? usersJson.result?.length ?? 0),
        orgCount: Number(orgsJson.details?.totalResult ?? orgsJson.result?.length ?? 0),
        projectCount,
        appCount,
        recentSessions: sessions.sessions ?? [],
        version,
      },
      error: null,
    };
  } catch (e) {
    console.error("Failed to load overview:", e);
    return {
      stats: {
        userCount: 0,
        orgCount: 0,
        projectCount: 0,
        appCount: 0,
        recentSessions: [],
        version: null,
      },
      error: e instanceof Error ? e.message : "Failed to connect to ZITADEL",
    };
  }
}

/**
 * Fetch instance version standalone (for components that only need version).
 */
export async function fetchInstanceVersion(): Promise<string | null> {
  try {
    const instanceService = getInstanceService();
    const resp = await instanceService.getInstance({});
    return resp.instance?.version ?? null;
  } catch {
    return null;
  }
}
