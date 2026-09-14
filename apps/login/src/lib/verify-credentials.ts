import { hasLoginClientKey, hasServiceUserToken, hasSystemUserCredentials } from "@/lib/deployment";
import { createLogger } from "@/lib/logger";
import { createServiceForHost } from "@/lib/service";
import { Code, ConnectError } from "@connectrpc/connect";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";

const logger = createLogger("startup");

/** gRPC codes that mean the API received the request but rejected the credentials. */
const CREDENTIAL_ERROR_CODES: ReadonlySet<Code> = new Set([Code.Unauthenticated, Code.PermissionDenied]);

export type CredentialCheckResult =
  /** The API accepted the configured credentials. */
  | "ok"
  /** The API rejected the configured credentials. */
  | "rejected"
  /** No credentials are configured at all. */
  | "missing"
  /** The API could not be reached, so the credentials could not be verified. */
  | "unreachable"
  /** ZITADEL_API_URL is not set, so no check was performed. */
  | "skipped";

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

  try {
    const settingsService = await createServiceForHost(SettingsService, { baseUrl: apiUrl });
    await settingsService.getGeneralSettings({});
    logger.info("API credentials verified", { apiUrl });
    return "ok";
  } catch (e) {
    if (e instanceof ConnectError && CREDENTIAL_ERROR_CODES.has(e.code)) {
      logger.error("The ZITADEL API rejected the configured credentials", { apiUrl, error: e });
      return "rejected";
    }
    logger.warn("Could not verify API credentials because the ZITADEL API is not reachable", { apiUrl, error: e });
    return "unreachable";
  }
}
