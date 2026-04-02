"use client"

import * as React from "react"
import { Check, ChevronRight, X } from "lucide-react"
import { cn } from "../../utils"
import { Button } from "./button"

export interface WizardStep {
  id: string
  title: string
  description?: string
  icon?: React.ReactNode
}

interface StepWizardContextValue {
  steps: WizardStep[]
  currentStepIndex: number
  currentStep: WizardStep
  goToStep: (index: number) => void
  nextStep: () => void
  prevStep: () => void
  isFirstStep: boolean
  isLastStep: boolean
  canGoBack: boolean
}

const StepWizardContext = React.createContext<StepWizardContextValue | null>(null)

export function useStepWizard() {
  const context = React.useContext(StepWizardContext)
  if (!context) {
    throw new Error("useStepWizard must be used within a StepWizard")
  }
  return context
}

interface StepWizardProps {
  steps: WizardStep[]
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  children: React.ReactNode
  onComplete?: () => void
  initialStep?: number
}

export function StepWizard({
  steps,
  open,
  onOpenChange,
  title,
  children,
  onComplete,
  initialStep = 0,
}: StepWizardProps) {
  const [currentStepIndex, setCurrentStepIndex] = React.useState(initialStep)

  React.useEffect(() => {
    if (open) {
      setCurrentStepIndex(initialStep)
    }
  }, [open, initialStep])

  const goToStep = React.useCallback((index: number) => {
    if (index >= 0 && index < steps.length) {
      setCurrentStepIndex(index)
    }
  }, [steps.length])

  const nextStep = React.useCallback(() => {
    if (currentStepIndex < steps.length - 1) {
      setCurrentStepIndex(prev => prev + 1)
    } else {
      onComplete?.()
    }
  }, [currentStepIndex, steps.length, onComplete])

  const prevStep = React.useCallback(() => {
    if (currentStepIndex > 0) {
      setCurrentStepIndex(prev => prev - 1)
    }
  }, [currentStepIndex])

  const contextValue: StepWizardContextValue = {
    steps,
    currentStepIndex,
    currentStep: steps[currentStepIndex],
    goToStep,
    nextStep,
    prevStep,
    isFirstStep: currentStepIndex === 0,
    isLastStep: currentStepIndex === steps.length - 1,
    canGoBack: currentStepIndex > 0,
  }

  if (!open) return null

  return (
    <StepWizardContext.Provider value={contextValue}>
      <div className="fixed inset-0 z-50">
        {/* Backdrop */}
        <div 
          className="absolute inset-0 bg-black/50 animate-in fade-in-0 duration-200"
          onClick={() => onOpenChange(false)}
        />
        
        {/* Panel */}
        <div className="absolute inset-y-0 right-0 w-full max-w-md bg-background border-l shadow-xl animate-in slide-in-from-right duration-300 flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between px-6 py-4 border-b">
            <h2 className="text-lg font-semibold">{title}</h2>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => onOpenChange(false)}
            >
              <X className="h-4 w-4" />
              <span className="sr-only">Close</span>
            </Button>
          </div>

          {/* Content */}
          <div className="flex-1 min-h-0 overflow-y-auto">
            <div className="p-6">
              {/* Stepper */}
              <div className="mb-8">
                {steps.map((step, index) => {
                  const isCompleted = index < currentStepIndex
                  const isCurrent = index === currentStepIndex
                  const isLast = index === steps.length - 1

                  return (
                    <div key={step.id} className="relative">
                      <div className="flex items-start gap-4">
                        {/* Step indicator */}
                        <div className="relative flex flex-col items-center">
                          <div
                            className={cn(
                              "flex h-8 w-8 items-center justify-center rounded-full border-2 transition-colors",
                              isCompleted && "bg-foreground border-foreground text-background",
                              isCurrent && "border-foreground bg-background text-foreground",
                              !isCompleted && !isCurrent && "border-muted-foreground/30 bg-background text-muted-foreground/50"
                            )}
                          >
                            {isCompleted ? (
                              <Check className="h-4 w-4" />
                            ) : isCurrent ? (
                              <ChevronRight className="h-4 w-4" />
                            ) : (
                              <span className="text-xs font-medium">{index + 1}</span>
                            )}
                          </div>
                          {/* Connector line */}
                          {!isLast && (
                            <div
                              className={cn(
                                "absolute top-8 w-0.5 h-8",
                                isCompleted ? "bg-foreground" : "bg-muted-foreground/30"
                              )}
                              style={{ left: "50%", transform: "translateX(-50%)" }}
                            />
                          )}
                        </div>

                        {/* Step content */}
                        <div className={cn("flex-1 pb-8", isLast && "pb-0")}>
                          <div
                            className={cn(
                              "font-medium transition-colors",
                              isCurrent ? "text-foreground" : "text-muted-foreground"
                            )}
                          >
                            {step.title}
                          </div>
                          {step.description && isCurrent && (
                            <p className="text-sm text-muted-foreground mt-0.5">
                              {step.description}
                            </p>
                          )}
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>

              {/* Step content */}
              <div className="space-y-6 pb-6">
                {children}
              </div>
            </div>
          </div>
        </div>
      </div>
    </StepWizardContext.Provider>
  )
}

interface StepContentProps {
  stepId: string
  children: React.ReactNode
}

export function StepContent({ stepId, children }: StepContentProps) {
  const { currentStep } = useStepWizard()
  
  if (currentStep.id !== stepId) return null
  
  return <>{children}</>
}

interface StepActionsProps {
  onNext?: () => void | Promise<void>
  onBack?: () => void
  nextLabel?: string
  backLabel?: string
  nextDisabled?: boolean
  showBack?: boolean
  isLoading?: boolean
}

export function StepActions({
  onNext,
  onBack,
  nextLabel,
  backLabel = "Back",
  nextDisabled = false,
  showBack = true,
  isLoading = false,
}: StepActionsProps) {
  const { isFirstStep, isLastStep, nextStep, prevStep } = useStepWizard()

  const handleNext = async () => {
    if (onNext) {
      await onNext()
    } else {
      nextStep()
    }
  }

  const handleBack = () => {
    if (onBack) {
      onBack()
    } else {
      prevStep()
    }
  }

  return (
    <div className="flex items-center justify-between pt-6 border-t mt-6">
      {showBack && !isFirstStep ? (
        <Button variant="outline" onClick={handleBack} disabled={isLoading}>
          {backLabel}
        </Button>
      ) : (
        <div />
      )}
      <Button onClick={handleNext} disabled={nextDisabled || isLoading}>
        {isLoading ? "Processing..." : nextLabel || (isLastStep ? "Complete" : "Continue")}
      </Button>
    </div>
  )
}

// Helper components for form sections
interface FormSectionProps {
  title?: string
  description?: string
  children: React.ReactNode
  className?: string
}

export function FormSection({ title, description, children, className }: FormSectionProps) {
  return (
    <div className={cn("space-y-4", className)}>
      {(title || description) && (
        <div>
          {title && <h3 className="font-medium">{title}</h3>}
          {description && <p className="text-sm text-muted-foreground">{description}</p>}
        </div>
      )}
      {children}
    </div>
  )
}

interface ParameterRowProps {
  label: string
  value: React.ReactNode
}

export function ParameterRow({ label, value }: ParameterRowProps) {
  return (
    <div className="flex items-center justify-between py-2.5 border-b last:border-0">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className="text-sm font-medium">{value}</span>
    </div>
  )
}

interface InfoBoxProps {
  title: string
  description?: string
  children?: React.ReactNode
  variant?: "default" | "success" | "warning"
  className?: string
}

export function InfoBox({ title, description, children, variant = "default", className }: InfoBoxProps) {
  return (
    <div className={cn(
      "rounded-lg border p-4",
      variant === "success" && "bg-green-50 border-green-200 dark:bg-green-950/20 dark:border-green-900",
      variant === "warning" && "bg-amber-50 border-amber-200 dark:bg-amber-950/20 dark:border-amber-900",
      variant === "default" && "bg-muted/50",
      className
    )}>
      <div className="font-medium text-sm">{title}</div>
      {description && <p className="text-sm text-muted-foreground mt-1">{description}</p>}
      {children && <div className="mt-3">{children}</div>}
    </div>
  )
}
