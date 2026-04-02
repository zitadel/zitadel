"use client"

import { useState, useEffect } from "react"
import { useConsoleRouter as useRouter } from "../../../hooks/use-console-router"
import { ConsoleLink as Link } from "../../../context/link-context"
import { Save, X } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../components/ui/card"
import { Badge } from "../../../components/ui/badge"
import { StatusBadge } from "../../../components/ui/status-badge"
import { Button } from "../../../components/ui/button"
import { Input } from "../../../components/ui/input"
import { Label } from "../../../components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../../components/ui/tabs"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../../../components/ui/dialog"
import {
  ArrowLeft, Building2, Users, FolderOpen, AppWindow, ChevronRight, Globe, Tag,
  Pencil, Trash2, Plus, Loader2, AlertTriangle,
  Power, PowerOff,
} from "lucide-react"
import { Avatar, AvatarFallback } from "../../../components/ui/avatar"
import {
  updateOrganization,
  deleteOrganization,
  deactivateOrganization,
  activateOrganization,
} from "../../../api/manage-organization"
import {
  addOrgDomain,
  deleteOrgDomain,
  setOrgMetadata,
  deleteOrgMetadata,
  type OrgDomain,
  type OrgMetadataEntry,
} from "../../../api/org-settings"

interface OrgDetailClientProps {
  organization: any
  orgId: string
  users: any[]
  projectCount: number
  applicationCount: number
  initialDomains: OrgDomain[]
  initialMetadata: OrgMetadataEntry[]
  error: string | null
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

function getOrgState(org: any): { label: string; variant: "active" | "inactive" | "destructive" | "warning" } {
  const state = org?.state ?? ""
  const labels: Record<string, { label: string; variant: "active" | "inactive" | "destructive" | "warning" }> = {
    ORGANIZATION_STATE_ACTIVE: { label: "Active", variant: "active" },
    ORG_STATE_ACTIVE: { label: "Active", variant: "active" },
    ORGANIZATION_STATE_INACTIVE: { label: "Inactive", variant: "inactive" },
    ORG_STATE_INACTIVE: { label: "Inactive", variant: "inactive" },
    ORGANIZATION_STATE_REMOVED: { label: "Removed", variant: "destructive" },
    ORG_STATE_REMOVED: { label: "Removed", variant: "destructive" },
  }
  return labels[state] ?? { label: state, variant: "inactive" }
}

function isActive(org: any): boolean {
  const state = org?.state ?? ""
  return state === "ORGANIZATION_STATE_ACTIVE" || state === "ORG_STATE_ACTIVE"
}


export function OrgDetailClient({ organization, orgId, users, projectCount, applicationCount, initialDomains, initialMetadata, error }: OrgDetailClientProps) {
  const router = useRouter()

  // Inline edit state — matches user detail page pattern
  const [isEditing, setIsEditing] = useState(false)
  const [editName, setEditName] = useState(organization?.name ?? "")
  const [isUpdating, setIsUpdating] = useState(false)
  const [actionFeedback, setActionFeedback] = useState<string | null>(null)

  // Domain state
  const [domains, setDomains] = useState<OrgDomain[]>(initialDomains)
  const [newDomain, setNewDomain] = useState("")
  const [isDomainAdding, setIsDomainAdding] = useState(false)

  // Metadata state
  const [metadata, setMetadata] = useState<OrgMetadataEntry[]>(initialMetadata)
  const [newMetaKey, setNewMetaKey] = useState("")
  const [newMetaValue, setNewMetaValue] = useState("")
  const [isMetaAdding, setIsMetaAdding] = useState(false)

  // Delete state
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState("")
  const [isDeleting, setIsDeleting] = useState(false)
  const [deleteError, setDeleteError] = useState<string | null>(null)

  // State toggle
  const [isToggling, setIsToggling] = useState(false)
  const [toggleError, setToggleError] = useState<string | null>(null)

  function startEditing() {
    setEditName(organization.name)
    setIsEditing(true)
  }

  function cancelEditing() {
    setIsEditing(false)
  }

  const handleSave = async () => {
    if (!editName.trim()) return
    setIsUpdating(true)
    try {
      if (editName.trim() !== organization.name) {
        await updateOrganization(orgId, editName.trim())
      }
      setIsEditing(false)
      setActionFeedback("Organization updated successfully")
      setTimeout(() => setActionFeedback(null), 3000)
      router.refresh()
    } catch (e) {
      setActionFeedback(e instanceof Error ? e.message : "Failed to update organization")
      setTimeout(() => setActionFeedback(null), 3000)
    } finally {
      setIsUpdating(false)
    }
  }

  const handleDelete = async () => {
    setIsDeleting(true)
    setDeleteError(null)
    try {
      await deleteOrganization(orgId)
      setShowDeleteDialog(false)
      router.push("/organizations")
      router.refresh()
    } catch (e) {
      setDeleteError(e instanceof Error ? e.message : "Failed to delete organization")
      setIsDeleting(false)
    }
  }

  const handleToggleState = async () => {
    setIsToggling(true)
    setToggleError(null)
    try {
      if (isActive(organization)) {
        await deactivateOrganization(orgId)
      } else {
        await activateOrganization(orgId)
      }
      router.refresh()
    } catch (e) {
      setToggleError(e instanceof Error ? e.message : "Failed to change organization state")
    } finally {
      setIsToggling(false)
    }
  }

  // ─── Domain handlers ─────────────────────────────────────────────────

  const handleAddDomain = async () => {
    if (!newDomain.trim()) return
    setIsDomainAdding(true)
    try {
      await addOrgDomain(orgId, newDomain.trim())
      setDomains((prev) => [...prev, { domain: newDomain.trim(), isVerified: false, isPrimary: false, validationType: "" }])
      setNewDomain("")
      setActionFeedback("Domain added successfully")
      setTimeout(() => setActionFeedback(null), 3000)
    } catch (e) {
      setActionFeedback(e instanceof Error ? e.message : "Failed to add domain")
      setTimeout(() => setActionFeedback(null), 3000)
    } finally {
      setIsDomainAdding(false)
    }
  }

  const handleDeleteDomain = async (domain: string) => {
    try {
      await deleteOrgDomain(orgId, domain)
      setDomains((prev) => prev.filter((d) => d.domain !== domain))
      setActionFeedback("Domain deleted")
      setTimeout(() => setActionFeedback(null), 3000)
    } catch (e) {
      setActionFeedback(e instanceof Error ? e.message : "Failed to delete domain")
      setTimeout(() => setActionFeedback(null), 3000)
    }
  }

  // ─── Metadata handlers ───────────────────────────────────────────────

  const handleAddMetadata = async () => {
    if (!newMetaKey.trim() || !newMetaValue.trim()) return
    setIsMetaAdding(true)
    try {
      await setOrgMetadata(orgId, newMetaKey.trim(), newMetaValue.trim())
      setMetadata((prev) => {
        const existing = prev.findIndex((m) => m.key === newMetaKey.trim())
        if (existing >= 0) {
          const updated = [...prev]
          updated[existing] = { key: newMetaKey.trim(), value: newMetaValue.trim() }
          return updated
        }
        return [...prev, { key: newMetaKey.trim(), value: newMetaValue.trim() }]
      })
      setNewMetaKey("")
      setNewMetaValue("")
      setActionFeedback("Metadata saved")
      setTimeout(() => setActionFeedback(null), 3000)
    } catch (e) {
      setActionFeedback(e instanceof Error ? e.message : "Failed to save metadata")
      setTimeout(() => setActionFeedback(null), 3000)
    } finally {
      setIsMetaAdding(false)
    }
  }

  const handleDeleteMetadata = async (key: string) => {
    try {
      await deleteOrgMetadata(orgId, [key])
      setMetadata((prev) => prev.filter((m) => m.key !== key))
      setActionFeedback("Metadata deleted")
      setTimeout(() => setActionFeedback(null), 3000)
    } catch (e) {
      setActionFeedback(e instanceof Error ? e.message : "Failed to delete metadata")
      setTimeout(() => setActionFeedback(null), 3000)
    }
  }

  if (error || !organization) {
    return (
      <div className="flex flex-col items-center justify-center h-[50vh] space-y-4">
        <h1 className="text-2xl font-bold">
          {error ? "Failed to load organization" : "Organization not found"}
        </h1>
        {error && <p className="text-sm text-muted-foreground">{error}</p>}
        <Button asChild>
          <Link href="/organizations">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Organizations
          </Link>
        </Button>
      </div>
    )
  }

  const stateInfo = getOrgState(organization)
  const details = organization.details ?? {}
  const orgIsActive = isActive(organization)

  return (
    <div className="space-y-6">
      {/* Header — aligned with user detail page */}
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/organizations">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <Avatar className="h-14 w-14">
            <AvatarFallback className="text-lg bg-muted">
              <Building2 className="h-6 w-6 text-muted-foreground" />
            </AvatarFallback>
          </Avatar>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold">{organization.name}</h1>
            </div>
            {organization.primaryDomain && (
              <p className="text-sm text-muted-foreground">{organization.primaryDomain}</p>
            )}
            <div className="flex items-center gap-2 mt-1">
              <StatusBadge variant={stateInfo.variant}>
                {stateInfo.label}
              </StatusBadge>
              {organization.isDefault && (
                <Badge variant="secondary" className="text-xs">default</Badge>
              )}
            </div>
          </div>
        </div>
        {/* Action buttons — matches user detail page pattern */}
        {isEditing ? (
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={cancelEditing} disabled={isUpdating}>
              <X className="mr-2 h-4 w-4" />
              Cancel
            </Button>
            <Button size="sm" onClick={handleSave} disabled={isUpdating}>
              {isUpdating ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Save className="mr-2 h-4 w-4" />}
              Save
            </Button>
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleToggleState}
              disabled={isToggling}
            >
              {isToggling ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : orgIsActive ? (
                <PowerOff className="mr-2 h-4 w-4" />
              ) : (
                <Power className="mr-2 h-4 w-4" />
              )}
              {orgIsActive ? "Deactivate" : "Activate"}
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={startEditing}
            >
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </Button>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => {
                setDeleteConfirm("")
                setDeleteError(null)
                setShowDeleteDialog(true)
              }}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Delete
            </Button>
          </div>
        )}
      </div>

      {/* Action feedback toast — matches user detail page */}
      {actionFeedback && (
        <div className="fixed top-4 right-4 z-50 bg-foreground text-background px-4 py-2.5 rounded-lg shadow-lg text-sm font-medium animate-in fade-in slide-in-from-top-2 duration-200">
          {actionFeedback}
        </div>
      )}

      {toggleError && (
        <div className="rounded-md border border-destructive/50 bg-destructive/10 p-3">
          <p className="text-sm text-destructive">{toggleError}</p>
        </div>
      )}

      {/* Tabs — aligned with user detail page (Overview, then resource tabs) */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="domains">Domains</TabsTrigger>
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
        </TabsList>

        {/* Overview — two-column layout like user detail */}
        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            {/* Organization Information */}
            <Card>
              <CardHeader>
                <CardTitle>Organization Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm text-muted-foreground">Name</Label>
                    {isEditing ? (
                      <Input
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                        onKeyDown={(e) => { if (e.key === "Enter" && editName.trim()) handleSave() }}
                        maxLength={200}
                        className="mt-1"
                        disabled={isUpdating}
                      />
                    ) : (
                      <p className="font-medium mt-1">{organization.name}</p>
                    )}
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Primary Domain</p>
                    <p className="font-medium">{organization.primaryDomain ?? "—"}</p>
                  </div>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">State</p>
                  <StatusBadge variant={stateInfo.variant} className="mt-0.5">
                    {stateInfo.label}
                  </StatusBadge>
                </div>
              </CardContent>
            </Card>

            {/* Account Details */}
            <Card>
              <CardHeader>
                <CardTitle>Account Details</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground">Organization ID</p>
                  <code className="text-sm font-mono bg-muted px-2 py-0.5 rounded">{orgId}</code>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Created</p>
                  <p className="font-medium">{formatDate(details.creationDate)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Last Changed</p>
                  <p className="font-medium">{formatDate(details.changeDate)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Resource Owner</p>
                  <code className="text-sm font-mono bg-muted px-2 py-0.5 rounded">
                    {details.resourceOwner ?? "—"}
                  </code>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Clickable resource cards — link to pages with preset org filter */}
          <div className="grid gap-4 md:grid-cols-3">
            <Link
              href={`/users?org=${orgId}`}
              className="block"
            >
              <Card className="hover:bg-muted/30 transition-colors cursor-pointer">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Users</CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent className="flex items-center justify-between">
                  <div className="text-2xl font-bold">{users.length}</div>
                  <ChevronRight className="h-4 w-4 text-muted-foreground" />
                </CardContent>
              </Card>
            </Link>
            <Link
              href={`/projects?org=${orgId}`}
              className="block"
            >
              <Card className="hover:bg-muted/30 transition-colors cursor-pointer">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Projects</CardTitle>
                  <FolderOpen className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent className="flex items-center justify-between">
                  <div className="text-2xl font-bold">{projectCount}</div>
                  <ChevronRight className="h-4 w-4 text-muted-foreground" />
                </CardContent>
              </Card>
            </Link>
            <Link
              href={`/applications?org=${orgId}`}
              className="block"
            >
              <Card className="hover:bg-muted/30 transition-colors cursor-pointer">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Applications</CardTitle>
                  <AppWindow className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent className="flex items-center justify-between">
                  <div className="text-2xl font-bold">{applicationCount}</div>
                  <ChevronRight className="h-4 w-4 text-muted-foreground" />
                </CardContent>
              </Card>
            </Link>
          </div>
        </TabsContent>

        {/* Domains tab */}
        <TabsContent value="domains" className="space-y-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Domains</CardTitle>
                <CardDescription>
                  Manage domains registered to this organization. Verified domains can be used for user login and domain discovery.
                </CardDescription>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Add domain form */}
              <div className="flex items-end gap-2">
                <div className="flex-1">
                  <Label htmlFor="org-domain" className="text-xs text-muted-foreground">Domain</Label>
                  <Input
                    id="org-domain"
                    value={newDomain}
                    onChange={(e) => setNewDomain(e.target.value)}
                    placeholder="example.com"
                    onKeyDown={(e) => {
                      if (e.key === "Enter" && newDomain.trim()) {
                        handleAddDomain()
                      }
                    }}
                    disabled={isDomainAdding}
                    className="h-8 text-sm"
                  />
                </div>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={handleAddDomain}
                  disabled={!newDomain.trim() || isDomainAdding}
                  className="h-8"
                >
                  {isDomainAdding ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
                </Button>
              </div>
              {/* Domain list */}
              {domains.length === 0 ? (
                <p className="text-muted-foreground text-sm">No domains configured</p>
              ) : (
                <div className="space-y-2">
                  {domains.map((d) => (
                    <div key={d.domain} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center gap-3">
                        <Globe className="h-4 w-4 text-muted-foreground" />
                        <div>
                          <p className="font-medium text-sm">{d.domain}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {d.isPrimary && (
                          <Badge variant="secondary" className="text-xs">primary</Badge>
                        )}
                        {d.isVerified ? (
                          <Badge variant="secondary" className="text-xs">verified</Badge>
                        ) : (
                          <Badge variant="outline" className="text-xs">unverified</Badge>
                        )}
                        {!d.isPrimary && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDeleteDomain(d.domain)}
                          >
                            <Trash2 className="h-3.5 w-3.5 text-destructive" />
                          </Button>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Metadata tab */}
        <TabsContent value="metadata" className="space-y-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Metadata</CardTitle>
                <CardDescription>
                  Key-value metadata attached to this organization.
                </CardDescription>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Add metadata form */}
              <div className="flex items-end gap-2">
                <div className="flex-1">
                  <Label htmlFor="org-meta-key" className="text-xs text-muted-foreground">Key</Label>
                  <Input
                    id="org-meta-key"
                    value={newMetaKey}
                    onChange={(e) => setNewMetaKey(e.target.value)}
                    placeholder="metadata.key"
                    disabled={isMetaAdding}
                    className="h-8 text-sm"
                  />
                </div>
                <div className="flex-1">
                  <Label htmlFor="org-meta-value" className="text-xs text-muted-foreground">Value</Label>
                  <Input
                    id="org-meta-value"
                    value={newMetaValue}
                    onChange={(e) => setNewMetaValue(e.target.value)}
                    placeholder="value"
                    disabled={isMetaAdding}
                    className="h-8 text-sm"
                    onKeyDown={(e) => {
                      if (e.key === "Enter" && newMetaKey.trim() && newMetaValue.trim()) {
                        handleAddMetadata()
                      }
                    }}
                  />
                </div>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={handleAddMetadata}
                  disabled={!newMetaKey.trim() || !newMetaValue.trim() || isMetaAdding}
                  className="h-8"
                >
                  {isMetaAdding ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
                </Button>
              </div>
              {/* Metadata list */}
              {metadata.length === 0 ? (
                <p className="text-muted-foreground text-sm">No metadata entries</p>
              ) : (
                <div className="space-y-2">
                  {metadata.map((m) => (
                    <div key={m.key} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center gap-3">
                        <Tag className="h-4 w-4 text-muted-foreground" />
                        <div>
                          <p className="font-medium text-sm">{m.key}</p>
                          <p className="text-xs text-muted-foreground">{m.value}</p>
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteMetadata(m.key)}
                      >
                        <Trash2 className="h-3.5 w-3.5 text-destructive" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Delete Confirmation Dialog */}
      <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-destructive">
              <AlertTriangle className="h-5 w-5" />
              Delete Organization
            </DialogTitle>
            <DialogDescription>
              This will permanently delete <strong>{organization.name}</strong> and all its resources
              including users, projects, and grants. This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-2">
              <Label htmlFor="delete-confirm">
                Type <strong>{organization.name}</strong> to confirm
              </Label>
              <Input
                id="delete-confirm"
                value={deleteConfirm}
                onChange={(e) => setDeleteConfirm(e.target.value)}
                placeholder={organization.name}
                disabled={isDeleting}
              />
            </div>
            {deleteError && (
              <div className="rounded-md border border-destructive/50 bg-destructive/10 p-3">
                <p className="text-sm text-destructive">{deleteError}</p>
              </div>
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDeleteDialog(false)} disabled={isDeleting}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteConfirm !== organization.name || isDeleting}
            >
              {isDeleting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Deleting...
                </>
              ) : (
                "Delete Organization"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
