"use client"

import * as React from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { FolderKanban, Shield, Users, Lock, Building2 } from "lucide-react"
import {
  StepWizard,
  StepContent,
  StepActions,
  FormSection,
  ParameterRow,
  InfoBox,
  type WizardStep,
} from "../ui/step-wizard"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Textarea } from "../ui/textarea"
import { RadioGroup, RadioGroupItem } from "../ui/radio-group"
import { Checkbox } from "../ui/checkbox"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select"
import { cn } from "../../utils"
import { useAppContext } from "../../context/app-context"
import { organizations } from "../../mock-data"

const steps: WizardStep[] = [
  { id: "details", title: "Project Details", description: "Basic information" },
  { id: "roles", title: "Role Configuration", description: "Set up access roles" },
  { id: "settings", title: "Settings", description: "Configure project" },
  { id: "confirmation", title: "Confirmation", description: "Review and create" },
]

const defaultRoles = [
  { id: "viewer", name: "Viewer", description: "Read-only access", permissions: ["view"] },
  { id: "editor", name: "Editor", description: "Can edit content", permissions: ["view", "edit"] },
  { id: "admin", name: "Admin", description: "Full project access", permissions: ["view", "edit", "admin"] },
]

interface CreateProjectWizardProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateProjectWizard({ open, onOpenChange }: CreateProjectWizardProps) {
  const router = useRouter()
  const { currentInstance, currentOrganization } = useAppContext()
  
  const [projectName, setProjectName] = React.useState("")
  const [description, setDescription] = React.useState("")
  const [selectedOrgId, setSelectedOrgId] = React.useState(currentOrganization?.id || "")
  const [selectedRoles, setSelectedRoles] = React.useState<string[]>(["viewer", "editor", "admin"])
  const [customRoleName, setCustomRoleName] = React.useState("")
  const [authRequired, setAuthRequired] = React.useState(true)
  const [checkAuthorization, setCheckAuthorization] = React.useState(true)
  const [privateProject, setPrivateProject] = React.useState(false)

  // Filter organizations by current instance
  const availableOrgs = React.useMemo(() => {
    if (!currentInstance) return organizations.slice(0, 10)
    return organizations.filter(org => org.instanceId === currentInstance.id).slice(0, 20)
  }, [currentInstance])

  React.useEffect(() => {
    if (currentOrganization) {
      setSelectedOrgId(currentOrganization.id)
    }
  }, [currentOrganization])

  const handleComplete = () => {
    onOpenChange(false)
    if (currentOrganization) {
      router.push("/org/projects")
    } else {
      router.push("/projects")
    }
  }

  const toggleRole = (roleId: string) => {
    setSelectedRoles(prev => 
      prev.includes(roleId) 
        ? prev.filter(r => r !== roleId)
        : [...prev, roleId]
    )
  }

  const selectedOrg = organizations.find(o => o.id === selectedOrgId)

  return (
    <StepWizard
      steps={steps}
      open={open}
      onOpenChange={onOpenChange}
      title="Create Project"
      onComplete={handleComplete}
    >
      {/* Step 1: Project Details */}
      <StepContent stepId="details">
        <FormSection
          title="Project Information"
          description="Enter basic details about your project"
        >
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="projectName">Project Name</Label>
              <Input
                id="projectName"
                placeholder="My Web Application"
                value={projectName}
                onChange={(e) => setProjectName(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">
                Used to identify this project across the platform
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description (Optional)</Label>
              <Textarea
                id="description"
                placeholder="A brief description of your project..."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
              />
            </div>
          </div>
        </FormSection>

        {!currentOrganization && (
          <FormSection title="Organization" description="Select which organization owns this project" className="mt-6">
            <Select value={selectedOrgId} onValueChange={setSelectedOrgId}>
              <SelectTrigger>
                <SelectValue placeholder="Select an organization" />
              </SelectTrigger>
              <SelectContent>
                {availableOrgs.map((org) => (
                  <SelectItem key={org.id} value={org.id}>
                    <div className="flex items-center gap-2">
                      <Building2 className="h-4 w-4 text-muted-foreground" />
                      {org.name}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </FormSection>
        )}

        {currentOrganization && (
          <InfoBox title="Organization" className="mt-6">
            <p className="text-sm">
              Project will be created in <strong>{currentOrganization.name}</strong>
            </p>
          </InfoBox>
        )}

        <StepActions nextDisabled={!projectName || (!selectedOrgId && !currentOrganization)} />
      </StepContent>

      {/* Step 2: Role Configuration */}
      <StepContent stepId="roles">
        <FormSection
          title="Project Roles"
          description="Define roles for user access control in this project"
        >
          <div className="space-y-3">
            {defaultRoles.map((role) => (
              <label
                key={role.id}
                className={cn(
                  "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                  selectedRoles.includes(role.id) ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
                )}
              >
                <Checkbox
                  checked={selectedRoles.includes(role.id)}
                  onCheckedChange={() => toggleRole(role.id)}
                />
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <Shield className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">{role.name}</span>
                  </div>
                  <p className="text-xs text-muted-foreground mt-1">{role.description}</p>
                  <div className="flex flex-wrap gap-1 mt-2">
                    {role.permissions.map((perm) => (
                      <span key={perm} className="text-xs bg-muted px-2 py-0.5 rounded">
                        {perm}
                      </span>
                    ))}
                  </div>
                </div>
              </label>
            ))}
          </div>
        </FormSection>

        <FormSection title="Custom Role (Optional)" className="mt-6">
          <div className="flex items-center gap-2">
            <Input
              placeholder="Enter custom role name..."
              value={customRoleName}
              onChange={(e) => setCustomRoleName(e.target.value)}
            />
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            You can add more custom roles after the project is created
          </p>
        </FormSection>

        <StepActions nextDisabled={selectedRoles.length === 0} />
      </StepContent>

      {/* Step 3: Settings */}
      <StepContent stepId="settings">
        <FormSection
          title="Authentication Settings"
          description="Configure how users access this project"
        >
          <div className="space-y-4">
            <div className="flex items-start gap-2">
              <Checkbox
                id="authRequired"
                checked={authRequired}
                onCheckedChange={(checked) => setAuthRequired(checked === true)}
              />
              <label htmlFor="authRequired" className="text-sm cursor-pointer leading-relaxed">
                <span className="font-medium">Require authentication</span>
                <p className="text-xs text-muted-foreground">
                  Users must be logged in to access this project
                </p>
              </label>
            </div>

            <div className="flex items-start gap-2">
              <Checkbox
                id="checkAuth"
                checked={checkAuthorization}
                onCheckedChange={(checked) => setCheckAuthorization(checked === true)}
              />
              <label htmlFor="checkAuth" className="text-sm cursor-pointer leading-relaxed">
                <span className="font-medium">Check authorization</span>
                <p className="text-xs text-muted-foreground">
                  Verify user has required roles before granting access
                </p>
              </label>
            </div>
          </div>
        </FormSection>

        <FormSection title="Visibility" className="mt-6">
          <RadioGroup 
            value={privateProject ? "private" : "public"} 
            onValueChange={(v) => setPrivateProject(v === "private")}
            className="space-y-2"
          >
            <label
              className={cn(
                "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                !privateProject ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="public" id="public" className="mt-0.5" />
              <div>
                <div className="flex items-center gap-2">
                  <Users className="h-4 w-4" />
                  <span className="font-medium text-sm">Organization Visible</span>
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  All organization members can see this project
                </p>
              </div>
            </label>

            <label
              className={cn(
                "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                privateProject ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="private" id="private" className="mt-0.5" />
              <div>
                <div className="flex items-center gap-2">
                  <Lock className="h-4 w-4" />
                  <span className="font-medium text-sm">Private</span>
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  Only users with explicit access can see this project
                </p>
              </div>
            </label>
          </RadioGroup>
        </FormSection>

        <StepActions />
      </StepContent>

      {/* Step 4: Confirmation */}
      <StepContent stepId="confirmation">
        <FormSection title="Review Project">
          <div className="rounded-lg border divide-y">
            <ParameterRow label="Name" value={projectName || "—"} />
            {description && <ParameterRow label="Description" value={description} />}
            <ParameterRow label="Organization" value={selectedOrg?.name || currentOrganization?.name || "—"} />
            <ParameterRow label="Roles" value={selectedRoles.length.toString()} />
            <ParameterRow label="Visibility" value={privateProject ? "Private" : "Organization"} />
            <ParameterRow label="Authentication" value={authRequired ? "Required" : "Optional"} />
          </div>
        </FormSection>

        <InfoBox title="Configured Roles" variant="default">
          <div className="flex flex-wrap gap-2 mt-2">
            {selectedRoles.map((roleId) => {
              const role = defaultRoles.find(r => r.id === roleId)
              return (
                <span key={roleId} className="text-xs bg-muted px-2 py-1 rounded font-medium">
                  {role?.name || roleId}
                </span>
              )
            })}
            {customRoleName && (
              <span className="text-xs bg-muted px-2 py-1 rounded font-medium">
                {customRoleName}
              </span>
            )}
          </div>
        </InfoBox>

        <InfoBox
          title="Ready to Create"
          description="Your project will be available immediately after creation."
          variant="success"
        />

        <StepActions nextLabel="Create Project" />
      </StepContent>
    </StepWizard>
  )
}
