import { fetchAllSessions } from "../../api/all-sessions"
import { SessionsClient } from "./sessions-client"

/**
 * Sessions list page — server component that fetches sessions.
 */
export default async function SessionsPage() {
  let initialSessions: any[] = []
  let totalSessions = 0
  let error: string | null = null

  try {
    const result = await fetchAllSessions(20, 0)
    initialSessions = result.sessions
    totalSessions = result.totalResult
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load sessions"
    console.error("Failed to load sessions:", e)
  }

  return (
    <SessionsClient
      initialSessions={initialSessions}
      totalSessions={totalSessions}
      error={error}
    />
  )
}
