import { fetchAllApplications } from "../../api/fetch-all-applications"
import { ApplicationsClient } from "./applications-client"

/**
 * Applications page — server component fetches initial data.
 * The client component handles pagination and search.
 */
export default async function ApplicationsPage() {
  let applications: any[] = []
  let totalResult = 0
  let error: string | null = null

  try {
    const result = await fetchAllApplications()
    applications = result.applications
    totalResult = result.totalResult
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load applications"
  }

  return (
    <ApplicationsClient
      applications={applications}
      totalResult={totalResult}
      error={error}
    />
  )
}
