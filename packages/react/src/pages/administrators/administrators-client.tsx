"use client"

import { useState, useMemo, useEffect } from "react"
import { useConsoleRouter } from "../../hooks/use-console-router"
import { Shield, Search, X, Plus, Building2, FolderKanban } from "lucide-react"
import { Badge } from "../../components/ui/badge"
import { Button } from "../../components/ui/button"
import { Input } from "../../components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../components/ui/table"
import { RequirePermission } from "../../components/auth/require-permission"
import { TableSkeleton } from "../../components/skeletons/table-skeleton"
import { TablePagination } from "../../components/ui/table-pagination"
import { useAppContext } from "../../context/app-context"
import { fetchAdministrators } from "../../api/fetch-administrators"

interface AdministratorsClientProps {
  administrators: any[]
  totalResult: number
  error: string | null
}

function getInitials(name: string): string {
  return name
    .split(" ")
    .map((s) => s[0])
    .filter(Boolean)
    .slice(0, 2)
    .join("")
    .toUpperCase()
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  })
}

function getResourceLabel(admin: any): { label: string; icon: React.ReactNode } {
  if (admin.instance) {
    return {
      label: "Instance",
      icon: <Shield className="h-3.5 w-3.5" />,
    }
  }
  if (admin.organization) {
    return {
      label: admin.organization.name || "Organization",
      icon: <Building2 className="h-3.5 w-3.5" />,
    }
  }
  if (admin.project) {
    return {
      label: admin.project.name || "Project",
      icon: <FolderKanban className="h-3.5 w-3.5" />,
    }
  }
  if (admin.projectGrant) {
    return {
      label: admin.projectGrant.projectName || "Project Grant",
      icon: <FolderKanban className="h-3.5 w-3.5" />,
    }
  }
  return { label: "Unknown", icon: <Shield className="h-3.5 w-3.5" /> }
}

export function AdministratorsClient({
  administrators: initialAdministrators,
  totalResult: initialTotalResult,
  error,
}: AdministratorsClientProps) {
  const router = useConsoleRouter()
  const { currentOrganization } = useAppContext()
  const [administrators, setAdministrators] = useState(initialAdministrators)
  const [totalResult, setTotalResult] = useState(initialTotalResult)
  const [searchQuery, setSearchQuery] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [page, setPage] = useState(0)
  const [pageSize, setPageSize] = useState(10)

  // Re-fetch when org or pagination changes
  useEffect(() => {
    let cancelled = false
    async function reload() {
      setIsLoading(true)
      try {
        const data = await fetchAdministrators(pageSize, page * pageSize)
        if (!cancelled) {
          setAdministrators(data.administrators)
          setTotalResult(data.totalResult)
        }
      } catch (e) {
        if (!cancelled) console.error("Failed to reload administrators:", e)
      } finally {
        if (!cancelled) setIsLoading(false)
      }
    }
    reload()
    return () => { cancelled = true }
  }, [currentOrganization?.id, page, pageSize])

  const filteredAdmins = useMemo(() => {
    if (!searchQuery) return administrators
    const q = searchQuery.toLowerCase()
    return administrators.filter((admin: any) => {
      const name = admin.user?.displayName?.toLowerCase() ?? ""
      const login = admin.user?.preferredLoginName?.toLowerCase() ?? ""
      const roles = (admin.roles ?? []).join(" ").toLowerCase()
      return name.includes(q) || login.includes(q) || roles.includes(q)
    })
  }, [administrators, searchQuery])

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">
            Administrators
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage administrators and their roles
          </p>
        </div>
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-6 text-center">
          <p className="text-sm font-medium text-destructive">
            Failed to load administrators
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
            Administrators
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage administrators and their roles ({administrators.length} total)
          </p>
        </div>
        <RequirePermission permission="iam.member.write">
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Add Administrator
          </Button>
        </RequirePermission>
      </div>

      {/* Search */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search administrators..."
            className="pl-9 pr-9"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <Button
              variant="ghost"
              size="sm"
              className="absolute right-1 top-1/2 -translate-y-1/2 h-6 w-6 p-0"
              onClick={() => setSearchQuery("")}
            >
              <X className="h-3 w-3" />
            </Button>
          )}
        </div>
      </div>

      {/* Table */}
      {isLoading ? (
        <TableSkeleton
          columns={["Administrator", "Roles", "Resource", "Created", "Updated"]}
          rows={5}
        />
      ) : filteredAdmins.length === 0 ? (
        <div className="rounded-lg border">
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <Shield className="h-12 w-12 text-muted-foreground/40 mb-4" />
            <p className="text-sm font-medium">
              {searchQuery
                ? "No administrators match your search"
                : "No administrators found"}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              {searchQuery
                ? "Try adjusting your search query"
                : "Add your first administrator to get started"}
            </p>
          </div>
        </div>
      ) : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead>Administrator</TableHead>
                <TableHead>Roles</TableHead>
                <TableHead className="w-[160px]">Resource</TableHead>
                <TableHead className="w-[120px]">Created</TableHead>
                <TableHead className="w-[120px]">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredAdmins.map((admin: any, index: number) => {
                const resource = getResourceLabel(admin)
                const userId = admin.user?.id
                return (
                  <TableRow
                    key={`${userId}-${index}`}
                    className="cursor-pointer"
                    onClick={() => {
                      if (userId) {
                        router.push(`/users/${userId}`)
                      }
                    }}
                  >
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-xs font-medium flex-shrink-0">
                          {getInitials(admin.user?.displayName || admin.user?.preferredLoginName || "?")}
                        </div>
                        <div className="min-w-0">
                          <p className="font-medium truncate">
                            {admin.user?.displayName || admin.user?.preferredLoginName || "Unknown"}
                          </p>
                          {admin.user?.preferredLoginName && admin.user?.displayName && (
                            <p className="text-xs text-muted-foreground truncate">
                              {admin.user.preferredLoginName}
                            </p>
                          )}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-wrap gap-1">
                        {(admin.roles ?? []).map((role: string) => (
                          <Badge
                            key={role}
                            variant="outline"
                            className="text-xs"
                          >
                            {role}
                          </Badge>
                        ))}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1.5 text-sm">
                        <span className="text-muted-foreground">
                          {resource.icon}
                        </span>
                        <span className="truncate">{resource.label}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(admin.creationDate)}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(admin.changeDate)}
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
