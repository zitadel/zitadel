"use client"

import * as React from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { Cloud, Server, Globe, Shield, Zap } from "lucide-react"
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
import { cn } from "../../utils"

const steps: WizardStep[] = [
  { id: "type", title: "Instance Type", description: "Choose your deployment method" },
  { id: "configuration", title: "Configuration", description: "Set up your instance" },
  { id: "plan", title: "Plan Selection", description: "Choose your plan" },
  { id: "confirmation", title: "Confirmation", description: "Review and create" },
]

const regions = [
  { id: "eu-frankfurt", name: "EU (Frankfurt)", flag: "DE" },
  { id: "us-virginia", name: "US (Virginia)", flag: "US" },
  { id: "us-oregon", name: "US (Oregon)", flag: "US" },
  { id: "asia-singapore", name: "Asia (Singapore)", flag: "SG" },
  { id: "asia-tokyo", name: "Asia (Tokyo)", flag: "JP" },
  { id: "au-sydney", name: "Australia (Sydney)", flag: "AU" },
]

const plans = [
  {
    id: "free",
    name: "Free",
    price: "$0",
    description: "For personal projects and testing",
    features: ["1,000 monthly active users", "Basic authentication", "Community support"],
  },
  {
    id: "pro",
    name: "Pro",
    price: "$99",
    description: "For growing teams and businesses",
    features: ["10,000 monthly active users", "Advanced authentication", "Priority support", "Custom branding"],
  },
  {
    id: "enterprise",
    name: "Enterprise",
    price: "Custom",
    description: "For large organizations",
    features: ["Unlimited users", "Enterprise SSO", "Dedicated support", "SLA guarantees", "Custom integrations"],
  },
]

interface CreateInstanceWizardProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateInstanceWizard({ open, onOpenChange }: CreateInstanceWizardProps) {
  const router = useRouter()
  const [instanceType, setInstanceType] = React.useState<"cloud" | "self-hosted">("cloud")
  const [instanceName, setInstanceName] = React.useState("")
  const [subdomain, setSubdomain] = React.useState("")
  const [region, setRegion] = React.useState("eu-frankfurt")
  const [plan, setPlan] = React.useState("free")
  const [acceptTerms, setAcceptTerms] = React.useState(false)

  const handleComplete = () => {
    // In a real app, this would create the instance via API
    onOpenChange(false)
    router.push("/overview")
  }

  const selectedPlan = plans.find(p => p.id === plan)
  const selectedRegion = regions.find(r => r.id === region)

  return (
    <StepWizard
      steps={steps}
      open={open}
      onOpenChange={onOpenChange}
      title="Create Instance"
      onComplete={handleComplete}
    >
      {/* Step 1: Instance Type */}
      <StepContent stepId="type">
        <FormSection
          title="Choose Deployment Type"
          description="Select how you want to deploy your ZITADEL instance"
        >
          <RadioGroup value={instanceType} onValueChange={(v) => setInstanceType(v as typeof instanceType)}>
            <label
              className={cn(
                "flex items-start gap-4 p-4 rounded-lg border cursor-pointer transition-colors",
                instanceType === "cloud" ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="cloud" id="cloud" className="mt-1" />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <Cloud className="h-5 w-5" />
                  <span className="font-medium">Cloud Hosted</span>
                </div>
                <p className="text-sm text-muted-foreground mt-1">
                  Fully managed by ZITADEL. No infrastructure to maintain.
                </p>
                <div className="flex flex-wrap gap-2 mt-3">
                  <span className="inline-flex items-center gap-1 text-xs bg-muted px-2 py-1 rounded">
                    <Zap className="h-3 w-3" /> Instant setup
                  </span>
                  <span className="inline-flex items-center gap-1 text-xs bg-muted px-2 py-1 rounded">
                    <Shield className="h-3 w-3" /> Auto-updates
                  </span>
                  <span className="inline-flex items-center gap-1 text-xs bg-muted px-2 py-1 rounded">
                    <Globe className="h-3 w-3" /> Global CDN
                  </span>
                </div>
              </div>
            </label>

            <label
              className={cn(
                "flex items-start gap-4 p-4 rounded-lg border cursor-pointer transition-colors",
                instanceType === "self-hosted" ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
              )}
            >
              <RadioGroupItem value="self-hosted" id="self-hosted" className="mt-1" />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <Server className="h-5 w-5" />
                  <span className="font-medium">Self-Hosted</span>
                </div>
                <p className="text-sm text-muted-foreground mt-1">
                  Deploy on your own infrastructure with full control.
                </p>
                <div className="flex flex-wrap gap-2 mt-3">
                  <span className="inline-flex items-center gap-1 text-xs bg-muted px-2 py-1 rounded">
                    Full control
                  </span>
                  <span className="inline-flex items-center gap-1 text-xs bg-muted px-2 py-1 rounded">
                    Data sovereignty
                  </span>
                  <span className="inline-flex items-center gap-1 text-xs bg-muted px-2 py-1 rounded">
                    Custom deployment
                  </span>
                </div>
              </div>
            </label>
          </RadioGroup>
        </FormSection>
        <StepActions />
      </StepContent>

      {/* Step 2: Configuration */}
      <StepContent stepId="configuration">
        <FormSection title="Instance Details">
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Instance Name</Label>
              <Input
                id="name"
                placeholder="My Production Instance"
                value={instanceName}
                onChange={(e) => setInstanceName(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">
                A friendly name to identify your instance
              </p>
            </div>

            {instanceType === "cloud" && (
              <div className="space-y-2">
                <Label htmlFor="subdomain">Subdomain</Label>
                <div className="flex items-center">
                  <Input
                    id="subdomain"
                    placeholder="my-company"
                    value={subdomain}
                    onChange={(e) => setSubdomain(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ""))}
                    className="rounded-r-none"
                  />
                  <span className="inline-flex items-center px-3 h-9 border border-l-0 rounded-r-md bg-muted text-sm text-muted-foreground">
                    .zitadel.cloud
                  </span>
                </div>
                <p className="text-xs text-muted-foreground">
                  Your instance will be available at {subdomain || "my-company"}.zitadel.cloud
                </p>
              </div>
            )}
          </div>
        </FormSection>

        {instanceType === "cloud" && (
          <FormSection title="Region" description="Select the region closest to your users" className="mt-6">
            <RadioGroup value={region} onValueChange={setRegion} className="grid grid-cols-2 gap-2">
              {regions.map((r) => (
                <label
                  key={r.id}
                  className={cn(
                    "flex items-center gap-2 p-3 rounded-lg border cursor-pointer transition-colors text-sm",
                    region === r.id ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
                  )}
                >
                  <RadioGroupItem value={r.id} id={r.id} />
                  <span>{r.name}</span>
                </label>
              ))}
            </RadioGroup>
          </FormSection>
        )}

        <StepActions nextDisabled={!instanceName || (instanceType === "cloud" && !subdomain)} />
      </StepContent>

      {/* Step 3: Plan Selection */}
      <StepContent stepId="plan">
        <FormSection title="Select Plan" description="Choose the plan that fits your needs">
          <RadioGroup value={plan} onValueChange={setPlan} className="space-y-3">
            {plans.map((p) => (
              <label
                key={p.id}
                className={cn(
                  "flex flex-col p-4 rounded-lg border cursor-pointer transition-colors",
                  plan === p.id ? "border-foreground bg-muted/50" : "border-border hover:bg-muted/30"
                )}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <RadioGroupItem value={p.id} id={p.id} />
                    <span className="font-medium">{p.name}</span>
                  </div>
                  <span className="font-semibold">{p.price}{p.price !== "Custom" && "/mo"}</span>
                </div>
                <p className="text-sm text-muted-foreground mt-2 ml-7">{p.description}</p>
                <ul className="mt-3 ml-7 space-y-1">
                  {p.features.map((feature, i) => (
                    <li key={i} className="text-xs text-muted-foreground flex items-center gap-1.5">
                      <span className="h-1 w-1 rounded-full bg-muted-foreground" />
                      {feature}
                    </li>
                  ))}
                </ul>
              </label>
            ))}
          </RadioGroup>
        </FormSection>
        <StepActions />
      </StepContent>

      {/* Step 4: Confirmation */}
      <StepContent stepId="confirmation">
        <FormSection title="Review Configuration">
          <div className="rounded-lg border divide-y">
            <ParameterRow label="Instance Name" value={instanceName || "—"} />
            <ParameterRow label="Type" value={instanceType === "cloud" ? "Cloud Hosted" : "Self-Hosted"} />
            {instanceType === "cloud" && (
              <>
                <ParameterRow label="Domain" value={`${subdomain || "—"}.zitadel.cloud`} />
                <ParameterRow label="Region" value={selectedRegion?.name || "—"} />
              </>
            )}
            <ParameterRow label="Plan" value={selectedPlan?.name || "—"} />
            {selectedPlan?.price !== "Custom" && (
              <ParameterRow label="Monthly Cost" value={selectedPlan?.price || "—"} />
            )}
          </div>
        </FormSection>

        <InfoBox
          title={selectedPlan?.name || "Free"}
          description={selectedPlan?.description}
          variant="success"
        >
          <ul className="space-y-1">
            {selectedPlan?.features.map((feature, i) => (
              <li key={i} className="text-xs flex items-center gap-1.5">
                <span className="h-1 w-1 rounded-full bg-green-600" />
                {feature}
              </li>
            ))}
          </ul>
        </InfoBox>

        <div className="flex items-start gap-2">
          <Checkbox
            id="terms"
            checked={acceptTerms}
            onCheckedChange={(checked) => setAcceptTerms(checked === true)}
          />
          <label htmlFor="terms" className="text-sm text-muted-foreground leading-relaxed cursor-pointer">
            I agree to the{" "}
            <a href="#" className="text-foreground underline">Terms of Service</a>
            {" "}and{" "}
            <a href="#" className="text-foreground underline">Privacy Policy</a>
          </label>
        </div>

        <StepActions nextLabel="Create Instance" nextDisabled={!acceptTerms} />
      </StepContent>
    </StepWizard>
  )
}
