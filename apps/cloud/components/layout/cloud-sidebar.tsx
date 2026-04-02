"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
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
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarHeader,
  SidebarFooter,
  SidebarSeparator,
} from "@zitadel/react/components/ui/sidebar"
import { Badge } from "@zitadel/react/components/ui/badge"
import { useAppContext } from "@zitadel/react/context/app-context"
import { useNavCounts } from "@zitadel/react/hooks/use-nav-counts"
import { coreNavItems, filterByContext } from "@zitadel/react/components/layout/nav-items"

interface InstanceInfo {
  id: string
  name: string
  url: string
}

/**
 * Cloud-specific sidebar — imports shared nav items from console,
 * prefixes links with /console/instances/{id}/, and adds cloud-only footer items.
 */
export function CloudSidebar({ instances }: { instances: InstanceInfo[] }) {
  const pathname = usePathname()
  const { currentOrganization } = useAppContext()

  // Detect instance from URL: /console/instances/{id}/...
  const instanceMatch = pathname.match(/^\/console\/instances\/([^/]+)/)
  const currentInstanceId = instanceMatch?.[1]
  const currentInstance = instances.find((i) => i.id === currentInstanceId)
  const instanceBase = currentInstanceId ? `/console/instances/${currentInstanceId}` : null

  // Only fetch/show counts when an instance is selected
  const navCounts = useNavCounts(currentInstanceId ? (currentOrganization?.id ?? null) : undefined)

  const hasOrgSelected = currentOrganization != null
  const visibleItems = filterByContext(coreNavItems, hasOrgSelected)

  return (
    <Sidebar>
      <SidebarHeader>
        <div className="flex items-center gap-2 px-3 py-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-foreground text-background text-sm font-bold">
            Z
          </div>
          <span className="font-semibold text-lg">ZITADEL</span>
        </div>
      </SidebarHeader>

      <SidebarContent>
        {/* Cloud top-level */}
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild isActive={pathname === "/console/getting-started"}>
                  <Link href="/console/getting-started">
                    <Sparkles className="h-4 w-4" />
                    <span>Getting Started</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild isActive={pathname === "/console"}>
                  <Link href="/console">
                    <Server className="h-4 w-4" />
                    <span>All Instances</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {/* Instance section — shared nav items from console */}
        {instances.length > 0 && (
          <>
            <SidebarSeparator />
            <SidebarGroup>
              <SidebarGroupLabel>
                {currentInstance?.name || instances[0]?.name || "Instance"}
              </SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  {visibleItems.map((item) => {
                    const base = instanceBase || `/console/instances/${instances[0]?.id}`
                    const href = `${base}/${item.path}`
                    return (
                      <SidebarMenuItem key={item.path}>
                        <SidebarMenuButton asChild isActive={pathname === href || pathname.startsWith(href + "/")}>
                          <Link href={href} className="flex items-center justify-between">
                            <span className="flex items-center gap-2">
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
                    )
                  })}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </>
        )}
      </SidebarContent>

      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild isActive={pathname === "/console/billing" || pathname.startsWith("/console/billing/")}>
              <Link href="/console/billing">
                <CreditCard className="h-4 w-4" />
                <span>Billing & Usage</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton asChild isActive={pathname === "/console/support"}>
              <Link href="/console/support">
                <LifeBuoy className="h-4 w-4" />
                <span>Support</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
