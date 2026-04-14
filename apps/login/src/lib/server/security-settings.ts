import { createLogger } from "@/lib/logger";

const logger = createLogger("security-settings");

/** Cache TTL: 1 hour. Security settings rarely change. */
const CACHE_TTL_MS = 60 * 60 * 1000;

interface CacheEntry {
  origins: string[] | undefined;
  expiresAt: number;
}

const cache = new Map<string, CacheEntry>();

/** In-flight promises per cache key to deduplicate concurrent fetches. */
const inflight = new Map<string, Promise<string[] | undefined>>();

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
 * Results are cached in-memory for 1 hour per instance host. Concurrent
 * requests that arrive while a fetch is in-flight share the same promise
 * to prevent thundering-herd stampedes on the backend.
 *
 * @param baseUrl - The ZITADEL API base URL (ZITADEL_API_URL)
 * @param instanceHost - Optional instance host for multi-tenant deployments
 * @returns An array of allowed iframe origins, or undefined if not configured
 */
export async function getIframeOrigins(baseUrl: string, instanceHost?: string): Promise<string[] | undefined> {
  const cacheKey = instanceHost || "__default__";
  const cached = cache.get(cacheKey);

  if (cached && cached.expiresAt > Date.now()) {
    return cached.origins;
  }

  // Deduplicate concurrent fetches for the same key
  const existing = inflight.get(cacheKey);
  if (existing) {
    return existing;
  }

  const promise = fetchIframeOrigins(baseUrl, instanceHost).then((origins) => {
    cache.set(cacheKey, { origins, expiresAt: Date.now() + CACHE_TTL_MS });
    inflight.delete(cacheKey);
    return origins;
  }).catch((err) => {
    inflight.delete(cacheKey);
    throw err;
  });

  inflight.set(cacheKey, promise);

  return promise;
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
