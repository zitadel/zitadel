"use client"

import * as React from "react"
import { SidebarProvider, SidebarInset } from "@zitadel/react/components/ui/sidebar"
import { CloudSidebar } from "@/components/layout/cloud-sidebar"
import { Header } from "@zitadel/react/components/layout/header"

interface InstanceInfo {
  id: string
  name: string
  url: string
}

/**
 * Cloud console layout — sidebar + console header.
 * Uses console's Header component (instance switcher, org switcher, search).
 * Only the sidebar is cloud-specific (multi-instance nav, billing, support).
 */
export function ConsoleLayout({
  children,
  instances,
}: {
  children: React.ReactNode
  instances: InstanceInfo[]
}) {
  return (
    <SidebarProvider defaultOpen={true} open={true}>
      <CloudSidebar instances={instances} />
      <SidebarInset>
        <Header />
        <main className="flex-1 p-6">
          {children}
        </main>
      </SidebarInset>
    </SidebarProvider>
  )
}
