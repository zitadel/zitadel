import {
  Users,
  FolderKanban,
  AppWindow,
  Building2,
  Zap,
  KeyRound,
  UserCog,
  Activity,
  LayoutDashboard,
  Shield,
} from "lucide-react"

/**
 * Shared nav item configuration used by both console and cloud sidebars.
 */
export interface NavItem {
  title: string
  /** Path segment (relative — the sidebar component decides how to prefix) */
  path: string
  icon: React.ComponentType<{ className?: string }>
  /** Required permission to show this nav item */
  permission?: string
  /** Alternative: any of these permissions */
  anyPermission?: string[]
  /** Only show in cloud mode */
  cloudOnly?: boolean
  /** Key into NavCounts for dynamic badge */
  countKey?: string
  /** 'instance' = only when no org selected, 'org' = only when org selected, 'both' = always */
  context?: "instance" | "org" | "both"
}

/**
 * Core nav items — shared between console standalone and cloud.
 */
export const coreNavItems: NavItem[] = [
  {
    title: "Overview",
    path: "/overview",
    icon: LayoutDashboard,
    context: "both",
  },
  {
    title: "Organizations",
    path: "/organizations",
    icon: Building2,
    permission: "org.read",
    context: "instance",
    countKey: "organizations",
  },
  {
    title: "Users",
    path: "/users",
    icon: Users,
    permission: "user.read",
    context: "both",
    countKey: "users",
  },
  {
    title: "Projects",
    path: "/projects",
    icon: FolderKanban,
    permission: "project.read",
    context: "both",
    countKey: "projects",
  },
  {
    title: "Applications",
    path: "/applications",
    icon: AppWindow,
    permission: "project.app.read",
    context: "both",
    countKey: "applications",
  },
  {
    title: "Actions",
    path: "/actions",
    icon: Zap,
    anyPermission: ["iam.action.read", "org.action.read"],
    context: "both",
  },
  {
    title: "Sessions",
    path: "/sessions",
    icon: KeyRound,
    permission: "session.read",
    context: "instance",
  },
  {
    title: "Administrators",
    path: "/administrators",
    icon: UserCog,
    anyPermission: ["iam.member.read", "org.member.read"],
    context: "both",
  },
  {
    title: "Activity Log",
    path: "/activity",
    icon: Activity,
    permission: "events.read",
    context: "both",
  },
  {
    title: "Settings & Policies",
    path: "/settings",
    icon: Shield,
    anyPermission: ["iam.policy.read", "policy.read"],
    context: "both",
  },
]

/**
 * Filter nav items based on org selection context.
 */
export function filterByContext(items: NavItem[], hasOrgSelected: boolean): NavItem[] {
  return items.filter((item) => {
    const ctx = item.context ?? "both"
    if (ctx === "instance" && hasOrgSelected) return false
    if (ctx === "org" && !hasOrgSelected) return false
    return true
  })
}
