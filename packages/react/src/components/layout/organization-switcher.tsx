"use client"

import * as React from "react"
import { Check, ChevronsUpDown, Plus, Building2, X, Loader2 } from "lucide-react"
import { cn } from "../../utils"
import { Button } from "../ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "../ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../ui/popover"
import { Badge } from "../ui/badge"
import { useAppContext } from "../../context/app-context"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { usePathname } from "next/navigation"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { createOrganization } from "../../api/create-organization"

/** Pages that are only visible at instance level (no org selected) */
const instanceOnlyPaths = ["/organizations", "/sessions"]

export function OrganizationSwitcher() {
  const { currentOrganization, availableOrganizations, setCurrentOrganization } = useAppContext()
  const router = useRouter()
  const pathname = usePathname()
  const [open, setOpen] = React.useState(false)
  const [showNewOrgDialog, setShowNewOrgDialog] = React.useState(false)
  const [newOrgName, setNewOrgName] = React.useState("")
  const [isCreating, setIsCreating] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const handleClearOrganization = (e: React.MouseEvent) => {
    e.stopPropagation()
    setCurrentOrganization(null)
  }

  const handleCreateOrganization = async () => {
    if (!newOrgName.trim()) return
    setIsCreating(true)
    setError(null)

    try {
      await createOrganization(newOrgName.trim())
      setShowNewOrgDialog(false)
      setNewOrgName("")
      // Refresh the page to reload org lists from the server
      router.refresh()
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to create organization")
    } finally {
      setIsCreating(false)
    }
  }

  return (
    <Dialog open={showNewOrgDialog} onOpenChange={(open) => {
      setShowNewOrgDialog(open)
      if (!open) {
        setNewOrgName("")
        setError(null)
      }
    }}>
      <div className="flex items-center gap-1">
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              role="combobox"
              aria-expanded={open}
              aria-label="Select an organization"
              className={cn(
                "justify-between",
                currentOrganization ? "w-auto max-w-[220px]" : "w-[200px]"
              )}
            >
              <Building2 className="mr-2 h-4 w-4 shrink-0" />
              <span className="truncate flex-1 text-left">
                {currentOrganization?.name ?? "Select organization"}
              </span>
              {currentOrganization?.isDefault && (
                <Badge variant="secondary" className="ml-2 shrink-0 text-xs">
                  default
                </Badge>
              )}
              <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </PopoverTrigger>
        <PopoverContent className="w-[280px] p-0">
          <Command>
            <CommandInput placeholder="Search organization..." />
            <CommandList>
              <CommandEmpty>No organization found.</CommandEmpty>
              <CommandGroup heading="Organizations">
                {availableOrganizations.map((org) => (
                  <CommandItem
                    key={org.id}
                    onSelect={() => {
                      setCurrentOrganization(org)
                      setOpen(false)
                      // If on an instance-only page (e.g. /organizations), navigate to overview
                      if (instanceOnlyPaths.some(p => pathname.endsWith(p) || pathname.endsWith(p + "/"))) {
                        const overviewPath = pathname.replace(/\/[^/]+\/?$/, '/overview')
                        router.push(overviewPath)
                      }
                    }}
                    className="text-sm"
                  >
                    <Building2 className="mr-2 h-4 w-4" />
                    <span className="truncate flex-1">{org.name}</span>
                    {org.isDefault && (
                      <Badge variant="secondary" className="ml-2 text-xs">
                        default
                      </Badge>
                    )}
                    <Check
                      className={cn(
                        "ml-2 h-4 w-4 shrink-0",
                        currentOrganization?.id === org.id
                          ? "opacity-100"
                          : "opacity-0"
                      )}
                    />
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
            <CommandSeparator />
            <CommandList>
              <CommandGroup>
                <CommandItem
                  onSelect={() => {
                    setOpen(false)
                    setShowNewOrgDialog(true)
                  }}
                >
                  <Plus className="mr-2 h-4 w-4" />
                  Add Organization
                </CommandItem>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Organization</DialogTitle>
          <DialogDescription>
            Create a new organization in this instance.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-2">
          <div className="space-y-2">
            <Label htmlFor="org-name">Organization Name</Label>
            <Input
              id="org-name"
              placeholder="Acme Corp"
              value={newOrgName}
              onChange={(e) => setNewOrgName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && newOrgName.trim()) {
                  handleCreateOrganization()
                }
              }}
              disabled={isCreating}
            />
          </div>
          {error && (
            <div className="rounded-md border border-destructive/50 bg-destructive/10 p-3">
              <p className="text-sm text-destructive">{error}</p>
            </div>
          )}
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => setShowNewOrgDialog(false)}
            disabled={isCreating}
          >
            Cancel
          </Button>
          <Button
            onClick={handleCreateOrganization}
            disabled={!newOrgName.trim() || isCreating}
          >
            {isCreating ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Creating...
              </>
            ) : (
              "Create Organization"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
      {currentOrganization && (
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 shrink-0"
          onClick={handleClearOrganization}
          aria-label="Clear organization selection"
        >
          <X className="h-4 w-4" />
        </Button>
      )}
      </div>
    </Dialog>
  )
}
