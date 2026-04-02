import { notFound } from "next/navigation"
import { getInstances } from "@/lib/instances"
import { ConsoleLinkProvider } from "@zitadel/react/context/link-context"

/**
 * Instance-scoped layout.
 * Validates the instanceId from the URL and configures the transport
 * to use this instance's credentials.
 *
 * Wraps children with ConsoleLinkProvider so all ConsoleLink components
 * in console pages resolve to /console/instances/{id}/...
 *
 * Cookie-setting is handled by middleware.ts (layouts can't set cookies).
 */
export default async function InstanceLayout({
  children,
  params,
}: {
  children: React.ReactNode
  params: Promise<{ instanceId: string }>
}) {
  const { instanceId } = await params
  const instances = getInstances()
  const instance = instances.find((i) => i.id === instanceId)

  if (!instance) {
    notFound()
  }

  // Set the instance URL and PAT for the transport layer
  process.env.ZITADEL_INSTANCE_URL = instance.url
  process.env.ZITADEL_PAT = instance.pat

  return (
    <ConsoleLinkProvider base={`/console/instances/${instanceId}`}>
      {children}
    </ConsoleLinkProvider>
  )
}
