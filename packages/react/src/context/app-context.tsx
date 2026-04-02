"use client"

import React, { createContext, useContext, useState, useCallback } from "react"

/**
 * Backward-compatible instance type for prototype pages that still
 * reference currentInstance. In self-hosted mode, we provide a
 * synthetic "local" instance. This will be removed once all pages
 * are migrated to use the API layer.
 */
interface SyntheticInstance {
  id: string
  name: string
  domain: string
  status: "active" | "inactive"
  hostingType: "cloud" | "self-hosted"
}

interface AppContextType {
  // Organization (will use proto org type once all pages are migrated)
  currentOrganization: any | null
  setCurrentOrganization: (org: any | null) => void

  // Backward-compatible instance fields for prototype pages
  currentInstance: SyntheticInstance | null
  setCurrentInstance: (instance: SyntheticInstance | null) => void
  availableInstances: SyntheticInstance[]

  // Backward-compatible org list for prototype pages
  availableOrganizations: any[]

  // Context switching (no-ops for now)
  setContextFromUser: (orgId: string) => void
  setContextFromProject: (orgId: string) => void
  setContextFromApplication: (orgId: string) => void

  // Loading state
  isLoading: boolean

  // Mounted state for hydration
  isMounted: boolean
}

const AppContext = createContext<AppContextType | undefined>(undefined)

/** Synthetic instance for self-hosted mode so prototype pages don't crash */
const SELF_HOSTED_INSTANCE: SyntheticInstance = {
  id: "self-hosted",
  name: "Local Instance",
  domain: process.env.ZITADEL_INSTANCE_URL ?? "localhost",
  status: "active",
  hostingType: "self-hosted",
}

export function AppProvider({ children, initialOrganizations = [] }: { children: React.ReactNode; initialOrganizations?: any[] }) {
  const [currentOrganization, setCurrentOrganizationState] = useState<any | null>(null)

  // In self-hosted mode, always provide a synthetic instance
  const [currentInstance] = useState<SyntheticInstance>(SELF_HOSTED_INSTANCE)

  const setCurrentOrganization = useCallback((org: any | null) => {
    setCurrentOrganizationState(org)
  }, [])

  // No-op setters for backward compatibility
  const setCurrentInstance = useCallback((_instance: SyntheticInstance | null) => {}, [])
  const noopContextSwitch = useCallback((_id: string) => {}, [])

  const value: AppContextType = {
    currentOrganization,
    setCurrentOrganization,
    currentInstance,
    setCurrentInstance,
    availableInstances: [SELF_HOSTED_INSTANCE],
    availableOrganizations: initialOrganizations,
    setContextFromUser: noopContextSwitch,
    setContextFromProject: noopContextSwitch,
    setContextFromApplication: noopContextSwitch,
    isLoading: false,
    isMounted: true,
  }

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>
}

export function useAppContext() {
  const context = useContext(AppContext)
  if (context === undefined) {
    throw new Error("useAppContext must be used within an AppProvider")
  }
  return context
}
