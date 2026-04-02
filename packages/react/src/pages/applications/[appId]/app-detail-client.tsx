"use client"

import { ConsoleLink as Link } from "../../../context/link-context"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../components/ui/card"
import { Badge } from "../../../components/ui/badge"
import { Button } from "../../../components/ui/button"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../../components/ui/tabs"
import { ArrowLeft, Trash2, Copy, AppWindow, Globe, Server, Shield, ExternalLink } from "lucide-react"
import { Input } from "../../../components/ui/input"
import { Label } from "../../../components/ui/label"

interface ApplicationDetailClientProps {
  app: any
  appId: string
  projectId: string
  error: string | null
}

function getAppType(app: any) {
  if (app?.oidcConfig) return { label: "OIDC", icon: Globe, className: "bg-blue-100 text-blue-700 dark:bg-blue-950 dark:text-blue-400" }
  if (app?.apiConfig) return { label: "API", icon: Server, className: "bg-purple-100 text-purple-700 dark:bg-purple-950 dark:text-purple-400" }
  if (app?.samlConfig) return { label: "SAML", icon: Shield, className: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-400" }
  return { label: "Unknown", icon: AppWindow, className: "bg-muted text-muted-foreground" }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleDateString()
}

function CopyField({ label, value }: { label: string; value: string }) {
  return (
    <div className="space-y-2">
      <Label>{label}</Label>
      <div className="flex gap-2">
        <Input value={value || "—"} readOnly className="font-mono text-sm" />
        {value && (
          <Button
            variant="outline"
            size="icon"
            onClick={() => navigator.clipboard.writeText(value)}
          >
            <Copy className="h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  )
}

export function ApplicationDetailClient({ app, appId, projectId, error }: ApplicationDetailClientProps) {
  if (error || !app) {
    return (
      <div className="flex flex-col items-center justify-center h-[50vh] space-y-4">
        <h1 className="text-2xl font-bold">
          {error ? "Failed to load application" : "Application not found"}
        </h1>
        {error && <p className="text-sm text-muted-foreground">{error}</p>}
        <Button asChild>
          <Link href="/applications">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Applications
          </Link>
        </Button>
      </div>
    )
  }

  const appType = getAppType(app)
  const TypeIcon = appType.icon
  const details = app.details ?? {}
  const oidcConfig = app.oidcConfig ?? null
  const apiConfig = app.apiConfig ?? null

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/applications">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
            <TypeIcon className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{app.name}</h1>
            <Link
              href={`/projects/${projectId}`}
              className="text-sm text-muted-foreground hover:underline flex items-center gap-1 mt-1"
            >
              Project {projectId}
              <ExternalLink className="h-3 w-3" />
            </Link>
            <div className="flex items-center gap-2 mt-2">
              <Badge variant="secondary" className={`${appType.className} border-0 gap-1`}>
                <TypeIcon className="h-3 w-3" />
                {appType.label}
              </Badge>
              <Badge variant="outline" className="font-mono text-xs">
                {appId}
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
          {oidcConfig && <TabsTrigger value="configuration">OIDC Configuration</TabsTrigger>}
          {apiConfig && <TabsTrigger value="api-config">API Configuration</TabsTrigger>}
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            {/* Credentials card */}
            {oidcConfig && (
              <Card>
                <CardHeader>
                  <CardTitle>Client Credentials</CardTitle>
                  <CardDescription>
                    Use these credentials to authenticate your application
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <CopyField label="Client ID" value={oidcConfig.clientId ?? appId} />
                </CardContent>
              </Card>
            )}

            {/* Details card */}
            <Card>
              <CardHeader>
                <CardTitle>Application Details</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground">Type</p>
                  <p className="font-medium">{appType.label}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Application ID</p>
                  <code className="text-sm font-mono bg-muted px-2 py-0.5 rounded">{appId}</code>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Created</p>
                  <p className="font-medium">{formatDate(details.creationDate)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Last Changed</p>
                  <p className="font-medium">{formatDate(details.changeDate)}</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {oidcConfig && (
          <TabsContent value="configuration" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Redirect URIs</CardTitle>
                <CardDescription>Allowed redirect URIs for this application</CardDescription>
              </CardHeader>
              <CardContent>
                {(oidcConfig.redirectUris ?? []).length > 0 ? (
                  <div className="space-y-2">
                    {oidcConfig.redirectUris.map((uri: string, i: number) => (
                      <div key={i} className="flex items-center gap-2">
                        <code className="text-sm bg-muted px-3 py-1.5 rounded font-mono flex-1">
                          {uri}
                        </code>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">No redirect URIs configured</p>
                )}
              </CardContent>
            </Card>

            {(oidcConfig.postLogoutRedirectUris ?? []).length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle>Post Logout Redirect URIs</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {oidcConfig.postLogoutRedirectUris.map((uri: string, i: number) => (
                      <div key={i} className="flex items-center gap-2">
                        <code className="text-sm bg-muted px-3 py-1.5 rounded font-mono flex-1">
                          {uri}
                        </code>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            <Card>
              <CardHeader>
                <CardTitle>Grant Types</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  {(oidcConfig.grantTypes ?? []).map((gt: string) => (
                    <Badge key={gt} variant="secondary">{gt}</Badge>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Response Types</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  {(oidcConfig.responseTypes ?? []).map((rt: string) => (
                    <Badge key={rt} variant="outline">{rt}</Badge>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Token & Auth Settings</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm">App Type</span>
                  <Badge variant="outline">{oidcConfig.appType ?? "—"}</Badge>
                </div>
                <div className="flex items-center justify-between border-t pt-3">
                  <span className="text-sm">Auth Method Type</span>
                  <Badge variant="outline">{oidcConfig.authMethodType ?? "—"}</Badge>
                </div>
                <div className="flex items-center justify-between border-t pt-3">
                  <span className="text-sm">Access Token Type</span>
                  <Badge variant="outline">{oidcConfig.accessTokenType ?? "—"}</Badge>
                </div>
                <div className="flex items-center justify-between border-t pt-3">
                  <span className="text-sm">Access Token Role Assertion</span>
                  <Badge variant={oidcConfig.accessTokenRoleAssertion ? "default" : "secondary"}>
                    {oidcConfig.accessTokenRoleAssertion ? "Yes" : "No"}
                  </Badge>
                </div>
                <div className="flex items-center justify-between border-t pt-3">
                  <span className="text-sm">ID Token Role Assertion</span>
                  <Badge variant={oidcConfig.idTokenRoleAssertion ? "default" : "secondary"}>
                    {oidcConfig.idTokenRoleAssertion ? "Yes" : "No"}
                  </Badge>
                </div>
                <div className="flex items-center justify-between border-t pt-3">
                  <span className="text-sm">ID Token Userinfo Assertion</span>
                  <Badge variant={oidcConfig.idTokenUserinfoAssertion ? "default" : "secondary"}>
                    {oidcConfig.idTokenUserinfoAssertion ? "Yes" : "No"}
                  </Badge>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        )}

        {apiConfig && (
          <TabsContent value="api-config" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>API Configuration</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm">Auth Method Type</span>
                  <Badge variant="outline">{apiConfig.authMethodType ?? "—"}</Badge>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        )}
      </Tabs>
    </div>
  )
}
