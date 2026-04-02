"use client"

import { useState, useMemo } from "react"
import { ColumnDef } from "@tanstack/react-table"
import { DataTable } from "../../../components/data-table/data-table"
import { getUsersByOrganization } from "../../../mock-data"
import type { User } from "../../../types"
import { Badge } from "../../../components/ui/badge"
import { Button } from "../../../components/ui/button"
import { Avatar, AvatarFallback } from "../../../components/ui/avatar"
import { Checkbox } from "../../../components/ui/checkbox"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../../../components/ui/dropdown-menu"
import { 
  MoreHorizontal, 
  Plus, 
  ArrowUpDown, 
  Eye, 
  Edit, 
  Lock, 
  Trash2,
  Download,
  Shield,
  Clock,
  UserPlus,
  AlertTriangle,
  Crown
} from "lucide-react"
import { ConsoleLink as Link } from "../../../context/link-context"
import { useAppContext } from "../../../context/app-context"
import { OrganizationSelectorPrompt } from "../../../components/organization-selector-prompt"
import { Tabs, TabsList, TabsTrigger } from "../../../components/ui/tabs"
import { toast } from "sonner"

const statusColors: Record<User["status"], string> = {
  active: "bg-foreground/10 text-foreground border-foreground/20",
  inactive: "bg-muted text-muted-foreground border-border",
  locked: "bg-muted text-muted-foreground border-border",
  pending: "bg-amber-500/10 text-amber-600 border-amber-500/20",
}

export default function OrgUsersPage() {
  const { currentOrganization } = useAppContext()
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set())
  const [activeTab, setActiveTab] = useState<"all" | "pending">("all")

  if (!currentOrganization) {
    return (
      <OrganizationSelectorPrompt 
        title="Select an Organization"
        description="Choose an organization to view its users"
        targetPath="/org/users"
      />
    )
  }
  
  const orgUsers = getUsersByOrganization(currentOrganization.id)

  // Stats for summary bar
  const stats = useMemo(() => {
    const pending = orgUsers.filter(u => u.status === "pending")
    const pendingOld = pending.filter(u => {
      const created = new Date(u.createdAt)
      const daysSince = (Date.now() - created.getTime()) / (1000 * 60 * 60 * 24)
      return daysSince > 7
    })
    const locked = orgUsers.filter(u => u.status === "locked")
    const lockedAdmins = locked.filter(u => u.role === "admin" || u.role === "owner")
    const neverLoggedIn = orgUsers.filter(u => u.lastLogin === null)
    const owners = orgUsers.filter(u => u.role === "owner")
    const admins = orgUsers.filter(u => u.role === "admin")
    
    return {
      total: orgUsers.length,
      pending: pending.length,
      pendingOld: pendingOld.length,
      neverLoggedIn: neverLoggedIn.length,
      locked: locked.length,
      lockedAdmins: lockedAdmins,
      owners: owners.length,
      admins: admins.length,
    }
  }, [orgUsers])

  // Filter users based on active tab
  const filteredByTab = activeTab === "pending" 
    ? orgUsers.filter(u => u.status === "pending")
    : orgUsers

  const statusOptions = [
    { label: "Active", value: "active" },
    { label: "Inactive", value: "inactive" },
    { label: "Locked", value: "locked" },
    { label: "Pending", value: "pending" },
  ]

  const roleOptions = [
    { label: "User", value: "user" },
    { label: "Admin", value: "admin" },
    { label: "Owner", value: "owner" },
  ]

  const loginOptions = [
    { label: "Has logged in", value: "logged_in" },
    { label: "Never logged in", value: "never" },
  ]

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedRows(new Set(filteredByTab.map(u => u.id)))
    } else {
      setSelectedRows(new Set())
    }
  }

  const handleSelectRow = (userId: string, checked: boolean) => {
    const newSelected = new Set(selectedRows)
    if (checked) {
      newSelected.add(userId)
    } else {
      newSelected.delete(userId)
    }
    setSelectedRows(newSelected)
  }

  const handleBulkAction = (action: string) => {
    toast.success(`${action} ${selectedRows.size} user(s)`)
    setSelectedRows(new Set())
  }

  const handleExport = () => {
    const csvContent = [
      ["Name", "Email", "Status", "Role", "Last Login", "Created At"].join(","),
      ...filteredByTab.map(u => [
        u.displayName,
        u.email,
        u.status,
        u.role,
        u.lastLogin ? new Date(u.lastLogin).toISOString() : "Never",
        new Date(u.createdAt).toISOString()
      ].join(","))
    ].join("\n")
    
    const blob = new Blob([csvContent], { type: "text/csv" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = `users-${currentOrganization.name}-${new Date().toISOString().split("T")[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
    toast.success("Users exported successfully")
  }

  const columns: ColumnDef<User>[] = [
    {
      id: "select",
      header: () => (
        <Checkbox
          checked={selectedRows.size === filteredByTab.length && filteredByTab.length > 0}
          onCheckedChange={(checked) => handleSelectAll(!!checked)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={selectedRows.has(row.original.id)}
          onCheckedChange={(checked) => handleSelectRow(row.original.id, !!checked)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "displayName",
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
        const user = row.original
        const isOwner = user.role === "owner"
        return (
          <div className="flex items-center gap-3">
            {isOwner && (
              <div className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
            )}
            <Avatar className={`h-8 w-8 ${isOwner ? "ring-2 ring-primary ring-offset-2 ring-offset-background" : ""}`}>
              <AvatarFallback className={`text-xs ${isOwner ? "bg-primary/10 text-primary" : ""}`}>
                {user.firstName.charAt(0)}{user.lastName.charAt(0)}
              </AvatarFallback>
            </Avatar>
            <div>
              <Link 
                href={`/users/${user.id}`} 
                className={`font-medium hover:underline ${isOwner ? "text-primary" : ""}`}
              >
                {user.displayName}
              </Link>
              <p className="text-sm text-muted-foreground">{user.email}</p>
            </div>
          </div>
        )
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.original.status
        return (
          <Badge variant="outline" className={statusColors[status]}>
            {status}
          </Badge>
        )
      },
      filterFn: (row, id, value) => {
        return value === row.getValue(id)
      },
    },
    {
      accessorKey: "role",
      header: "Role",
      cell: ({ row }) => {
        const role = row.original.role
        const isOwner = role === "owner"
        const isAdmin = role === "admin"
        return (
          <div className={`flex items-center gap-1.5 ${isOwner ? "text-primary font-medium" : ""}`}>
            {isOwner && <Crown className="h-3.5 w-3.5 text-primary" />}
            {isAdmin && <Shield className="h-3.5 w-3.5 text-muted-foreground" />}
            <span className="capitalize">{role}</span>
          </div>
        )
      },
      filterFn: (row, id, value) => {
        return value === row.getValue(id)
      },
    },
    {
      accessorKey: "lastLogin",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="-ml-4"
        >
          Last Login
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const lastLogin = row.original.lastLogin
        return lastLogin 
          ? new Date(lastLogin).toLocaleDateString() 
          : <span className="text-muted-foreground flex items-center gap-1">
              <Clock className="h-3.5 w-3.5" />
              Never
            </span>
      },
      filterFn: (row, id, value) => {
        const lastLogin = row.original.lastLogin
        if (value === "never") return lastLogin === null
        if (value === "logged_in") return lastLogin !== null
        return true
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const user = row.original
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem asChild>
                <Link href={`/users/${user.id}`}>
                  <Eye className="mr-2 h-4 w-4" />
                  View Details
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Edit className="mr-2 h-4 w-4" />
                Edit User
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Lock className="mr-2 h-4 w-4" />
                {user.status === "locked" ? "Unlock" : "Lock"} User
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-destructive">
                <Trash2 className="mr-2 h-4 w-4" />
                Delete User
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Users</h1>
          <p className="text-muted-foreground">
            Users in {currentOrganization?.name} ({orgUsers.length} total)
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={handleExport}>
            <Download className="mr-2 h-4 w-4" />
            Export
          </Button>
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Add User
          </Button>
        </div>
      </div>

      {/* Summary Bar */}
      <div className="rounded-lg border bg-card p-4 space-y-3">
        <div className="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm">
          <button 
            onClick={() => setActiveTab("pending")}
            className={`flex items-center gap-2 hover:text-foreground transition-colors ${activeTab === "pending" ? "text-foreground font-medium" : "text-muted-foreground"}`}
          >
            <UserPlus className="h-4 w-4" />
            <span>
              {stats.pending} pending
              {stats.pendingOld > 0 && (
                <span className="text-amber-600 ml-1">({stats.pendingOld} older than 7 days)</span>
              )}
            </span>
          </button>
          
          <span className="text-border hidden sm:inline">|</span>
          
          <button 
            onClick={() => setActiveTab("all")}
            className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
          >
            <Clock className="h-4 w-4" />
            <span>{stats.neverLoggedIn} never logged in</span>
          </button>
          
          <span className="text-border hidden sm:inline">|</span>
          
          <button 
            onClick={() => setActiveTab("all")}
            className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
          >
            <Crown className="h-4 w-4 text-primary" />
            <span>{stats.owners} Owners</span>
          </button>
          
          <button 
            onClick={() => setActiveTab("all")}
            className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
          >
            <Shield className="h-4 w-4" />
            <span>{stats.admins} Admins</span>
          </button>
          
          {stats.locked > 0 && (
            <>
              <span className="text-border hidden sm:inline">|</span>
              <span className="flex items-center gap-2 text-muted-foreground">
                <Lock className="h-4 w-4" />
                <span>{stats.locked} locked</span>
              </span>
            </>
          )}
        </div>
        
        {stats.lockedAdmins.length > 0 && (
          <div className="flex items-start gap-2 p-2 rounded-md bg-destructive/10 text-destructive text-sm">
            <AlertTriangle className="h-4 w-4 mt-0.5 flex-shrink-0" />
            <div>
              <span className="font-medium">Locked privileged accounts: </span>
              {stats.lockedAdmins.map((u, i) => (
                <span key={u.id}>
                  <Link href={`/users/${u.id}`} className="underline hover:no-underline">
                    {u.displayName}
                  </Link>
                  <span className="text-destructive/70"> ({u.role})</span>
                  {i < stats.lockedAdmins.length - 1 && ", "}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as "all" | "pending")}>
        <TabsList>
          <TabsTrigger value="all">All Users ({orgUsers.length})</TabsTrigger>
          <TabsTrigger value="pending">Pending Invites ({stats.pending})</TabsTrigger>
        </TabsList>
      </Tabs>

      {/* Bulk Actions */}
      {selectedRows.size > 0 && (
        <div className="flex items-center gap-3 p-3 bg-muted rounded-lg">
          <span className="text-sm font-medium">{selectedRows.size} selected</span>
          <div className="flex items-center gap-2 ml-auto">
            <Button variant="outline" size="sm" onClick={() => handleBulkAction("Activated")}>
              Activate
            </Button>
            <Button variant="outline" size="sm" onClick={() => handleBulkAction("Deactivated")}>
              Deactivate
            </Button>
            <Button variant="outline" size="sm" onClick={() => handleBulkAction("Locked")}>
              <Lock className="mr-1.5 h-3.5 w-3.5" />
              Lock
            </Button>
            <Button variant="outline" size="sm" className="text-destructive" onClick={() => handleBulkAction("Deleted")}>
              <Trash2 className="mr-1.5 h-3.5 w-3.5" />
              Delete
            </Button>
            <Button variant="ghost" size="sm" onClick={() => setSelectedRows(new Set())}>
              Clear
            </Button>
          </div>
        </div>
      )}

      <DataTable
        columns={columns}
        data={filteredByTab}
        searchKey="displayName"
        searchPlaceholder="Search users by name or email..."
        filterOptions={[
          { key: "status", label: "Status", options: statusOptions },
          { key: "role", label: "Role", options: roleOptions },
          { key: "lastLogin", label: "Login Status", options: loginOptions },
        ]}
      />
    </div>
  )
}
