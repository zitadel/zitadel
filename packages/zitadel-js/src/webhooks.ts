import { timingSafeEqual, createHmac } from "node:crypto";
import {
  createRemoteJWKSet,
  jwtVerify,
  compactDecrypt,
  importPKCS8,
  type JWTPayload,
} from "jose";

// ---------------------------------------------------------------------------
// Payload types matching ZITADEL Actions v2 target configuration
// ---------------------------------------------------------------------------

/** The three payload/signing modes ZITADEL supports for action targets. */
export type PayloadType = "json" | "jwt" | "jwe";

// ---------------------------------------------------------------------------
// JSON (HMAC-SHA256) verification — existing behaviour
// ---------------------------------------------------------------------------

export interface WebhookVerifyOptions {
  /** The raw request body as a string or Buffer. */
  body: string | Buffer;
  /** The signature header value from the incoming request. */
  signature: string;
  /** The shared webhook signing key. */
  signingKey: string;
}

/**
 * Verifies a ZITADEL webhook signature using HMAC-SHA256.
 * Returns true if the signature is valid.
 */
export function verifyZitadelWebhook(options: WebhookVerifyOptions): boolean {
  const { body, signature, signingKey } = options;
  const expected = createHmac("sha256", signingKey)
    .update(body)
    .digest("hex");

  const sigBuffer = Buffer.from(signature, "hex");
  const expectedBuffer = Buffer.from(expected, "hex");

  if (sigBuffer.length !== expectedBuffer.length) {
    return false;
  }

  return timingSafeEqual(sigBuffer, expectedBuffer);
}

// ---------------------------------------------------------------------------
// JWT verification — PAYLOAD_TYPE_JWT
// ---------------------------------------------------------------------------

export interface JwtWebhookOptions {
  /** The raw JWT string received as the request body. */
  token: string;
  /**
   * JWKS endpoint of the ZITADEL instance, e.g.
   * `https://my.zitadel.cloud/oauth/v2/keys` or the webkeys endpoint.
   */
  jwksEndpoint: string;
  /** Expected JWT issuer (the ZITADEL instance URL). */
  issuer?: string;
  /** Expected JWT audience. */
  audience?: string;
}

/**
 * Verifies a JWT-signed webhook payload against ZITADEL's JWKS.
 * Returns the decoded payload on success; throws on failure.
 */
export async function verifyJwtWebhook<T = Record<string, unknown>>(
  options: JwtWebhookOptions,
): Promise<T & JWTPayload> {
  const { token, jwksEndpoint, issuer, audience } = options;
  const JWKS = createRemoteJWKSet(new URL(jwksEndpoint));

  const { payload } = await jwtVerify(token, JWKS, { issuer, audience });
  return payload as T & JWTPayload;
}

// ---------------------------------------------------------------------------
// JWE decryption + JWT verification — PAYLOAD_TYPE_JWE
// ---------------------------------------------------------------------------

export interface JweWebhookOptions {
  /** The raw JWE compact-serialised string received as the request body. */
  token: string;
  /**
   * PEM-encoded private key used to decrypt the JWE outer layer.
   * Must correspond to the public key uploaded to the ZITADEL target.
   */
  privateKey: string;
  /**
   * JWKS endpoint of the ZITADEL instance for verifying the inner JWT signature.
   */
  jwksEndpoint: string;
  /** Expected JWT issuer (the ZITADEL instance URL). */
  issuer?: string;
  /** Expected JWT audience. */
  audience?: string;
}

/**
 * Decrypts a JWE-encrypted webhook payload and verifies the inner JWT.
 * Returns the decoded payload on success; throws on failure.
 */
export async function decryptJweWebhook<T = Record<string, unknown>>(
  options: JweWebhookOptions,
): Promise<T & JWTPayload> {
  const { token, privateKey, jwksEndpoint, issuer, audience } = options;

  // Decrypt the outer JWE layer
  const key = await importPKCS8(privateKey, "RSA-OAEP-256");
  const { plaintext } = await compactDecrypt(token, key);

  // The plaintext is the inner JWT
  const innerJwt = new TextDecoder().decode(plaintext);

  // Verify the inner JWT against ZITADEL's JWKS
  return verifyJwtWebhook<T>({
    token: innerJwt,
    jwksEndpoint,
    issuer,
    audience,
  });
}

// ---------------------------------------------------------------------------
// Unified webhook handler
// ---------------------------------------------------------------------------

export interface WebhookHandlerOptions {
  /**
   * The payload type. Defaults to `"json"` (HMAC-SHA256 signature).
   * Set to `"jwt"` or `"jwe"` to use token-based verification.
   */
  payloadType?: PayloadType;

  // ---- JSON (HMAC) options ----
  /** The shared webhook signing key (required when payloadType is "json"). */
  signingKey?: string;

  // ---- JWT / JWE options ----
  /** JWKS endpoint of the ZITADEL instance (required when payloadType is "jwt" or "jwe"). */
  jwksEndpoint?: string;
  /** Expected JWT issuer. */
  issuer?: string;
  /** Expected JWT audience. */
  audience?: string;
  /** PEM-encoded private key for JWE decryption (required when payloadType is "jwe"). */
  privateKey?: string;

  /** Callback invoked with the parsed webhook payload when verification succeeds. */
  onEvent: (event: unknown) => void | Promise<void>;
}

/**
 * Creates a webhook handler that verifies the payload and invokes the callback.
 *
 * Supports all three ZITADEL Actions v2 payload types:
 * - `"json"` (default) — HMAC-SHA256 signature in `x-zitadel-signature` header
 * - `"jwt"` — JWT signed by the instance, verified via JWKS
 * - `"jwe"` — JWE-encrypted JWT, decrypted with your private key then verified via JWKS
 */
export function createWebhookHandler(options: WebhookHandlerOptions) {
  const payloadType = options.payloadType ?? "json";

  return async (request: {
    body: string;
    headers: Record<string, string>;
  }) => {
    try {
      let event: unknown;

      switch (payloadType) {
        case "json": {
          if (!options.signingKey) {
            throw new Error("signingKey is required for JSON payload type");
          }
          const signature =
            request.headers["x-zitadel-signature"] ?? "";
          const isValid = verifyZitadelWebhook({
            body: request.body,
            signature,
            signingKey: options.signingKey,
          });
          if (!isValid) {
            return { status: 401, body: "Invalid signature" };
          }
          event = JSON.parse(request.body);
          break;
        }

        case "jwt": {
          if (!options.jwksEndpoint) {
            throw new Error("jwksEndpoint is required for JWT payload type");
          }
          event = await verifyJwtWebhook({
            token: request.body.trim(),
            jwksEndpoint: options.jwksEndpoint,
            issuer: options.issuer,
            audience: options.audience,
          });
          break;
        }

        case "jwe": {
          if (!options.jwksEndpoint || !options.privateKey) {
            throw new Error(
              "jwksEndpoint and privateKey are required for JWE payload type",
            );
          }
          event = await decryptJweWebhook({
            token: request.body.trim(),
            privateKey: options.privateKey,
            jwksEndpoint: options.jwksEndpoint,
            issuer: options.issuer,
            audience: options.audience,
          });
          break;
        }

        default:
          return { status: 400, body: `Unknown payload type: ${payloadType}` };
      }

      await options.onEvent(event);
      return { status: 200, body: "OK" };
    } catch {
      return { status: 401, body: "Verification failed" };
    }
  };
}
