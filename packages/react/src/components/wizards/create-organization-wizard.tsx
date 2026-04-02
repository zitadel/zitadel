"use client"

import * as React from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { Building2, Plus, Trash2, User } from "lucide-react"
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
import { Button } from "../ui/button"
import { Checkbox } from "../ui/checkbox"
import { useAppContext } from "../../context/app-context"
import { createOrganization } from "../../api/create-organization"
import { UserSearch } from "../ui/user-search"

/** Available org member roles from ListOrgMemberRoles */
const ORG_ROLES = [
  { role: "ORG_OWNER", label: "Owner", description: "Full control over the organization" },
  { role: "ORG_OWNER_VIEWER", label: "Owner Viewer", description: "Read-only view of org management" },
  { role: "ORG_USER_MANAGER", label: "User Manager", description: "Manage users in the org" },
  { role: "ORG_SETTINGS_MANAGER", label: "Settings Manager", description: "Manage org settings" },
  { role: "ORG_USER_PERMISSION_EDITOR", label: "User Permission Editor", description: "Edit user permissions" },
  { role: "ORG_PROJECT_PERMISSION_EDITOR", label: "Project Permission Editor", description: "Edit project permissions" },
  { role: "ORG_PROJECT_CREATOR", label: "Project Creator", description: "Create projects" },
  { role: "ORG_USER_SELF_MANAGER", label: "User Self Manager", description: "Self-manage user account" },
  { role: "ORG_ADMIN_IMPERSONATOR", label: "Admin Impersonator", description: "Impersonate org admins" },
  { role: "ORG_END_USER_IMPERSONATOR", label: "End User Impersonator", description: "Impersonate end users" },
] as const

interface AdminEntry {
  userId: string
  displayName: string
  username: string
  roles: string[]
}

const steps: WizardStep[] = [
  { id: "details", title: "Organization Details", description: "Name your organization" },
  { id: "admins", title: "Administrators", description: "Add users and assign roles (optional)" },
  { id: "confirmation", title: "Confirmation", description: "Review and create" },
]

interface CreateOrganizationWizardProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateOrganizationWizard({ open, onOpenChange }: CreateOrganizationWizardProps) {
  const router = useRouter()
  const { currentInstance } = useAppContext()

  const [orgName, setOrgName] = React.useState("")
  const [admins, setAdmins] = React.useState<AdminEntry[]>([])
  const [isCreating, setIsCreating] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const addAdmin = () => {
    setAdmins([...admins, { userId: "", displayName: "", username: "", roles: ["ORG_OWNER"] }])
  }

  const removeAdmin = (index: number) => {
    setAdmins(admins.filter((_, i) => i !== index))
  }

  const updateAdminUserId = (index: number, userId: string) => {
    const updated = [...admins]
    updated[index] = { ...updated[index], userId }
    setAdmins(updated)
  }

  const toggleAdminRole = (index: number, role: string) => {
    const updated = [...admins]
    const current = updated[index].roles
    if (current.includes(role)) {
      updated[index] = { ...updated[index], roles: current.filter(r => r !== role) }
    } else {
      updated[index] = { ...updated[index], roles: [...current, role] }
    }
    setAdmins(updated)
  }

  const handleComplete = async () => {
    if (!orgName.trim()) return
    setIsCreating(true)
    setError(null)
    try {
      // TODO: pass admins to createOrganization when API wrapper supports it
      await createOrganization(orgName.trim())
      onOpenChange(false)
      setOrgName("")
      setAdmins([])
      router.refresh()
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to create organization")
    } finally {
      setIsCreating(false)
    }
  }

  return (
    <StepWizard
      steps={steps}
      open={open}
      onOpenChange={onOpenChange}
      title="Create Organization"
      onComplete={handleComplete}
    >
      {/* Step 1: Organization Name */}
      <StepContent stepId="details">
        <FormSection
          title="Organization Information"
          description="The name must be unique across the instance."
        >
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="orgName">Organization Name</Label>
              <Input
                id="orgName"
                placeholder="Acme Corporation"
                value={orgName}
                onChange={(e) => setOrgName(e.target.value)}
                maxLength={200}
              />
              <p className="text-xs text-muted-foreground">
                Max 200 characters. This will be displayed across the platform.
              </p>
            </div>
          </div>
        </FormSection>

        <InfoBox
          title="Instance"
          description={`Organization will be created in ${currentInstance?.name || "the current instance"}`}
          variant="default"
        />

        <StepActions nextDisabled={!orgName.trim()} />
      </StepContent>

      {/* Step 2: Administrators (optional) */}
      <StepContent stepId="admins">
        <FormSection
          title="Organization Administrators"
          description="Optionally add users as administrators. If no admins are specified, the organization can still be managed by any instance administrator."
        >
          <div className="space-y-4">
            {admins.map((admin, index) => (
              <div key={index} className="rounded-lg border p-4 space-y-3">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 space-y-2">
                    {admin.userId ? (
                      <div className="flex items-center gap-2 rounded-md border px-3 py-2">
                        <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 shrink-0">
                          <User className="h-3 w-3 text-primary" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <span className="font-medium text-sm">{admin.displayName}</span>
                          <span className="text-xs text-muted-foreground ml-1.5">@{admin.username}</span>
                        </div>
                      </div>
                    ) : (
                      <UserSearch
                        placeholder="Search for a user..."
                        onSelect={(user) => {
                          const updated = [...admins]
                          updated[index] = {
                            ...updated[index],
                            userId: user.userId,
                            displayName: user.displayName,
                            username: user.username,
                          }
                          setAdmins(updated)
                        }}
                      />
                    )}
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className={`h-9 w-9 shrink-0 text-destructive hover:text-destructive ${admin.userId ? 'mt-0.5' : 'mt-0'}`}
                    onClick={() => removeAdmin(index)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>

                <div className="space-y-2">
                  <Label className="text-xs">Roles (defaults to ORG_OWNER if none selected)</Label>
                  <div className="grid grid-cols-1 gap-1.5">
                    {ORG_ROLES.map(({ role, label, description }) => (
                      <div key={role} className="flex items-start gap-2">
                        <Checkbox
                          id={`${index}-${role}`}
                          checked={admin.roles.includes(role)}
                          onCheckedChange={() => toggleAdminRole(index, role)}
                        />
                        <label htmlFor={`${index}-${role}`} className="text-sm cursor-pointer leading-relaxed">
                          <span className="font-medium">{label}</span>
                          <span className="text-muted-foreground ml-1">— {description}</span>
                        </label>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}

            <Button variant="outline" onClick={addAdmin} className="w-full">
              <Plus className="mr-2 h-4 w-4" />
              Add Administrator
            </Button>
          </div>
        </FormSection>

        <InfoBox
          title="Note"
          description="Instance administrators (IAM_OWNER) can always manage organizations regardless of org membership."
          variant="default"
        />

        <StepActions />
      </StepContent>

      {/* Step 3: Confirmation */}
      <StepContent stepId="confirmation">
        <FormSection title="Review Organization">
          <div className="rounded-lg border divide-y">
            <ParameterRow label="Name" value={orgName || "—"} />
            <ParameterRow
              label="Administrators"
              value={
                admins.length === 0
                  ? "None — managed by instance admins"
                  : `${admins.length} user${admins.length > 1 ? "s" : ""}`
              }
            />
          </div>
        </FormSection>

        {admins.length > 0 && (
          <FormSection title="Admin Details" className="mt-4">
            <div className="rounded-lg border divide-y">
              {admins.map((admin, i) => (
                <ParameterRow
                  key={i}
                  label={admin.displayName ? `${admin.displayName} (@${admin.username})` : "No user selected"}
                  value={admin.roles.length > 0 ? admin.roles.join(", ") : "ORG_OWNER (default)"}
                />
              ))}
            </div>
          </FormSection>
        )}

        {error && (
          <div className="rounded-md border border-destructive/50 bg-destructive/10 p-3">
            <p className="text-sm text-destructive">{error}</p>
          </div>
        )}

        <InfoBox
          title="What happens next"
          description="The organization will be immediately available after creation."
          variant="success"
        >
          <ul className="space-y-1 text-xs">
            <li className="flex items-center gap-1.5">
              <span className="h-1 w-1 rounded-full bg-green-600" />
              Organization is created with a generated domain
            </li>
            {admins.length > 0 && (
              <li className="flex items-center gap-1.5">
                <span className="h-1 w-1 rounded-full bg-green-600" />
                {admins.length} admin{admins.length > 1 ? "s" : ""} will be assigned
              </li>
            )}
            <li className="flex items-center gap-1.5">
              <span className="h-1 w-1 rounded-full bg-green-600" />
              Custom domains can be added afterwards
            </li>
          </ul>
        </InfoBox>

        <StepActions nextLabel={isCreating ? "Creating..." : "Create Organization"} />
      </StepContent>
    </StepWizard>
  )
}
