"use client"

import { useState, useMemo } from "react"
import { ColumnDef } from "@tanstack/react-table"
import { DataTable } from "../../components/data-table/data-table"
import { roleAssignments, organizations, projects } from "../../mock-data"
import type { RoleAssignment } from "../../types"
import { Badge } from "../../components/ui/badge"
import { Button } from "../../components/ui/button"
import { Input } from "../../components/ui/input"
import { Checkbox } from "../../components/ui/checkbox"
import { Avatar, AvatarFallback } from "../../components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem,
} from "../../components/ui/dropdown-menu"
import { 
  MoreHorizontal, 
  Plus, 
  ArrowUpDown, 
  Edit, 
  Trash2, 
  Search, 
  Filter, 
  X, 
  Building2,
  FolderKanban,
  Shield,
  Crown,
  UserCog
} from "lucide-react"
import { ConsoleLink as Link } from "../../context/link-context"
import { useAppContext } from "../../context/app-context"
import { InstanceSelectorPrompt } from "../../components/instance-selector-prompt"

// Role badge colors by privilege level
const roleBadgeStyles: Record<string, string> = {
  owner: "bg-primary text-primary-foreground",
  admin: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
  editor: "bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-200",
  developer: "bg-violet-100 text-violet-800 dark:bg-violet-900 dark:text-violet-200",
  viewer: "bg-muted text-muted-foreground",
  tester: "bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200",
  analyst: "bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-200",
}

const privilegedRoles = ["owner", "admin"]

export default function InstanceRolesPage() {
  const { currentInstance } = useAppContext()
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set())
  const [searchQuery, setSearchQuery] = useState("")
  const [roleFilter, setRoleFilter] = useState<string[]>([])
  const [orgFilter, setOrgFilter] = useState<string[]>([])
  const [projectFilter, setProjectFilter] = useState<string[]>([])

  if (!currentInstance) {
    return (
      <InstanceSelectorPrompt 
        title="Continue to Role Assignments"
        description="Choose an instance to view its role assignments"
        icon={<UserCog className="h-6 w-6 text-muted-foreground" />}
        targetPath="/roles"
      />
    )
  }

  // Get unique values for filters
  const allRoles = useMemo(() => {
    const roles = new Set<string>()
    roleAssignments.forEach(ra => ra.roles.forEach(r => roles.add(r)))
    return Array.from(roles).sort()
  }, [])

  const instanceOrgs = useMemo(() => {
    return organizations.filter(o => o.instanceId === currentInstance.id)
  }, [currentInstance.id])

  const instanceProjects = useMemo(() => {
    const orgIds = new Set(instanceOrgs.map(o => o.id))
    return projects.filter(p => orgIds.has(p.orgId))
  }, [instanceOrgs])

  // Filter role assignments
  const filteredAssignments = useMemo(() => {
    return roleAssignments.filter(ra => {
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase()
        const matchesSearch = 
          ra.userName.toLowerCase().includes(query) ||
          ra.userEmail.toLowerCase().includes(query) ||
          ra.projectName.toLowerCase().includes(query) ||
          ra.orgName.toLowerCase().includes(query) ||
          ra.roles.some(r => r.toLowerCase().includes(query)) ||
          ra.grantedBy.toLowerCase().includes(query)
        if (!matchesSearch) return false
      }

      // Role filter
      if (roleFilter.length > 0 && !ra.roles.some(r => roleFilter.includes(r))) {
        return false
      }

      // Org filter
      if (orgFilter.length > 0 && !orgFilter.includes(ra.orgId)) {
        return false
      }

      // Project filter
      if (projectFilter.length > 0 && !projectFilter.includes(ra.projectId)) {
        return false
      }

      return true
    })
  }, [searchQuery, roleFilter, orgFilter, projectFilter])

  // Stats
  const stats = useMemo(() => ({
    total: roleAssignments.length,
    privileged: roleAssignments.filter(ra => ra.roles.some(r => privilegedRoles.includes(r))).length,
    uniqueUsers: new Set(roleAssignments.map(ra => ra.userId)).size,
  }), [])

  const hasFilters = searchQuery || roleFilter.length > 0 || orgFilter.length > 0 || projectFilter.length > 0

  const clearFilters = () => {
    setSearchQuery("")
    setRoleFilter([])
    setOrgFilter([])
    setProjectFilter([])
  }

  const toggleRowSelection = (id: string) => {
    setSelectedRows(prev => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  const toggleAllRows = () => {
    if (selectedRows.size === filteredAssignments.length) {
      setSelectedRows(new Set())
    } else {
      setSelectedRows(new Set(filteredAssignments.map(ra => ra.id)))
    }
  }

  const columns: ColumnDef<RoleAssignment>[] = [
    {
      id: "select",
      header: () => (
        <Checkbox
          checked={selectedRows.size === filteredAssignments.length && filteredAssignments.length > 0}
          onCheckedChange={toggleAllRows}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={selectedRows.has(row.original.id)}
          onCheckedChange={() => toggleRowSelection(row.original.id)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "userName",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="-ml-4"
        >
          User
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const assignment = row.original
        const hasPrivilegedRole = assignment.roles.some(r => privilegedRoles.includes(r))
        return (
          <div className="flex items-center gap-3">
            <Avatar className={`h-8 w-8 ${hasPrivilegedRole ? "ring-2 ring-primary ring-offset-2" : ""}`}>
              <AvatarFallback className={`text-xs ${hasPrivilegedRole ? "bg-primary/10 text-primary" : ""}`}>
                {assignment.userName.split(" ").map(n => n[0]).join("")}
              </AvatarFallback>
            </Avatar>
            <div>
              <Link 
                href={`/users/${assignment.userId}`} 
                className={`font-medium hover:underline ${hasPrivilegedRole ? "text-primary" : ""}`}
              >
                {assignment.userName}
              </Link>
              <p className="text-sm text-muted-foreground">{assignment.userEmail}</p>
            </div>
          </div>
        )
      },
    },
    {
      accessorKey: "orgName",
      header: "Organization",
      cell: ({ row }) => (
        <div className="flex items-center gap-2 text-muted-foreground">
          <Building2 className="h-4 w-4" />
          <span>{row.original.orgName}</span>
        </div>
      ),
    },
    {
      accessorKey: "projectName",
      header: "Project",
      cell: ({ row }) => (
        <Link 
          href={`/projects/${row.original.projectId}`} 
          className="flex items-center gap-2 hover:underline"
        >
          <FolderKanban className="h-4 w-4 text-muted-foreground" />
          {row.original.projectName}
        </Link>
      ),
    },
    {
      accessorKey: "roles",
      header: "Roles",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1">
          {row.original.roles.map((role) => (
            <Badge 
              key={role} 
              variant="secondary"
              className={roleBadgeStyles[role] || ""}
            >
              {role === "owner" && <Crown className="h-3 w-3 mr-1" />}
              {role === "admin" && <Shield className="h-3 w-3 mr-1" />}
              {role}
            </Badge>
          ))}
        </div>
      ),
    },
    {
      accessorKey: "grantedAt",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="-ml-4"
        >
          Granted
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => (
        <div className="space-y-1">
          <p className="text-sm">{new Date(row.original.grantedAt).toLocaleDateString()}</p>
          <div className="flex items-center gap-1.5">
            <span className="text-xs text-muted-foreground">by</span>
            <Link 
              href={`/users/${row.original.grantedById}`}
              className="text-sm font-medium hover:underline"
            >
              {row.original.grantedBy}
            </Link>
          </div>
        </div>
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <Edit className="mr-2 h-4 w-4" />
              Edit Roles
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-destructive">
              <Trash2 className="mr-2 h-4 w-4" />
              Remove Assignment
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Role Assignments</h1>
          <p className="text-muted-foreground">
            {stats.total} assignments across {stats.uniqueUsers} users ({stats.privileged} with privileged access)
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Assign Roles
        </Button>
      </div>

      {/* Bulk Actions Bar */}
      {selectedRows.size > 0 && (
        <div className="flex items-center gap-4 p-3 bg-muted rounded-lg">
          <span className="text-sm font-medium">{selectedRows.size} selected</span>
          <div className="flex items-center gap-2">
            <Button size="sm" variant="outline">
              <Edit className="mr-2 h-4 w-4" />
              Bulk Edit Roles
            </Button>
            <Button size="sm" variant="destructive">
              <Trash2 className="mr-2 h-4 w-4" />
              Revoke Selected
            </Button>
          </div>
          <Button 
            size="sm" 
            variant="ghost" 
            onClick={() => setSelectedRows(new Set())}
            className="ml-auto"
          >
            Clear selection
          </Button>
        </div>
      )}

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3">
        {/* Search */}
        <div className="relative flex-1 min-w-[250px] max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search user, project, role, or grantor..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9"
          />
        </div>

        {/* Role Filter */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm">
              <Shield className="mr-2 h-4 w-4" />
              Role
              {roleFilter.length > 0 && (
                <Badge variant="secondary" className="ml-2">{roleFilter.length}</Badge>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-48">
            <DropdownMenuLabel>Filter by Role</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {allRoles.map(role => (
              <DropdownMenuCheckboxItem
                key={role}
                checked={roleFilter.includes(role)}
                onCheckedChange={(checked) => 
                  setRoleFilter(prev => checked ? [...prev, role] : prev.filter(r => r !== role))
                }
              >
                <Badge variant="secondary" className={`mr-2 ${roleBadgeStyles[role] || ""}`}>
                  {role}
                </Badge>
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>

        {/* Organization Filter */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm">
              <Building2 className="mr-2 h-4 w-4" />
              Organization
              {orgFilter.length > 0 && (
                <Badge variant="secondary" className="ml-2">{orgFilter.length}</Badge>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-56 max-h-64 overflow-y-auto">
            <DropdownMenuLabel>Filter by Organization</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {instanceOrgs.map(org => (
              <DropdownMenuCheckboxItem
                key={org.id}
                checked={orgFilter.includes(org.id)}
                onCheckedChange={(checked) => 
                  setOrgFilter(prev => checked ? [...prev, org.id] : prev.filter(o => o !== org.id))
                }
              >
                {org.name}
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>

        {/* Project Filter */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm">
              <FolderKanban className="mr-2 h-4 w-4" />
              Project
              {projectFilter.length > 0 && (
                <Badge variant="secondary" className="ml-2">{projectFilter.length}</Badge>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-56 max-h-64 overflow-y-auto">
            <DropdownMenuLabel>Filter by Project</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {instanceProjects.map(project => (
              <DropdownMenuCheckboxItem
                key={project.id}
                checked={projectFilter.includes(project.id)}
                onCheckedChange={(checked) => 
                  setProjectFilter(prev => checked ? [...prev, project.id] : prev.filter(p => p !== project.id))
                }
              >
                {project.name}
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>

        {/* Clear filters */}
        {hasFilters && (
          <Button variant="ghost" size="sm" onClick={clearFilters}>
            <X className="mr-2 h-4 w-4" />
            Clear filters
          </Button>
        )}
      </div>

      {/* Results count when filtered */}
      {hasFilters && (
        <p className="text-sm text-muted-foreground">
          Showing {filteredAssignments.length} of {roleAssignments.length} assignments
        </p>
      )}

      <DataTable
        columns={columns}
        data={filteredAssignments}
      />
    </div>
  )
}
