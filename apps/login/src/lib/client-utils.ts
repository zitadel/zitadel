import { sanitizeUrl } from "@braintree/sanitize-url";

export type ServerActionResponse =
  | { redirect: string }
  | { error: string }
  | { samlData: { url: string; fields: Record<string, string> } }
  | undefined
  | null;

export function handleServerActionResponse(
  response: ServerActionResponse,
  router: { push: (url: string) => void },
  setError: (error: string) => void,
): boolean {
  if (!response) {
    return false;
  }

  if ("redirect" in response && response.redirect) {
    if (isSafeRedirectUri(response.redirect)) {
      if (typeof window !== "undefined" && isExternalUrl(response.redirect)) {
        // External/custom-protocol URLs: use window.location for full navigation.
        // router.push() would trigger an RSC prefetch fetch() that gets blocked
        // by CSP connect-src 'self' for non-same-origin URLs.
        // CodeQL: This is safe — the URL is validated by isSafeRedirectUri above,
        // which blocks javascript:, data:, file:, blob:, and about: schemes.
        window.location.href = response.redirect; // lgtm[js/client-side-unvalidated-url-redirection]
      } else {
        router.push(response.redirect);
      }
      return true;
    } else {
      console.warn("handleServerActionResponse: Blocked unsafe redirect URI:", response.redirect);
      setError("Unsafe redirect URI was blocked");
      return true;
    }
  }

  if ("samlData" in response && response.samlData) {
    if (isSafeRedirectUri(response.samlData.url)) {
      submitSamlForm(response.samlData);
      return true;
    } else {
      console.warn("handleServerActionResponse: Blocked unsafe SAML post URL:", response.samlData.url);
      setError("Unsafe redirect URI was blocked");
      return true;
    }
  }

  if ("error" in response && response.error) {
    setError(response.error);
    return true;
  }

  return false;
}

/**
 * Posts a SAML response to the service provider (POST binding) and leaves the page.
 *
 * The form is built on the document, not rendered by a component: the form that started
 * the login can be gone by the time the response is handled (the page re-rendered into
 * something else), and a state update on an unmounted component would post nothing.
 */
export function submitSamlForm({ url, fields }: { url: string; fields: Record<string, string> }): void {
  const form = document.createElement("form");
  form.method = "post";
  form.action = url;
  form.hidden = true;

  for (const [name, value] of Object.entries(fields)) {
    const input = document.createElement("input");
    input.type = "hidden";
    input.name = name;
    input.value = value;
    form.appendChild(input);
  }

  document.body.appendChild(form);
  form.submit();
}

/**
 * Returns true if the URL is external (absolute URL or custom protocol).
 * Relative paths (starting with "/" but not "//") are considered internal.
 */
export function isExternalUrl(url: string): boolean {
  if (url.startsWith("/") && !url.startsWith("//")) {
    return false;
  }
  return true;
}

const SANITIZE_BLANK = "about:blank";
const EXTRA_BLOCKED = new Set(["file:", "blob:", "about:"]);

/** Allowlist for the back link on the IdP error pages. */
export function isAllowedPostErrorRedirectUrl(url: string | undefined): url is "/loginname" | "/idp" {
  return url === "/loginname" || url === "/idp";
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

  const sanitized = sanitizeUrl(uri);
  if (sanitized === SANITIZE_BLANK) return false;

  try {
    const parsed = new URL(sanitized);
    return !EXTRA_BLOCKED.has(parsed.protocol);
  } catch {
    return false;
  }
}
