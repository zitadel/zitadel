import { createConnectTransport } from "@connectrpc/connect-node";
import { NewAuthorizationBearerInterceptor } from "@zitadel/client";

let _transport: ReturnType<typeof createConnectTransport> | null = null;
let _lastBaseUrl: string | undefined;

/**
 * Resolve the instance URL and PAT from environment.
 * Priority: ZITADEL_INSTANCES (first entry) > ZITADEL_INSTANCE_URL + ZITADEL_PAT
 */
function resolveInstanceConfig(): { baseUrl: string; pat: string } | null {
  // Try ZITADEL_INSTANCES first (set by debug page)
  try {
    const raw = process.env.ZITADEL_INSTANCES;
    if (raw) {
      const instances = JSON.parse(raw);
      if (instances.length > 0 && instances[0].url && instances[0].pat) {
        return { baseUrl: instances[0].url, pat: instances[0].pat };
      }
    }
  } catch {}

  // Fallback to explicit env vars
  const baseUrl = process.env.ZITADEL_INSTANCE_URL;
  const pat = process.env.ZITADEL_PAT;
  if (baseUrl && pat) {
    return { baseUrl, pat };
  }

  return null;
}

/**
 * Get or create the connectRPC transport for server-side API calls.
 * Reads instance config from ZITADEL_INSTANCES or ZITADEL_INSTANCE_URL/PAT.
 */
export function getTransport() {
  const config = resolveInstanceConfig();

  if (!config) {
    throw new Error(
      "No ZITADEL instance configured. Add one via the debug page or set ZITADEL_INSTANCES in .env.local."
    );
  }

  if (_transport && _lastBaseUrl === config.baseUrl) {
    return _transport;
  }

  _transport = createConnectTransport({
    baseUrl: config.baseUrl,
    httpVersion: "2",
    interceptors: [
      NewAuthorizationBearerInterceptor(config.pat),
    ],
  });
  _lastBaseUrl = config.baseUrl;

  return _transport;
}

/** Check if an instance is configured */
export function isInstanceConfigured(): boolean {
  return resolveInstanceConfig() !== null;
}

