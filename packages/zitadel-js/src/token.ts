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
