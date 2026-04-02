"use client"

import { useState, useTransition } from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import { Loader2, UserPlus, Eye, EyeOff, Check } from "lucide-react"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select"
import { Checkbox } from "../ui/checkbox"
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from "../ui/sheet"
import { createUser } from "../../api/create-user"
import { cn } from "../../utils"

interface AddUserSheetProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  organizations: any[]
}

const STEPS = [
  { id: "profile", label: "User Details", description: "Basic information" },
  { id: "email", label: "Email & Login", description: "Contact details" },
  { id: "password", label: "Password", description: "Initial credentials" },
  { id: "confirm", label: "Confirmation", description: "Review & create" },
] as const

type Step = typeof STEPS[number]["id"]

export function AddUserSheet({ open, onOpenChange, organizations }: AddUserSheetProps) {
  const router = useRouter()
  const [isPending, startTransition] = useTransition()
  const [error, setError] = useState<string | null>(null)
  const [currentStep, setCurrentStep] = useState<Step>("profile")
  const [showPassword, setShowPassword] = useState(false)

  // Form state
  const [organizationId, setOrganizationId] = useState(organizations[0]?.id ?? "")
  const [givenName, setGivenName] = useState("")
  const [familyName, setFamilyName] = useState("")
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [isEmailVerified, setIsEmailVerified] = useState(false)
  const [password, setPassword] = useState("")
  const [requirePasswordChange, setRequirePasswordChange] = useState(true)

  const currentStepIndex = STEPS.findIndex(s => s.id === currentStep)
  const selectedOrg = organizations.find((o: any) => o.id === organizationId)

  function resetForm() {
    setCurrentStep("profile")
    setOrganizationId(organizations[0]?.id ?? "")
    setGivenName("")
    setFamilyName("")
    setUsername("")
    setEmail("")
    setIsEmailVerified(false)
    setPassword("")
    setRequirePasswordChange(true)
    setError(null)
    setShowPassword(false)
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  function canContinue(): boolean {
    switch (currentStep) {
      case "profile":
        return !!(organizationId && givenName && familyName)
      case "email":
        return !!email
      case "password":
        return true // password is optional
      case "confirm":
        return true
      default:
        return false
    }
  }

  function handleNext() {
    const idx = currentStepIndex
    if (idx < STEPS.length - 1) {
      setCurrentStep(STEPS[idx + 1].id)
    }
  }

  function handleBack() {
    const idx = currentStepIndex
    if (idx > 0) {
      setCurrentStep(STEPS[idx - 1].id)
    }
  }

  function handleCreate() {
    setError(null)
    startTransition(async () => {
      try {
        const result = await createUser({
          organizationId,
          givenName,
          familyName,
          email,
          username: username || undefined,
          isEmailVerified,
          password: password || undefined,
          requirePasswordChange,
        })
        handleOpenChange(false)
        router.refresh()
      } catch (e) {
        setError(e instanceof Error ? e.message : "Failed to create user")
      }
    })
  }

  return (
    <Sheet open={open} onOpenChange={handleOpenChange}>
      <SheetContent side="right" className="sm:max-w-md w-full flex flex-col">
        <SheetHeader>
          <SheetTitle>Create User</SheetTitle>
          <SheetDescription className="sr-only">
            Add a new human user to your ZITADEL instance
          </SheetDescription>
        </SheetHeader>

        {/* Stepper */}
        <div className="px-4 pb-2">
          <div className="relative flex flex-col gap-0">
            {STEPS.map((step, i) => {
              const isActive = step.id === currentStep
              const isCompleted = i < currentStepIndex
              return (
                <div key={step.id} className="flex items-start gap-3">
                  {/* Step indicator line + circle */}
                  <div className="flex flex-col items-center">
                    <button
                      type="button"
                      onClick={() => {
                        if (isCompleted) setCurrentStep(step.id)
                      }}
                      disabled={!isCompleted && !isActive}
                      className={cn(
                        "flex h-7 w-7 items-center justify-center rounded-full border-2 text-xs font-medium transition-colors shrink-0",
                        isActive && "border-foreground bg-foreground text-background",
                        isCompleted && "border-foreground bg-transparent text-foreground cursor-pointer hover:bg-muted",
                        !isActive && !isCompleted && "border-muted-foreground/30 text-muted-foreground/50",
                      )}
                    >
                      {isCompleted ? <Check className="h-3.5 w-3.5" /> : i + 1}
                    </button>
                    {i < STEPS.length - 1 && (
                      <div className={cn(
                        "w-px h-4",
                        i < currentStepIndex ? "bg-foreground" : "bg-border",
                      )} />
                    )}
                  </div>
                  {/* Step label */}
                  <div className="pb-4">
                    <p className={cn(
                      "text-sm font-medium leading-7",
                      isActive && "text-foreground",
                      !isActive && !isCompleted && "text-muted-foreground",
                    )}>
                      {step.label}
                    </p>
                    {isActive && (
                      <p className="text-xs text-muted-foreground">{step.description}</p>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        </div>

        {/* Step Content */}
        <div className="flex-1 overflow-y-auto px-4 space-y-4">
          {/* Step 1: Profile */}
          {currentStep === "profile" && (
            <>
              <div className="space-y-2">
                <Label htmlFor="givenName">First Name *</Label>
                <Input
                  id="givenName"
                  placeholder="Jane"
                  value={givenName}
                  onChange={(e) => setGivenName(e.target.value)}
                  autoFocus
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="familyName">Last Name *</Label>
                <Input
                  id="familyName"
                  placeholder="Doe"
                  value={familyName}
                  onChange={(e) => setFamilyName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  placeholder="Optional — defaults to email"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                />
                <p className="text-xs text-muted-foreground">
                  If not set, the email will be used as the username
                </p>
              </div>
              <div className="rounded-lg border bg-muted/30 p-3 mt-4">
                <Label className="text-xs text-muted-foreground mb-1.5 block">Organization</Label>
                <Select value={organizationId} onValueChange={setOrganizationId}>
                  <SelectTrigger className="bg-background">
                    <SelectValue placeholder="Select organization" />
                  </SelectTrigger>
                  <SelectContent>
                    {organizations.map((org: any) => (
                      <SelectItem key={org.id} value={org.id}>
                        {org.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </>
          )}

          {/* Step 2: Email */}
          {currentStep === "email" && (
            <>
              <div className="space-y-2">
                <Label htmlFor="email">Email Address *</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="jane.doe@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  autoFocus
                />
              </div>
              <div className="flex items-center space-x-2 pt-2">
                <Checkbox
                  id="isEmailVerified"
                  checked={isEmailVerified}
                  onCheckedChange={(checked) => setIsEmailVerified(checked === true)}
                />
                <Label htmlFor="isEmailVerified" className="text-sm font-normal">
                  Mark email as already verified
                </Label>
              </div>
              {!isEmailVerified && (
                <p className="text-xs text-muted-foreground">
                  A verification email will be sent to the user
                </p>
              )}
            </>
          )}

          {/* Step 3: Password */}
          {currentStep === "password" && (
            <>
              <div className="space-y-2">
                <Label htmlFor="password">Initial Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Leave blank to skip"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="pr-10"
                    autoFocus
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOff className="h-4 w-4 text-muted-foreground" />
                    ) : (
                      <Eye className="h-4 w-4 text-muted-foreground" />
                    )}
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  If left blank, the user will set their password via email
                </p>
              </div>
              {password && (
                <div className="flex items-center space-x-2 pt-2">
                  <Checkbox
                    id="requirePasswordChange"
                    checked={requirePasswordChange}
                    onCheckedChange={(checked) => setRequirePasswordChange(checked === true)}
                  />
                  <Label htmlFor="requirePasswordChange" className="text-sm font-normal">
                    Require password change on first login
                  </Label>
                </div>
              )}
            </>
          )}

          {/* Step 4: Confirmation */}
          {currentStep === "confirm" && (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">
                Review the details below before creating the user.
              </p>
              <div className="rounded-lg border divide-y">
                <div className="p-3">
                  <p className="text-xs text-muted-foreground">Name</p>
                  <p className="text-sm font-medium">{givenName} {familyName}</p>
                </div>
                {username && (
                  <div className="p-3">
                    <p className="text-xs text-muted-foreground">Username</p>
                    <p className="text-sm font-medium">@{username}</p>
                  </div>
                )}
                <div className="p-3">
                  <p className="text-xs text-muted-foreground">Email</p>
                  <p className="text-sm font-medium">{email}</p>
                  <p className="text-xs text-muted-foreground mt-0.5">
                    {isEmailVerified ? "Pre-verified" : "Verification email will be sent"}
                  </p>
                </div>
                <div className="p-3">
                  <p className="text-xs text-muted-foreground">Organization</p>
                  <p className="text-sm font-medium">{selectedOrg?.name ?? organizationId}</p>
                </div>
                <div className="p-3">
                  <p className="text-xs text-muted-foreground">Password</p>
                  <p className="text-sm font-medium">
                    {password ? "Set (hidden)" : "Not set — user will set via email"}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Error */}
          {error && (
            <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-3">
              <p className="text-sm font-medium text-destructive">
                Failed to create user
              </p>
              <p className="text-xs text-muted-foreground mt-1">{error}</p>
            </div>
          )}
        </div>

        {/* Footer Actions */}
        <SheetFooter className="flex-row justify-between border-t pt-4">
          {currentStepIndex > 0 ? (
            <Button variant="outline" onClick={handleBack} disabled={isPending}>
              Back
            </Button>
          ) : (
            <div />
          )}
          {currentStep === "confirm" ? (
            <Button onClick={handleCreate} disabled={isPending}>
              {isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Creating...
                </>
              ) : (
                <>
                  <UserPlus className="mr-2 h-4 w-4" />
                  Create User
                </>
              )}
            </Button>
          ) : (
            <Button onClick={handleNext} disabled={!canContinue()}>
              Continue
            </Button>
          )}
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}
