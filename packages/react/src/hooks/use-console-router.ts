"use client"

import { useRouter } from "next/navigation"
import { useConsoleBase } from "../context/link-context"
import { useCallback, useMemo } from "react"

/**
 * A wrapper around next/navigation's useRouter that prepends the console base path.
 * Use this instead of useRouter() for all programmatic navigation in console pages.
 */
export function useConsoleRouter() {
  const router = useRouter()
  const base = useConsoleBase()

  const push = useCallback(
    (href: string, options?: Parameters<typeof router.push>[1]) => {
      const resolved = href.startsWith("/") ? `${base}${href}` : href
      return router.push(resolved, options)
    },
    [router, base],
  )

  const replace = useCallback(
    (href: string, options?: Parameters<typeof router.replace>[1]) => {
      const resolved = href.startsWith("/") ? `${base}${href}` : href
      return router.replace(resolved, options)
    },
    [router, base],
  )

  return useMemo(
    () => ({
      ...router,
      push,
      replace,
    }),
    [router, push, replace],
  )
}
