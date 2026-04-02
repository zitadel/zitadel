import { ConsoleLayout } from '@/components/layout/console-layout'
import { getInstances } from '@/lib/instances'
import { isInstanceConfigured } from '@zitadel/react/api/transport'
import { AppProvider } from '@zitadel/react/context/app-context'
import { PermissionProvider } from '@zitadel/react/context/permissions'
import { DeploymentProvider } from '@zitadel/react/context/deployment'
import { Toaster } from '@zitadel/react/components/ui/toaster'
import { discoverUserRoles } from '@zitadel/react/api/auth'
import { listOrganizations } from '@zitadel/react/api/organizations'
import { toJson } from '@zitadel/client'
import { ListOrganizationsResponseSchema } from '@zitadel/proto/zitadel/org/v2/org_service_pb'

/**
 * Console layout — wraps /console/* routes.
 * IMPORTANT: All imports use @/ (cloud's own modules) because the
 * re-exported console page components also resolve @/ to cloud.
 * Using @zitadel/react/* here would create separate React Context objects.
 */
export default async function ConsoleRouteLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const instances = getInstances()
  const configured = isInstanceConfigured()

  let roles: string[] = []
  let orgs: any[] = []

  if (configured) {
    try {
      const [userRoles, orgsResponse] = await Promise.all([
        discoverUserRoles(),
        listOrganizations({ pageSize: 10 })
          .then((res) => {
            const json = toJson(ListOrganizationsResponseSchema, res) as any
            return (json.result ?? []).map((org: any) => ({
              id: org.id ?? "",
              name: org.name ?? "",
              primaryDomain: org.primaryDomain ?? "",
              isDefault: false,
            }))
          })
          .catch((e) => {
            console.error("Failed to load organizations:", e)
            return []
          }),
      ])
      roles = userRoles
      orgs = orgsResponse
    } catch (e) {
      console.error("Failed to initialize console context:", e)
    }
  }

  return (
    <DeploymentProvider>
      <PermissionProvider initialRoles={roles}>
        <AppProvider initialOrganizations={orgs}>
          <ConsoleLayout instances={instances.map(i => ({ id: i.id, name: i.name, url: i.url }))}>
            {children}
          </ConsoleLayout>
          <Toaster />
        </AppProvider>
      </PermissionProvider>
    </DeploymentProvider>
  )
}
