"use client"

import * as React from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { Globe, Smartphone, Server, Puzzle, FolderKanban, Copy, Check } from "lucide-react"
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
import { Button } from "../ui/button"
import { cn } from "../../utils"
import { useAppContext } from "../../context/app-context"
import { projects } from "../../mock-data"

const steps: WizardStep[] = [
  { id: "type", title: "Application Type", description: "Choose your app type" },
  { id: "details", title: "Application Details", description: "Configure your app" },
  { id: "auth", title: "Authentication", description: "Set up auth flow" },
  { id: "confirmation", title: "Confirmation", description: "Review and create" },
]

const appTypes = [
  {
    id: "web",
    name: "Web Application",
    description: "Single-page app or traditional web app",
    icon: Globe,
    authMethods: ["PKCE", "Code Flow"],
  },
  {
    id: "native",
    name: "Native / Mobile",
    description: "iOS, Android, or desktop application",
    icon: Smartphone,
    authMethods: ["PKCE"],
  },
  {
    id: "api",
    name: "API / Machine-to-Machine",
    description: "Backend service or API client",
    icon: Server,
    authMethods: ["Client Credentials"],
  },
  {
    id: "extension",
    name: "Browser Extension",
    description: "Chrome, Firefox, or other browser extension",
    icon: Puzzle,
    authMethods: ["PKCE"],
  },
]

const authFlows = [
  { id: "pkce", name: "Authorization Code with PKCE", recommended: true, description: "Most secure for public clients" },
  { id: "code", name: "Authorization Code", description: "Traditional OAuth flow for confidential clients" },
  { id: "implicit", name: "Implicit Flow (Legacy)", description: "Not recommended for new applications" },
  { id: "client_credentials", name: "Client Credentials", description: "For machine-to-machine communication" },
]

interface CreateApplicationWizardProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateApplicationWizard({ open, onOpenChange }: CreateApplicationWizardProps) {
  const router = useRouter()
  const { currentInstance, currentOrganization } = useAppContext()
  
  const [appType, setAppType] = React.useState("web")
  const [appName, setAppName] = React.useState("")
  const [selectedProjectId, setSelectedProjectId] = React.useState("")
  const [redirectUris, setRedirectUris] = React.useState("")
  const [postLogoutUris, setPostLogoutUris] = React.useState("")
  const [authFlow, setAuthFlow] = React.useState("pkce")
  const [devMode, setDevMode] = React.useState(true)
  const [copied, setCopied] = React.useState(false)

  // Filter projects by current organization
  const availableProjects = React.useMemo(() => {
    if (!currentOrganization) return projects.slice(0, 10)
    return projects.filter(p => p.orgId === currentOrganization.id).slice(0, 20)
  }, [currentOrganization])

  const handleComplete = () => {
    onOpenChange(false)
    if (currentOrganization) {
      router.push("/org/applications")
    } else {
      router.push("/applications")
    }
  }

  const selectedAppType = appTypes.find(t => t.id === appType)
  const selectedProject = projects.find(p => p.id === selectedProjectId)
  const selectedAuthFlow = authFlows.find(f => f.id === authFlow)
  
  // Generate a mock client ID
  const clientId = React.useMemo(() => {
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
    let result = ""
    for (let i = 0; i < 24; i++) {
      result += chars[Math.floor(Math.random() * chars.length)]
    }
    return result
  }, [])

  const copyClientId = () => {
    navigator.clipboard.writeText(clientId)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <StepWizard
      steps={steps}
      open={open}
      onOpenChange={onOpenChange}
      title="Create Application"
      onComplete={handleComplete}
    >
      {/* Step 1: Application Type */}
      <StepContent stepId="type">
        <FormSection
          title="Select Application Type"
          description="Choose the type that best matches your application"
        >
          <RadioGroup value={appType} onValueChange={setAppType} className="space-y-3">
            {appTypes.map((type) => {
              const Icon = type.icon
              return (
                <label
                  key={type.id}
                  className={cn(
                    "flex items-start gap-4 p-4 rounded-lg border cursor-pointer transition-colors",
                    appType === type.id ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
                  )}
                >
                  <RadioGroupItem value={type.id} id={type.id} className="mt-1" />
                  <div className="flex items-start gap-3 flex-1">
                    <div className="h-10 w-10 rounded-lg bg-muted flex items-center justify-center">
                      <Icon className="h-5 w-5 text-muted-foreground" />
                    </div>
                    <div className="flex-1">
                      <span className="font-medium">{type.name}</span>
                      <p className="text-sm text-muted-foreground mt-0.5">{type.description}</p>
                      <div className="flex flex-wrap gap-1 mt-2">
                        {type.authMethods.map((method) => (
                          <span key={method} className="text-xs bg-muted px-2 py-0.5 rounded">
                            {method}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>
                </label>
              )
            })}
          </RadioGroup>
        </FormSection>
        <StepActions />
      </StepContent>

      {/* Step 2: Application Details */}
      <StepContent stepId="details">
        <FormSection
          title="Application Information"
          description="Enter basic details about your application"
        >
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="appName">Application Name</Label>
              <Input
                id="appName"
                placeholder="My Web App"
                value={appName}
                onChange={(e) => setAppName(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="project">Project</Label>
              <Select value={selectedProjectId} onValueChange={setSelectedProjectId}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a project" />
                </SelectTrigger>
                <SelectContent>
                  {availableProjects.map((project) => (
                    <SelectItem key={project.id} value={project.id}>
                      <div className="flex items-center gap-2">
                        <FolderKanban className="h-4 w-4 text-muted-foreground" />
                        {project.name}
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground">
                Applications inherit roles and permissions from their project
              </p>
            </div>
          </div>
        </FormSection>

        {(appType === "web" || appType === "extension") && (
          <FormSection title="Redirect URIs" className="mt-6">
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="redirectUris">Redirect URIs</Label>
                <Input
                  id="redirectUris"
                  placeholder="https://myapp.com/callback"
                  value={redirectUris}
                  onChange={(e) => setRedirectUris(e.target.value)}
                />
                <p className="text-xs text-muted-foreground">
                  Where to redirect after authentication (one per line)
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="postLogoutUris">Post-Logout Redirect URIs (Optional)</Label>
                <Input
                  id="postLogoutUris"
                  placeholder="https://myapp.com"
                  value={postLogoutUris}
                  onChange={(e) => setPostLogoutUris(e.target.value)}
                />
              </div>
            </div>
          </FormSection>
        )}

        <FormSection className="mt-6">
          <div className="flex items-start gap-2">
            <Checkbox
              id="devMode"
              checked={devMode}
              onCheckedChange={(checked) => setDevMode(checked === true)}
            />
            <label htmlFor="devMode" className="text-sm cursor-pointer leading-relaxed">
              <span className="font-medium">Enable development mode</span>
              <p className="text-xs text-muted-foreground">
                Allows localhost redirects and relaxed security for testing
              </p>
            </label>
          </div>
        </FormSection>

        <StepActions nextDisabled={!appName || !selectedProjectId} />
      </StepContent>

      {/* Step 3: Authentication */}
      <StepContent stepId="auth">
        <FormSection
          title="Authentication Flow"
          description="Select the OAuth flow for your application"
        >
          <RadioGroup value={authFlow} onValueChange={setAuthFlow} className="space-y-2">
            {authFlows
              .filter(flow => {
                if (appType === "api") return flow.id === "client_credentials"
                if (appType === "native" || appType === "extension") return flow.id === "pkce"
                return flow.id !== "client_credentials"
              })
              .map((flow) => (
                <label
                  key={flow.id}
                  className={cn(
                    "flex items-start gap-3 p-4 rounded-lg border cursor-pointer transition-colors",
                    authFlow === flow.id ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
                  )}
                >
                  <RadioGroupItem value={flow.id} id={flow.id} className="mt-0.5" />
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-sm">{flow.name}</span>
                      {flow.recommended && (
                        <span className="text-xs bg-foreground text-background px-1.5 py-0.5 rounded">
                          Recommended
                        </span>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">{flow.description}</p>
                  </div>
                </label>
              ))}
          </RadioGroup>
        </FormSection>

        <InfoBox
          title={selectedAppType?.name || "Application"}
          variant="default"
        >
          <p className="text-xs text-muted-foreground mt-1">
            {appType === "web" && "Web applications should use PKCE for public clients or Authorization Code for confidential clients with a backend."}
            {appType === "native" && "Native applications must use PKCE to securely authenticate without storing secrets."}
            {appType === "api" && "API clients use Client Credentials for server-to-server authentication."}
            {appType === "extension" && "Browser extensions should use PKCE for secure authentication."}
          </p>
        </InfoBox>

        <StepActions />
      </StepContent>

      {/* Step 4: Confirmation */}
      <StepContent stepId="confirmation">
        <FormSection title="Review Application">
          <div className="rounded-lg border divide-y">
            <ParameterRow label="Name" value={appName || "—"} />
            <ParameterRow label="Type" value={selectedAppType?.name || "—"} />
            <ParameterRow label="Project" value={selectedProject?.name || "—"} />
            <ParameterRow label="Auth Flow" value={selectedAuthFlow?.name || "—"} />
            {redirectUris && <ParameterRow label="Redirect URI" value={redirectUris} />}
            <ParameterRow label="Dev Mode" value={devMode ? "Enabled" : "Disabled"} />
          </div>
        </FormSection>

        <FormSection title="Client Credentials" className="mt-4">
          <div className="rounded-lg border p-4 bg-muted/30">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs text-muted-foreground">Client ID</p>
                <code className="text-sm font-mono">{clientId}</code>
              </div>
              <Button variant="ghost" size="sm" onClick={copyClientId}>
                {copied ? (
                  <Check className="h-4 w-4 text-green-600" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>
            {appType === "api" && (
              <div className="mt-3 pt-3 border-t">
                <p className="text-xs text-muted-foreground">Client Secret</p>
                <code className="text-sm font-mono">••••••••••••••••••••••••</code>
                <p className="text-xs text-muted-foreground mt-1">
                  Secret will be shown once after creation
                </p>
              </div>
            )}
          </div>
        </FormSection>

        <InfoBox
          title="Ready to Create"
          description="Your application will be created and ready for integration."
          variant="success"
        />

        <StepActions nextLabel="Create Application" />
      </StepContent>
    </StepWizard>
  )
}
