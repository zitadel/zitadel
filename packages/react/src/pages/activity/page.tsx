"use client"

import { useState, useMemo } from "react"
import { ColumnDef } from "@tanstack/react-table"
import { DataTable } from "../../components/data-table/data-table"
import { getActivityLogByInstance, organizations } from "../../mock-data"
import type { ActivityLogEntry } from "../../types"
import { Badge } from "../../components/ui/badge"
import { Button } from "../../components/ui/button"
import { Input } from "../../components/ui/input"
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../components/ui/select"
import { Calendar } from "../../components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../../components/ui/popover"
import { 
  ArrowUpDown, 
  User, 
  FolderKanban, 
  AppWindow, 
  Building2, 
  Settings, 
  Activity,
  Download,
  CalendarIcon,
  Search,
  X,
  AlertTriangle,
  Info,
  UserCog,
  Key,
  ExternalLink,
  Server
} from "lucide-react"
import { ConsoleLink as Link } from "../../context/link-context"
import { useAppContext } from "../../context/app-context"
import { InstanceSelectorPrompt } from "../../components/instance-selector-prompt"
import { format, subHours, subDays, isWithinInterval } from "date-fns"

const resourceIcons: Record<ActivityLogEntry["resourceType"], React.ComponentType<{ className?: string }>> = {
  user: User,
  project: FolderKanban,
  application: AppWindow,
  organization: Building2,
  settings: Settings,
  role_assignment: UserCog,
  session: Key,
}

const severityConfig = {
  routine: { 
    label: "Routine", 
    color: "text-muted-foreground", 
    bg: "bg-muted",
    icon: null 
  },
  important: { 
    label: "Important", 
    color: "text-blue-600", 
    bg: "bg-blue-50 dark:bg-blue-950",
    icon: Info 
  },
  sensitive: { 
    label: "Sensitive", 
    color: "text-amber-600", 
    bg: "bg-amber-50 dark:bg-amber-950",
    icon: AlertTriangle 
  },
}

const actionTypeColors: Record<string, string> = {
  created: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400",
  updated: "bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-400",
  deleted: "bg-red-100 text-red-700 dark:bg-red-950 dark:text-red-400",
  activated: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400",
  deactivated: "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400",
  revoked: "bg-red-100 text-red-700 dark:bg-red-950 dark:text-red-400",
  assigned: "bg-purple-100 text-purple-700 dark:bg-purple-950 dark:text-purple-400",
  unassigned: "bg-orange-100 text-orange-700 dark:bg-orange-950 dark:text-orange-400",
}

function getResourceLink(entry: ActivityLogEntry): string | null {
  switch (entry.resourceType) {
    case "user": return `/users/${entry.resourceId}`
    case "project": return `/projects/${entry.resourceId}`
    case "application": return `/applications/${entry.resourceId}`
    case "session": return `/sessions`
    case "role_assignment": return `/roles`
    default: return null
  }
}

// Time range presets
const timeRangePresets = [
  { label: "Last hour", value: "1h", getRange: () => ({ start: subHours(new Date(), 1), end: new Date() }) },
  { label: "Last 24 hours", value: "24h", getRange: () => ({ start: subHours(new Date(), 24), end: new Date() }) },
  { label: "Last 7 days", value: "7d", getRange: () => ({ start: subDays(new Date(), 7), end: new Date() }) },
  { label: "Last 30 days", value: "30d", getRange: () => ({ start: subDays(new Date(), 30), end: new Date() }) },
  { label: "All time", value: "all", getRange: () => null },
  { label: "Custom", value: "custom", getRange: () => null },
]

export default function InstanceActivityPage() {
  const { currentInstance } = useAppContext()
  
  // Filter states
  const [timeRange, setTimeRange] = useState("30d")
  const [customDateRange, setCustomDateRange] = useState<{ start: Date | undefined; end: Date | undefined }>({ start: undefined, end: undefined })
  const [actorFilter, setActorFilter] = useState("")
  const [actionTypeFilter, setActionTypeFilter] = useState<string>("all")
  const [resourceTypeFilter, setResourceTypeFilter] = useState<string>("all")
  const [severityFilter, setSeverityFilter] = useState<string>("all")
  const [orgFilter, setOrgFilter] = useState<string>("all")

  if (!currentInstance) {
    return (
      <InstanceSelectorPrompt 
        title="Continue to Activity Log"
        description="Choose an instance to view its activity log"
        icon={<Activity className="h-6 w-6 text-muted-foreground" />}
        targetPath="/activity"
      />
    )
  }
  
  const instanceActivity = getActivityLogByInstance(currentInstance.id)
  const instanceOrgs = organizations.filter(o => o.instanceId === currentInstance.id)

  // Apply filters
  const filteredActivity = useMemo(() => {
    return instanceActivity.filter(entry => {
      // Time range filter
      if (timeRange !== "all" && timeRange !== "custom") {
        const preset = timeRangePresets.find(p => p.value === timeRange)
        const range = preset?.getRange()
        if (range && !isWithinInterval(entry.timestamp, { start: range.start, end: range.end })) {
          return false
        }
      }
      if (timeRange === "custom" && customDateRange.start && customDateRange.end) {
        if (!isWithinInterval(entry.timestamp, { start: customDateRange.start, end: customDateRange.end })) {
          return false
        }
      }
      
      // Org filter
      if (orgFilter !== "all" && entry.orgId !== orgFilter) {
        return false
      }
      
      // Actor filter
      if (actorFilter && !entry.actorName.toLowerCase().includes(actorFilter.toLowerCase())) {
        return false
      }
      
      // Action type filter
      if (actionTypeFilter !== "all" && entry.actionType !== actionTypeFilter) {
        return false
      }
      
      // Resource type filter
      if (resourceTypeFilter !== "all" && entry.resourceType !== resourceTypeFilter) {
        return false
      }
      
      // Severity filter
      if (severityFilter !== "all" && entry.severity !== severityFilter) {
        return false
      }
      
      return true
    })
  }, [instanceActivity, timeRange, customDateRange, actorFilter, actionTypeFilter, resourceTypeFilter, severityFilter, orgFilter])

  const hasActiveFilters = timeRange !== "30d" || actorFilter || actionTypeFilter !== "all" || resourceTypeFilter !== "all" || severityFilter !== "all" || orgFilter !== "all"

  const clearFilters = () => {
    setTimeRange("30d")
    setCustomDateRange({ start: undefined, end: undefined })
    setActorFilter("")
    setActionTypeFilter("all")
    setResourceTypeFilter("all")
    setSeverityFilter("all")
    setOrgFilter("all")
  }

  // Export function
  const handleExport = () => {
    const headers = ["Timestamp", "Action", "Action Type", "Resource Type", "Resource", "Resource ID", "Actor", "Severity", "Organization"]
    const rows = filteredActivity.map(entry => {
      const org = instanceOrgs.find(o => o.id === entry.orgId)
      return [
        entry.timestamp.toISOString(),
        entry.action,
        entry.actionType,
        entry.resourceType,
        entry.resourceName,
        entry.resourceId,
        entry.actorName,
        entry.severity,
        org?.name || "Unknown",
      ]
    })
    
    const csv = [headers, ...rows].map(row => row.map(cell => `"${cell}"`).join(",")).join("\n")
    const blob = new Blob([csv], { type: "text/csv" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = `activity-log-${currentInstance.name}-${new Date().toISOString().split("T")[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  const columns: ColumnDef<ActivityLogEntry>[] = [
    {
      accessorKey: "timestamp",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="-ml-4"
        >
          Time
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const date = new Date(row.original.timestamp)
        return (
          <div>
            <p className="font-medium">{format(date, "MMM d, yyyy")}</p>
            <p className="text-xs text-muted-foreground">{format(date, "HH:mm:ss")}</p>
          </div>
        )
      },
    },
    {
      accessorKey: "action",
      header: "Action",
      cell: ({ row }) => {
        const entry = row.original
        return (
          <Badge 
            variant="secondary" 
            className={actionTypeColors[entry.actionType] || ""}
          >
            {entry.action}
          </Badge>
        )
      },
    },
    {
      accessorKey: "resourceName",
      header: "Resource",
      cell: ({ row }) => {
        const entry = row.original
        const Icon = resourceIcons[entry.resourceType]
        const href = getResourceLink(entry)
        
        return (
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
              <Icon className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              {href ? (
                <Link href={href} className="font-medium hover:underline flex items-center gap-1">
                  {entry.resourceName}
                  <ExternalLink className="h-3 w-3 text-muted-foreground" />
                </Link>
              ) : (
                <span className="font-medium">{entry.resourceName}</span>
              )}
              <div className="flex items-center gap-2">
                <p className="text-xs text-muted-foreground capitalize">{entry.resourceType.replace("_", " ")}</p>
                <span className="text-xs text-muted-foreground/50">|</span>
                <p className="text-xs text-muted-foreground font-mono">{entry.resourceId}</p>
              </div>
            </div>
          </div>
        )
      },
    },
    {
      id: "organization",
      header: "Organization",
      cell: ({ row }) => {
        const org = instanceOrgs.find(o => o.id === row.original.orgId)
        return (
          <span className="text-sm text-muted-foreground">
            {org?.name || "Unknown"}
          </span>
        )
      },
    },
    {
      accessorKey: "actorName",
      header: "Actor",
      cell: ({ row }) => (
        <Link 
          href={`/users/${row.original.actorId}`}
          className="hover:underline"
        >
          {row.original.actorName}
        </Link>
      ),
    },
    {
      accessorKey: "severity",
      header: "Severity",
      cell: ({ row }) => {
        const severity = row.original.severity
        const config = severityConfig[severity]
        const SeverityIcon = config.icon
        
        return (
          <div className={`flex items-center gap-1.5 px-2 py-1 rounded-md text-xs font-medium ${config.bg} ${config.color}`}>
            {SeverityIcon && <SeverityIcon className="h-3 w-3" />}
            {config.label}
          </div>
        )
      },
    },
  ]

  // Unique actors for autocomplete hint
  const uniqueActors = [...new Set(instanceActivity.map(e => e.actorName))]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Instance Activity Log</h1>
          <p className="text-muted-foreground">
            All activity across {currentInstance?.name} ({filteredActivity.length} of {instanceActivity.length} events)
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={handleExport}>
          <Download className="h-4 w-4 mr-2" />
          Export
        </Button>
      </div>

      {/* Filter Strip */}
      <div className="flex flex-wrap items-center gap-3 p-4 rounded-lg border bg-card">
        {/* Time Range */}
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground">Time:</span>
          <Select value={timeRange} onValueChange={setTimeRange}>
            <SelectTrigger className="w-[140px] h-9">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {timeRangePresets.map(preset => (
                <SelectItem key={preset.value} value={preset.value}>
                  {preset.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {timeRange === "custom" && (
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline" size="sm" className="h-9">
                  <CalendarIcon className="h-4 w-4 mr-2" />
                  {customDateRange.start && customDateRange.end 
                    ? `${format(customDateRange.start, "MMM d")} - ${format(customDateRange.end, "MMM d")}`
                    : "Pick dates"
                  }
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                  mode="range"
                  selected={{ from: customDateRange.start, to: customDateRange.end }}
                  onSelect={(range) => setCustomDateRange({ start: range?.from, end: range?.to })}
                  numberOfMonths={2}
                />
              </PopoverContent>
            </Popover>
          )}
        </div>

        <span className="text-border">|</span>

        {/* Organization Filter */}
        <Select value={orgFilter} onValueChange={setOrgFilter}>
          <SelectTrigger className="w-[160px] h-9">
            <SelectValue placeholder="Organization" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All organizations</SelectItem>
            {instanceOrgs.map(org => (
              <SelectItem key={org.id} value={org.id}>
                {org.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <span className="text-border">|</span>

        {/* Actor Search */}
        <div className="relative">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input 
            placeholder="Filter by actor..."
            value={actorFilter}
            onChange={(e) => setActorFilter(e.target.value)}
            className="pl-8 h-9 w-[160px]"
            list="actors"
          />
          <datalist id="actors">
            {uniqueActors.slice(0, 10).map(actor => (
              <option key={actor} value={actor} />
            ))}
          </datalist>
        </div>

        <span className="text-border">|</span>

        {/* Action Type Filter */}
        <Select value={actionTypeFilter} onValueChange={setActionTypeFilter}>
          <SelectTrigger className="w-[130px] h-9">
            <SelectValue placeholder="Action type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All actions</SelectItem>
            <SelectItem value="created">Created</SelectItem>
            <SelectItem value="updated">Updated</SelectItem>
            <SelectItem value="deleted">Deleted</SelectItem>
            <SelectItem value="activated">Activated</SelectItem>
            <SelectItem value="deactivated">Deactivated</SelectItem>
            <SelectItem value="revoked">Revoked</SelectItem>
            <SelectItem value="assigned">Assigned</SelectItem>
            <SelectItem value="unassigned">Unassigned</SelectItem>
          </SelectContent>
        </Select>

        {/* Resource Type Filter */}
        <Select value={resourceTypeFilter} onValueChange={setResourceTypeFilter}>
          <SelectTrigger className="w-[140px] h-9">
            <SelectValue placeholder="Resource type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All resources</SelectItem>
            <SelectItem value="user">User</SelectItem>
            <SelectItem value="project">Project</SelectItem>
            <SelectItem value="application">Application</SelectItem>
            <SelectItem value="organization">Organization</SelectItem>
            <SelectItem value="role_assignment">Role Assignment</SelectItem>
            <SelectItem value="session">Session</SelectItem>
            <SelectItem value="settings">Settings</SelectItem>
          </SelectContent>
        </Select>

        {/* Severity Filter */}
        <Select value={severityFilter} onValueChange={setSeverityFilter}>
          <SelectTrigger className="w-[120px] h-9">
            <SelectValue placeholder="Severity" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All severity</SelectItem>
            <SelectItem value="routine">Routine</SelectItem>
            <SelectItem value="important">Important</SelectItem>
            <SelectItem value="sensitive">Sensitive</SelectItem>
          </SelectContent>
        </Select>

        {hasActiveFilters && (
          <Button variant="ghost" size="sm" onClick={clearFilters} className="h-9">
            <X className="h-4 w-4 mr-1" />
            Clear
          </Button>
        )}
      </div>

      {/* Data Table or Empty State */}
      {filteredActivity.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 px-4 rounded-lg border border-dashed bg-muted/30">
          <Activity className="h-12 w-12 text-muted-foreground/50 mb-4" />
          <h3 className="text-lg font-medium mb-1">No activity found</h3>
          <p className="text-sm text-muted-foreground text-center max-w-md mb-4">
            No events match your current filters. Try expanding the time range or clearing some filters to see more results.
          </p>
          <Button variant="outline" onClick={clearFilters}>
            Clear all filters
          </Button>
        </div>
      ) : (
        <DataTable
          columns={columns}
          data={filteredActivity}
          searchKey="resourceName"
          searchPlaceholder="Search resources..."
        />
      )}
    </div>
  )
}
