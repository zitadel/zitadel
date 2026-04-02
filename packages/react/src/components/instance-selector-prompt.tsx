"use client"

import * as React from "react"
import { Server, Plus, Cloud, HardDrive } from "lucide-react"
import { Input } from "./ui/input"
import { Button } from "./ui/button"
import { Badge } from "./ui/badge"
import { useAppContext } from "../context/app-context"
import { useConsoleRouter as useRouter } from "../hooks/use-console-router"

interface InstanceSelectorPromptProps {
  title: string
  description?: string
  icon?: React.ReactNode
  targetPath?: string
}

export function InstanceSelectorPrompt({ 
  title, 
  description = "Choose an instance to continue",
  icon,
  targetPath
}: InstanceSelectorPromptProps) {
  const router = useRouter()
  const { availableInstances, setCurrentInstance, currentInstance } = useAppContext()
  const [search, setSearch] = React.useState("")

  const filteredInstances = availableInstances.filter(instance =>
    instance.name.toLowerCase().includes(search.toLowerCase()) ||
    instance.domain.toLowerCase().includes(search.toLowerCase())
  )

  const handleSelectInstance = (instance: typeof availableInstances[0]) => {
    setCurrentInstance(instance)
    if (targetPath) {
      router.push(targetPath)
    }
  }

  // If instance is already selected, don't show the prompt
  if (currentInstance) {
    return null
  }

  return (
    <div className="flex flex-col items-center justify-center py-16">
      <div className="flex flex-col items-center gap-4 max-w-md w-full">
        {/* Icon */}
        <div className="flex h-12 w-12 items-center justify-center rounded-lg border bg-muted">
          {icon || <Server className="h-6 w-6 text-muted-foreground" />}
        </div>

        {/* Title */}
        <div className="text-center">
          <h2 className="text-lg font-semibold">{title}</h2>
          <p className="text-sm text-muted-foreground">{description}</p>
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
            {filteredInstances.map((instance) => (
              <button
                key={instance.id}
                onClick={() => handleSelectInstance(instance)}
                className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted transition-colors border-b last:border-b-0 text-left"
              >
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
                  {instance.hostingType === "cloud" ? (
                    <Cloud className="h-4 w-4 text-foreground" />
                  ) : (
                    <HardDrive className="h-4 w-4 text-foreground" />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium truncate">{instance.name}</span>
                    <Badge 
                      variant={instance.status === "active" ? "default" : "secondary"}
                      className="text-xs shrink-0"
                    >
                      {instance.status}
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground truncate">{instance.domain}</p>
                </div>
              </button>
            ))}
          </div>

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
