export type ServerActionResponse =
  | { redirect: string }
  | { error: string }
  | { samlData: { url: string; fields: Record<string, string> } }
  | undefined
  | null;

const PROXY_PATH_PREFIXES = ["/.well-known/", "/oauth/", "/oidc/", "/idps/callback/", "/saml/", "/assets/"];

function isHttpUrl(url: string): boolean {
  return /^https?:\/\//i.test(url);
}

function getRedirectPath(url: string): string {
  if (!url) {
    return "";
  }

  if (url.startsWith("/")) {
    const queryIndex = url.indexOf("?");
    return queryIndex >= 0 ? url.slice(0, queryIndex) : url;
  }

  if (!isHttpUrl(url)) {
    return "";
  }

  try {
    return new URL(url).pathname;
  } catch {
    return "";
  }
}

export function isCrossOrigin(url: string, currentOrigin?: string): boolean {
  if (!isHttpUrl(url)) {
    return false;
  }

  try {
    const targetOrigin = new URL(url).origin;
    const origin = currentOrigin ?? (typeof window !== "undefined" ? window.location.origin : undefined);
    if (!origin) {
      return true;
    }
    return targetOrigin !== origin;
  } catch {
    return false;
  }
}

export function shouldUseHardNavigation(url: string, currentOrigin?: string): boolean {
  const path = getRedirectPath(url);
  if (PROXY_PATH_PREFIXES.some((prefix) => path.startsWith(prefix))) {
    return true;
  }

  return isCrossOrigin(url, currentOrigin);
}

export function navigateHard(url: string): void {
  if (typeof window !== "undefined") {
    window.location.assign(url);
  }
}

export function handleServerActionResponse(
  response: ServerActionResponse,
  router: { push: (url: string) => void },
  setSamlData: (data: { url: string; fields: Record<string, string> }) => void,
  setError: (error: string) => void,
  hardNavigate: (url: string) => void = navigateHard,
  onRedirectStart?: () => void,
): boolean {
  if (!response) {
    return false;
  }

  if ("redirect" in response && response.redirect) {
    onRedirectStart?.();
    if (shouldUseHardNavigation(response.redirect)) {
      hardNavigate(response.redirect);
    } else {
      router.push(response.redirect);
    }
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
