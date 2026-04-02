"use client"

import { usePathname } from "next/navigation"
import {
  BarChart3,
  Sparkles,
  CreditCard,
  LifeBuoy,
  Server,
} from "lucide-react"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarHeader,
  SidebarFooter,
  SidebarSeparator,
} from "../ui/sidebar"
import { useAppContext } from "../../context/app-context"
import { usePermissions } from "../../context/permissions"
import { useDeployment } from "../../context/deployment"
import { Badge } from "../ui/badge"
import { AccountDropdown } from "./account-dropdown"
import { useNavCounts } from "../../hooks/use-nav-counts"
import { type NavItem, coreNavItems, filterByContext } from "./nav-items"
import { ConsoleLink as Link } from "../../context/link-context"

const cloudOnlyItems: NavItem[] = [
  {
    title: "Instances",
    path: "/instances",
    icon: Server,
    cloudOnly: true,
  },
  {
    title: "Analytics",
    path: "/analytics",
    icon: BarChart3,
    cloudOnly: true,
  },
  {
    title: "Billing",
    path: "/billing",
    icon: CreditCard,
    cloudOnly: true,
  },
  {
    title: "Support",
    path: "/support",
    icon: LifeBuoy,
    cloudOnly: true,
  },
]

export function AppSidebar() {
  const pathname = usePathname()
  const { currentOrganization } = useAppContext()
  const { can, canAny } = usePermissions()
  const { isCloud } = useDeployment()
  const navCounts = useNavCounts(currentOrganization?.id)

  const hasOrgSelected = currentOrganization != null

  const isVisible = (item: NavItem): boolean => {
    if (item.cloudOnly && !isCloud) return false
    if (item.permission && !can(item.permission)) return false
    if (item.anyPermission && !canAny(item.anyPermission)) return false
    return true
  }

  const visibleItems = filterByContext(coreNavItems, hasOrgSelected).filter(isVisible)
  const visibleCloudItems = cloudOnlyItems.filter(isVisible)

  return (
    <Sidebar className="border-r-0">
      <SidebarHeader className="px-4 py-4">
        <Link href="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground font-bold text-sm">
            Z
          </div>
          <span className="font-semibold text-lg tracking-tight">ZITADEL</span>
        </Link>
      </SidebarHeader>

      <SidebarContent className="px-2">
        {/* Getting Started Link */}
        <SidebarGroup className="py-1">
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  asChild
                  isActive={pathname === "/getting-started"}
                  className="h-9"
                >
                  <Link href="/getting-started" className="flex items-center gap-2.5">
                    <Sparkles className="h-4 w-4" />
                    <span className="font-medium">Getting Started</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator className="my-2" />

        {/* Main Navigation */}
        <SidebarGroup className="py-1">
          <SidebarGroupContent>
            <SidebarMenu>
              {visibleItems.map((item) => (
                <SidebarMenuItem key={item.path}>
                  <SidebarMenuButton
                    asChild
                    isActive={pathname === item.path || pathname.startsWith(item.path + "/")}
                    className="h-9"
                  >
                    <Link href={item.path} className="flex items-center justify-between">
                      <span className="flex items-center gap-2.5">
                        <item.icon className="h-4 w-4" />
                        <span>{item.title}</span>
                      </span>
                      {item.countKey && navCounts && navCounts[item.countKey as keyof typeof navCounts] > 0 && (
                        <Badge variant="secondary" className="text-[10px] px-1.5 py-0 h-5 font-normal tabular-nums">
                          {navCounts[item.countKey as keyof typeof navCounts]}
                        </Badge>
                      )}
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {/* Cloud-only items */}
        {isCloud && visibleCloudItems.length > 0 && (
          <>
            <SidebarSeparator className="my-2" />
            <SidebarGroup className="py-1">
              <div className="px-2 py-1.5 text-xs font-medium uppercase tracking-wider text-muted-foreground">
                Cloud
              </div>
              <SidebarGroupContent className="mt-1">
                <SidebarMenu>
                  {visibleCloudItems.map((item) => (
                    <SidebarMenuItem key={item.path}>
                      <SidebarMenuButton
                        asChild
                        isActive={pathname === item.path}
                        className="h-9"
                      >
                        <Link href={item.path} className="flex items-center gap-2.5">
                          <item.icon className="h-4 w-4" />
                          <span>{item.title}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </>
        )}
      </SidebarContent>

      <SidebarFooter className="p-3">
        <AccountDropdown />
      </SidebarFooter>
    </Sidebar>
  )
}
