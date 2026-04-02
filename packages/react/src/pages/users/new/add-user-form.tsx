"use client"

import { useState, useTransition } from "react"
import { useConsoleRouter } from "../../../hooks/use-console-router"
import { ConsoleLink as Link } from "../../../context/link-context"
import { ArrowLeft, Loader2, UserPlus, Eye, EyeOff } from "lucide-react"
import { Button } from "../../../components/ui/button"
import { Input } from "../../../components/ui/input"
import { Label } from "../../../components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../../../components/ui/select"
import { Checkbox } from "../../../components/ui/checkbox"
import { createUser } from "../../../api/create-user"

interface AddUserFormProps {
  organizations: any[]
}

export function AddUserForm({ organizations }: AddUserFormProps) {
  const router = useConsoleRouter()
  const [isPending, startTransition] = useTransition()
  const [error, setError] = useState<string | null>(null)
  const [showPassword, setShowPassword] = useState(false)

  // Form state
  const [organizationId, setOrganizationId] = useState(organizations[0]?.id ?? "")
  const [givenName, setGivenName] = useState("")
  const [familyName, setFamilyName] = useState("")
  const [email, setEmail] = useState("")
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [isEmailVerified, setIsEmailVerified] = useState(false)
  const [requirePasswordChange, setRequirePasswordChange] = useState(true)

  const canSubmit = organizationId && givenName && familyName && email

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!canSubmit) return
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
        // Redirect to the newly created user
        router.push(`/users/${result.userId}`)
      } catch (e) {
        setError(e instanceof Error ? e.message : "Failed to create user")
      }
    })
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/users">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Add User</h1>
          <p className="text-sm text-muted-foreground">
            Create a new human user in your ZITADEL instance
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Organization */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Organization</CardTitle>
            <CardDescription>
              Select the organization this user belongs to
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Select value={organizationId} onValueChange={setOrganizationId}>
              <SelectTrigger>
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
          </CardContent>
        </Card>

        {/* Profile */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Profile</CardTitle>
            <CardDescription>
              Basic user information
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="givenName">First Name *</Label>
                <Input
                  id="givenName"
                  placeholder="Jane"
                  value={givenName}
                  onChange={(e) => setGivenName(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="familyName">Last Name *</Label>
                <Input
                  id="familyName"
                  placeholder="Doe"
                  value={familyName}
                  onChange={(e) => setFamilyName(e.target.value)}
                  required
                />
              </div>
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
                If not set, the email address will be used as the username
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Email */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Email</CardTitle>
            <CardDescription>
              The user&apos;s email address for login and notifications
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email Address *</Label>
              <Input
                id="email"
                type="email"
                placeholder="jane.doe@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox
                id="isEmailVerified"
                checked={isEmailVerified}
                onCheckedChange={(checked) => setIsEmailVerified(checked === true)}
              />
              <Label htmlFor="isEmailVerified" className="text-sm font-normal">
                Mark email as already verified (skip verification email)
              </Label>
            </div>
          </CardContent>
        </Card>

        {/* Password */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Password</CardTitle>
            <CardDescription>
              Optionally set an initial password for the user
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
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
                If left blank, the user will need to set their password via email
              </p>
            </div>
            {password && (
              <div className="flex items-center space-x-2">
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
          </CardContent>
        </Card>

        {/* Error */}
        {error && (
          <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4">
            <p className="text-sm font-medium text-destructive">
              Failed to create user
            </p>
            <p className="text-xs text-muted-foreground mt-1">{error}</p>
          </div>
        )}

        {/* Actions */}
        <div className="flex items-center justify-end gap-3">
          <Button variant="outline" asChild>
            <Link href="/users">Cancel</Link>
          </Button>
          <Button type="submit" disabled={!canSubmit || isPending}>
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
        </div>
      </form>
    </div>
  )
}
