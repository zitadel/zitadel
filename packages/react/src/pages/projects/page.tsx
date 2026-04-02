import { fetchProjects } from "../../api/projects"
import { listOrganizations } from "../../api/organizations"
import { toJson } from "@zitadel/client"
import { ListOrganizationsResponseSchema } from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { ProjectsClient } from "./projects-client"

/**
 * Projects list page — server component that fetches projects and organizations.
 */
export default async function ProjectsPage() {
  let projects: any[] = []
  let organizations: any[] = []
  let totalResult = 0
  let error: string | null = null

  try {
    const [projectsResult, orgsResponse] = await Promise.all([
      fetchProjects(10),
      listOrganizations({ pageSize: 100 }),
    ])
    projects = projectsResult.projects
    totalResult = projectsResult.totalResult
    const orgsJson = toJson(ListOrganizationsResponseSchema, orgsResponse) as any
    organizations = orgsJson.result ?? []
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load projects"
    console.error("Failed to load projects:", e)
  }

  return (
    <ProjectsClient
      projects={projects}
      organizations={organizations}
      totalResult={totalResult}
      error={error}
    />
  )
}
