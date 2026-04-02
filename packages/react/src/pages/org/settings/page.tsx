"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../components/ui/card"
import { Button } from "../../../components/ui/button"
import { Input } from "../../../components/ui/input"
import { Label } from "../../../components/ui/label"
import { Switch } from "../../../components/ui/switch"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../../components/ui/tabs"
import { Separator } from "../../../components/ui/separator"
import { useAppContext } from "../../../context/app-context"
import { Badge } from "../../../components/ui/badge"
import { OrganizationSelectorPrompt } from "../../../components/organization-selector-prompt"

export default function OrgSettingsPage() {
  const { currentOrganization } = useAppContext()

  if (!currentOrganization) {
    return (
      <OrganizationSelectorPrompt 
        title="Select an Organization"
        description="Choose an organization to manage its settings"
        targetPath="/org/settings"
      />
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Settings & Policies</h1>
        <p className="text-muted-foreground">
          Organization settings for {currentOrganization?.name}
        </p>
      </div>

      <Tabs defaultValue="general" className="space-y-4">
        <TabsList>
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="login">Login Policy</TabsTrigger>
          <TabsTrigger value="password">Password Policy</TabsTrigger>
          <TabsTrigger value="branding">Branding</TabsTrigger>
          <TabsTrigger value="domains">Domains</TabsTrigger>
        </TabsList>

        <TabsContent value="general">
          <Card>
            <CardHeader>
              <CardTitle>Organization Details</CardTitle>
              <CardDescription>Basic organization information</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="name">Organization Name</Label>
                <Input id="name" defaultValue={currentOrganization?.name} />
              </div>
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Default Organization</Label>
                  <p className="text-sm text-muted-foreground">
                    Set this as the default organization for new users
                  </p>
                </div>
                <Switch defaultChecked={currentOrganization?.isDefault} />
              </div>
              <Separator />
              <Button>Save Changes</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="login">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Login Policy</CardTitle>
                  <CardDescription>Override instance login settings for this organization</CardDescription>
                </div>
                <Badge variant="secondary">Inherits from Instance</Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Override Instance Policy</Label>
                  <p className="text-sm text-muted-foreground">
                    Enable to customize login settings for this organization
                  </p>
                </div>
                <Switch />
              </div>
              <Separator />
              <div className="flex items-center justify-between opacity-50">
                <div className="space-y-0.5">
                  <Label>Allow Registration</Label>
                  <p className="text-sm text-muted-foreground">
                    Allow new users to self-register
                  </p>
                </div>
                <Switch disabled defaultChecked />
              </div>
              <Separator />
              <div className="flex items-center justify-between opacity-50">
                <div className="space-y-0.5">
                  <Label>External Login</Label>
                  <p className="text-sm text-muted-foreground">
                    Allow login with external identity providers
                  </p>
                </div>
                <Switch disabled defaultChecked />
              </div>
              <Separator />
              <Button disabled>Save Changes</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="password">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Password Policy</CardTitle>
                  <CardDescription>Override instance password requirements</CardDescription>
                </div>
                <Badge variant="secondary">Inherits from Instance</Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Override Instance Policy</Label>
                  <p className="text-sm text-muted-foreground">
                    Enable to customize password settings for this organization
                  </p>
                </div>
                <Switch />
              </div>
              <Separator />
              <div className="grid grid-cols-2 gap-4 opacity-50">
                <div className="space-y-2">
                  <Label htmlFor="minLength">Minimum Length</Label>
                  <Input id="minLength" type="number" defaultValue="8" disabled />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="maxLength">Maximum Length</Label>
                  <Input id="maxLength" type="number" defaultValue="72" disabled />
                </div>
              </div>
              <Separator />
              <Button disabled>Save Changes</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="branding">
          <Card>
            <CardHeader>
              <CardTitle>Branding</CardTitle>
              <CardDescription>Customize the look and feel for this organization</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="logo">Logo URL</Label>
                <Input id="logo" placeholder="https://example.com/logo.png" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="primaryColor">Primary Color</Label>
                <div className="flex gap-2">
                  <Input id="primaryColor" defaultValue="#6366f1" className="w-32" />
                  <div className="w-10 h-10 rounded border bg-primary" />
                </div>
              </div>
              <Separator />
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Custom Login Screen</Label>
                  <p className="text-sm text-muted-foreground">
                    Use organization branding on login pages
                  </p>
                </div>
                <Switch />
              </div>
              <Separator />
              <Button>Save Changes</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="domains">
          <Card>
            <CardHeader>
              <CardTitle>Verified Domains</CardTitle>
              <CardDescription>Domains verified for this organization</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-2">
                <Label>Add Domain</Label>
                <div className="flex gap-2">
                  <Input placeholder="example.com" />
                  <Button>Add</Button>
                </div>
              </div>
              <Separator />
              <div className="space-y-2">
                <Label>Verified Domains</Label>
                <p className="text-sm text-muted-foreground">No domains verified yet</p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
