"use client"

import * as React from "react"
import { Check, ChevronsUpDown, Plus, Server, X } from "lucide-react"
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

export function InstanceSwitcher() {
  const { currentInstance, availableInstances, setCurrentInstance } = useAppContext()
  const [open, setOpen] = React.useState(false)
  const [showNewInstanceDialog, setShowNewInstanceDialog] = React.useState(false)

  const handleClearInstance = (e: React.MouseEvent) => {
    e.stopPropagation()
    setCurrentInstance(null)
  }

  return (
    <Dialog open={showNewInstanceDialog} onOpenChange={setShowNewInstanceDialog}>
      <div className="flex items-center gap-1">
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              role="combobox"
              aria-expanded={open}
              aria-label="Select an instance"
              className={cn(
                "justify-between",
                currentInstance ? "w-auto max-w-[200px]" : "w-[180px]"
              )}
            >
              <Server className="mr-2 h-4 w-4 shrink-0" />
              <span className="truncate">{currentInstance?.name ?? "Select instance"}</span>
              <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0">
          <Command>
            <CommandInput placeholder="Search instance..." />
            <CommandList>
              <CommandEmpty>No instance found.</CommandEmpty>
              <CommandGroup heading="Instances">
                {availableInstances.map((instance) => (
                  <CommandItem
                    key={instance.id}
                    onSelect={() => {
                      setCurrentInstance(instance)
                      setOpen(false)
                    }}
                    className="text-sm"
                  >
                    <Server className="mr-2 h-4 w-4" />
                    <span className="truncate">{instance.name}</span>
                    <Check
                      className={cn(
                        "ml-auto h-4 w-4",
                        currentInstance?.id === instance.id
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
                    setShowNewInstanceDialog(true)
                  }}
                >
                  <Plus className="mr-2 h-4 w-4" />
                  Add Instance
                </CommandItem>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Instance</DialogTitle>
          <DialogDescription>
            Add a new ZITADEL instance to manage.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-2">
          <div className="space-y-2">
            <Label htmlFor="name">Instance Name</Label>
            <Input id="name" placeholder="Production" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="domain">Domain</Label>
            <Input id="domain" placeholder="example.zitadel.cloud" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setShowNewInstanceDialog(false)}>
            Cancel
          </Button>
          <Button onClick={() => setShowNewInstanceDialog(false)}>
            Add Instance
          </Button>
        </DialogFooter>
      </DialogContent>
      {currentInstance && (
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 shrink-0"
          onClick={handleClearInstance}
          aria-label="Clear instance selection"
        >
          <X className="h-4 w-4" />
        </Button>
      )}
      </div>
    </Dialog>
  )
}
