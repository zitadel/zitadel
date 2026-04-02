import Link from "next/link"
import {
  Server,
  Plus,
  Cloud,
  ChevronRight,
  CircleDot,
  Building2,
  Globe,
} from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@zitadel/react/components/ui/card"
import { Badge } from "@zitadel/react/components/ui/badge"
import { getInstances } from "@/lib/instances"

/**
 * Console root page — All Instances dashboard.
 * Uses the same Card/Stats pattern as the org overview.
 */
export default function ConsolePage() {
  const instances = getInstances()
  const total = instances.length
  const active = instances.length

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Instances</h1>
          <p className="text-muted-foreground">
            Manage your ZITADEL instances across all environments
          </p>
        </div>
        <Link
          href="/debug"
          className="inline-flex items-center gap-2 rounded-lg bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add Instance
        </Link>
      </div>

      {instances.length > 0 ? (
        <>
          {/* Stats Grid — same pattern as org overview */}
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <Card className="transition-all hover:border-foreground hover:shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  Total Instances
                </CardTitle>
                <div className="rounded-lg p-2 bg-muted text-foreground">
                  <Server className="h-4 w-4" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{total}</div>
                <p className="text-xs text-muted-foreground">Configured instances</p>
              </CardContent>
            </Card>

            <Card className="transition-all hover:border-foreground hover:shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  Active
                </CardTitle>
                <div className="rounded-lg p-2 bg-muted text-foreground">
                  <CircleDot className="h-4 w-4" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{active}</div>
                <p className="text-xs text-muted-foreground">Currently running</p>
              </CardContent>
            </Card>

            <Card className="transition-all hover:border-foreground hover:shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  Environment
                </CardTitle>
                <div className="rounded-lg p-2 bg-muted text-foreground">
                  <Cloud className="h-4 w-4" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{total}</div>
                <p className="text-xs text-muted-foreground">Cloud hosted</p>
              </CardContent>
            </Card>
          </div>

          {/* Instance List */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-lg">All Instances</CardTitle>
                  <CardDescription>{total} instance{total !== 1 && "s"} configured</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {instances.map((inst, i) => {
                  let hostname = inst.url
                  try { hostname = new URL(inst.url).hostname } catch {}
                  return (
                    <Link
                      key={i}
                      href={`/console/instances/${inst.id}/overview`}
                      className="flex items-center gap-4 p-3 rounded-lg border hover:bg-muted/50 transition-colors group"
                    >
                      <div className="rounded-lg p-2 bg-muted">
                        <Globe className="h-4 w-4 text-muted-foreground" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium group-hover:text-primary transition-colors">
                          {inst.name || "Unnamed"}
                        </p>
                        <p className="text-xs text-muted-foreground font-mono truncate">{hostname}</p>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge
                          variant="outline"
                          className="bg-emerald-500/10 text-emerald-700 border-emerald-500/30 text-xs"
                        >
                          Active
                        </Badge>
                        <Badge variant="secondary" className="text-xs">Local</Badge>
                      </div>
                      <ChevronRight className="h-4 w-4 text-muted-foreground/50 group-hover:text-foreground transition-colors flex-shrink-0" />
                    </Link>
                  )
                })}
              </div>
            </CardContent>
          </Card>
        </>
      ) : (
        <Card className="border-dashed">
          <CardContent className="py-12 text-center">
            <Server className="h-10 w-10 text-muted-foreground/40 mx-auto mb-3" />
            <h3 className="font-semibold text-lg mb-1">No instances configured</h3>
            <p className="text-sm text-muted-foreground mb-4 max-w-md mx-auto">
              Add a ZITADEL instance to get started. You can connect to a local instance
              for development or a cloud-hosted one.
            </p>
            <Link
              href="/debug"
              className="inline-flex items-center gap-2 rounded-lg bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors"
            >
              Configure Instance
            </Link>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
