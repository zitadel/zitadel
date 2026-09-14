import { hasLoginClientKey, hasServiceUserToken, hasSystemUserCredentials } from "@/lib/deployment";
import { createLogger } from "@/lib/logger";
import { createServiceForHost } from "@/lib/service";
import { Code, ConnectError } from "@connectrpc/connect";
import { Client } from "@zitadel/client";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";

const logger = createLogger("startup");

/** Upper bound for the startup RPC so a hanging API cannot block process startup. */
const CREDENTIAL_CHECK_TIMEOUT_MS = 10_000;

/** gRPC codes that mean the API received the request but rejected the credentials. */
const CREDENTIAL_ERROR_CODES: ReadonlySet<Code> = new Set([Code.Unauthenticated, Code.PermissionDenied]);

export type CredentialCheckResult =
  /** The API accepted the configured credentials. */
  | "ok"
  /** The API rejected the configured credentials. */
  | "rejected"
  /** Credentials are configured but could not be loaded or signed (e.g. missing or malformed key file). */
  | "invalid"
  /** No credentials are configured at all. */
  | "missing"
  /** The API could not be reached (or did not answer in time), so the credentials could not be verified. */
  | "unreachable"
  /** No check was performed: ZITADEL_API_URL is not set, or it does not resolve to an instance. */
  | "skipped";

/** Results that indicate a configuration error the process cannot recover from. */
export const FATAL_CREDENTIAL_CHECK_RESULTS: ReadonlySet<CredentialCheckResult> = new Set<CredentialCheckResult>([
  "rejected",
  "invalid",
  "missing",
]);

/**
 * Verifies once, at process startup, that the configured API credentials are
 * accepted by the ZITADEL API, using the same auth and transport path as real
 * requests.
 *
 * This deliberately runs at startup and not in the readiness probe: probes are
 * executed every few seconds per replica and every authenticated API call is
 * recorded as a billable request of the instance. Verifying once per process
 * catches misconfigured credentials early (fail fast) without generating
 * continuous authenticated traffic.
 *
 * The check is scoped to the instance that ZITADEL_API_URL resolves to. In
 * multi-tenant deployments, where the instance is only known per request, the
 * API answers with NotFound and the check is skipped.
 *
 * The function never throws; the caller decides what to do with the result.
 */
export async function verifyApiCredentials(): Promise<CredentialCheckResult> {
  const apiUrl = process.env.ZITADEL_API_URL;
  if (!apiUrl) {
    logger.warn("ZITADEL_API_URL is not set, skipping API credential check");
    return "skipped";
  }

  if (!hasSystemUserCredentials() && !hasLoginClientKey() && !hasServiceUserToken()) {
    logger.error(
      "No API credentials configured. Set ZITADEL_LOGINCLIENT_KEYFILE, system user credentials (AUDIENCE, SYSTEM_USER_ID, SYSTEM_USER_PRIVATE_KEY or SYSTEM_USER_PRIVATE_KEY_FILE), or ZITADEL_SERVICE_USER_TOKEN",
    );
    return "missing";
  }

  // Loading and signing the credentials happens locally, before any request is
  // sent. Failures here are configuration errors, not connectivity problems.
  let settingsService: Client<typeof SettingsService>;
  try {
    settingsService = await createServiceForHost(SettingsService, { baseUrl: apiUrl });
  } catch (e) {
    logger.error("The configured API credentials could not be loaded", { apiUrl, error: e });
    return "invalid";
  }

  try {
    await settingsService.getGeneralSettings({}, { timeoutMs: CREDENTIAL_CHECK_TIMEOUT_MS });
    logger.info("API credentials verified", { apiUrl });
    return "ok";
  } catch (e) {
    if (e instanceof ConnectError) {
      if (CREDENTIAL_ERROR_CODES.has(e.code)) {
        logger.error("The ZITADEL API rejected the configured credentials", { apiUrl, error: e });
        return "rejected";
      }
      if (e.code === Code.NotFound) {
        logger.warn(
          "ZITADEL_API_URL does not resolve to an instance, skipping API credential check (expected for multi-tenant deployments)",
          { apiUrl, error: e },
        );
        return "skipped";
      }
    }
    logger.warn("Could not verify API credentials because the ZITADEL API is not reachable", { apiUrl, error: e });
    return "unreachable";
  }
}
