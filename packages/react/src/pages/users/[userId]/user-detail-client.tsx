"use client"

import { useState, useTransition, useEffect } from "react"
import { ConsoleLink as Link } from "../../../context/link-context"
import { useConsoleRouter } from "../../../hooks/use-console-router"
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card"
import { Badge } from "../../../components/ui/badge"
import { StatusBadge } from "../../../components/ui/status-badge"
import { Button } from "../../../components/ui/button"
import { Avatar, AvatarFallback } from "../../../components/ui/avatar"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../../components/ui/tabs"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "../../../components/ui/alert-dialog"
import { ArrowLeft, Edit, Lock, Unlock, Trash2, Mail, KeyRound, Shield, Activity, User, Monitor, Globe, Clock, Fingerprint, ChevronLeft, ChevronRight, Loader2, Save, X, Power, Plus, Check } from "lucide-react"
import { Input } from "../../../components/ui/input"
import { Label } from "../../../components/ui/label"
import { listUserSessions, deleteSession } from "../../../api/sessions"
import { lockUser, unlockUser, deactivateUser, reactivateUser, deleteUser, resetPassword, updateUser, type UpdateUserData } from "../../../api/user-actions"
import { listAuthMethodTypes, listAuthFactors, listPasskeys, removePasskey, removeTOTP, removeOTPSMS, removeOTPEmail, removeU2F } from "../../../api/user-security"
import { listUserMetadata, setUserMetadata, deleteUserMetadata, type UserMetadataEntry } from "../../../api/user-metadata"

interface UserDetailClientProps {
  user: any
  userId: string
  initialSessions: any[]
  totalSessions: number
  initialAuthMethods: string[]
  initialAuthFactors: any[]
  initialPasskeys: any[]
  initialMetadata: UserMetadataEntry[]
  error: string | null
}

const stateLabels: Record<string, { label: string; variant: "active" | "inactive" | "destructive" | "warning" }> = {
  USER_STATE_ACTIVE: {
    label: "Active",
    variant: "active",
  },
  USER_STATE_INACTIVE: {
    label: "Inactive",
    variant: "inactive",
  },
  USER_STATE_LOCKED: {
    label: "Locked",
    variant: "destructive",
  },
  USER_STATE_INITIAL: {
    label: "Initial",
    variant: "warning",
  },
}

function getUserInfo(user: any) {
  if (user.human) {
    const human = user.human
    const profile = human?.profile ?? {}
    const email = human?.email ?? {}
    const phone = human?.phone ?? {}
    return {
      kind: "human" as const,
      givenName: profile.givenName ?? "",
      familyName: profile.familyName ?? "",
      displayName: profile.displayName || `${profile.givenName ?? ""} ${profile.familyName ?? ""}`.trim() || user.username || "Unknown",
      email: email.email ?? "",
      isEmailVerified: email.isVerified ?? false,
      phone: phone.phone ?? "",
      isPhoneVerified: phone.isVerified ?? false,
      avatarUrl: profile.avatarUrl ?? "",
      nickName: profile.nickName ?? "",
      preferredLanguage: profile.preferredLanguage ?? "",
      gender: profile.gender ?? "",
    }
  }
  if (user.machine) {
    return {
      kind: "machine" as const,
      givenName: "",
      familyName: "",
      displayName: user.machine?.name || user.username || "Machine User",
      email: "",
      isEmailVerified: false,
      phone: "",
      isPhoneVerified: false,
      avatarUrl: "",
      nickName: "",
      preferredLanguage: "",
      gender: "",
      description: user.machine?.description ?? "",
      accessTokenType: user.machine?.accessTokenType ?? "",
    }
  }
  return {
    kind: "unknown" as const,
    givenName: "",
    familyName: "",
    displayName: user.username || "Unknown",
    email: "",
    isEmailVerified: false,
    phone: "",
    isPhoneVerified: false,
    avatarUrl: "",
    nickName: "",
    preferredLanguage: "",
    gender: "",
  }
}

function getInitials(info: ReturnType<typeof getUserInfo>) {
  if (info.givenName && info.familyName) {
    return `${info.givenName.charAt(0)}${info.familyName.charAt(0)}`.toUpperCase()
  }
  return info.displayName.charAt(0).toUpperCase()
}

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  try {
    return new Date(dateStr).toLocaleDateString(undefined, {
      year: "numeric", month: "short", day: "numeric", hour: "2-digit", minute: "2-digit",
    })
  } catch {
    return dateStr
  }
}

function getSessionFactors(session: any) {
  const factors = session.factors ?? {}
  const items: { label: string; verifiedAt?: string }[] = []
  if (factors.user) items.push({ label: "User", verifiedAt: factors.user.verifiedAt })
  if (factors.password) items.push({ label: "Password", verifiedAt: factors.password.verifiedAt })
  if (factors.webAuthN) items.push({ label: "WebAuthn", verifiedAt: factors.webAuthN.verifiedAt })
  if (factors.totp) items.push({ label: "TOTP", verifiedAt: factors.totp.verifiedAt })
  if (factors.otpSms) items.push({ label: "OTP SMS", verifiedAt: factors.otpSms.verifiedAt })
  if (factors.otpEmail) items.push({ label: "OTP Email", verifiedAt: factors.otpEmail.verifiedAt })
  if (factors.intent) items.push({ label: "Intent", verifiedAt: factors.intent.verifiedAt })
  return items
}

const PAGE_SIZE = 10

export function UserDetailClient({ user, userId, initialSessions, totalSessions, initialAuthMethods, initialAuthFactors, initialPasskeys, initialMetadata, error }: UserDetailClientProps) {
  const router = useConsoleRouter()
  const [sessions, setSessions] = useState(initialSessions)
  const [total, setTotal] = useState(totalSessions)
  const [page, setPage] = useState(0)
  const [isLoadingSessions, startSessionTransition] = useTransition()
  const [isActing, startActing] = useTransition()
  const [actionFeedback, setActionFeedback] = useState<string | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [deleteConfirmation, setDeleteConfirmation] = useState("")
  const [isEditing, setIsEditing] = useState(false)
  const [isSaving, startSaving] = useTransition()

  // Security state
  const [authMethods, setAuthMethods] = useState<string[]>(initialAuthMethods)
  const [authFactors, setAuthFactors] = useState<any[]>(initialAuthFactors)
  const [passkeys, setPasskeys] = useState<any[]>(initialPasskeys)
  const [isSecurityLoading, startSecurityTransition] = useTransition()

  // Metadata state
  const [metadata, setMetadata] = useState<UserMetadataEntry[]>(initialMetadata)
  const [newMetaKey, setNewMetaKey] = useState("")
  const [newMetaValue, setNewMetaValue] = useState("")
  const [isMetadataLoading, startMetadataTransition] = useTransition()

  // Edit form state — initialized from current user
  const currentInfo = user ? getUserInfo(user) : null
  const [editForm, setEditForm] = useState({
    givenName: currentInfo?.givenName ?? "",
    familyName: currentInfo?.familyName ?? "",
    nickName: currentInfo?.nickName ?? "",
    displayName: currentInfo?.displayName ?? "",
    email: currentInfo?.email ?? "",
    phone: currentInfo?.phone ?? "",
    preferredLanguage: currentInfo?.preferredLanguage ?? "",
    username: user?.username ?? "",
  })

  function startEditing() {
    const info = getUserInfo(user)
    setEditForm({
      givenName: info.givenName,
      familyName: info.familyName,
      nickName: info.nickName,
      displayName: info.displayName,
      email: info.email,
      phone: info.phone,
      preferredLanguage: info.preferredLanguage,
      username: user.username ?? "",
    })
    setIsEditing(true)
  }

  function cancelEditing() {
    setIsEditing(false)
  }

  function handleSave() {
    startSaving(async () => {
      try {
        const info = getUserInfo(user)
        const data: UpdateUserData = {}

        // Detect changed profile fields
        const profileChanges: UpdateUserData["profile"] = {}
        let hasProfileChanges = false
        if (editForm.givenName !== info.givenName) { profileChanges.givenName = editForm.givenName; hasProfileChanges = true }
        if (editForm.familyName !== info.familyName) { profileChanges.familyName = editForm.familyName; hasProfileChanges = true }
        if (editForm.nickName !== info.nickName) { profileChanges.nickName = editForm.nickName; hasProfileChanges = true }
        if (editForm.displayName !== info.displayName) { profileChanges.displayName = editForm.displayName; hasProfileChanges = true }
        if (editForm.preferredLanguage !== info.preferredLanguage) { profileChanges.preferredLanguage = editForm.preferredLanguage; hasProfileChanges = true }
        if (hasProfileChanges) data.profile = profileChanges

        // Detect changed email
        if (editForm.email !== info.email) data.email = editForm.email

        // Detect changed phone
        if (editForm.phone !== info.phone) data.phone = editForm.phone

        // Detect changed username
        if (editForm.username !== (user.username ?? "")) data.username = editForm.username

        // Only call if something changed
        if (Object.keys(data).length > 0) {
          await updateUser(userId, data)
        }

        setIsEditing(false)
        setActionFeedback("User updated successfully")
        setTimeout(() => setActionFeedback(null), 3000)
        router.refresh()
      } catch (e) {
        console.error("Failed to update user:", e)
        setActionFeedback("Failed to update user")
        setTimeout(() => setActionFeedback(null), 3000)
      }
    })
  }

  function handleActivateDeactivate() {
    startActing(async () => {
      try {
        if (user.state === "USER_STATE_INACTIVE") {
          await reactivateUser(userId)
        } else {
          await deactivateUser(userId)
        }
        router.refresh()
      } catch (e) {
        console.error("Failed to activate/deactivate user:", e)
      }
    })
  }

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

  function goToPage(newPage: number) {
    startSessionTransition(async () => {
      try {
        const result = await listUserSessions(userId, PAGE_SIZE, newPage * PAGE_SIZE)
        setSessions(result.sessions)
        setTotal(result.totalResult)
        setPage(newPage)
      } catch (e) {
        console.error("Failed to load sessions page:", e)
      }
    })
  }

  function handleRevokeSession(sessionId: string) {
    startActing(async () => {
      try {
        await deleteSession(sessionId)
        // Refresh session list
        const result = await listUserSessions(userId, PAGE_SIZE, page * PAGE_SIZE)
        setSessions(result.sessions)
        setTotal(result.totalResult)
        setActionFeedback("Session revoked successfully")
        setTimeout(() => setActionFeedback(null), 3000)
      } catch (e) {
        console.error("Failed to revoke session:", e)
        setActionFeedback("Failed to revoke session")
        setTimeout(() => setActionFeedback(null), 3000)
      }
    })
  }

  function handleLockUnlock() {
    startActing(async () => {
      try {
        if (user.state === "USER_STATE_LOCKED") {
          await unlockUser(userId)
        } else {
          await lockUser(userId)
        }
        router.refresh()
      } catch (e) {
        console.error("Failed to lock/unlock user:", e)
      }
    })
  }

  function handleDelete() {
    startActing(async () => {
      try {
        await deleteUser(userId)
        router.push("/users")
      } catch (e) {
        console.error("Failed to delete user:", e)
      }
    })
  }

  function handleResetPassword() {
    startActing(async () => {
      try {
        await resetPassword(userId)
        setActionFeedback("Password reset email sent successfully")
        setTimeout(() => setActionFeedback(null), 3000)
      } catch (e) {
        console.error("Failed to reset password:", e)
      }
    })
  }

  if (error || !user) {
    return (
      <div className="flex flex-col items-center justify-center h-[50vh] space-y-4">
        <h1 className="text-2xl font-bold">
          {error ? "Failed to load user" : "User not found"}
        </h1>
        {error && <p className="text-sm text-muted-foreground">{error}</p>}
        <Button asChild>
          <Link href="/users">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Users
          </Link>
        </Button>
      </div>
    )
  }

  const info = getUserInfo(user)
  const stateInfo = stateLabels[user.state] ?? { label: user.state, variant: "inactive" as const }
  const details = user.details ?? {}

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-start gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/users">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <Avatar className="h-16 w-16">
            <AvatarFallback className="text-lg">
              {getInitials(info)}
            </AvatarFallback>
          </Avatar>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{info.displayName}</h1>
            {info.email && (
              <p className="text-muted-foreground">{info.email}</p>
            )}
            <div className="flex items-center gap-2 mt-2">
              <StatusBadge variant={stateInfo.variant}>
                {stateInfo.label}
              </StatusBadge>
              {info.kind === "machine" && (
                <Badge variant="secondary">Machine User</Badge>
              )}
              {user.username && (
                <Badge variant="outline">@{user.username}</Badge>
              )}
            </div>
          </div>
        </div>
        {isEditing ? (
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={cancelEditing} disabled={isSaving}>
              <X className="mr-2 h-4 w-4" />
              Cancel
            </Button>
            <Button size="sm" onClick={handleSave} disabled={isSaving}>
              {isSaving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Save className="mr-2 h-4 w-4" />}
              Save
            </Button>
          </div>
        ) : (
          <div className="flex items-center gap-2">
            {info.email && (
              <Button variant="outline" size="sm" asChild>
                <a href={`mailto:${info.email}`}>
                  <Mail className="mr-2 h-4 w-4" />
                  Email
                </a>
              </Button>
            )}
            <Button variant="outline" size="sm" onClick={handleLockUnlock} disabled={isActing}>
              {isActing ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : user.state === "USER_STATE_LOCKED" ? (
                <Unlock className="mr-2 h-4 w-4" />
              ) : (
                <Lock className="mr-2 h-4 w-4" />
              )}
              {user.state === "USER_STATE_LOCKED" ? "Unlock" : "Lock"}
            </Button>
            <Button variant="outline" size="sm" onClick={startEditing}>
              <Edit className="mr-2 h-4 w-4" />
              Edit
            </Button>
            <AlertDialog open={deleteDialogOpen} onOpenChange={(open) => { setDeleteDialogOpen(open); if (!open) setDeleteConfirmation("") }}>
              <AlertDialogTrigger asChild>
                <Button variant="destructive" size="sm" disabled={isActing}>
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete User</AlertDialogTitle>
                  <AlertDialogDescription>
                    This action cannot be undone. This will permanently delete <strong>{info.displayName}</strong> and remove all associated data.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <div className="py-2">
                  <Label htmlFor="delete-confirm" className="text-sm">
                    Type <span className="font-mono font-semibold">{user.username}</span> to confirm
                  </Label>
                  <Input
                    id="delete-confirm"
                    className="mt-2"
                    placeholder={user.username}
                    value={deleteConfirmation}
                    onChange={(e) => setDeleteConfirmation(e.target.value)}
                    autoComplete="off"
                  />
                </div>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    onClick={handleDelete}
                    disabled={deleteConfirmation !== user.username || isActing}
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                  >
                    {isActing ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                    Delete User
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        )}
      </div>

      {/* Action feedback toast */}
      {actionFeedback && (
        <div className="fixed top-4 right-4 z-50 bg-foreground text-background px-4 py-2.5 rounded-lg shadow-lg text-sm font-medium animate-in fade-in slide-in-from-top-2 duration-200">
          {actionFeedback}
        </div>
      )}

      {/* Tabs — anchored via URL hash */}
      <Tabs defaultValue={typeof window !== "undefined" && window.location.hash ? window.location.hash.slice(1) : "overview"} className="space-y-4" onValueChange={(v) => { window.history.replaceState(null, "", `#${v}`) }}>
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="sessions">Sessions ({total})</TabsTrigger>
          <TabsTrigger value="security">Security</TabsTrigger>
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
          <TabsTrigger value="activity">Activity</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            {/* User Information Card */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">User Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {info.kind === "human" ? (
                  <>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label className="text-sm text-muted-foreground">First Name</Label>
                        {isEditing ? (
                          <Input
                            value={editForm.givenName}
                            onChange={(e) => setEditForm(prev => ({ ...prev, givenName: e.target.value }))}
                            className="mt-1"
                          />
                        ) : (
                          <p className="font-medium mt-1">{info.givenName || "—"}</p>
                        )}
                      </div>
                      <div>
                        <Label className="text-sm text-muted-foreground">Last Name</Label>
                        {isEditing ? (
                          <Input
                            value={editForm.familyName}
                            onChange={(e) => setEditForm(prev => ({ ...prev, familyName: e.target.value }))}
                            className="mt-1"
                          />
                        ) : (
                          <p className="font-medium mt-1">{info.familyName || "—"}</p>
                        )}
                      </div>
                    </div>
                    <div>
                      <Label className="text-sm text-muted-foreground">Nickname</Label>
                      {isEditing ? (
                        <Input
                          value={editForm.nickName}
                          onChange={(e) => setEditForm(prev => ({ ...prev, nickName: e.target.value }))}
                          placeholder="Optional"
                          className="mt-1"
                        />
                      ) : (
                        <p className="font-medium mt-1">{info.nickName || "—"}</p>
                      )}
                    </div>
                    <div>
                      <Label className="text-sm text-muted-foreground">Display Name</Label>
                      {isEditing ? (
                        <Input
                          value={editForm.displayName}
                          onChange={(e) => setEditForm(prev => ({ ...prev, displayName: e.target.value }))}
                          className="mt-1"
                        />
                      ) : (
                        <p className="font-medium mt-1">{info.displayName}</p>
                      )}
                    </div>
                    <div>
                      <Label className="text-sm text-muted-foreground">Email</Label>
                      {isEditing ? (
                        <Input
                          type="email"
                          value={editForm.email}
                          onChange={(e) => setEditForm(prev => ({ ...prev, email: e.target.value }))}
                          className="mt-1"
                        />
                      ) : (
                        <div className="flex items-center gap-2 mt-1">
                          <p className="font-medium">{info.email || "—"}</p>
                          {info.email && (
                            <Badge variant={info.isEmailVerified ? "secondary" : "outline"} className="text-xs">
                              {info.isEmailVerified ? "Verified" : "Unverified"}
                            </Badge>
                          )}
                        </div>
                      )}
                    </div>
                    <div>
                      <Label className="text-sm text-muted-foreground">Phone</Label>
                      {isEditing ? (
                        <Input
                          type="tel"
                          value={editForm.phone}
                          onChange={(e) => setEditForm(prev => ({ ...prev, phone: e.target.value }))}
                          placeholder="Optional"
                          className="mt-1"
                        />
                      ) : info.phone ? (
                        <div className="flex items-center gap-2 mt-1">
                          <p className="font-medium">{info.phone}</p>
                          <Badge variant={info.isPhoneVerified ? "secondary" : "outline"} className="text-xs">
                            {info.isPhoneVerified ? "Verified" : "Unverified"}
                          </Badge>
                        </div>
                      ) : (
                        <p className="font-medium mt-1 text-muted-foreground">—</p>
                      )}
                    </div>
                    <div>
                      <Label className="text-sm text-muted-foreground">Preferred Language</Label>
                      {isEditing ? (
                        <Input
                          value={editForm.preferredLanguage}
                          onChange={(e) => setEditForm(prev => ({ ...prev, preferredLanguage: e.target.value }))}
                          placeholder="e.g. en, de"
                          className="mt-1"
                        />
                      ) : (
                        <p className="font-medium mt-1">{info.preferredLanguage || "—"}</p>
                      )}
                    </div>
                  </>
                ) : (
                  <>
                    <div>
                      <p className="text-sm text-muted-foreground">Name</p>
                      <p className="font-medium">{info.displayName}</p>
                    </div>
                    {"description" in info && info.description && (
                      <div>
                        <p className="text-sm text-muted-foreground">Description</p>
                        <p className="font-medium">{info.description}</p>
                      </div>
                    )}
                    {"accessTokenType" in info && (
                      <div>
                        <p className="text-sm text-muted-foreground">Access Token Type</p>
                        <p className="font-medium">{info.accessTokenType || "Bearer"}</p>
                      </div>
                    )}
                  </>
                )}
              </CardContent>
            </Card>

            {/* Account Details Card */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Account Details</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground">User ID</p>
                  <p className="font-medium font-mono text-sm">{user.userId}</p>
                </div>
                <div>
                  <Label className="text-sm text-muted-foreground">Username</Label>
                  {isEditing ? (
                    <Input
                      value={editForm.username}
                      onChange={(e) => setEditForm(prev => ({ ...prev, username: e.target.value }))}
                      className="mt-1 font-mono text-sm"
                    />
                  ) : (
                    <p className="font-medium">{user.username || "—"}</p>
                  )}
                </div>
                {user.loginNames && user.loginNames.length > 0 && (
                  <div>
                    <p className="text-sm text-muted-foreground">Login Names</p>
                    <div className="flex flex-wrap gap-1 mt-1">
                      {user.loginNames.map((name: string) => (
                        <Badge key={name} variant="outline" className="text-xs">
                          {name}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
                <div>
                  <p className="text-sm text-muted-foreground">State</p>
                  <StatusBadge variant={stateInfo.variant}>
                    {stateInfo.label}
                  </StatusBadge>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Created</p>
                  <p className="font-medium text-sm">{formatDate(details.creationDate)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Last Changed</p>
                  <p className="font-medium text-sm">{formatDate(details.changeDate)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Resource Owner</p>
                  <p className="font-medium font-mono text-sm">{details.resourceOwner || "—"}</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Sessions Tab */}
        <TabsContent value="sessions" className="space-y-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0">
              <div>
                <CardTitle className="text-base">Sessions</CardTitle>
                <p className="text-sm text-muted-foreground mt-1">Login sessions for this user</p>
              </div>
              {total > 0 && (
                <Badge variant="secondary" className="tabular-nums">{total}</Badge>
              )}
            </CardHeader>
            <CardContent>
              {sessions.length === 0 ? (
                <p className="text-muted-foreground text-sm">No sessions found for this user.</p>
              ) : (
                <div className="space-y-3">
                  {sessions.map((session: any) => {
                    const factors = getSessionFactors(session)
                    const userAgent = session.userAgent ?? {}
                    const headerName = userAgent.header?.name ?? ""
                    const fingerprintId = userAgent.fingerprintId ?? ""
                    const ip = userAgent.ip ?? ""
                    const displayTitle = ip || fingerprintId || session.id
                    const isExpired = session.expirationDate ? new Date(session.expirationDate) < new Date() : false
                    return (
                      <div
                        key={session.id}
                        className={`flex items-center justify-between p-4 border rounded-lg${isExpired ? " opacity-60" : ""}`}
                      >
                        <div className="space-y-1.5 flex-1 min-w-0">
                          <div className="flex items-center gap-2">
                            <Monitor className="h-4 w-4 text-muted-foreground shrink-0" />
                            <span className="font-medium text-sm">
                              {displayTitle}
                            </span>
                          </div>
                          {headerName && (
                            <p className="text-sm text-muted-foreground pl-6 truncate">
                              {headerName}
                            </p>
                          )}
                          <p className="text-sm text-muted-foreground pl-6">
                            Created: {formatDate(session.creationDate)}
                          </p>
                          {session.expirationDate && (
                            <p className="text-sm text-muted-foreground pl-6">
                              {isExpired ? "Expired" : "Expires"}: {formatDate(session.expirationDate)}
                            </p>
                          )}
                        </div>
                        <div className="flex items-center gap-2 shrink-0 ml-4">
                          {isExpired ? (
                            <StatusBadge variant="inactive">Expired</StatusBadge>
                          ) : (
                            <>
                              <StatusBadge variant="active">Active</StatusBadge>
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleRevokeSession(session.id)}
                                disabled={isActing}
                              >
                                Revoke
                              </Button>
                            </>
                          )}
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
              {/* Pagination */}
              {total > PAGE_SIZE && (
                <div className="flex items-center justify-between pt-4 border-t mt-4">
                  <p className="text-sm text-muted-foreground">
                    Page {page + 1} of {totalPages} · {total} sessions
                  </p>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => goToPage(page - 1)}
                      disabled={page === 0 || isLoadingSessions}
                    >
                      <ChevronLeft className="h-4 w-4 mr-1" />
                      Previous
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => goToPage(page + 1)}
                      disabled={page >= totalPages - 1 || isLoadingSessions}
                    >
                      Next
                      <ChevronRight className="h-4 w-4 ml-1" />
                    </Button>
                  </div>
                </div>
              )}
              {isLoadingSessions && (
                <div className="flex items-center justify-center py-4">
                  <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security Tab */}
        <TabsContent value="security" className="space-y-4">
          {/* Password */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <KeyRound className="h-5 w-5 text-muted-foreground" />
                  <div>
                    <CardTitle className="text-base">Password</CardTitle>
                    <p className="text-sm text-muted-foreground mt-0.5">
                      {authMethods.includes("AUTHENTICATION_METHOD_TYPE_PASSWORD")
                        ? "Password authentication is configured"
                        : "No password set for this user"}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <StatusBadge variant={authMethods.includes("AUTHENTICATION_METHOD_TYPE_PASSWORD") ? "active" : "inactive"}>
                    {authMethods.includes("AUTHENTICATION_METHOD_TYPE_PASSWORD") ? "Active" : "Inactive"}
                  </StatusBadge>
                  <Button variant="outline" size="sm" onClick={handleResetPassword} disabled={isActing}>
                    {isActing ? <Loader2 className="h-4 w-4 animate-spin" /> : "Reset"}
                  </Button>
                </div>
              </div>
            </CardHeader>
          </Card>

          {/* Multi-Factor Authentication */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Shield className="h-5 w-5 text-muted-foreground" />
                  <CardTitle className="text-base">Multi-Factor Authentication</CardTitle>
                  {authFactors.length > 0 && (
                    <Badge variant="outline" className="bg-foreground/10 text-foreground border-foreground/20">
                      {authFactors.length} factor{authFactors.length !== 1 ? "s" : ""}
                    </Badge>
                  )}
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {authFactors.length === 0 ? (
                <p className="text-sm text-muted-foreground">
                  No multi-factor authentication methods configured
                </p>
              ) : (
                <div className="space-y-2">
                  {authFactors.map((factor: any, i: number) => {
                    const factorType = factor.otp ? "otp" : factor.u2f ? "u2f" : factor.otpSms ? "otpSms" : factor.otpEmail ? "otpEmail" : "unknown"
                    const factorLabels: Record<string, string> = {
                      otp: "TOTP (Authenticator App)",
                      u2f: "U2F (Security Key)",
                      otpSms: "OTP via SMS",
                      otpEmail: "OTP via Email",
                    }
                    const factorIcons: Record<string, typeof Shield> = {
                      otp: Shield,
                      u2f: Fingerprint,
                      otpSms: Monitor,
                      otpEmail: Mail,
                    }
                    const FactorIcon = factorIcons[factorType] ?? Shield
                    const isReady = factor.state === "AUTH_FACTOR_STATE_READY"

                    return (
                      <div key={`${factorType}-${i}`} className="flex items-center justify-between p-3 rounded-md border">
                        <div className="flex items-center gap-3">
                          <FactorIcon className="h-4 w-4 text-muted-foreground" />
                          <div>
                            <p className="text-sm font-medium">{factorLabels[factorType] ?? factorType}</p>
                            {factorType === "u2f" && factor.u2f?.name && (
                              <p className="text-xs text-muted-foreground">{factor.u2f.name}</p>
                            )}
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <StatusBadge variant={isReady ? "active" : "inactive"}>
                            {isReady ? "Active" : "Inactive"}
                          </StatusBadge>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="text-destructive hover:text-destructive hover:bg-destructive/10 h-8"
                            disabled={isSecurityLoading}
                            onClick={() => {
                              startSecurityTransition(async () => {
                                try {
                                  if (factorType === "otp") await removeTOTP(userId)
                                  else if (factorType === "otpSms") await removeOTPSMS(userId)
                                  else if (factorType === "otpEmail") await removeOTPEmail(userId)
                                  else if (factorType === "u2f" && factor.u2f?.id) await removeU2F(userId, factor.u2f.id)
                                  const [updatedMethods, updatedFactors] = await Promise.all([
                                    listAuthMethodTypes(userId),
                                    listAuthFactors(userId),
                                  ])
                                  setAuthMethods(updatedMethods)
                                  setAuthFactors(updatedFactors)
                                  setActionFeedback("Factor removed")
                                  setTimeout(() => setActionFeedback(null), 3000)
                                } catch (e) {
                                  setActionFeedback("Failed to remove factor")
                                  setTimeout(() => setActionFeedback(null), 3000)
                                }
                              })
                            }}
                          >
                            {isSecurityLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                          </Button>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Passkeys */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Fingerprint className="h-5 w-5 text-muted-foreground" />
                  <CardTitle className="text-base">Passkeys</CardTitle>
                  {passkeys.length > 0 && (
                    <Badge variant="outline" className="bg-foreground/10 text-foreground border-foreground/20">
                      {passkeys.length}
                    </Badge>
                  )}
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {passkeys.length === 0 ? (
                <p className="text-sm text-muted-foreground">
                  No passkeys registered
                </p>
              ) : (
                <div className="space-y-2">
                  {passkeys.map((pk: any) => {
                    const isReady = pk.state === "AUTH_FACTOR_STATE_READY"
                    return (
                      <div key={pk.id} className="flex items-center justify-between p-3 rounded-md border">
                        <div>
                          <p className="text-sm font-medium">{pk.name || "Unnamed passkey"}</p>
                          <p className="text-xs font-mono text-muted-foreground">{pk.id}</p>
                        </div>
                        <div className="flex items-center gap-2">
                          <StatusBadge variant={isReady ? "active" : "inactive"}>
                            {isReady ? "Active" : pk.state?.replace("AUTH_FACTOR_STATE_", "") ?? "Unknown"}
                          </StatusBadge>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="text-destructive hover:text-destructive hover:bg-destructive/10 h-8"
                            disabled={isSecurityLoading}
                            onClick={() => {
                              startSecurityTransition(async () => {
                                try {
                                  await removePasskey(userId, pk.id)
                                  const [updatedMethods, updatedPasskeys] = await Promise.all([
                                    listAuthMethodTypes(userId),
                                    listPasskeys(userId),
                                  ])
                                  setAuthMethods(updatedMethods)
                                  setPasskeys(updatedPasskeys)
                                  setActionFeedback("Passkey removed")
                                  setTimeout(() => setActionFeedback(null), 3000)
                                } catch (e) {
                                  setActionFeedback("Failed to remove passkey")
                                  setTimeout(() => setActionFeedback(null), 3000)
                                }
                              })
                            }}
                          >
                            {isSecurityLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
                          </Button>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Recovery Codes & IDP */}
          <Card>
            <CardContent className="pt-6 space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <KeyRound className="h-5 w-5 text-muted-foreground" />
                  <div>
                    <p className="font-medium text-sm">Recovery Codes</p>
                    <p className="text-sm text-muted-foreground">
                      {authMethods.includes("AUTHENTICATION_METHOD_TYPE_RECOVERY_CODE")
                        ? "Recovery codes are configured"
                        : "No recovery codes set"}
                    </p>
                  </div>
                </div>
                <StatusBadge variant={authMethods.includes("AUTHENTICATION_METHOD_TYPE_RECOVERY_CODE") ? "active" : "inactive"}>
                  {authMethods.includes("AUTHENTICATION_METHOD_TYPE_RECOVERY_CODE") ? "Active" : "Inactive"}
                </StatusBadge>
              </div>

              {authMethods.includes("AUTHENTICATION_METHOD_TYPE_IDP") && (
                <div className="flex items-center justify-between border-t pt-4">
                  <div className="flex items-center gap-3">
                    <Globe className="h-5 w-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium text-sm">Identity Provider</p>
                      <p className="text-sm text-muted-foreground">Linked via external identity provider</p>
                    </div>
                  </div>
                  <StatusBadge variant="active">
                    Linked
                  </StatusBadge>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Metadata Tab */}
        <TabsContent value="metadata" className="space-y-4">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-base">
                  User Metadata
                  {metadata.length > 0 && (
                    <Badge variant="secondary" className="ml-2 text-xs">{metadata.length}</Badge>
                  )}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Add new metadata form */}
              <div className="flex items-end gap-2">
                <div className="flex-1">
                  <Label htmlFor="meta-key" className="text-xs text-muted-foreground">Key</Label>
                  <Input
                    id="meta-key"
                    value={newMetaKey}
                    onChange={(e) => setNewMetaKey(e.target.value)}
                    placeholder="metadata.key"
                    className="h-8 text-sm"
                  />
                </div>
                <div className="flex-1">
                  <Label htmlFor="meta-value" className="text-xs text-muted-foreground">Value</Label>
                  <Input
                    id="meta-value"
                    value={newMetaValue}
                    onChange={(e) => setNewMetaValue(e.target.value)}
                    placeholder="value"
                    className="h-8 text-sm"
                  />
                </div>
                <Button
                  size="sm"
                  variant="outline"
                  disabled={!newMetaKey.trim() || !newMetaValue.trim() || isMetadataLoading}
                  onClick={() => {
                    startMetadataTransition(async () => {
                      try {
                        await setUserMetadata(userId, newMetaKey.trim(), newMetaValue.trim())
                        const updated = await listUserMetadata(userId)
                        setMetadata(updated)
                        setNewMetaKey("")
                        setNewMetaValue("")
                        setActionFeedback("Metadata added")
                        setTimeout(() => setActionFeedback(null), 3000)
                      } catch (e) {
                        setActionFeedback("Failed to add metadata")
                        setTimeout(() => setActionFeedback(null), 3000)
                      }
                    })
                  }}
                  className="h-8"
                >
                  {isMetadataLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
                </Button>
              </div>

              {/* Metadata table */}
              {metadata.length === 0 ? (
                <div className="py-8 text-center">
                  <KeyRound className="h-10 w-10 text-muted-foreground mx-auto mb-3" />
                  <p className="text-sm text-muted-foreground">No metadata entries</p>
                  <p className="text-xs text-muted-foreground mt-1">Add custom key-value pairs using the form above</p>
                </div>
              ) : (
                <div className="border rounded-md overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b bg-muted/50">
                        <th className="text-left px-3 py-2 font-medium text-muted-foreground">Key</th>
                        <th className="text-left px-3 py-2 font-medium text-muted-foreground">Value</th>
                        <th className="text-left px-3 py-2 font-medium text-muted-foreground">Last Modified</th>
                        <th className="w-10"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {metadata.map((entry) => (
                        <tr key={entry.key} className="border-b last:border-b-0 hover:bg-muted/30">
                          <td className="px-3 py-2 font-mono text-xs">{entry.key}</td>
                          <td className="px-3 py-2 text-xs max-w-[200px] truncate">{entry.value}</td>
                          <td className="px-3 py-2 text-xs text-muted-foreground">{formatDate(entry.changeDate)}</td>
                          <td className="px-3 py-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-6 w-6 p-0 text-destructive hover:text-destructive"
                              disabled={isMetadataLoading}
                              onClick={() => {
                                startMetadataTransition(async () => {
                                  try {
                                    await deleteUserMetadata(userId, [entry.key])
                                    const updated = await listUserMetadata(userId)
                                    setMetadata(updated)
                                    setActionFeedback("Metadata deleted")
                                    setTimeout(() => setActionFeedback(null), 3000)
                                  } catch (e) {
                                    setActionFeedback("Failed to delete metadata")
                                    setTimeout(() => setActionFeedback(null), 3000)
                                  }
                                })
                              }}
                            >
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Activity Tab */}
        <TabsContent value="activity" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Activity Log</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground text-sm">
                Activity log coming soon. This will show recent authentication events,
                profile changes, and other user activity.
              </p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
