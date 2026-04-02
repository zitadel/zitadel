"use client"

import { ColumnDef } from "@tanstack/react-table"
import { DataTable } from "../../../components/data-table/data-table"
import { getProjectsByOrganization, roleAssignments } from "../../../mock-data"
import type { Project } from "../../../types"
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
import { MoreHorizontal, Plus, ArrowUpDown, Eye, Edit, Trash2, FolderKanban, Users, AppWindow } from "lucide-react"
import { ConsoleLink as Link } from "../../../context/link-context"
import { useAppContext } from "../../../context/app-context"
import { OrganizationSelectorPrompt } from "../../../components/organization-selector-prompt"

// Helper to get unique user count per project
function getProjectUserCount(projectId: string): number {
  const userIds = new Set(roleAssignments.filter(ra => ra.projectId === projectId).map(ra => ra.userId))
  return userIds.size
}

const columns: ColumnDef<Project>[] = [
  {
    accessorKey: "name",
    header: ({ column }) => (
      <Button
        variant="ghost"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        className="-ml-4"
      >
        Project
        <ArrowUpDown className="ml-2 h-4 w-4" />
      </Button>
    ),
    cell: ({ row }) => {
      const project = row.original
      return (
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10">
            <FolderKanban className="h-4 w-4 text-primary" />
          </div>
          <div>
            <Link 
              href={`/projects/${project.id}`} 
              className="font-medium hover:underline"
            >
              {project.name}
            </Link>
            <p className="text-sm text-muted-foreground line-clamp-1">{project.description}</p>
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
        <Badge 
          variant="outline" 
          className={status === "active" 
            ? "bg-foreground/10 text-foreground border-foreground/20" 
            : "bg-muted text-muted-foreground border-border"
          }
        >
          {status}
        </Badge>
      )
    },
  },
  {
    accessorKey: "applicationCount",
    header: "Applications",
    cell: ({ row }) => {
      const project = row.original
      return (
        <Link 
          href={`/projects/${project.id}?tab=applications`}
          className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors group"
        >
          <AppWindow className="h-3.5 w-3.5" />
          <span className="group-hover:underline">{project.applicationCount}</span>
        </Link>
      )
    },
  },
  {
    id: "userGrants",
    header: "User Grants",
    cell: ({ row }) => {
      const project = row.original
      const userCount = getProjectUserCount(project.id)
      return (
        <Link 
          href={`/projects/${project.id}?tab=roles`}
          className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors group"
        >
          <Users className="h-3.5 w-3.5" />
          <span className="group-hover:underline">{userCount}</span>
        </Link>
      )
    },
  },
  {
    accessorKey: "updatedAt",
    header: ({ column }) => (
      <Button
        variant="ghost"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        className="-ml-4"
      >
        Last Updated
        <ArrowUpDown className="ml-2 h-4 w-4" />
      </Button>
    ),
    cell: ({ row }) => new Date(row.original.updatedAt).toLocaleDateString(),
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
          <DropdownMenuItem asChild>
            <Link href={`/projects/${row.original.id}`}>
              <Eye className="mr-2 h-4 w-4" />
              View Details
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Edit className="mr-2 h-4 w-4" />
            Edit Project
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem className="text-destructive">
            <Trash2 className="mr-2 h-4 w-4" />
            Delete Project
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    ),
  },
]

export default function OrgProjectsPage() {
  const { currentOrganization } = useAppContext()

  if (!currentOrganization) {
    return (
      <OrganizationSelectorPrompt 
        title="Select an Organization"
        description="Choose an organization to view its projects"
        targetPath="/org/projects"
      />
    )
  }
  
  const orgProjects = getProjectsByOrganization(currentOrganization.id)

  const statusOptions = [
    { label: "Active", value: "active" },
    { label: "Inactive", value: "inactive" },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Projects</h1>
          <p className="text-muted-foreground">
            Projects in {currentOrganization?.name} ({orgProjects.length} total)
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Create Project
        </Button>
      </div>

      <DataTable
        columns={columns}
        data={orgProjects}
        searchKey="name"
        searchPlaceholder="Search projects..."
        filterOptions={[
          { key: "status", label: "Status", options: statusOptions },
        ]}
      />
    </div>
  )
}
