import { fetchOverviewStats } from "../../api/fetch-overview"
import { OverviewClient } from "./overview-client"

/**
 * Overview dashboard — server component fetches initial stats.
 * The client component re-fetches when org context changes.
 */
export default async function OverviewPage() {
  const { stats, error } = await fetchOverviewStats()

  return (
    <OverviewClient
      initialStats={stats}
      initialError={error}
    />
  )
}
