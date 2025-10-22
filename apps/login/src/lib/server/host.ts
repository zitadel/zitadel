import { headers } from "next/headers";

/**
 * Gets the original host that the user sees in their browser URL.
 * When using rewrites this function prioritizes forwarded headers that preserve the original host.
 *
 * ⚠️ SERVER-SIDE ONLY: This function can only be used in:
 * - Server Actions (functions with "use server")
 * - Server Components (React components that run on the server)
 * - Route Handlers (API routes)
 * - Middleware
 *
 * @returns The host string (e.g., "zitadel.com")
 * @throws Error if no host is found
 */
export async function getOriginalHost(): Promise<string> {
  const _headers = await headers();

  // Priority order:
  // 1. x-forwarded-host - Set by proxies/CDNs with the original host
  // 2. x-original-host - Alternative header sometimes used
  // 3. host - Fallback to the current host header
  const host = _headers.get("x-forwarded-host") || _headers.get("x-original-host") || _headers.get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found in headers");
  }

  return host;
}

/**
 * Gets the original host with protocol prefix.
 * Automatically detects if localhost should use http:// or https://
 *
 * ⚠️ SERVER-SIDE ONLY: This function can only be used in:
 * - Server Actions (functions with "use server")
 * - Server Components (React components that run on the server)
 * - Route Handlers (API routes)
 * - Middleware
 *
 * @returns The full URL prefix (e.g., "https://zitadel.com")
 */
export async function getOriginalHostWithProtocol(): Promise<string> {
  const host = await getOriginalHost();
  const protocol = host.includes("localhost") ? "http://" : "https://";
  return `${protocol}${host}`;
}
