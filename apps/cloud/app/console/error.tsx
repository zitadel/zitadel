"use client"

import Link from "next/link"
import { Server } from "lucide-react"
import { useEffect } from "react"

/**
 * Error boundary for /console/* routes.
 * Catches transport errors (missing ZITADEL_INSTANCE_URL) and shows
 * a configuration prompt instead of a crash page.
 */
export default function ConsoleError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error("Console error:", error)
  }, [error])

  const isConfigError = error.message?.includes("ZITADEL_INSTANCE_URL") ||
    error.message?.includes("ZITADEL_PAT")

  if (isConfigError) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="max-w-md w-full text-center space-y-4">
          <Server className="h-10 w-10 text-muted-foreground mx-auto" />
          <h2 className="text-xl font-bold">Instance not configured</h2>
          <p className="text-sm text-muted-foreground">
            Set <code className="px-1 py-0.5 rounded bg-muted text-xs">ZITADEL_INSTANCE_URL</code> and{" "}
            <code className="px-1 py-0.5 rounded bg-muted text-xs">ZITADEL_PAT</code> in your{" "}
            <code className="px-1 py-0.5 rounded bg-muted text-xs">.env.local</code> file,
            or configure an instance via the debug page.
          </p>
          <div className="flex items-center justify-center gap-3">
            <Link
              href="/debug"
              className="inline-flex items-center gap-2 rounded-md bg-foreground text-background px-4 py-2 text-sm font-medium"
            >
              Configure Instance
            </Link>
            <button
              onClick={reset}
              className="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex items-center justify-center min-h-[60vh]">
      <div className="max-w-md w-full text-center space-y-4">
        <h2 className="text-xl font-bold">Something went wrong</h2>
        <p className="text-sm text-muted-foreground">{error.message}</p>
        <button
          onClick={reset}
          className="inline-flex items-center gap-2 rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent"
        >
          Try again
        </button>
      </div>
    </div>
  )
}
