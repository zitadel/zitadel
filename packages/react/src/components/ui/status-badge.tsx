import { Badge } from "./badge"
import { cn } from "../../utils"

type StatusVariant = "active" | "inactive" | "destructive" | "warning"

const variantStyles: Record<StatusVariant, string> = {
  active:
    "bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/50 dark:text-emerald-400 dark:border-emerald-800",
  inactive:
    "bg-muted text-muted-foreground border-border",
  destructive:
    "bg-destructive/10 text-destructive border-destructive/20",
  warning:
    "bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/50 dark:text-amber-400 dark:border-amber-800",
}

interface StatusBadgeProps {
  variant: StatusVariant
  children: React.ReactNode
  className?: string
}

export function StatusBadge({ variant, children, className }: StatusBadgeProps) {
  return (
    <Badge variant="outline" className={cn(variantStyles[variant], "text-xs", className)}>
      {children}
    </Badge>
  )
}
