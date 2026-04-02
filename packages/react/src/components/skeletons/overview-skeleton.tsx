import { Skeleton } from "../ui/skeleton"
import { Card, CardContent, CardHeader } from "../ui/card"

interface OverviewSkeletonProps {
  /** Number of stat cards (3 for org-scoped, 4 for instance) */
  cardCount?: number
}

/**
 * Skeleton loading state for the Overview dashboard.
 * Renders placeholder stat cards, activity list, and quick actions.
 */
export function OverviewSkeleton({ cardCount = 4 }: OverviewSkeletonProps) {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <Skeleton className="h-7 w-[120px]" />
        <Skeleton className="h-4 w-[200px] mt-2" />
      </div>

      {/* Stats Grid */}
      <div
        className={`grid gap-4 ${
          cardCount === 3
            ? "sm:grid-cols-3"
            : "sm:grid-cols-2 lg:grid-cols-4"
        }`}
      >
        {Array.from({ length: cardCount }).map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <Skeleton className="h-4 w-[90px]" />
              <Skeleton className="h-8 w-8 rounded-lg" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-7 w-[60px]" />
              <Skeleton className="h-3 w-[120px] mt-2" />
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Middle Row */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Recent Activity skeleton */}
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-[160px]" />
            <Skeleton className="h-3.5 w-[200px] mt-1" />
          </CardHeader>
          <CardContent className="space-y-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton className="h-8 w-8 rounded-full flex-shrink-0" />
                <div className="flex-1 space-y-1.5">
                  <Skeleton className="h-3.5 w-[180px]" />
                  <Skeleton className="h-3 w-[120px]" />
                </div>
                <Skeleton className="h-5 w-[60px] rounded-full" />
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Quick Actions skeleton */}
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-[130px]" />
            <Skeleton className="h-3.5 w-[180px] mt-1" />
          </CardHeader>
          <CardContent className="space-y-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-10 w-full rounded-md" />
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
