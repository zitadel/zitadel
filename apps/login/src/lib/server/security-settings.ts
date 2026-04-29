import { PromiseCache } from "@/lib/cache";
import { applyCustomHeaders } from "@/lib/custom-headers";
import { createLogger } from "@/lib/logger";
import type { GetSecuritySettingsResponseJson } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";

const logger = createLogger("security-settings");

let cacheConfig: Record<string, number> = {};
try {
  if (process.env.API_CACHE_CONFIG) {
    cacheConfig = JSON.parse(process.env.API_CACHE_CONFIG);
  }
} catch (e) {
  console.error("Failed to parse API_CACHE_CONFIG", e);
}

/** Cache TTL: defaults to longMinutes (1 hour). Security settings rarely change. */
const CACHE_TTL_MS = (cacheConfig.longMinutes ?? 60) * 60 * 1000;

/**
 * Bounded LRU cache for security settings per instance host.
 *
 * Uses the shared PromiseCache which provides:
 * - LRU eviction (capped by API_CACHE_CONFIG maxSize, default 10 000)
 * - Automatic TTL expiry (API_CACHE_CONFIG longMinutes, default 60 min)
 * - Stale-while-revalidate (serves stale data while refreshing in background)
 * - Built-in in-flight request deduplication
 */
const cache = new PromiseCache(Number(cacheConfig.maxSize) || 10_000);

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
 * via the Connect protocol (POST + JSON). This uses raw fetch (no connectRPC
 * node transport) so it stays compatible with the Next.js proxy runtime.
 *
 * Results are cached in-memory for 1 hour per instance host using a bounded
 * LRU cache. Concurrent requests for the same key share a single in-flight
 * promise to prevent thundering-herd stampedes on the backend.
 *
 * @param baseUrl - The ZITADEL API base URL (ZITADEL_API_URL)
 * @param instanceHost - Optional instance host for multi-tenant deployments
 * @returns An array of allowed iframe origins, or undefined if not configured
 */
export async function getIframeOrigins(baseUrl: string, instanceHost?: string): Promise<string[] | undefined> {
  const cacheKey = instanceHost || "__default__";

  return cache.getOrFetch<string[] | undefined>(cacheKey, () => fetchIframeOrigins(baseUrl, instanceHost), CACHE_TTL_MS);
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

  // Apply custom headers from environment
  applyCustomHeaders({
    set: (key, value) => {
      reqHeaders[key] = value;
    },
    remove: (key) => {
      delete reqHeaders[key];
    },
  });

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

  const data: GetSecuritySettingsResponseJson = await response.json();
  const settings = data.settings;

  const origins = settings?.embeddedIframe?.enabled ? settings.embeddedIframe.allowedOrigins : undefined;

  return origins && origins.length > 0 ? origins : undefined;
}
