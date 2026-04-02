"use client"

import { useState } from "react"
import { ColumnDef } from "@tanstack/react-table"
import { DataTable } from "../../../components/data-table/data-table"
import { getApplicationsByOrganization, roleAssignments, applications } from "../../../mock-data"
import type { Application } from "../../../types"
import { Badge } from "../../../components/ui/badge"
import { Button } from "../../../components/ui/button"
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
  Trash2, 
  Globe, 
  Smartphone, 
  Server, 
  Monitor,
  Copy,
  Check,
  Users,
  Settings,
  Key
} from "lucide-react"
import { ConsoleLink as Link } from "../../../context/link-context"
import { useAppContext } from "../../../context/app-context"
import { OrganizationSelectorPrompt } from "../../../components/organization-selector-prompt"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../../../components/ui/tooltip"

// Type configuration with icons, colors, and security info
const typeConfig: Record<Application["type"], {
  icon: React.ComponentType<{ className?: string }>
  label: string
  color: string
  bgColor: string
  description: string
}> = {
  web: {
    icon: Globe,
    label: "Web App",
    color: "text-blue-600",
    bgColor: "bg-blue-100 dark:bg-blue-950",
    description: "Browser-based application with redirect flows",
  },
  native: {
    icon: Smartphone,
    label: "Native App",
    color: "text-green-600",
    bgColor: "bg-green-100 dark:bg-green-950",
    description: "Mobile or desktop app with PKCE flow",
  },
  api: {
    icon: Server,
    label: "API",
    color: "text-purple-600",
    bgColor: "bg-purple-100 dark:bg-purple-950",
    description: "Machine-to-machine with client credentials",
  },
  "user-agent": {
    icon: Monitor,
    label: "User Agent",
    color: "text-amber-600",
    bgColor: "bg-amber-100 dark:bg-amber-950",
    description: "SPA without backend - implicit flow",
  },
}

// Helper to get user grant count per application
function getAppUserCount(appId: string): number {
  const app = applications.find(a => a.id === appId)
  if (!app) return 0
  const projectGrants = roleAssignments.filter(ra => ra.projectId === app.projectId)
  const uniqueUsers = new Set(projectGrants.map(ra => ra.userId))
  return uniqueUsers.size
}

// Copy button component with feedback
function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async (e: React.MouseEvent) => {
    e.stopPropagation()
    await navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={handleCopy}
          >
            {copied ? (
              <Check className="h-3 w-3 text-green-600" />
            ) : (
              <Copy className="h-3 w-3" />
            )}
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          <p>{copied ? "Copied!" : "Copy Client ID"}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}

const columns: ColumnDef<Application>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && "indeterminate")
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "name",
    header: ({ column }) => (
      <Button
        variant="ghost"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        className="-ml-4"
      >
        Application
        <ArrowUpDown className="ml-2 h-4 w-4" />
      </Button>
    ),
    cell: ({ row }) => {
      const app = row.original
      const config = typeConfig[app.type]
      const TypeIcon = config.icon
      return (
        <div className="flex items-center gap-3">
          <div className={`flex h-8 w-8 items-center justify-center rounded-md ${config.bgColor}`}>
            <TypeIcon className={`h-4 w-4 ${config.color}`} />
          </div>
          <div>
            <Link 
              href={`/applications/${app.id}`} 
              className="font-medium hover:underline"
            >
              {app.name}
            </Link>
            <p className="text-sm text-muted-foreground">{app.projectName}</p>
          </div>
        </div>
      )
    },
  },
  {
    accessorKey: "type",
    header: "Type",
    cell: ({ row }) => {
      const type = row.original.type
      const config = typeConfig[type]
      const TypeIcon = config.icon
      return (
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Badge 
                variant="secondary" 
                className={`${config.bgColor} ${config.color} border-0 gap-1.5`}
              >
                <TypeIcon className="h-3 w-3" />
                {config.label}
              </Badge>
            </TooltipTrigger>
            <TooltipContent>
              <p>{config.description}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )
    },
  },
  {
    accessorKey: "clientId",
    header: "Client ID",
    cell: ({ row }) => (
      <div className="flex items-center gap-1 group">
        <code className="text-xs bg-muted px-2 py-1 rounded font-mono">
          {row.original.clientId}
        </code>
        <CopyButton value={row.original.clientId} />
      </div>
    ),
  },
  {
    id: "userGrants",
    header: "Users",
    cell: ({ row }) => {
      const userCount = getAppUserCount(row.original.id)
      return (
        <div className="flex items-center gap-1.5 text-muted-foreground">
          <Users className="h-3.5 w-3.5" />
          <span>{userCount}</span>
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
    id: "actions",
    cell: ({ row }) => {
      const app = row.original
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
              <Link href={`/applications/${app.id}`}>
                <Eye className="mr-2 h-4 w-4" />
                View Details
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link href={`/applications/${app.id}?tab=token`}>
                <Key className="mr-2 h-4 w-4" />
                Token Configuration
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Edit className="mr-2 h-4 w-4" />
              Edit Application
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Settings className="mr-2 h-4 w-4" />
              Security Settings
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-destructive">
              <Trash2 className="mr-2 h-4 w-4" />
              Delete Application
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]

export default function OrgApplicationsPage() {
  const { currentOrganization } = useAppContext()
  const [rowSelection, setRowSelection] = useState({})

  if (!currentOrganization) {
    return (
      <OrganizationSelectorPrompt 
        title="Select an Organization"
        description="Choose an organization to view its applications"
        targetPath="/org/applications"
      />
    )
  }
  
  const orgApps = getApplicationsByOrganization(currentOrganization.id)
  const selectedCount = Object.keys(rowSelection).length

  const typeOptions = [
    { label: "Web App", value: "web" },
    { label: "Native App", value: "native" },
    { label: "API", value: "api" },
    { label: "User Agent", value: "user-agent" },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Applications</h1>
          <p className="text-muted-foreground">
            Applications in {currentOrganization?.name} ({orgApps.length} total)
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Create Application
        </Button>
      </div>

      {/* Bulk Actions Bar */}
      {selectedCount > 0 && (
        <div className="flex items-center gap-4 p-3 rounded-lg border bg-muted/50">
          <span className="text-sm font-medium">{selectedCount} selected</span>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              <Settings className="mr-2 h-4 w-4" />
              Bulk Configure
            </Button>
            <Button variant="outline" size="sm" className="text-destructive hover:text-destructive">
              <Trash2 className="mr-2 h-4 w-4" />
              Delete
            </Button>
          </div>
          <Button 
            variant="ghost" 
            size="sm" 
            className="ml-auto"
            onClick={() => setRowSelection({})}
          >
            Clear selection
          </Button>
        </div>
      )}

      <DataTable
        columns={columns}
        data={orgApps}
        searchKey="name"
        searchPlaceholder="Search applications..."
        filterOptions={[
          { key: "type", label: "Type", options: typeOptions },
        ]}
        rowSelection={rowSelection}
        onRowSelectionChange={setRowSelection}
      />
    </div>
  )
}
