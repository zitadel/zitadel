"use client"

import * as React from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { User, Mail, Lock, Shield, Building2 } from "lucide-react"
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
import { organizations } from "../../mock-data"
import { useAppContext } from "../../context/app-context"

const steps: WizardStep[] = [
  { id: "organization", title: "Organization", description: "Select target organization" },
  { id: "profile", title: "User Profile", description: "Enter user details" },
  { id: "authentication", title: "Authentication", description: "Set up login method" },
  { id: "confirmation", title: "Confirmation", description: "Review and create" },
]

const roles = [
  { id: "user", name: "User", description: "Standard access to applications" },
  { id: "admin", name: "Admin", description: "Can manage organization settings" },
  { id: "owner", name: "Owner", description: "Full control over the organization" },
]

interface CreateUserWizardProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateUserWizard({ open, onOpenChange }: CreateUserWizardProps) {
  const router = useRouter()
  const { currentInstance, currentOrganization } = useAppContext()
  
  const [selectedOrgId, setSelectedOrgId] = React.useState(currentOrganization?.id || "")
  const [firstName, setFirstName] = React.useState("")
  const [lastName, setLastName] = React.useState("")
  const [email, setEmail] = React.useState("")
  const [username, setUsername] = React.useState("")
  const [role, setRole] = React.useState("user")
  const [authMethod, setAuthMethod] = React.useState<"password" | "passwordless" | "invite">("invite")
  const [password, setPassword] = React.useState("")
  const [sendWelcomeEmail, setSendWelcomeEmail] = React.useState(true)
  const [requirePasswordChange, setRequirePasswordChange] = React.useState(true)

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
      router.push("/org/users")
    } else {
      router.push("/users")
    }
  }

  const selectedOrg = organizations.find(o => o.id === selectedOrgId)
  const selectedRole = roles.find(r => r.id === role)

  return (
    <StepWizard
      steps={steps}
      open={open}
      onOpenChange={onOpenChange}
      title="Create User"
      onComplete={handleComplete}
    >
      {/* Step 1: Organization */}
      <StepContent stepId="organization">
        <FormSection
          title="Select Organization"
          description="Choose which organization this user will belong to"
        >
          {currentOrganization ? (
            <InfoBox title="Organization Pre-selected">
              <p className="text-sm">
                User will be created in <strong>{currentOrganization.name}</strong>
              </p>
            </InfoBox>
          ) : (
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
          )}
        </FormSection>

        <FormSection title="User Role" description="Set the user's permissions level" className="mt-6">
          <RadioGroup value={role} onValueChange={setRole} className="space-y-2">
            {roles.map((r) => (
              <label
                key={r.id}
                className={cn(
                  "flex items-start gap-3 p-3 rounded-lg border cursor-pointer transition-colors",
                  role === r.id ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
                )}
              >
                <RadioGroupItem value={r.id} id={r.id} className="mt-0.5" />
                <div>
                  <span className="font-medium text-sm">{r.name}</span>
                  <p className="text-xs text-muted-foreground">{r.description}</p>
                </div>
              </label>
            ))}
          </RadioGroup>
        </FormSection>

        <StepActions nextDisabled={!selectedOrgId && !currentOrganization} />
      </StepContent>

      {/* Step 2: Profile */}
      <StepContent stepId="profile">
        <FormSection title="User Information">
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">First Name</Label>
                <Input
                  id="firstName"
                  placeholder="John"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Last Name</Label>
                <Input
                  id="lastName"
                  placeholder="Doe"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                type="email"
                placeholder="john.doe@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">
                Primary contact and login identifier
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="username">Username (Optional)</Label>
              <Input
                id="username"
                placeholder="johndoe"
                value={username}
                onChange={(e) => setUsername(e.target.value.toLowerCase().replace(/[^a-z0-9._-]/g, ""))}
              />
              <p className="text-xs text-muted-foreground">
                Alternative login identifier
              </p>
            </div>
          </div>
        </FormSection>

        <StepActions nextDisabled={!firstName || !lastName || !email} />
      </StepContent>

      {/* Step 3: Authentication */}
      <StepContent stepId="authentication">
        <FormSection
          title="Authentication Method"
          description="Choose how the user will sign in"
        >
          <RadioGroup value={authMethod} onValueChange={(v) => setAuthMethod(v as typeof authMethod)} className="space-y-2">
            <label
              className={cn(
                "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                authMethod === "invite" ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="invite" id="invite" className="mt-0.5" />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <Mail className="h-4 w-4" />
                  <span className="font-medium text-sm">Send Invite Email</span>
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  User will receive an email to set up their own password
                </p>
              </div>
            </label>

            <label
              className={cn(
                "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                authMethod === "password" ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="password" id="password" className="mt-0.5" />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <Lock className="h-4 w-4" />
                  <span className="font-medium text-sm">Set Password</span>
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  Create a password for the user now
                </p>
              </div>
            </label>

            <label
              className={cn(
                "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                authMethod === "passwordless" ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="passwordless" id="passwordless" className="mt-0.5" />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4" />
                  <span className="font-medium text-sm">Passwordless Only</span>
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  User will use passkeys or magic links to sign in
                </p>
              </div>
            </label>
          </RadioGroup>
        </FormSection>

        {authMethod === "password" && (
          <FormSection className="mt-6">
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="password">Initial Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Enter a strong password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>
              <div className="flex items-center gap-2">
                <Checkbox
                  id="requireChange"
                  checked={requirePasswordChange}
                  onCheckedChange={(checked) => setRequirePasswordChange(checked === true)}
                />
                <label htmlFor="requireChange" className="text-sm cursor-pointer">
                  Require password change on first login
                </label>
              </div>
            </div>
          </FormSection>
        )}

        {authMethod === "invite" && (
          <FormSection className="mt-6">
            <div className="flex items-center gap-2">
              <Checkbox
                id="welcomeEmail"
                checked={sendWelcomeEmail}
                onCheckedChange={(checked) => setSendWelcomeEmail(checked === true)}
              />
              <label htmlFor="welcomeEmail" className="text-sm cursor-pointer">
                Include welcome message in invitation email
              </label>
            </div>
          </FormSection>
        )}

        <StepActions nextDisabled={authMethod === "password" && !password} />
      </StepContent>

      {/* Step 4: Confirmation */}
      <StepContent stepId="confirmation">
        <FormSection title="Review User Details">
          <div className="rounded-lg border divide-y">
            <ParameterRow label="Organization" value={selectedOrg?.name || currentOrganization?.name || "—"} />
            <ParameterRow label="Name" value={`${firstName} ${lastName}`} />
            <ParameterRow label="Email" value={email || "—"} />
            {username && <ParameterRow label="Username" value={username} />}
            <ParameterRow label="Role" value={selectedRole?.name || "—"} />
            <ParameterRow 
              label="Authentication" 
              value={
                authMethod === "invite" ? "Email Invitation" :
                authMethod === "password" ? "Password" : "Passwordless"
              } 
            />
          </div>
        </FormSection>

        <InfoBox
          title="User Access"
          description={selectedRole?.description}
          variant="default"
        />

        {authMethod === "invite" && (
          <InfoBox
            title="Invitation Email"
            description={`An invitation will be sent to ${email} with instructions to set up their account.`}
            variant="success"
          />
        )}

        <StepActions nextLabel="Create User" />
      </StepContent>
    </StepWizard>
  )
}
