"use client"

import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@zitadel/react/components/ui/card"
import { Button } from "@zitadel/react/components/ui/button"
import { Input } from "@zitadel/react/components/ui/input"
import { Textarea } from "@zitadel/react/components/ui/textarea"
import { Label } from "@zitadel/react/components/ui/label"
import { Badge } from "@zitadel/react/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@zitadel/react/components/ui/select"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@zitadel/react/components/ui/alert-dialog"
import {
  LifeBuoy,
  BookOpen,
  MessageCircle,
  Bug,
  ExternalLink,
  ArrowRight,
  Send,
  AlertTriangle,
  Clock,
  CheckCircle2,
  Loader2,
  Ticket,
} from "lucide-react"

const TOPICS = [
  { value: "authentication", label: "Authentication & Login" },
  { value: "authorization", label: "Authorization & Roles" },
  { value: "api", label: "API & Integration" },
  { value: "console", label: "Console & UI" },
  { value: "performance", label: "Performance & Availability" },
  { value: "billing", label: "Billing & Subscription" },
  { value: "data", label: "Data & Migration" },
  { value: "security", label: "Security & Compliance" },
  { value: "feature", label: "Feature Request" },
  { value: "other", label: "Other" },
] as const

const PRIORITIES = [
  {
    value: "low",
    label: "Low",
    description: "General question or minor issue",
    icon: Clock,
    color: "text-muted-foreground",
    badgeVariant: "secondary" as const,
  },
  {
    value: "medium",
    label: "Medium",
    description: "Issue affecting workflow but has workaround",
    icon: AlertTriangle,
    color: "text-amber-500",
    badgeVariant: "outline" as const,
  },
  {
    value: "high",
    label: "High",
    description: "Significant impact, no workaround available",
    icon: AlertTriangle,
    color: "text-orange-500",
    badgeVariant: "outline" as const,
  },
  {
    value: "urgent",
    label: "Urgent",
    description: "System outage or critical functionality loss",
    icon: AlertTriangle,
    color: "text-destructive",
    badgeVariant: "destructive" as const,
  },
] as const

// Mock instances — in production, fetched from ZITADEL_INSTANCES env
const MOCK_INSTANCES = [
  { name: "Production", url: "https://auth.example.com" },
  { name: "Staging", url: "https://auth-staging.example.com" },
  { name: "Development", url: "https://auth-dev.example.com" },
]

// Mock previous tickets
const MOCK_TICKETS = [
  {
    id: "TKT-1024",
    subject: "OIDC token refresh failing intermittently",
    topic: "authentication",
    priority: "high",
    status: "open",
    instance: "https://auth.example.com",
    createdAt: "2026-03-15T10:30:00Z",
  },
  {
    id: "TKT-1021",
    subject: "Need to increase rate limits for API",
    topic: "api",
    priority: "medium",
    status: "closed",
    instance: "https://auth.example.com",
    createdAt: "2026-03-10T14:20:00Z",
  },
]

function getStatusBadge(status: string) {
  switch (status) {
    case "open":
      return <Badge className="bg-blue-500/10 text-blue-500 border-blue-500/20 hover:bg-blue-500/20">Open</Badge>
    case "in_progress":
      return <Badge className="bg-amber-500/10 text-amber-500 border-amber-500/20 hover:bg-amber-500/20">In Progress</Badge>
    case "closed":
      return <Badge variant="secondary">Closed</Badge>
    default:
      return <Badge variant="outline">{status}</Badge>
  }
}

function getPriorityBadge(priority: string) {
  const p = PRIORITIES.find((pr) => pr.value === priority)
  if (!p) return <Badge variant="outline">{priority}</Badge>
  return <Badge variant={p.badgeVariant}>{p.label}</Badge>
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  })
}

export default function SupportPage() {
  const [subject, setSubject] = useState("")
  const [description, setDescription] = useState("")
  const [topic, setTopic] = useState("")
  const [priority, setPriority] = useState("")
  const [instance, setInstance] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [showSuccess, setShowSuccess] = useState(false)
  const [ticketFilter, setTicketFilter] = useState<"all" | "open" | "closed">("all")

  const canSubmit = subject && description && topic && priority && instance

  const filteredTickets = MOCK_TICKETS.filter((t) => {
    if (ticketFilter === "all") return true
    return t.status === ticketFilter
  })

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!canSubmit) return

    setIsSubmitting(true)
    // Mock submission — will connect to HubSpot
    setTimeout(() => {
      setIsSubmitting(false)
      setShowSuccess(true)
      setSubject("")
      setDescription("")
      setTopic("")
      setPriority("")
      setInstance("")
    }, 1500)
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Support</h1>
        <p className="text-muted-foreground">
          Get help with ZITADEL and manage your support requests
        </p>
      </div>

      {/* Support Requests */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-semibold tracking-tight">Support Requests</h2>
          <div className="flex items-center gap-2">
            <div className="flex rounded-lg border p-0.5">
              {(["all", "open", "closed"] as const).map((f) => (
                <button
                  key={f}
                  onClick={() => setTicketFilter(f)}
                  className={`px-3 py-1 text-xs font-medium rounded-md transition-colors ${
                    ticketFilter === f
                      ? "bg-primary text-primary-foreground"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  {f.charAt(0).toUpperCase() + f.slice(1)}
                </button>
              ))}
            </div>
          </div>
        </div>

        {filteredTickets.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Ticket className="h-10 w-10 text-muted-foreground/40 mx-auto mb-3" />
              <p className="text-sm text-muted-foreground">
                No {ticketFilter === "all" ? "" : ticketFilter} tickets found.
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-2">
            {filteredTickets.map((ticket) => (
              <Card
                key={ticket.id}
                className="hover:bg-muted/30 transition-colors cursor-pointer"
              >
                <CardContent className="p-4">
                  <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-3 min-w-0">
                      <span className="text-xs font-mono text-muted-foreground shrink-0">
                        {ticket.id}
                      </span>
                      <p className="font-medium truncate">{ticket.subject}</p>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      {getPriorityBadge(ticket.priority)}
                      {getStatusBadge(ticket.status)}
                      <span className="text-xs text-muted-foreground">
                        {formatDate(ticket.createdAt)}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      <hr />

      {/* New Ticket Form */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <LifeBuoy className="h-5 w-5" />
            New Support Ticket
          </CardTitle>
          <CardDescription>
            Submit a request and our team will get back to you. Fields marked * are required.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Row 1: Topic + Priority */}
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="topic">Topic *</Label>
                <Select value={topic} onValueChange={setTopic}>
                  <SelectTrigger id="topic">
                    <SelectValue placeholder="Select a topic" />
                  </SelectTrigger>
                  <SelectContent>
                    {TOPICS.map((t) => (
                      <SelectItem key={t.value} value={t.value}>
                        {t.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="priority">Priority *</Label>
                <Select value={priority} onValueChange={setPriority}>
                  <SelectTrigger id="priority">
                    <SelectValue placeholder="Select priority" />
                  </SelectTrigger>
                  <SelectContent>
                    {PRIORITIES.map((p) => (
                      <SelectItem key={p.value} value={p.value}>
                        <div className="flex items-center gap-2">
                          <p.icon className={`h-3.5 w-3.5 ${p.color}`} />
                          <span>{p.label}</span>
                          <span className="text-muted-foreground text-xs">
                            — {p.description}
                          </span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Priority warning for urgent */}
            {priority === "urgent" && (
              <div className="rounded-lg border border-destructive/50 bg-destructive/5 p-4 flex items-start gap-3">
                <AlertTriangle className="h-5 w-5 text-destructive shrink-0 mt-0.5" />
                <div className="text-sm">
                  <p className="font-medium text-destructive">Priority Policy</p>
                  <p className="text-muted-foreground mt-1">
                    Urgent tickets are reserved for system outages, critical functionality loss,
                    or issues causing major business impact. Urgent tickets immediately page our
                    on-call team — please use this status thoughtfully.
                  </p>
                </div>
              </div>
            )}

            {/* Row 2: Affected Instance */}
            <div className="space-y-2">
              <Label htmlFor="instance">Affected Instance *</Label>
              <Select value={instance} onValueChange={setInstance}>
                <SelectTrigger id="instance">
                  <SelectValue placeholder="Select the affected instance" />
                </SelectTrigger>
                <SelectContent>
                  {MOCK_INSTANCES.map((inst) => (
                    <SelectItem key={inst.url} value={inst.url}>
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{inst.name}</span>
                        <span className="text-muted-foreground text-xs">{inst.url}</span>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Row 3: Subject */}
            <div className="space-y-2">
              <Label htmlFor="subject">Subject *</Label>
              <Input
                id="subject"
                placeholder="Brief description of your issue"
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
              />
            </div>

            {/* Row 4: Description */}
            <div className="space-y-2">
              <Label htmlFor="description">Description *</Label>
              <Textarea
                id="description"
                placeholder="Please describe your issue in detail. Include steps to reproduce, expected behavior, and any error messages."
                className="min-h-[160px]"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>

            {/* Submit */}
            <div className="flex items-center justify-between pt-2">
              <p className="text-xs text-muted-foreground">
                Tickets are sent to our support team via HubSpot
              </p>
              <Button type="submit" disabled={!canSubmit || isSubmitting}>
                {isSubmitting ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Submitting...
                  </>
                ) : (
                  <>
                    <Send className="mr-2 h-4 w-4" />
                    Submit Ticket
                  </>
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Support Resources */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card className="group hover:border-primary/30 transition-colors">
          <CardContent className="p-5">
            <div className="flex items-start justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <BookOpen className="h-4 w-4 text-primary" />
                  <h3 className="font-semibold">Documentation</h3>
                </div>
                <p className="text-sm text-muted-foreground">
                  Comprehensive guides, API references, and tutorials.
                </p>
              </div>
              <a
                href="https://zitadel.com/docs"
                target="_blank"
                rel="noopener noreferrer"
                className="shrink-0"
              >
                <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              </a>
            </div>
          </CardContent>
        </Card>

        <Card className="group hover:border-primary/30 transition-colors">
          <CardContent className="p-5">
            <div className="flex items-start justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <MessageCircle className="h-4 w-4 text-primary" />
                  <h3 className="font-semibold">Community</h3>
                </div>
                <p className="text-sm text-muted-foreground">
                  Connect with other developers on Discord.
                </p>
              </div>
              <a
                href="https://discord.gg/zitadel"
                target="_blank"
                rel="noopener noreferrer"
                className="shrink-0"
              >
                <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              </a>
            </div>
          </CardContent>
        </Card>

        <Card className="group hover:border-primary/30 transition-colors">
          <CardContent className="p-5">
            <div className="flex items-start justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Bug className="h-4 w-4 text-primary" />
                  <h3 className="font-semibold">Bug Reports</h3>
                </div>
                <p className="text-sm text-muted-foreground">
                  Report bugs or request features on GitHub.
                </p>
              </div>
              <a
                href="https://github.com/zitadel/zitadel/issues"
                target="_blank"
                rel="noopener noreferrer"
                className="shrink-0"
              >
                <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              </a>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Success Dialog */}
      <AlertDialog open={showSuccess} onOpenChange={setShowSuccess}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <CheckCircle2 className="h-5 w-5 text-green-500" />
              Ticket Submitted
            </AlertDialogTitle>
            <AlertDialogDescription>
              Your support ticket has been created. Our team will review it and
              respond based on your selected priority. You&apos;ll receive updates
              via email.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction>OK</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
