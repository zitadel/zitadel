"use client"

import * as React from "react"
import { SidebarProvider, SidebarInset } from "../ui/sidebar"
import { AppSidebar } from "./app-sidebar"
import { Header } from "./header"
import { useAppContext } from "../../context/app-context"

export function ConsoleLayout({ children }: { children: React.ReactNode }) {
  const { isMounted } = useAppContext()
  
  // Show a minimal loading state during SSR and initial hydration
  // This ensures consistent rendering between server and client
  if (!isMounted) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-muted border-t-foreground" />
      </div>
    )
  }

  return (
    <SidebarProvider defaultOpen={true} open={true}>
      <AppSidebar />
      <SidebarInset>
        <Header />
        <main className="flex-1 p-6">
          {children}
        </main>
      </SidebarInset>
    </SidebarProvider>
  )
}
