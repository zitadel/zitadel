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