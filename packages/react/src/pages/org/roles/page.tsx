"use client"

import { ColumnDef } from "@tanstack/react-table"
import { DataTable } from "../../../components/data-table/data-table"
import { getRoleAssignmentsByOrganization } from "../../../mock-data"
import type { RoleAssignment } from "../../../types"
import { Badge } from "../../../components/ui/badge"
import { Button } from "../../../components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../../../components/ui/dropdown-menu"
import { MoreHorizontal, Plus, ArrowUpDown, Edit, Trash2, UserCog } from "lucide-react"
import { ConsoleLink as Link } from "../../../context/link-context"
import { useAppContext } from "../../../context/app-context"
import { OrganizationSelectorPrompt } from "../../../components/organization-selector-prompt"

const columns: ColumnDef<RoleAssignment>[] = [
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
      return (
        <div>
          <Link 
            href={`/users/${assignment.userId}`} 
            className="font-medium hover:underline"
          >
            {assignment.userName}
          </Link>
          <p className="text-sm text-muted-foreground">{assignment.userEmail}</p>
        </div>
      )
    },
  },
  {
    accessorKey: "projectName",
    header: "Project",
    cell: ({ row }) => (
      <Link 
        href={`/projects/${row.original.projectId}`} 
        className="hover:underline"
      >
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
          <Badge key={role} variant="secondary">{role}</Badge>
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
      <div>
        <p>{new Date(row.original.grantedAt).toLocaleDateString()}</p>
        <p className="text-xs text-muted-foreground">by {row.original.grantedBy}</p>
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

export default function OrgRolesPage() {
  const { currentOrganization } = useAppContext()

  if (!currentOrganization) {
    return (
      <OrganizationSelectorPrompt 
        title="Select an Organization"
        description="Choose an organization to view its role assignments"
        targetPath="/org/roles"
      />
    )
  }
  
  const orgRoles = getRoleAssignmentsByOrganization(currentOrganization.id)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Role Assignments</h1>
          <p className="text-muted-foreground">
            Role assignments in {currentOrganization?.name} ({orgRoles.length} total)
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Assign Roles
        </Button>
      </div>

      <DataTable
        columns={columns}
        data={orgRoles}
      />
    </div>
  )
}
