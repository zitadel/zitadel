"use client"

/**
 * Debug banner — shows a small persistent bar at the top of every page
 * when running in development or preview mode.
 * 
 * Appears across all routes (console, docs, login, debug).
 */
export function DebugBanner() {
  const vercelEnv = process.env.NEXT_PUBLIC_VERCEL_ENV
  const nodeEnv = process.env.NODE_ENV

  const isPreview = vercelEnv === "preview"
  const isDev = nodeEnv === "development"
  
  if (!isDev && !isPreview) return null

  const label = isPreview ? "Preview" : "Local Dev"

  return (
    <div className="bg-amber-500 text-black px-4 py-1 text-xs flex items-center gap-3 font-medium sticky top-0 z-50">
      <span className="flex items-center gap-1.5">
        <span className="inline-block w-1.5 h-1.5 rounded-full bg-black/40 animate-pulse" />
        {label}
      </span>
      <a
        href="/debug"
        className="ml-auto text-black/60 hover:text-black underline"
      >
        Configure
      </a>
    </div>
  )
}
