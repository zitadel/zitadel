"use client"

import { ConsoleLink as Link } from "../../../context/link-context"
import { useConsoleRouter } from "../../../hooks/use-console-router"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../components/ui/card"
import { Badge } from "../../../components/ui/badge"
import { StatusBadge } from "../../../components/ui/status-badge"
import { Button } from "../../../components/ui/button"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../../components/ui/tabs"
import { ArrowLeft, Trash2, FolderKanban, AppWindow, ChevronRight } from "lucide-react"

interface ProjectDetailClientProps {
  project: any
  projectId: string
  applications: any[]
  error: string | null
}

function getProjectState(project: any): { label: string; variant: "active" | "inactive" | "destructive" | "warning" } {
  const state = project?.state ?? "PROJECT_STATE_UNSPECIFIED"
  const labels: Record<string, { label: string; variant: "active" | "inactive" | "destructive" | "warning" }> = {
    PROJECT_STATE_ACTIVE: { label: "Active", variant: "active" },
    PROJECT_STATE_INACTIVE: { label: "Inactive", variant: "inactive" },
    PROJECT_STATE_UNSPECIFIED: { label: "Unknown", variant: "inactive" },
  }
  return labels[state] ?? { label: state, variant: "inactive" }
}

function getAppType(app: any) {
  if (app.oidcConfig) return { label: "OIDC", className: "bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-400" }
  if (app.apiConfig) return { label: "API", className: "bg-purple-100 text-purple-700 dark:bg-purple-950 dark:text-purple-400" }
  if (app.samlConfig) return { label: "SAML", className: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-400" }
  return { label: "Unknown", className: "bg-muted text-muted-foreground" }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleDateString()
}

export function ProjectDetailClient({ project, projectId, applications, error }: ProjectDetailClientProps) {
  const router = useConsoleRouter()

  if (error || !project) {
    return (
      <div className="flex flex-col items-center justify-center h-[50vh] space-y-4">
        <h1 className="text-2xl font-bold">
          {error ? "Failed to load project" : "Project not found"}
        </h1>
        {error && <p className="text-sm text-muted-foreground">{error}</p>}
        <Button asChild>
          <Link href="/projects">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Projects
          </Link>
        </Button>
      </div>
    )
  }

  const stateInfo = getProjectState(project)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/projects">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
            <FolderKanban className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{project.name}</h1>
            <div className="flex items-center gap-2 mt-2">
              <StatusBadge variant={stateInfo.variant}>
                {stateInfo.label}
              </StatusBadge>
              <Badge variant="outline" className="font-mono text-xs">
                {projectId}
              </Badge>
            </div>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="destructive" size="sm">
            <Trash2 className="mr-2 h-4 w-4" />
            Delete
          </Button>
        </div>
      </div>

      {/* Content Tabs */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="applications">Applications ({applications.length})</TabsTrigger>
          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Applications</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{applications.length}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Created</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatDate(project.creationDate)}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Last Changed</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatDate(project.changeDate)}</div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Project Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <p className="text-sm text-muted-foreground">Project ID</p>
                <code className="text-sm font-mono bg-muted px-2 py-0.5 rounded">{projectId}</code>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Resource Owner</p>
                <code className="text-sm font-mono bg-muted px-2 py-0.5 rounded">{project.organizationId ?? "—"}</code>
              </div>
              {project.privateLabelingSetting && (
                <div>
                  <p className="text-sm text-muted-foreground">Private Labeling</p>
                  <p className="font-medium">{project.privateLabelingSetting}</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="applications" className="space-y-4">
          {applications.length === 0 ? (
            <Card>
              <CardContent className="py-12 text-center">
                <AppWindow className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                <p className="text-muted-foreground">No applications in this project</p>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-2">
              {applications.map((app: any) => {
                const appType = getAppType(app)
                return (
                  <Card
                    key={app.applicationId}
                    className="cursor-pointer hover:bg-muted/30 transition-colors"
                    onClick={() => router.push(`/applications/${app.applicationId}?projectId=${projectId}`)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                            <AppWindow className="h-5 w-5 text-primary" />
                          </div>
                          <div>
                            <p className="font-medium">{app.name}</p>
                            <p className="text-xs text-muted-foreground">
                              Created {formatDate(app.creationDate)}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-3">
                          <Badge variant="secondary" className={`${appType.className} border-0`}>
                            {appType.label}
                          </Badge>
                          <ChevronRight className="h-4 w-4 text-muted-foreground" />
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                )
              })}
            </div>
          )}
        </TabsContent>

        <TabsContent value="settings" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Project Settings</CardTitle>
              <CardDescription>Configure project options</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Project Role Assertion</p>
                  <p className="text-sm text-muted-foreground">
                    Assert project roles in tokens
                  </p>
                </div>
                <Badge variant={project.projectRoleAssertion ? "default" : "secondary"}>
                  {project.projectRoleAssertion ? "Enabled" : "Disabled"}
                </Badge>
              </div>
              <div className="flex items-center justify-between border-t pt-4">
                <div>
                  <p className="font-medium">Project Role Check</p>
                  <p className="text-sm text-muted-foreground">
                    Check that user has project role before login
                  </p>
                </div>
                <Badge variant={project.projectRoleCheck ? "default" : "secondary"}>
                  {project.projectRoleCheck ? "Enabled" : "Disabled"}
                </Badge>
              </div>
              <div className="flex items-center justify-between border-t pt-4">
                <div>
                  <p className="font-medium">Has Project Check</p>
                  <p className="text-sm text-muted-foreground">
                    Check that application exists in project
                  </p>
                </div>
                <Badge variant={project.hasProjectCheck ? "default" : "secondary"}>
                  {project.hasProjectCheck ? "Enabled" : "Disabled"}
                </Badge>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
