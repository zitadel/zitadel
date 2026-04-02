"use client"

import { useState } from "react"
import { actions } from "../../mock-data"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card"
import { Badge } from "../../components/ui/badge"
import { Button } from "../../components/ui/button"
import { Plus, Zap, Edit, Trash2, Play, Pause } from "lucide-react"
import { useAppContext } from "../../context/app-context"
import { InstanceSelectorPrompt } from "../../components/instance-selector-prompt"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../../components/ui/dialog"
import { Input } from "../../components/ui/input"
import { Label } from "../../components/ui/label"
import { Textarea } from "../../components/ui/textarea"
import { Switch } from "../../components/ui/switch"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../../components/ui/dropdown-menu"
import { MoreHorizontal } from "lucide-react"

function CreateActionDialog() {
  const [open, setOpen] = useState(false)

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Create Action
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>Create Action</DialogTitle>
          <DialogDescription>
            Create a new action to execute custom logic during authentication flows.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="name">Action Name</Label>
            <Input id="name" placeholder="My Custom Action" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Input id="description" placeholder="What does this action do?" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="script">Script</Label>
            <Textarea 
              id="script" 
              placeholder="function execute(ctx) {&#10;  // Your code here&#10;}"
              className="font-mono min-h-[200px]"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="timeout">Timeout (ms)</Label>
              <Input id="timeout" type="number" placeholder="5000" />
            </div>
            <div className="flex items-center justify-between space-y-2 pt-6">
              <Label htmlFor="allowFail">Allow to Fail</Label>
              <Switch id="allowFail" />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button onClick={() => setOpen(false)}>Create Action</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default function ActionsPage() {
  const { currentInstance } = useAppContext()

  if (!currentInstance) {
    return (
      <InstanceSelectorPrompt 
        title="Continue to Actions"
        description="Choose an instance to view actions"
        icon={<Zap className="h-6 w-6 text-muted-foreground" />}
        targetPath="/actions"
      />
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Actions</h1>
          <p className="text-muted-foreground">
            Custom actions to execute during authentication flows
          </p>
        </div>
        <CreateActionDialog />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {actions.map((action) => (
          <Card key={action.id}>
            <CardHeader className="flex flex-row items-start justify-between space-y-0">
              <div className="flex items-start gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-md bg-primary/10">
                  <Zap className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <CardTitle className="text-lg">{action.name}</CardTitle>
                  <CardDescription>{action.description}</CardDescription>
                </div>
              </div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem>
                    <Edit className="mr-2 h-4 w-4" />
                    Edit
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    {action.status === "active" ? (
                      <>
                        <Pause className="mr-2 h-4 w-4" />
                        Deactivate
                      </>
                    ) : (
                      <>
                        <Play className="mr-2 h-4 w-4" />
                        Activate
                      </>
                    )}
                  </DropdownMenuItem>
                  <DropdownMenuItem className="text-destructive">
                    <Trash2 className="mr-2 h-4 w-4" />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-4 text-sm">
                <Badge 
                  variant="outline"
                  className={action.status === "active" 
                    ? "bg-foreground/10 text-foreground border-foreground/20"
                    : "bg-muted text-muted-foreground border-border"
                  }
                >
                  {action.status}
                </Badge>
                <span className="text-muted-foreground">
                  Timeout: {action.timeout}ms
                </span>
                {action.allowedToFail && (
                  <Badge variant="secondary">Can fail</Badge>
                )}
              </div>
              <div className="mt-4">
                <pre className="text-xs bg-muted p-3 rounded-md overflow-x-auto">
                  <code>{action.script}</code>
                </pre>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
