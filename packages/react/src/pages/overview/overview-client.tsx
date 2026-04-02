"use client"

import { useState, useEffect } from "react"
import {
  Users,
  Building2,
  Shield,
  ArrowRight,
  FolderKanban,
  AppWindow,
  Activity,
  AlertTriangle,
  Lock,
  UserPlus,
  Clock,
} from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card"
import { Badge } from "../../components/ui/badge"
import { ConsoleLink as Link } from "../../context/link-context"
import { useDeployment } from "../../context/deployment"
import { useAppContext } from "../../context/app-context"
import { fetchOverviewStats, type OverviewStats } from "../../api/fetch-overview"
import { OverviewSkeleton } from "../../components/skeletons/overview-skeleton"

interface OverviewClientProps {
  initialStats: OverviewStats
  initialError: string | null
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffHours < 1) return "Just now"
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`
  return date.toLocaleDateString()
}

function getAuthMethod(session: any): string {
  const factors = session.factors ?? {}
  if (factors.webAuthN) return "Passkey"
  if (factors.otpSms || factors.otpEmail) return "OTP"
  if (factors.intent) return "Intent"
  if (factors.password) return "Password"
  if (factors.user) return "Session"
  return "Unknown"
}

export function OverviewClient({ initialStats, initialError }: OverviewClientProps) {
  const { isSelfHosted } = useDeployment()
  const { currentOrganization } = useAppContext()
  const [stats, setStats] = useState<OverviewStats>(initialStats)
  const [error, setError] = useState<string | null>(initialError)
  const [loading, setLoading] = useState(false)

  // Re-fetch when org context changes
  useEffect(() => {
    let cancelled = false

    async function refetch() {
      setLoading(true)
      try {
        const result = await fetchOverviewStats(currentOrganization?.id ?? null)
        if (!cancelled) {
          setStats(result.stats)
          setError(result.error)
        }
      } catch (e) {
        if (!cancelled) {
          setError(e instanceof Error ? e.message : "Failed to load")
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    refetch()
    return () => { cancelled = true }
  }, [currentOrganization?.id])

  if (error && !stats.userCount) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Overview</h1>
          <p className="text-muted-foreground">
            Your ZITADEL dashboard
          </p>
        </div>
        <Card className="border-destructive/50 bg-destructive/5">
          <CardHeader>
            <CardTitle className="text-lg text-destructive flex items-center gap-2">
              <Shield className="h-5 w-5" />
              Connection Error
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <p className="text-sm">{error}</p>
            <div className="rounded-lg bg-muted p-4 text-sm space-y-2">
              <p className="font-medium">To connect, set the following in your <code>.env</code> file:</p>
              <pre className="text-xs bg-background rounded p-3 overflow-x-auto">
{`ZITADEL_INSTANCE_URL=https://your-instance.zitadel.cloud
ZITADEL_PAT=your-personal-access-token`}
              </pre>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  const recentSessions = stats.recentSessions

  // Contextual subtitle
  const subtitle = currentOrganization
    ? currentOrganization.name
    : isSelfHosted
      ? "Self-hosted ZITADEL instance"
      : "ZITADEL Cloud instance"

  const statCards = [
    {
      title: "Total Users",
      value: stats.userCount,
      description: currentOrganization ? `In ${currentOrganization.name}` : `${stats.userCount} active`,
      icon: Users,
      href: "/users",
    },
    // Show org card only when no org is selected (instance-level view)
    ...(!currentOrganization ? [{
      title: "Organizations",
      value: stats.orgCount,
      description: `${stats.orgCount} with active users`,
      icon: Building2,
      href: "/organizations",
    }] : []),
    {
      title: "Projects",
      value: stats.projectCount,
      description: "With apps configured",
      icon: FolderKanban,
      href: "/projects",
    },
    {
      title: "Applications",
      value: stats.appCount,
      description: "OIDC, API & SAML apps",
      icon: AppWindow,
      href: "/applications",
    },
  ]

  if (loading && !stats.userCount) {
    return <OverviewSkeleton cardCount={currentOrganization ? 3 : 4} />
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Overview</h1>
        <p className="text-muted-foreground">{subtitle}</p>
      </div>

      {/* Stats Grid — always 4 cols (or 3 with org), never stack */}
      <div className="grid gap-4 grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => (
          <Link key={stat.title} href={stat.href}>
            <Card className="transition-all hover:border-foreground hover:shadow-sm h-full">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.title}
                </CardTitle>
                <div className="rounded-lg p-2 bg-muted text-foreground">
                  <stat.icon className="h-4 w-4" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">{stat.description}</p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Alerts + Recent Activity — side by side like the prototype */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Alerts */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg flex items-center gap-2">
                  <AlertTriangle className="h-5 w-5 text-amber-500" />
                  Alerts
                </CardTitle>
                <CardDescription>Items requiring attention</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {/* Alert items — when we have real data these will be dynamic */}
            <div className="flex flex-col items-center justify-center py-6 text-center">
              <div className="rounded-full bg-emerald-100 p-3 mb-3">
                <Shield className="h-5 w-5 text-emerald-600" />
              </div>
              <p className="text-sm font-medium">All clear</p>
              <p className="text-xs text-muted-foreground">No items require attention</p>
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-lg">Recent Activity</CardTitle>
            <CardDescription>Latest authentication events</CardDescription>
          </CardHeader>
          <CardContent>
            {recentSessions.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-4">
                No recent activity
              </p>
            ) : (
              <div className="space-y-0">
                {recentSessions.map((session: any, idx: number) => {
                  const factors = session.factors ?? {}
                  const user = factors.user ?? {}
                  const userId = user.id
                  const displayName = user.displayName || user.loginName || "Unknown"
                  const authMethod = getAuthMethod(session)

                  const content = (
                    <div className="flex items-center gap-3 py-2 border-b last:border-0 hover:bg-muted/50 transition-colors cursor-pointer rounded-md px-2 -mx-2 group">
                      <Activity className="h-4 w-4 text-muted-foreground shrink-0" />
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-medium truncate">{displayName}</p>
                        <p className="text-xs text-muted-foreground">Signed in with {authMethod}</p>
                      </div>
                      <div className="text-right shrink-0">
                        <Badge
                          variant="outline"
                          className="bg-emerald-500/10 text-emerald-700 border-emerald-500/30 text-xs"
                        >
                          active
                        </Badge>
                        <p className="text-xs text-muted-foreground mt-0.5 flex items-center justify-end gap-1">
                          <Clock className="h-3 w-3" />
                          {formatDate(session.creationDate)}
                        </p>
                      </div>
                    </div>
                  )

                  return userId ? (
                    <Link key={session.id || idx} href={`/users/${userId}`}>
                      {content}
                    </Link>
                  ) : (
                    <div key={session.id || idx}>{content}</div>
                  )
                })}
              </div>
            )}
            {recentSessions.length > 0 && (
              <Link
                href="/sessions"
                className="flex items-center justify-center gap-1 text-sm text-primary hover:text-primary/80 font-medium pt-2 border-t mt-2"
              >
                View all sessions
                <ArrowRight className="h-4 w-4" />
              </Link>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
