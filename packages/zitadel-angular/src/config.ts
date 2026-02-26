export interface ZitadelConfig {
  /** The ZITADEL issuer URL. */
  issuerUrl: string;
  /** The OIDC client ID. */
  clientId: string;
  /** The callback URL registered in ZITADEL. */
  callbackUrl: string;
  /** Scopes to request. Defaults to ["openid", "profile", "email"]. */
  scopes?: string[];
}
