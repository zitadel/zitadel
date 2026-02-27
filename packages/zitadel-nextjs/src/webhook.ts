import {
  createWebhookHandler,
  type PayloadType,
} from "@zitadel/zitadel-js/webhooks";

export interface WebhookHandlerOptions {
  /**
   * The payload type. Defaults to `"json"` (HMAC-SHA256 signature).
   * Set to `"jwt"` or `"jwe"` for token-based verification.
   */
  payloadType?: PayloadType;

  // ---- JSON (HMAC) options ----
  /**
   * The shared webhook signing key (required when payloadType is "json").
   * Falls back to the `ZITADEL_WEBHOOK_SECRET` environment variable.
   */
  signingKey?: string;

  // ---- JWT / JWE options ----
  /**
   * JWKS endpoint of the ZITADEL instance (required for "jwt" / "jwe").
   * Falls back to the `ZITADEL_WEBHOOK_JWKS_ENDPOINT` environment variable.
   */
  jwksEndpoint?: string;
  /** Expected JWT issuer. */
  issuer?: string;
  /** Expected JWT audience. */
  audience?: string;
  /**
   * PEM-encoded private key for JWE decryption (required for "jwe").
   * Falls back to the `ZITADEL_WEBHOOK_JWE_PRIVATE_KEY` environment variable.
   */
  privateKey?: string;

  /** Callback invoked with the parsed webhook payload when verification succeeds. */
  onEvent: (event: unknown) => void | Promise<void>;
}

/**
 * Creates a Next.js Route Handler for ZITADEL Actions v2 webhooks.
 *
 * Returns a `POST` handler suitable for use in `app/api/.../route.ts`.
 * Supports all three payload types: JSON (HMAC), JWT, and JWE.
 *
 * @example
 * ```ts
 * // app/api/webhook/route.ts
 * import { createZitadelWebhookHandler } from "@zitadel/nextjs";
 *
 * export const POST = createZitadelWebhookHandler({
 *   signingKey: process.env.ZITADEL_WEBHOOK_SECRET!,
 *   onEvent: async (event) => {
 *     console.log("Received webhook:", event);
 *   },
 * });
 * ```
 */
export function createZitadelWebhookHandler(options: WebhookHandlerOptions) {
  const handler = createWebhookHandler({
    payloadType: options.payloadType,
    signingKey: options.signingKey ?? process.env.ZITADEL_WEBHOOK_SECRET,
    jwksEndpoint: options.jwksEndpoint ?? process.env.ZITADEL_WEBHOOK_JWKS_ENDPOINT,
    issuer: options.issuer,
    audience: options.audience,
    privateKey: options.privateKey ?? process.env.ZITADEL_WEBHOOK_JWE_PRIVATE_KEY,
    onEvent: options.onEvent,
  });

  return async function POST(request: Request): Promise<Response> {
    const body = await request.text();

    // Convert Next.js Request headers to a plain object
    const headers: Record<string, string> = {};
    request.headers.forEach((value, key) => {
      headers[key] = value;
    });

    const result = await handler({ body, headers });

    return new Response(result.body, {
      status: result.status,
      headers: { "Content-Type": "text/plain" },
    });
  };
}
