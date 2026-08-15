import { createLogger } from "@/lib/logger";

const logger = createLogger("custom-headers");

/**
 * Parses the CUSTOM_REQUEST_HEADERS environment variable and applies the
 * headers to a mutable target via the provided `set` and `remove` callbacks.
 *
 * Format: comma-separated "key:value" pairs.
 * - Non-empty values call `set(key, value)`.
 * - Empty values (e.g., "X-Remove-Me:") call `remove(key)` to delete the header.
 * - Malformed entries (no colon) are logged and skipped.
 */
export function applyCustomHeaders(actions: {
  set: (key: string, value: string) => void;
  remove: (key: string) => void;
}): void {
  const raw = process.env.CUSTOM_REQUEST_HEADERS;
  if (!raw) return;

  raw.split(",").forEach((header) => {
    const kv = header.indexOf(":");
    if (kv > 0) {
      const key = header.slice(0, kv).trim();
      const value = header.slice(kv + 1).trim();
      if (value) {
        actions.set(key, value);
      } else {
        actions.remove(key);
      }
    } else {
      logger.warn("Skipping malformed CUSTOM_REQUEST_HEADERS entry", { header });
    }
  });
}
