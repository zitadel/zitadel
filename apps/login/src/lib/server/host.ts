import { ReadonlyHeaders } from "next/dist/server/web/spec-extension/adapters/headers";

const INVALID_PUBLIC_URL_ERROR =
  "ZITADEL_PUBLIC_URL must be a valid absolute URL, for example: https://login.example.com";

export function getConfiguredPublicURL(): URL | null {
  const configuredPublicURL = process.env.ZITADEL_PUBLIC_URL?.trim();

  if (!configuredPublicURL) {
    return null;
  }

  try {
    const publicURL = new URL(configuredPublicURL);

    const isValidOrigin =
      (publicURL.protocol === "http:" || publicURL.protocol === "https:") &&
      publicURL.host !== "" &&
      publicURL.pathname === "/" &&
      publicURL.search === "" &&
      publicURL.hash === "" &&
      publicURL.username === "" &&
      publicURL.password === "";

    if (!isValidOrigin) {
      throw new Error(INVALID_PUBLIC_URL_ERROR);
    }

    return publicURL;
  } catch {
    throw new Error(INVALID_PUBLIC_URL_ERROR);
  }
}

function getPublicHostFromHeaders(headers: ReadonlyHeaders): string {
  const publicHost =
    headers.get("x-zitadel-public-host") ||
    headers.get("x-zitadel-forward-host") ||
    headers.get("x-forwarded-host") ||
    headers.get("host");

  if (!publicHost || typeof publicHost !== "string") {
    throw new Error("No host found in headers");
  }

  return publicHost;
}

/**
 * Gets the original host that the user sees in their browser URL.
 * When using rewrites this function prioritizes forwarded headers that preserve the original host.
 *
 * @returns The host string (e.g., "zitadel.com")
 * @throws Error if no host is found
 */
export function getInstanceHost(headers: ReadonlyHeaders): string | null {
  // use standard proxy headers (x-forwarded-host → host) for both multi-tenant and self-hosted, do not use x-zitadel-instance-host
  const instanceHost = headers.get("x-zitadel-instance-host") || headers.get("x-zitadel-forward-host");

  return instanceHost;
}

/**
 * Gets the public host that the user sees in their browser URL.
 * Only considers standard proxy headers (x-forwarded-host and host).
 * Does NOT include x-zitadel-instance-host.
 *
 * Use this when you need the public-facing host that the user actually sees,
 * not the internal instance host used for API routing.
 *
 * @returns The public host string (e.g., "accounts.company.com")
 * @throws Error if no host is found
 */
export function getPublicHost(headers: ReadonlyHeaders): string {
  const configuredPublicURL = getConfiguredPublicURL();
  if (configuredPublicURL) {
    return configuredPublicURL.host;
  }

  // Only use standard proxy headers (x-zitadel-public-host → x-zitadel-forward-host → x-forwarded-host → host)
  // Do NOT use x-zitadel-instance-host as it may differ from what the user sees
  return getPublicHostFromHeaders(headers);
}

export function getPublicOrigin(headers: ReadonlyHeaders, fallbackProtocol?: string): string {
  const configuredPublicURL = getConfiguredPublicURL();
  if (configuredPublicURL) {
    return configuredPublicURL.origin;
  }

  const host = getPublicHostFromHeaders(headers);
  const protocol = fallbackProtocol ?? (host.includes("localhost") ? "http:" : "https:");
  const normalizedProtocol = protocol.endsWith(":") ? protocol : `${protocol}:`;

  return `${normalizedProtocol}//${host}`;
}

export function getPublicHostWithProtocol(headers: ReadonlyHeaders): string {
  return getPublicOrigin(headers);
}
