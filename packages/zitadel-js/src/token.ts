import * as jose from "jose";

export interface SystemTokenOptions {
  /** The key ID (kid) for the JWT header. */
  keyId: string;
  /** The PEM-encoded private key used to sign the JWT. */
  key: string;
  /** The issuer of the token (typically the service user's client ID). */
  issuer: string;
  /** The audience for the token (typically the ZITADEL issuer URL). */
  audience: string;
  /** Token lifetime in seconds. Defaults to 3600 (1 hour). */
  expiresInSeconds?: number;
}

/**
 * Creates a signed JWT for ZITADEL system-level authentication using jose.
 */
export async function newSystemToken(options: SystemTokenOptions): Promise<string> {
  const { keyId, key, issuer, audience, expiresInSeconds = 3600 } = options;
  const privateKey = await jose.importPKCS8(key, "RS256");

  return new jose.SignJWT({})
    .setProtectedHeader({ alg: "RS256", kid: keyId })
    .setIssuedAt()
    .setIssuer(issuer)
    .setSubject(issuer)
    .setAudience(audience)
    .setExpirationTime(`${expiresInSeconds}s`)
    .sign(privateKey);
}

// ---------------------------------------------------------------------------
// JWT verification
// ---------------------------------------------------------------------------

export interface VerifyJwtOptions {
  /** Expected JWT issuer (ZITADEL instance URL). */
  issuer?: string;
  /** Expected JWT audience. */
  audience?: string;
  /**
   * Value for the `x-zitadel-instance-host` header sent with the JWKS request.
   * Required in multi-tenant setups where the JWKS endpoint is behind a proxy.
   */
  instanceHost?: string;
  /**
   * Value for the `x-zitadel-public-host` header sent with the JWKS request.
   */
  publicHost?: string;
}

/**
 * Verifies a JWT against a remote JWKS endpoint.
 *
 * Useful for validating ZITADEL-issued access tokens in server-side code
 * (API routes, middleware, server actions).
 *
 * @example
 * ```ts
 * import { verifyJwt } from "@zitadel/zitadel-js/node";
 *
 * const payload = await verifyJwt(
 *   accessToken,
 *   "https://my.zitadel.cloud/oauth/v2/keys",
 *   { issuer: "https://my.zitadel.cloud", audience: "my-client-id" },
 * );
 * ```
 */
export async function verifyJwt<T = jose.JWTPayload>(
  token: string,
  keysEndpoint: string,
  options?: VerifyJwtOptions,
): Promise<T & jose.JWTPayload> {
  const headers: Record<string, string> = {};
  if (options?.instanceHost) {
    headers["x-zitadel-instance-host"] = options.instanceHost;
  }
  if (options?.publicHost) {
    headers["x-zitadel-public-host"] = options.publicHost;
  }

  const JWKS = jose.createRemoteJWKSet(new URL(keysEndpoint), { headers });

  const { payload } = await jose.jwtVerify(token, JWKS, {
    issuer: options?.issuer,
    audience: options?.audience,
  });

  return payload as T & jose.JWTPayload;
}
