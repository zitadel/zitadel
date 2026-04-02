import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { AppProvider } from '@zitadel/react/context/app-context'
import { PermissionProvider } from '@zitadel/react/context/permissions'
import { DeploymentProvider } from '@zitadel/react/context/deployment'
import { ConsoleLayout } from '@zitadel/react/components/layout/console-layout'
import { ErrorBoundary } from '@zitadel/react/components/error-boundary'
import { Toaster } from '@zitadel/react/components/ui/toaster'
import { ConsoleLinkProvider } from '@zitadel/react/context/link-context'
import { discoverUserRoles } from '@zitadel/react/api/auth'
import { listOrganizations } from '@zitadel/react/api/organizations'
import { toJson } from '@zitadel/client'
import { ListOrganizationsResponseSchema } from '@zitadel/proto/zitadel/org/v2/org_service_pb'

const _geist = Geist({ subsets: ["latin"] })
const _geistMono = Geist_Mono({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: 'ZITADEL Console',
  description: 'IAM Admin & Management Console for ZITADEL',
}

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  // Fetch roles and organizations in parallel
  const [roles, orgs] = await Promise.all([
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

  return (
    <html lang="en" suppressHydrationWarning>
      <body className="font-sans antialiased" suppressHydrationWarning>
        <ErrorBoundary>
          <DeploymentProvider>
            <PermissionProvider initialRoles={roles}>
              <AppProvider initialOrganizations={orgs}>
                <ConsoleLayout>
                  <ConsoleLinkProvider base="">
                    {children}
                  </ConsoleLinkProvider>
                </ConsoleLayout>
                <Toaster />
              </AppProvider>
            </PermissionProvider>
          </DeploymentProvider>
        </ErrorBoundary>
      </body>
    </html>
  )
}
