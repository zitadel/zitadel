export type ServerActionResponse =
  | { redirect: string }
  | { error: string }
  | { samlData: { url: string; fields: Record<string, string> } }
  | undefined
  | null;

export function handleServerActionResponse(
  response: ServerActionResponse,
  router: { push: (url: string) => void },
  setSamlData: (data: { url: string; fields: Record<string, string> }) => void,
  setError: (error: string) => void,
): boolean {
  if (!response) {
    return false;
  }

  if ("redirect" in response && response.redirect) {
    router.push(response.redirect);
    return true;
  }

  if ("samlData" in response && response.samlData) {
    setSamlData(response.samlData);
    return true;
  }

  if ("error" in response && response.error) {
    setError(response.error);
    return true;
  }

  return false;
}

/**
 * Validates whether a given redirect URI is safe.
 * Safe URIs are either relative paths or absolute URLs matching the current host.
 * This prevents open redirect vulnerabilities and XSS via javascript:/data: URIs.
 */
export function isSafeRedirectUri(uri: string): boolean {
  if (!uri) return false;

  // 1. Relative paths are generally safe
  if (uri.startsWith("/") && !uri.startsWith("//")) {
    return true;
  }

  // 2. Check absolute URLs for safe protocols (http/https)
  // We allow external domains, but strictly forbid javascript:/data: etc.
  try {
    const parsedUri = new URL(uri);

    // Only allow http(s) protocols
    if (parsedUri.protocol !== "http:" && parsedUri.protocol !== "https:") {
      return false;
    }

    return true;
  } catch {
    // If it can't be parsed as a URL and didn't start with /, it's unsafe
    return false;
  }
}

/**
 * Sanitizes a redirect URI by parsing and reconstructing it.
 * Returns the reconstructed URI if safe, or undefined if unsafe.
 *
 * Security policy:
 * - Always allow relative paths starting with a single "/".
 * - For absolute URLs:
 *   - Only allow http/https protocols.
 *   - Reject URLs with embedded credentials (userinfo).
 *   - When trustOrigin is true, allow any http/https URL (for admin-controlled env vars).
 *   - Otherwise, enforce same-origin using allowedOrigin or window.location.origin.
 */
export function sanitizeRedirectUri(uri: string, allowedOrigin?: string, trustOrigin = false): string | undefined {
  if (!uri) return undefined;

  // Relative paths: reconstruct from parsed components to break taint chain
  if (uri.startsWith("/") && !uri.startsWith("//")) {
    try {
      const parsed = new URL(uri, "http://n");
      return parsed.pathname + parsed.search + parsed.hash;
    } catch {
      return undefined;
    }
  }

  // Absolute URLs: parse and validate
  try {
    const parsed = new URL(uri);
    if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
      return undefined;
    }

    // Reject URLs with embedded credentials
    if (parsed.username || parsed.password) {
      return undefined;
    }

    // Reconstruct from components to break taint chain
    const reconstructed = parsed.origin + parsed.pathname + parsed.search + parsed.hash;

    // Admin-controlled origins (e.g., DEFAULT_REDIRECT_URI env var) are trusted
    if (trustOrigin) {
      return reconstructed;
    }

    // Determine the trusted origin: prefer explicit allowedOrigin, fall back to browser origin
    let trustedOrigin: string | undefined;
    if (allowedOrigin) {
      try {
        trustedOrigin = new URL(allowedOrigin).origin;
      } catch {
        // Invalid allowedOrigin, ignore
      }
    } else if (typeof window !== "undefined" && window.location?.origin && window.location.origin !== "null") {
      trustedOrigin = window.location.origin;
    }

    if (trustedOrigin && parsed.origin === trustedOrigin) {
      return reconstructed;
    }

    return undefined;
  } catch {
    return undefined;
  }
}