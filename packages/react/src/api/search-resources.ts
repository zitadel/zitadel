"use server"

import { create, toJson } from "@zitadel/client"
import {
  ListUsersRequestSchema,
  ListUsersResponseSchema,
  UserService,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb"
import {
  ListOrganizationsRequestSchema,
  ListOrganizationsResponseSchema,
  OrganizationService,
} from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import {
  ListProjectsRequestSchema,
  ListProjectsResponseSchema,
  ProjectService,
} from "@zitadel/proto/zitadel/project/v2/project_service_pb"
import {
  ListApplicationsRequestSchema,
  ListApplicationsResponseSchema,
  ApplicationService,
} from "@zitadel/proto/zitadel/application/v2/application_service_pb"
import { createClient } from "@connectrpc/connect"
import { getTransport } from "./transport"

// ─── Types ───────────────────────────────────────────────────────────────────

export interface SearchResult {
  id: string
  name: string
  description: string
  type: "user" | "org" | "project" | "app"
}

export interface SearchResponse {
  users: SearchResult[]
  organizations: SearchResult[]
  projects: SearchResult[]
  applications: SearchResult[]
}

// TEXT_QUERY_METHOD_CONTAINS = 2
const CONTAINS = 2

// ─── Users ───────────────────────────────────────────────────────────────────

export async function searchUsersOmni(query: string): Promise<SearchResult[]> {
  const transport = getTransport()
  const client = createClient(UserService, transport)

  const queries: any[] = []
  if (query.trim()) {
    queries.push({
      query: {
        case: "orQuery",
        value: {
          queries: [
            { query: { case: "userNameQuery", value: { userName: query, method: CONTAINS } } },
            { query: { case: "displayNameQuery", value: { displayName: query, method: CONTAINS } } },
            { query: { case: "emailQuery", value: { emailAddress: query, method: CONTAINS } } },
          ],
        },
      },
    })
  }

  const request = create(ListUsersRequestSchema, {
    query: { limit: 5, asc: true },
    queries,
  })

  const response = await client.listUsers(request)
  const json = toJson(ListUsersResponseSchema, response) as any
  return (json.result ?? []).map((u: any) => ({
    id: u.userId ?? "",
    name: u.human?.profile?.displayName ?? u.username ?? "",
    description: u.human?.email?.email ?? u.username ?? "",
    type: "user" as const,
  }))
}

// ─── Organizations ───────────────────────────────────────────────────────────

export async function searchOrgsOmni(query: string): Promise<SearchResult[]> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)

  const queries: any[] = []
  if (query.trim()) {
    queries.push({
      query: {
        case: "nameQuery",
        value: { name: query, method: CONTAINS },
      },
    })
  }

  const request = create(ListOrganizationsRequestSchema, {
    query: { limit: 5, asc: true },
    queries,
  })

  const response = await client.listOrganizations(request)
  const json = toJson(ListOrganizationsResponseSchema, response) as any
  return (json.result ?? []).map((o: any) => ({
    id: o.id ?? "",
    name: o.name ?? "",
    description: o.primaryDomain ?? "",
    type: "org" as const,
  }))
}

// ─── Projects ────────────────────────────────────────────────────────────────

export async function searchProjectsOmni(query: string): Promise<SearchResult[]> {
  const transport = getTransport()
  const client = createClient(ProjectService, transport)

  const filters: any[] = []
  if (query.trim()) {
    filters.push({
      filter: {
        case: "nameFilter",
        value: { name: query, method: CONTAINS },
      },
    })
  }

  const request = create(ListProjectsRequestSchema, {
    pagination: { limit: 5 },
    filters,
  })

  const response = await client.listProjects(request)
  const json = toJson(ListProjectsResponseSchema, response) as any
  return (json.projects ?? []).map((p: any) => ({
    id: p.projectId ?? "",
    name: p.name ?? "",
    description: p.organizationId ? `Org: ${p.organizationId}` : "",
    type: "project" as const,
  }))
}

// ─── Applications ────────────────────────────────────────────────────────────

export async function searchAppsOmni(query: string): Promise<SearchResult[]> {
  const transport = getTransport()
  const client = createClient(ApplicationService, transport)

  const filters: any[] = []
  if (query.trim()) {
    filters.push({
      filter: {
        case: "nameFilter",
        value: { name: query, method: CONTAINS },
      },
    })
  }

  const request = create(ListApplicationsRequestSchema, {
    pagination: { limit: 5 },
    filters,
  })

  const response = await client.listApplications(request)
  const json = toJson(ListApplicationsResponseSchema, response) as any
  return (json.applications ?? []).map((a: any) => ({
    id: a.applicationId ?? "",
    name: a.name ?? "",
    description: a.oidcConfiguration ? "OIDC" : a.apiConfiguration ? "API" : a.samlConfiguration ? "SAML" : "",
    type: "app" as const,
  }))
}

// ─── Unified search ──────────────────────────────────────────────────────────

/**
 * Search across all resource types in parallel.
 * If `scope` is provided, only searches that resource type.
 */
export async function searchAll(
  query: string,
  scope?: "user" | "org" | "project" | "app"
): Promise<SearchResponse> {
  const empty: SearchResult[] = []

  if (!query.trim() && !scope) {
    return { users: empty, organizations: empty, projects: empty, applications: empty }
  }

  const [users, organizations, projects, applications] = await Promise.all([
    !scope || scope === "user" ? searchUsersOmni(query).catch(() => empty) : Promise.resolve(empty),
    !scope || scope === "org" ? searchOrgsOmni(query).catch(() => empty) : Promise.resolve(empty),
    !scope || scope === "project" ? searchProjectsOmni(query).catch(() => empty) : Promise.resolve(empty),
    !scope || scope === "app" ? searchAppsOmni(query).catch(() => empty) : Promise.resolve(empty),
  ])

  return { users, organizations, projects, applications }
}
