import { applyCustomHeaders } from "@/lib/custom-headers";
import { createLogger } from "@/lib/logger";
import { NextResponse } from "next/server";

const logger = createLogger("readiness");

const READINESS_TIMEOUT_MS = 5_000;

/**
 * Readiness probe.
 *
 * Reports whether the ZITADEL API is reachable and ready by calling its own
 * unauthenticated readiness endpoint. This endpoint is served outside of the
 * API's access-log and quota interceptors, so probe traffic is never counted
 * as authenticated requests of the instance.
 *
 * The configured API credentials are validated once at process startup
 * (see instrumentation.ts), not on every probe.
 */
export async function GET() {
  const apiUrl = process.env.ZITADEL_API_URL;
  if (!apiUrl) {
    return unavailable();
  }

  try {
    const headers = new Headers();
    applyCustomHeaders({
      set: (key, value) => headers.set(key, value),
      remove: (key) => headers.delete(key),
    });

    const response = await fetch(new URL("/debug/ready", apiUrl), {
      method: "GET",
      headers,
      cache: "no-store",
      signal: AbortSignal.timeout(READINESS_TIMEOUT_MS),
    });
    // consume the body so the connection can be reused
    await response.arrayBuffer().catch(() => undefined);

    if (!response.ok) {
      throw new Error(`ZITADEL API readiness check returned HTTP ${response.status}`);
    }

    return new NextResponse("OK", {
      status: 200,
      headers: { "Content-Type": "text/plain", "Cache-Control": "no-store" },
    });
  } catch (e) {
    logger.error("Readiness check failed", { error: e });
    return unavailable();
  }
}

function unavailable() {
  return new NextResponse("Service unavailable", {
    status: 503,
    headers: { "Content-Type": "text/plain", "Cache-Control": "no-store" },
  });
}
