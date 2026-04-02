"use server"

import { create } from "@zitadel/client"
import {
  ListProjectsRequestSchema,
  ProjectService,
} from "@zitadel/proto/zitadel/project/v2/project_service_pb"
import {
  ProjectSearchFilterSchema,
  ProjectOrganizationIDFilterSchema,
  ProjectOrganizationIDFilter_Type,
} from "@zitadel/proto/zitadel/project/v2/query_pb"
import {
  ListApplicationsRequestSchema,
  ApplicationService,
} from "@zitadel/proto/zitadel/application/v2/application_service_pb"
import {
  ApplicationSearchFilterSchema,
  ProjectIDFilterSchema,
} from "@zitadel/proto/zitadel/application/v2/application_pb"
import { PaginationRequestSchema } from "@zitadel/proto/zitadel/filter/v2/filter_pb"
import { createClient } from "@connectrpc/connect"
import { getTransport } from "./transport"

/**
 * Count projects owned by a given organization.
 * Uses ListProjects with ProjectOrganizationIDFilter (OWNED type) and limit 0
 * to get just the total count from pagination.
 */
export async function countOrgProjects(organizationId: string): Promise<number> {
  const transport = getTransport()
  const client = createClient(ProjectService, transport)

  const orgFilter = create(ProjectOrganizationIDFilterSchema, {
    organizationId,
    type: ProjectOrganizationIDFilter_Type.OWNED,
  })

  const filter = create(ProjectSearchFilterSchema, {
    filter: {
      case: "organizationIdFilter",
      value: orgFilter,
    },
  })

  const request = create(ListProjectsRequestSchema, {
    pagination: create(PaginationRequestSchema, { limit: 0 }),
    filters: [filter],
  })

  const response = await client.listProjects(request)
  return Number(response.pagination?.totalResult ?? 0)
}

/**
 * Count applications across all projects owned by a given organization.
 * First lists projects for the org, then counts apps for each.
 * If there are no projects, returns 0.
 */
export async function countOrgApplications(organizationId: string): Promise<number> {
  const transport = getTransport()
  const projectClient = createClient(ProjectService, transport)
  const appClient = createClient(ApplicationService, transport)

  // First get all project IDs for this org
  const orgFilter = create(ProjectOrganizationIDFilterSchema, {
    organizationId,
    type: ProjectOrganizationIDFilter_Type.OWNED,
  })

  const projectFilter = create(ProjectSearchFilterSchema, {
    filter: {
      case: "organizationIdFilter",
      value: orgFilter,
    },
  })

  const projectRequest = create(ListProjectsRequestSchema, {
    filters: [projectFilter],
  })

  const projectResponse = await projectClient.listProjects(projectRequest)
  const projects = projectResponse.projects ?? []

  if (projects.length === 0) return 0

  // Count apps across all projects
  let totalApps = 0
  for (const project of projects) {
    const pidFilter = create(ProjectIDFilterSchema, {
      projectId: project.projectId,
    })
    const appFilter = create(ApplicationSearchFilterSchema, {
      filter: {
        case: "projectIdFilter",
        value: pidFilter,
      },
    })
    const appRequest = create(ListApplicationsRequestSchema, {
      pagination: create(PaginationRequestSchema, { limit: 0 }),
      filters: [appFilter],
    })
    const appResponse = await appClient.listApplications(appRequest)
    totalApps += Number(appResponse.pagination?.totalResult ?? 0)
  }

  return totalApps
}
