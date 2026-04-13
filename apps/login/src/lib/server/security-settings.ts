import { PromiseCache } from "@/lib/cache";
import { createLogger } from "@/lib/logger";

const logger = createLogger("security-settings");

/**
 * Cache TTL for security settings: 15 minutes (matching default in zitadel.ts).
 * Uses the same PromiseCache (LRU + stale-while-revalidate) as the rest of
 * the app for consistent caching behavior and deduplication of concurrent requests.
 */
const SECURITY_SETTINGS_TTL_MS = 15 * 60 * 1000;

const securityCache = new PromiseCache(100);

/**
 * Resolves an authentication token from available credential sources.
 * Priority: system user JWT > login client key > service account token.
 */
async function resolveAuthToken(): Promise<string> {
  const { hasSystemUserCredentials, hasLoginClientKey, hasServiceUserToken } = await import("@/lib/deployment");

  if (hasSystemUserCredentials()) {
    const { systemAPIToken } = await import("@/lib/api");
    return systemAPIToken();
  }
  if (hasLoginClientKey()) {
    const { loginClientKeyToken } = await import("@/lib/api");
    return loginClientKeyToken();
  }
  if (hasServiceUserToken()) {
    return process.env.ZITADEL_SERVICE_USER_TOKEN!;
  }
  throw new Error("No authentication credentials found");
}

/**
 * Fetches iframe origins from security settings using the ZITADEL API directly
 * via the Connect protocol (POST + JSON). This avoids the HTTPS self-loopback
 * through the load balancer that caused TLS errors on Cloud Run.
 *
 * Results are cached using the LRU PromiseCache with stale-while-revalidate
 * semantics — stale data is served immediately while a background refresh runs.
 *
 * @param baseUrl - The ZITADEL API base URL (ZITADEL_API_URL)
 * @param instanceHost - Optional instance host for multi-tenant deployments
 * @returns An array of allowed iframe origins, or undefined if not configured
 */
export async function getIframeOrigins(baseUrl: string, instanceHost?: string): Promise<string[] | undefined> {
  const cacheKey = `security-settings:${instanceHost || "__default__"}`;

  return securityCache.getOrFetch(cacheKey, () => fetchIframeOrigins(baseUrl, instanceHost), SECURITY_SETTINGS_TTL_MS);
}

async function fetchIframeOrigins(baseUrl: string, instanceHost?: string): Promise<string[] | undefined> {
  const token = await resolveAuthToken();

  const reqHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  if (instanceHost) {
    reqHeaders["x-zitadel-instance-host"] = instanceHost;
  }

  // Apply custom headers from environment, matching the pattern in zitadel.ts
  if (process.env.CUSTOM_REQUEST_HEADERS) {
    process.env.CUSTOM_REQUEST_HEADERS.split(",").forEach((header) => {
      const kv = header.indexOf(":");
      if (kv > 0) {
        const key = header.slice(0, kv).trim();
        const value = header.slice(kv + 1).trim();
        if (value) {
          reqHeaders[key] = value;
        }
      }
    });
  }

  const response = await fetch(`${baseUrl}/zitadel.settings.v2.SettingsService/GetSecuritySettings`, {
    method: "POST",
    headers: reqHeaders,
    body: "{}",
  });

  if (!response.ok) {
    logger.error("Failed to fetch security settings from API", {
      status: response.status,
      statusText: response.statusText,
    });
    return undefined;
  }

  const data = await response.json();
  const settings = data?.settings;

  return settings?.embeddedIframe?.enabled && settings.embeddedIframe.allowedOrigins?.length > 0
    ? settings.embeddedIframe.allowedOrigins
    : undefined;
}
