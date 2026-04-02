"use client"

import { useState, useMemo, useEffect, useTransition } from "react"
import { useConsoleRouter } from "../../hooks/use-console-router"
import { StatusBadge } from "../../components/ui/status-badge"
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
import { FolderKanban, Plus } from "lucide-react"
import { RequirePermission } from "../../components/auth/require-permission"
import { TablePagination } from "../../components/ui/table-pagination"
import { TableSkeleton } from "../../components/skeletons/table-skeleton"
import { fetchProjects } from "../../api/projects"

interface ProjectsClientProps {
  projects: any[]
  organizations: any[]
  totalResult: number
  error: string | null
}

function getProjectState(project: any): { label: string; variant: "active" | "inactive" | "destructive" | "warning" } {
  const state = project.state ?? "PROJECT_STATE_UNSPECIFIED"
  const labels: Record<string, { label: string; variant: "active" | "inactive" | "destructive" | "warning" }> = {
    PROJECT_STATE_ACTIVE: { label: "Active", variant: "active" },
    PROJECT_STATE_INACTIVE: { label: "Inactive", variant: "inactive" },
    PROJECT_STATE_UNSPECIFIED: { label: "Unknown", variant: "inactive" },
  }
  return labels[state] ?? { label: state, variant: "inactive" }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  })
}

export function ProjectsClient({
  projects: initialProjects,
  organizations,
  totalResult: initialTotalResult,
  error,
}: ProjectsClientProps) {
  const router = useConsoleRouter()
  const [searchQuery, setSearchQuery] = useState("")
  const [activeFilters, setActiveFilters] = useState<Record<string, string>>({})
  const [projects, setProjects] = useState(initialProjects)
  const [totalResult, setTotalResult] = useState(initialTotalResult)
  const [page, setPage] = useState(0)
  const [pageSize, setPageSize] = useState(10)
  const [isRefetching, startTransition] = useTransition()

  useEffect(() => {
    startTransition(async () => {
      try {
        const result = await fetchProjects(pageSize, page * pageSize)
        setProjects(result.projects)
        setTotalResult(result.totalResult)
      } catch (e) {
        console.error("Failed to refresh projects:", e)
      }
    })
  }, [page, pageSize])

  // Deduplicate projects by projectId (API returns duplicates for granted projects)
  const uniqueProjects = useMemo(() => {
    const seen = new Set<string>()
    return projects.filter((p: any) => {
      if (seen.has(p.projectId)) return false
      seen.add(p.projectId)
      return true
    })
  }, [projects])

  // Build org ID -> name lookup map
  const orgNameMap = useMemo(() => {
    const map: Record<string, string> = {}
    for (const org of organizations) {
      const orgId = org.organizationId ?? org.id
      if (orgId && org.name) {
        map[orgId] = org.name
      }
    }
    return map
  }, [organizations])

  const filteredProjects = useMemo(() => {
    let result = uniqueProjects
    if (searchQuery) {
      const q = searchQuery.toLowerCase()
      result = result.filter((p: any) => {
        const orgName = orgNameMap[p.organizationId] ?? ""
        return (
          (p.name ?? "").toLowerCase().includes(q) ||
          orgName.toLowerCase().includes(q)
        )
      })
    }
    if (activeFilters.state) {
      result = result.filter((p: any) => (p.state ?? "PROJECT_STATE_UNSPECIFIED") === activeFilters.state)
    }
    if (activeFilters.name) {
      const v = activeFilters.name.toLowerCase()
      result = result.filter((p: any) => (p.name ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.org) {
      const v = activeFilters.org.toLowerCase()
      result = result.filter((p: any) => {
        const orgName = orgNameMap[p.organizationId] ?? ""
        return orgName.toLowerCase().includes(v) ||
          (p.organizationId ?? "").toLowerCase().includes(v)
      })
    }
    return result
  }, [uniqueProjects, searchQuery, orgNameMap, activeFilters])

  const projectFilters: FilterDef[] = [
    { key: "name", label: "name" },
    { key: "org", label: "org" },
    {
      key: "state",
      label: "state",
      options: [
        { value: "PROJECT_STATE_ACTIVE", label: "active" },
        { value: "PROJECT_STATE_INACTIVE", label: "inactive" },
      ],
    },
  ]

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Projects</h1>
          <p className="text-sm text-muted-foreground">
            Manage your ZITADEL projects
          </p>
        </div>
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-6 text-center">
          <p className="text-sm font-medium text-destructive">
            Failed to load projects
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
          <h1 className="text-2xl font-semibold tracking-tight">Projects</h1>
          <p className="text-sm text-muted-foreground">
            Manage projects across all organizations ({uniqueProjects.length}{" "}
            total)
          </p>
        </div>
        <RequirePermission permission="project.write">
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Create Project
          </Button>
        </RequirePermission>
      </div>

      {/* Filters */}
      <FilterBar
        searchPlaceholder="Search projects..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
        filters={projectFilters}
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

      {/* Table — Columns: Project | Organization | Status | Created | Updated */}
      {isRefetching ? (
        <div className="rounded-lg border">
          <TableSkeleton
            columns={["Project", "Organization", "Status", "Created", "Updated"]}
            rows={pageSize}
          />
        </div>
      ) : filteredProjects.length === 0 ? (
        <div className="rounded-lg border">
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <FolderKanban className="h-12 w-12 text-muted-foreground/40 mb-4" />
            <p className="text-sm font-medium">
              {searchQuery
                ? "No projects match your search"
                : "No projects found"}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              {searchQuery
                ? "Try adjusting your search query"
                : "Create your first project to get started"}
            </p>
          </div>
        </div>
      ) : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead>Project</TableHead>
                <TableHead className="w-[140px]">Organization</TableHead>
                <TableHead className="w-[80px]">Status</TableHead>
                <TableHead className="w-[120px]">Created</TableHead>
                <TableHead className="w-[120px]">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredProjects.map((project: any) => {
                const stateInfo = getProjectState(project)
                const orgName = orgNameMap[project.organizationId]
                return (
                  <TableRow
                    key={project.projectId}
                    className="cursor-pointer"
                    onClick={() =>
                      router.push(`/projects/${project.projectId}`)
                    }
                  >
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 flex-shrink-0">
                          <FolderKanban className="h-4 w-4 text-primary" />
                        </div>
                        <p className="font-medium truncate">{project.name}</p>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm">
                      {orgName ? (
                        <span className="text-foreground">{orgName}</span>
                      ) : (
                        <span className="text-muted-foreground">—</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <StatusBadge variant={stateInfo.variant}>
                        {stateInfo.label}
                      </StatusBadge>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(project.creationDate)}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(project.changeDate)}
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
