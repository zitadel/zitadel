import { fetchProject } from "../../../api/projects"
import { fetchApplications } from "../../../api/applications"
import { ProjectDetailClient } from "./project-detail-client"

interface Params {
  projectId: string
}

/**
 * Project detail page — server component.
 */
export default async function ProjectDetailPage({ params }: { params: Promise<Params> }) {
  const { projectId } = await params
  let project: any = null
  let applications: any[] = []
  let error: string | null = null

  try {
    const [proj, apps] = await Promise.all([
      fetchProject(projectId),
      fetchApplications(projectId, 100),
    ])
    project = proj
    applications = apps.applications
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load project"
    console.error("Failed to load project:", e)
  }

  return (
    <ProjectDetailClient
      project={project}
      projectId={projectId}
      applications={applications}
      error={error}
    />
  )
}
