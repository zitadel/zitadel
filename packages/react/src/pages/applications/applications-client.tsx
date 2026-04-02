"use client"

import { useState, useMemo, useEffect, useTransition } from "react"
import { useConsoleRouter } from "../../hooks/use-console-router"
import { StatusBadge } from "../../components/ui/status-badge"
import { Badge } from "../../components/ui/badge"
import { Button } from "../../components/ui/button"
import { FilterBar, type FilterDef } from "../../components/ui/filter-bar"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../components/ui/table"
import {
  AppWindow,
  Globe,
  Server,
  Shield,
} from "lucide-react"
import { TablePagination } from "../../components/ui/table-pagination"
import { TableSkeleton } from "../../components/skeletons/table-skeleton"
import { fetchAllApplications } from "../../api/fetch-all-applications"

interface ApplicationsClientProps {
  applications: any[]
  totalResult: number
  error: string | null
}

function getAppType(app: any) {
  if (app.oidcConfiguration)
    return {
      label: "OIDC",
      icon: Globe,
      className:
        "bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/50 dark:text-blue-400 dark:border-blue-800",
    }
  if (app.apiConfiguration)
    return {
      label: "API",
      icon: Server,
      className:
        "bg-purple-50 text-purple-700 border-purple-200 dark:bg-purple-950/50 dark:text-purple-400 dark:border-purple-800",
    }
  if (app.samlConfiguration)
    return {
      label: "SAML",
      icon: Shield,
      className:
        "bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/50 dark:text-amber-400 dark:border-amber-800",
    }
  // v2beta fallback
  if (app.oidcConfig) return getAppType({ oidcConfiguration: app.oidcConfig })
  if (app.apiConfig) return getAppType({ apiConfiguration: app.apiConfig })
  if (app.samlConfig) return getAppType({ samlConfiguration: app.samlConfig })
  return {
    label: "Unknown",
    icon: AppWindow,
    className: "bg-muted text-muted-foreground",
  }
}

function getClientId(app: any): string {
  return (
    app.oidcConfiguration?.clientId ??
    app.apiConfiguration?.clientId ??
    app.oidcConfig?.clientId ??
    app.apiConfig?.clientId ??
    ""
  )
}

function getStateInfo(state?: string): { label: string; variant: "active" | "inactive" | "destructive" | "warning" } {
  switch (state) {
    case "APPLICATION_STATE_ACTIVE":
      return { label: "Active", variant: "active" }
    case "APPLICATION_STATE_INACTIVE":
      return { label: "Inactive", variant: "inactive" }
    default:
      return { label: "Unknown", variant: "inactive" }
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  })
}

export function ApplicationsClient({
  applications: initialApplications,
  totalResult: initialTotalResult,
  error,
}: ApplicationsClientProps) {
  const router = useConsoleRouter()
  const [searchQuery, setSearchQuery] = useState("")
  const [activeFilters, setActiveFilters] = useState<Record<string, string>>({})
  const [applications, setApplications] = useState(initialApplications)
  const [totalResult, setTotalResult] = useState(initialTotalResult)
  const [page, setPage] = useState(0)
  const [pageSize, setPageSize] = useState(10)
  const [isRefetching, startTransition] = useTransition()

  useEffect(() => {
    startTransition(async () => {
      try {
        const result = await fetchAllApplications(pageSize, page * pageSize)
        setApplications(result.applications)
        setTotalResult(result.totalResult)
      } catch (e) {
        console.error("Failed to refresh applications:", e)
      }
    })
  }, [page, pageSize])

  const filteredApps = useMemo(() => {
    let result = applications
    if (searchQuery) {
      const q = searchQuery.toLowerCase()
      result = result.filter(
        (app: any) =>
          (app.name ?? "").toLowerCase().includes(q) ||
          (getClientId(app) ?? "").toLowerCase().includes(q)
      )
    }
    if (activeFilters.type) {
      result = result.filter((app: any) => {
        const t = getAppType(app).label.toLowerCase()
        return t === activeFilters.type
      })
    }
    if (activeFilters.state) {
      result = result.filter((app: any) => app.state === activeFilters.state)
    }
    if (activeFilters.name) {
      const v = activeFilters.name.toLowerCase()
      result = result.filter((app: any) => (app.name ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.projectid) {
      const v = activeFilters.projectid.toLowerCase()
      result = result.filter((app: any) => (app.projectId ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.clientid) {
      const v = activeFilters.clientid.toLowerCase()
      result = result.filter((app: any) => (getClientId(app) ?? "").toLowerCase().includes(v))
    }
    return result
  }, [applications, searchQuery, activeFilters])

  const appFilters: FilterDef[] = [
    { key: "name", label: "name" },
    { key: "projectid", label: "projectid" },
    { key: "clientid", label: "clientid" },
    {
      key: "type",
      label: "type",
      options: [
        { value: "oidc", label: "oidc" },
        { value: "api", label: "api" },
        { value: "saml", label: "saml" },
      ],
    },
    {
      key: "state",
      label: "state",
      options: [
        { value: "APPLICATION_STATE_ACTIVE", label: "active" },
        { value: "APPLICATION_STATE_INACTIVE", label: "inactive" },
      ],
    },
  ]

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">
            Applications
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage applications across your projects
          </p>
        </div>
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-6 text-center">
          <p className="text-sm font-medium text-destructive">
            Failed to load applications
          </p>
          <p className="text-xs text-muted-foreground mt-1">{error}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">
            Applications
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage applications across all organizations ({applications.length}{" "}
            total)
          </p>
        </div>
      </div>

      {/* Filters */}
      <FilterBar
        searchPlaceholder="Search applications..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
        filters={appFilters}
        activeFilters={activeFilters}
        onFilterChange={(key, value) => {
          setActiveFilters((prev) => {
            if (value === null) {
              const { [key]: _, ...rest } = prev
              return rest
            }
            return { ...prev, [key]: value }
          })
        }}
      />

      {/* Table */}
      {isRefetching ? (
        <div className="rounded-lg border">
          <TableSkeleton
            columns={["Application", "Type", "Client ID", "Status", "Created", "Updated"]}
            rows={pageSize}
          />
        </div>
      ) : filteredApps.length === 0 ? (
        <div className="rounded-lg border">
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <AppWindow className="h-12 w-12 text-muted-foreground/40 mb-4" />
            <p className="text-sm font-medium">
              {searchQuery
                ? "No applications match your search"
                : "No applications found"}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              {searchQuery
                ? "Try adjusting your search query"
                : "Create an application in one of your projects to get started"}
            </p>
          </div>
        </div>
      ) : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead>Application</TableHead>
                <TableHead className="w-[80px]">Type</TableHead>
                <TableHead>Client ID</TableHead>
                <TableHead className="w-[80px]">Status</TableHead>
                <TableHead className="w-[120px]">Created</TableHead>
                <TableHead className="w-[120px]">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredApps.map((app: any) => {
                const appType = getAppType(app)
                const TypeIcon = appType.icon
                const clientId = getClientId(app)
                const stateInfo = getStateInfo(app.state)

                return (
                  <TableRow
                    key={app.applicationId}
                    className="cursor-pointer"
                    onClick={() =>
                      router.push(
                        `/applications/${app.applicationId}?projectId=${app.projectId}`
                      )
                    }
                  >
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 flex-shrink-0">
                          <TypeIcon className="h-4 w-4 text-primary" />
                        </div>
                        <div className="min-w-0">
                          <p className="font-medium truncate">{app.name}</p>
                          {app.projectName && (
                            <p className="text-xs text-muted-foreground truncate">
                              {app.projectName}
                            </p>
                          )}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={`${appType.className} text-xs gap-1`}
                      >
                        {appType.label}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {clientId ? (
                        <code className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded font-mono">
                          {clientId}
                        </code>
                      ) : (
                        <span className="text-xs text-muted-foreground">—</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <StatusBadge variant={stateInfo.variant}>
                        {stateInfo.label}
                      </StatusBadge>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(app.creationDate)}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(app.changeDate)}
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
          <TablePagination
            page={page}
            pageSize={pageSize}
            totalResult={totalResult}
            onPageChange={setPage}
            onPageSizeChange={(size) => { setPageSize(size); setPage(0) }}
          />
        </div>
      )}
    </div>
  )
}
