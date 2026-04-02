import Link from "next/link"

/**
 * Cloud app root — navigation hub.
 * Links to all sections: Console, Docs, Debug, etc.
 * Will eventually become the marketing landing page.
 */
export default function CloudHomePage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-8 p-8">
      <div className="text-center space-y-2">
        <h1 className="text-4xl font-bold tracking-tight">ZITADEL Cloud</h1>
        <p className="text-muted-foreground text-lg">
          Multi-instance management, billing, and more.
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 max-w-2xl w-full">
        <Link
          href="/console"
          className="rounded-lg border p-6 hover:bg-accent transition-colors"
        >
          <h2 className="font-semibold mb-1">Console</h2>
          <p className="text-sm text-muted-foreground">
            Instance admin — users, organizations, billing
          </p>
        </Link>

        <Link
          href="/debug"
          className="rounded-lg border p-6 hover:bg-accent transition-colors"
        >
          <h2 className="font-semibold mb-1">Debug</h2>
          <p className="text-sm text-muted-foreground">
            Configure test instances for preview testing
          </p>
        </Link>
      </div>
    </div>
  )
}
