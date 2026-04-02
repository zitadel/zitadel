"use client"

import React, { createContext, useContext } from "react"
import Link from "next/link"

const ConsoleLinkContext = createContext<string>("")

/**
 * Returns the current console base path.
 * Empty string in standalone console, "/console/instances/{id}" in cloud.
 */
export function useConsoleBase() {
  return useContext(ConsoleLinkContext)
}

/**
 * Provides the base path for all ConsoleLink components.
 * - Standalone console: <ConsoleLinkProvider base="">
 * - Cloud wrapper:      <ConsoleLinkProvider base="/console/instances/{id}">
 */
export function ConsoleLinkProvider({
  base,
  children,
}: {
  base: string
  children: React.ReactNode
}) {
  return (
    <ConsoleLinkContext value={base}>
      {children}
    </ConsoleLinkContext>
  )
}

/**
 * Link component that prepends the console base path to href.
 * Use this instead of next/link's Link for all internal console navigation.
 */
export function ConsoleLink({
  href,
  children,
  ...props
}: React.ComponentProps<typeof Link>) {
  const base = useConsoleBase()
  const resolvedHref =
    typeof href === "string" && href.startsWith("/")
      ? `${base}${href}`
      : href
  return (
    <Link href={resolvedHref} {...props}>
      {children}
    </Link>
  )
}
