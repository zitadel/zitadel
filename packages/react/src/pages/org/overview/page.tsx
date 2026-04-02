"use client"

import { 
  Users, 
  FolderKanban, 
  AppWindow, 
  UserCog,
  Activity,
  Clock,
  Shield,
  Lock,
  UserPlus,
  Timer,
  ArrowRight,
  AlertTriangle,
  Key,
  Fingerprint,
  KeyRound,
  ShieldCheck
} from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../components/ui/card"
import { Badge } from "../../../components/ui/badge"
import { StatusBadge } from "../../../components/ui/status-badge"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../../../components/ui/tooltip"
import { useAppContext } from "../../../context/app-context"
import { users, projects, applications, sessions, getRoleAssignmentsByOrganization, getActivityLogByOrganization } from "../../../mock-data"
import { ConsoleLink as Link } from "../../../context/link-context"
import { OrganizationSelectorPrompt } from "../../../components/organization-selector-prompt"
import { useMemo } from "react"

// Format relative time
function formatRelativeTime(date: Date): string {
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffMins < 1) return "Just now"
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`
  return date.toLocaleDateString()
}

// Auth method labels and icons
const authMethodConfig = {
  password: { label: "Password", icon: Key },
  passkey: { label: "Passkey", icon: Fingerprint },
  sso: { label: "SSO", icon: KeyRound },
  mfa: { label: "MFA", icon: ShieldCheck },
}

export default function OrgOverviewPage() {
  const { currentOrganization } = useAppContext()

  if (!currentOrganization) {
    return (
      <OrganizationSelectorPrompt 
        title="Select an Organization"
        description="Choose an organization to view its overview"
        targetPath="/org/overview"
      />
    )
  }

  // Filter data for current organization
  const orgUsers = users.filter(user => user.orgId === currentOrganization.id)
  const orgProjects = projects.filter(project => project.orgId === currentOrganization.id)
  const orgApps = applications.filter(app => app.orgId === currentOrganization.id)
  const orgSessions = sessions.filter(session => session.orgId === currentOrganization.id)
  const orgRoles = getRoleAssignmentsByOrganization(currentOrganization.id)

  // Meaningful secondary stats
  const activeUsers = orgUsers.filter(u => u.status === "active").length
  const projectsWithApps = orgProjects.filter(p => 
    orgApps.some(a => a.projectId === p.id)
  ).length
  const appsWithActiveSessions = new Set(
    orgSessions.filter(s => s.status === "active").map(s => {
      const userApps = orgApps.filter(a => 
        orgProjects.some(p => p.orgId === orgUsers.find(u => u.id === s.userId)?.orgId)
      )
      return userApps[0]?.id
    }).filter(Boolean)
  ).size

  // Alerts data
  const alerts = useMemo(() => {
    const lockedUsers = orgUsers.filter(u => u.status === "locked")
    const pendingOld = orgUsers.filter(u => {
      if (u.status !== "pending") return false
      const created = new Date(u.createdAt)
      const daysSince = (Date.now() - created.getTime()) / (1000 * 60 * 60 * 24)
      return daysSince > 14
    })
    const sessionsExpiringToday = orgSessions.filter(s => {
      if (s.status !== "active") return false
      const expiresAt = new Date(s.expiresAt)
      const today = new Date()
      return expiresAt.toDateString() === today.toDateString()
    })
    
    return { lockedUsers, pendingOld, sessionsExpiringToday }
  }, [orgUsers, orgSessions])

  const totalAlerts = alerts.lockedUsers.length + alerts.pendingOld.length + alerts.sessionsExpiringToday.length

  const stats = [
    {
      title: "Users",
      value: orgUsers.length,
      description: `${activeUsers} active`,
      icon: Users,
      href: "/org/users",
    },
    {
      title: "Projects",
      value: orgProjects.length,
      description: `${projectsWithApps} with apps configured`,
      icon: FolderKanban,
      href: "/org/projects",
    },
    {
      title: "Applications",
      value: orgApps.length,
      description: `${appsWithActiveSessions} with active sessions`,
      icon: AppWindow,
      href: "/org/applications",
    },
    {
      title: "Role Assignments",
      value: orgRoles.length,
      description: `${new Set(orgRoles.map(r => r.userId)).size} unique users`,
      icon: UserCog,
      href: "/org/roles",
    },
  ]

  // Recent activity from sessions
  const recentActivity = orgSessions
    .sort((a, b) => new Date(b.lastActivity).getTime() - new Date(a.lastActivity).getTime())
    .slice(0, 5)
    .map(session => {
      const user = orgUsers.find(u => u.id === session.userId)
      const authConfig = authMethodConfig[session.authMethod]
      return {
        id: session.id,
        user: user?.displayName || session.userName,
        action: `Signed in with ${authConfig.label}`,
        authIcon: authConfig.icon,
        time: formatRelativeTime(new Date(session.lastActivity)),
        status: session.status,
      }
    })

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold tracking-tight">Organization Overview</h1>
          {currentOrganization.isDefault && (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Badge variant="secondary" className="cursor-help">Default</Badge>
                </TooltipTrigger>
                <TooltipContent className="max-w-xs">
                  <p>Users without an explicit organization assignment will be placed in this organization. New users and SSO logins without org mapping will join here.</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          )}
        </div>
        <p className="text-muted-foreground">
          {currentOrganization.name}
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Link key={stat.title} href={stat.href}>
            <Card className="transition-all hover:border-foreground hover:shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.title}
                </CardTitle>
                <div className="rounded-lg p-2 bg-muted text-foreground">
                  <stat.icon className="h-4 w-4" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className="text-xs text-muted-foreground">{stat.description}</p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Alerts Panel */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg flex items-center gap-2">
                  <AlertTriangle className="h-5 w-5 text-amber-500" />
                  Alerts
                </CardTitle>
                <CardDescription>Items requiring attention</CardDescription>
              </div>
              {totalAlerts > 0 && (
                <Badge variant="secondary" className="text-amber-600 bg-amber-100">
                  {totalAlerts}
                </Badge>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {totalAlerts === 0 ? (
              <div className="flex flex-col items-center justify-center py-6 text-center">
                <div className="rounded-full bg-emerald-100 p-3 mb-3">
                  <Shield className="h-5 w-5 text-emerald-600" />
                </div>
                <p className="text-sm font-medium">All clear</p>
                <p className="text-xs text-muted-foreground">No items require attention</p>
              </div>
            ) : (
              <div className="space-y-3">
                {alerts.lockedUsers.length > 0 && (
                  <Link 
                    href="/org/users?status=locked"
                    className="flex items-center justify-between p-3 rounded-lg border border-destructive/20 bg-destructive/5 hover:bg-destructive/10 transition-colors group"
                  >
                    <div className="flex items-center gap-3">
                      <Lock className="h-4 w-4 text-destructive" />
                      <div>
                        <p className="text-sm font-medium">{alerts.lockedUsers.length} locked user{alerts.lockedUsers.length !== 1 && "s"}</p>
                        <p className="text-xs text-muted-foreground">
                          {alerts.lockedUsers.slice(0, 2).map(u => u.displayName).join(", ")}
                          {alerts.lockedUsers.length > 2 && ` +${alerts.lockedUsers.length - 2} more`}
                        </p>
                      </div>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                  </Link>
                )}
                
                {alerts.pendingOld.length > 0 && (
                  <Link 
                    href="/org/users?status=pending"
                    className="flex items-center justify-between p-3 rounded-lg border border-amber-500/20 bg-amber-500/5 hover:bg-amber-500/10 transition-colors group"
                  >
                    <div className="flex items-center gap-3">
                      <UserPlus className="h-4 w-4 text-amber-600" />
                      <div>
                        <p className="text-sm font-medium">{alerts.pendingOld.length} pending invite{alerts.pendingOld.length !== 1 && "s"} older than 14 days</p>
                        <p className="text-xs text-muted-foreground">
                          {alerts.pendingOld.slice(0, 2).map(u => u.displayName).join(", ")}
                          {alerts.pendingOld.length > 2 && ` +${alerts.pendingOld.length - 2} more`}
                        </p>
                      </div>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                  </Link>
                )}
                
                {alerts.sessionsExpiringToday.length > 0 && (
                  <Link 
                    href="/sessions"
                    className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors group"
                  >
                    <div className="flex items-center gap-3">
                      <Timer className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="text-sm font-medium">{alerts.sessionsExpiringToday.length} session{alerts.sessionsExpiringToday.length !== 1 && "s"} expiring today</p>
                        <p className="text-xs text-muted-foreground">May require re-authentication</p>
                      </div>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                  </Link>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Recent Activity</CardTitle>
            <CardDescription>Latest authentication events</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-3">
              {recentActivity.length > 0 ? recentActivity.map((activity) => (
                <div key={activity.id} className="flex items-center justify-between border-b pb-3 last:border-0 last:pb-0">
                  <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted">
                      <activity.authIcon className="h-4 w-4 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="text-sm font-medium">{activity.user}</p>
                      <p className="text-xs text-muted-foreground">{activity.action}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <StatusBadge variant={activity.status === "active" ? "active" : "inactive"}>
                      {activity.status === "active" ? "Active" : "Expired"}
                    </StatusBadge>
                    <p className="mt-1 text-xs text-muted-foreground flex items-center justify-end gap-1">
                      <Clock className="h-3 w-3" />
                      {activity.time}
                    </p>
                  </div>
                </div>
              )) : (
                <p className="text-sm text-muted-foreground text-center py-4">No recent activity</p>
              )}
            </div>
            <Link 
              href="/org/activity" 
              className="flex items-center justify-center gap-1 text-sm text-primary hover:text-primary/80 font-medium pt-2 border-t"
            >
              View all activity
              <ArrowRight className="h-4 w-4" />
            </Link>
          </CardContent>
        </Card>
      </div>

      {/* Organization Details */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Organization Details</CardTitle>
          <CardDescription>Key information about this organization</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-3">
            <div className="rounded-lg border p-4">
              <p className="text-sm font-medium text-muted-foreground">Organization ID</p>
              <p className="mt-1 font-mono text-sm">{currentOrganization.id}</p>
            </div>
            <div className="rounded-lg border p-4">
              <p className="text-sm font-medium text-muted-foreground">Created</p>
              <p className="mt-1 text-sm">
                {new Date(currentOrganization.createdAt).toLocaleDateString()}
              </p>
            </div>
            <div className="rounded-lg border p-4">
              <p className="text-sm font-medium text-muted-foreground">Status</p>
              <div className="mt-1">
                <StatusBadge variant={currentOrganization.status === "active" ? "active" : "inactive"}>
                  {currentOrganization.status === "active" ? "Active" : "Inactive"}
                </StatusBadge>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
