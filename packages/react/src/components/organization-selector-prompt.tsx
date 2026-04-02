"use client"

import * as React from "react"
import { Building2, Plus, Cloud, Server } from "lucide-react"
import { Input } from "./ui/input"
import { Badge } from "./ui/badge"
import { useAppContext } from "../context/app-context"
import { useConsoleRouter as useRouter } from "../hooks/use-console-router"
import { getOrganizationsByInstance, instances} from "../mock-data"

interface OrganizationSelectorPromptProps {
  title: string
  description?: string
  icon?: React.ReactNode
  targetPath?: string
}

export function OrganizationSelectorPrompt({ 
  title, 
  description = "Choose an organization to continue",
  icon,
  targetPath
}: OrganizationSelectorPromptProps) {
  const router = useRouter()
  const { currentInstance, currentOrganization, setCurrentOrganization, setCurrentInstance } = useAppContext()
  const [search, setSearch] = React.useState("")

  // Get organizations for the current instance
  const availableOrganizations = currentInstance 
    ? getOrganizationsByInstance(currentInstance.id)
    : []

  const filteredOrganizations = availableOrganizations.filter(org =>
    org.name.toLowerCase().includes(search.toLowerCase())
  )

  // Filter instances for when no instance is selected
  const allFilteredInstances = instances.filter(instance =>
    instance.name.toLowerCase().includes(search.toLowerCase()) ||
    instance.domain.toLowerCase().includes(search.toLowerCase())
  )
  // Show max 20 instances for better performance
  const filteredInstances = allFilteredInstances.slice(0, 20)
  const hasMoreInstances = allFilteredInstances.length > 20

  const handleSelectOrganization = (organization: typeof availableOrganizations[0]) => {
    setCurrentOrganization(organization)
    if (targetPath) {
      router.push(targetPath)
    }
  }

  const handleSelectInstance = (instance: typeof instances[0]) => {
    setCurrentInstance(instance)
    setSearch("") // Reset search when switching to org selection
  }

  // If organization is already selected, don't show the prompt
  if (currentOrganization) {
    return null
  }

  // If no instance is selected, show instance selection first
  if (!currentInstance) {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <div className="flex flex-col items-center gap-4 max-w-md w-full">
          {/* Icon */}
          <div className="flex h-12 w-12 items-center justify-center rounded-lg border bg-muted">
            <Building2 className="h-6 w-6 text-muted-foreground" />
          </div>

          {/* Title */}
          <div className="text-center">
            <h2 className="text-lg font-semibold">Continue to Organizations</h2>
            <p className="text-sm text-muted-foreground">Choose an instance to view organizations</p>
          </div>

          {/* Search */}
          <div className="w-full">
            <Input
              placeholder="Find Instance..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full"
            />
          </div>

          {/* Instance List */}
          <div className="w-full border rounded-lg overflow-hidden bg-background">
            <div className="max-h-80 overflow-y-auto">
              {filteredInstances.length === 0 ? (
                <div className="px-4 py-8 text-center text-sm text-muted-foreground">
                  No instances found
                </div>
              ) : (
                filteredInstances.map((instance) => (
                  <button
                    key={instance.id}
                    onClick={() => handleSelectInstance(instance)}
                    className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted transition-colors border-b last:border-b-0 text-left"
                  >
                    <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
                      {instance.hostingType === "cloud" ? (
                        <Cloud className="h-4 w-4 text-muted-foreground" />
                      ) : (
                        <Server className="h-4 w-4 text-muted-foreground" />
                      )}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="font-medium truncate">{instance.name}</span>
                        <Badge 
                          variant={instance.status === "active" ? "outline" : "secondary"} 
                          className="text-xs shrink-0"
                        >
                          {instance.status}
                        </Badge>
                      </div>
                      <p className="text-xs text-muted-foreground truncate">
                        {instance.domain}
                      </p>
                    </div>
                  </button>
                ))
              )}
            </div>

            {/* Show more indicator */}
            {hasMoreInstances && (
              <div className="px-4 py-2 text-xs text-muted-foreground text-center border-t bg-muted/50">
                Showing 20 of {allFilteredInstances.length} instances. Type to search for more.
              </div>
            )}

            {/* Add Instance */}
            <button
              onClick={() => router.push("/instances/new")}
              className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted transition-colors border-t text-left text-muted-foreground hover:text-foreground"
            >
              <Plus className="h-4 w-4" />
              <span>Add Instance</span>
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center justify-center py-16">
      <div className="flex flex-col items-center gap-4 max-w-md w-full">
        {/* Icon */}
        <div className="flex h-12 w-12 items-center justify-center rounded-lg border bg-muted">
          {icon || <Building2 className="h-6 w-6 text-muted-foreground" />}
        </div>

        {/* Title */}
        <div className="text-center">
          <h2 className="text-lg font-semibold">{title}</h2>
          <p className="text-sm text-muted-foreground">{description}</p>
        </div>

        {/* Search */}
        <div className="w-full">
          <Input
            placeholder="Find Organization..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full"
          />
        </div>

        {/* Organization List */}
        <div className="w-full border rounded-lg overflow-hidden bg-background">
          <div className="max-h-80 overflow-y-auto">
            {filteredOrganizations.length === 0 ? (
              <div className="px-4 py-8 text-center text-sm text-muted-foreground">
                No organizations found
              </div>
            ) : (
              filteredOrganizations.map((org) => (
                <button
                  key={org.id}
                  onClick={() => handleSelectOrganization(org)}
                  className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted transition-colors border-b last:border-b-0 text-left"
                >
                  <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
                    <Building2 className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-medium truncate">{org.name}</span>
                      {org.isDefault && (
                        <Badge variant="outline" className="text-xs shrink-0">
                          Default
                        </Badge>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground truncate">
                      {org.userCount} users · {org.projectCount} projects
                    </p>
                  </div>
                </button>
              ))
            )}
          </div>

          {/* Create Organization */}
          <button
            onClick={() => router.push("/organizations/new")}
            className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted transition-colors border-t text-left text-muted-foreground hover:text-foreground"
          >
            <Plus className="h-4 w-4" />
            <span>Create Organization</span>
          </button>
        </div>
      </div>
    </div>
  )
}
