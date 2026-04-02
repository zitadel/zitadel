"use client"

import { useState, useMemo, useEffect, useTransition } from "react"
import { useConsoleRouter } from "../../hooks/use-console-router"
import {
  Users,
  Plus,
  Search,
  X,
  Mail,
  Shield,
  User,
  Loader2,
  Bot,
} from "lucide-react"
import { Button } from "../../components/ui/button"
import { StatusBadge } from "../../components/ui/status-badge"
import { FilterBar, type FilterDef } from "../../components/ui/filter-bar"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../components/ui/table"
import { RequirePermission } from "../../components/auth/require-permission"
import { AddUserSheet } from "../../components/users/add-user-sheet"
import { TablePagination } from "../../components/ui/table-pagination"
import { TableSkeleton } from "../../components/skeletons/table-skeleton"
import { useAppContext } from "../../context/app-context"
import { fetchUsers } from "../../api/fetch-users"

interface UsersClientProps {
  users: any[]
  organizations: any[]
  totalResult: number
  error: string | null
}

/**
 * Extract display info from a toJson()-converted User object.
 * In proto3 JSON, oneof `type` serializes as `human` or `machine` key.
 */
function getUserDisplayInfo(user: any) {
  if (user.human) {
    const human = user.human
    return {
      displayName:
        human?.profile?.displayName ||
        `${human?.profile?.givenName ?? ""} ${human?.profile?.familyName ?? ""}`.trim() ||
        user.username ||
        "Unknown",
      email: human?.email?.email ?? "",
      kind: "human" as const,
    }
  }

  if (user.machine) {
    return {
      displayName: user.machine?.name || user.username || "Machine User",
      email: "",
      kind: "machine" as const,
    }
  }

  return {
    displayName: user.username || "Unknown",
    email: "",
    kind: "unknown" as const,
  }
}

function getUserState(state?: string): { label: string; variant: "active" | "inactive" | "destructive" | "warning" } {
  switch (state) {
    case "USER_STATE_ACTIVE":
      return { label: "Active", variant: "active" }
    case "USER_STATE_INACTIVE":
      return { label: "Inactive", variant: "inactive" }
    case "USER_STATE_LOCKED":
      return { label: "Locked", variant: "destructive" }
    case "USER_STATE_INITIAL":
      return { label: "Initial", variant: "warning" }
    default:
      return { label: "Unknown", variant: "inactive" }
  }
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

export function UsersClient({
  users: initialUsers,
  organizations: initialOrgs,
  totalResult: initialTotalResult,
  error,
}: UsersClientProps) {
  const router = useConsoleRouter()
  const { currentOrganization } = useAppContext()
  const [searchQuery, setSearchQuery] = useState("")
  const [activeFilters, setActiveFilters] = useState<Record<string, string>>({})
  const [addUserOpen, setAddUserOpen] = useState(false)
  const [users, setUsers] = useState(initialUsers)
  const [organizations, setOrganizations] = useState(initialOrgs)
  const [totalResult, setTotalResult] = useState(initialTotalResult)
  const [page, setPage] = useState(0)
  const [pageSize, setPageSize] = useState(10)
  const [isRefetching, startTransition] = useTransition()

  // Re-fetch users when the selected organization or pagination changes
  useEffect(() => {
    startTransition(async () => {
      try {
        const result = await fetchUsers(currentOrganization?.id ?? null, pageSize, page * pageSize)
        setUsers(result.users)
        setOrganizations(result.organizations)
        setTotalResult(result.totalResult)
      } catch (e) {
        console.error("Failed to refresh users:", e)
      }
    })
  }, [currentOrganization?.id, page, pageSize])

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

  const filteredUsers = useMemo(() => {
    let result = users
    // Free-text search
    if (searchQuery) {
      const q = searchQuery.toLowerCase()
      result = result.filter((user: any) => {
        const info = getUserDisplayInfo(user)
        const orgName = orgNameMap[user.details?.resourceOwner] ?? ""
        return (
          info.displayName.toLowerCase().includes(q) ||
          info.email.toLowerCase().includes(q) ||
          (user.username ?? "").toLowerCase().includes(q) ||
          orgName.toLowerCase().includes(q)
        )
      })
    }
    // Filter tokens
    if (activeFilters.state) {
      result = result.filter((user: any) => user.state === activeFilters.state)
    }
    if (activeFilters.type) {
      result = result.filter((user: any) => {
        if (activeFilters.type === "human") return !!user.human
        if (activeFilters.type === "machine") return !!user.machine
        return true
      })
    }
    if (activeFilters.username) {
      const v = activeFilters.username.toLowerCase()
      result = result.filter((u: any) => (u.username ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.firstname) {
      const v = activeFilters.firstname.toLowerCase()
      result = result.filter((u: any) => (u.human?.profile?.givenName ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.lastname) {
      const v = activeFilters.lastname.toLowerCase()
      result = result.filter((u: any) => (u.human?.profile?.familyName ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.displayname) {
      const v = activeFilters.displayname.toLowerCase()
      result = result.filter((u: any) => {
        const info = getUserDisplayInfo(u)
        return info.displayName.toLowerCase().includes(v)
      })
    }
    if (activeFilters.email) {
      const v = activeFilters.email.toLowerCase()
      result = result.filter((u: any) => (u.human?.email?.email ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.phone) {
      const v = activeFilters.phone.toLowerCase()
      result = result.filter((u: any) => (u.human?.phone?.phone ?? "").toLowerCase().includes(v))
    }
    if (activeFilters.loginname) {
      const v = activeFilters.loginname.toLowerCase()
      result = result.filter((u: any) =>
        (u.loginNames ?? []).some((ln: string) => ln.toLowerCase().includes(v))
      )
    }
    if (activeFilters.org) {
      const v = activeFilters.org.toLowerCase()
      result = result.filter((u: any) => {
        const orgName = orgNameMap[u.details?.resourceOwner] ?? ""
        return orgName.toLowerCase().includes(v) ||
          (u.details?.resourceOwner ?? "").toLowerCase().includes(v)
      })
    }
    return result
  }, [users, searchQuery, orgNameMap, activeFilters])

  const userFilters: FilterDef[] = [
    { key: "username", label: "username" },
    { key: "firstname", label: "firstname" },
    { key: "lastname", label: "lastname" },
    { key: "displayname", label: "displayname" },
    { key: "email", label: "email" },
    { key: "phone", label: "phone" },
    { key: "loginname", label: "loginname" },
    { key: "org", label: "org" },
    {
      key: "state",
      label: "state",
      options: [
        { value: "USER_STATE_ACTIVE", label: "active" },
        { value: "USER_STATE_INACTIVE", label: "inactive" },
        { value: "USER_STATE_LOCKED", label: "locked" },
        { value: "USER_STATE_INITIAL", label: "initial" },
      ],
    },
    {
      key: "type",
      label: "type",
      options: [
        { value: "human", label: "human" },
        { value: "machine", label: "machine" },
      ],
    },
  ]

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Users</h1>
          <p className="text-sm text-muted-foreground">
            Manage users in your ZITADEL instance
          </p>
        </div>
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-6 text-center">
          <p className="text-sm font-medium text-destructive">
            Failed to load users
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
            Users{" "}
            {isRefetching && (
              <Loader2 className="inline h-5 w-5 animate-spin ml-2" />
            )}
          </h1>
          <p className="text-sm text-muted-foreground">
            {currentOrganization
              ? `Users in ${currentOrganization.name}`
              : `Manage users across all organizations (${users.length} total)`}
          </p>
        </div>
        <RequirePermission permission="user.write">
          <Button onClick={() => setAddUserOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add User
          </Button>
        </RequirePermission>
      </div>

      {/* Filters */}
      <FilterBar
        searchPlaceholder="Search users by name or email..."
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
        filters={userFilters}
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
            columns={["User", "Organization", "Status", "Created", "Updated"]}
            rows={pageSize}
          />
        </div>
      ) : filteredUsers.length === 0 ? (
        <div className="rounded-lg border">
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <Users className="h-12 w-12 text-muted-foreground/40 mb-4" />
            <p className="text-sm font-medium">
              {searchQuery
                ? "No users match your search"
                : "No users found"}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              {searchQuery
                ? "Try adjusting your search query"
                : "Add your first user to get started"}
            </p>
          </div>
        </div>
      ) : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead>User</TableHead>
                <TableHead className="w-[140px]">Organization</TableHead>
                <TableHead className="w-[80px]">Status</TableHead>
                <TableHead className="w-[120px]">Created</TableHead>
                <TableHead className="w-[120px]">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredUsers.map((user: any) => {
                const info = getUserDisplayInfo(user)
                const stateInfo = getUserState(user.state)

                return (
                  <TableRow
                    key={user.userId ?? user.username}
                    className="cursor-pointer"
                    onClick={() => {
                      if (user.userId) {
                        router.push(`/users/${user.userId}`)
                      }
                    }}
                  >
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-xs font-medium flex-shrink-0">
                          {info.kind === "machine" ? (
                            <Bot className="h-4 w-4 text-muted-foreground" />
                          ) : (
                            getInitials(info.displayName)
                          )}
                        </div>
                        <div className="min-w-0">
                          <p className="font-medium truncate">
                            {info.displayName}
                          </p>
                          <p className="text-xs text-muted-foreground truncate">
                            {info.email || user.username || ""}
                          </p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm">
                      {orgNameMap[user.details?.resourceOwner] ? (
                        <span className="text-foreground">
                          {orgNameMap[user.details.resourceOwner]}
                        </span>
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
                      {formatDate(user.details?.creationDate)}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(user.details?.changeDate)}
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

      {/* Add User Sheet */}
      <AddUserSheet
        open={addUserOpen}
        onOpenChange={setAddUserOpen}
        organizations={organizations}
      />
    </div>
  )
}
