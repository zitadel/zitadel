import { fetchApplication } from "../../../api/applications"
import { ApplicationDetailClient } from "./app-detail-client"

interface Params {
  appId: string
}

/**
 * Application detail page — server component.
 */
export default async function ApplicationDetailPage({
  params,
  searchParams,
}: {
  params: Promise<Params>
  searchParams: Promise<{ projectId?: string }>
}) {
  const { appId } = await params
  const { projectId } = await searchParams
  let app: any = null
  let error: string | null = null

  if (!projectId) {
    error = "Missing projectId query parameter"
  } else {
    try {
      app = await fetchApplication(projectId, appId)
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load application"
      console.error("Failed to load application:", e)
    }
  }

  return (
    <ApplicationDetailClient
      app={app}
      appId={appId}
      projectId={projectId ?? ""}
      error={error}
    />
  )
}
