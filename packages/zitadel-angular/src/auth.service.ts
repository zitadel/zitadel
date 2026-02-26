import type { ZitadelConfig } from "./config.js";

/**
 * Authentication service for ZITADEL.
 * Placeholder — to be implemented as an Angular Injectable.
 */
export class ZitadelAuthService {
  private _config: ZitadelConfig;
  private _accessToken: string | undefined;
  private _isAuthenticated = false;

  constructor(config: ZitadelConfig) {
    this._config = config;
  }

  get config(): ZitadelConfig {
    return this._config;
  }

  get accessToken(): string | undefined {
    return this._accessToken;
  }

  get isAuthenticated(): boolean {
    return this._isAuthenticated;
  }

  /** Initiates the OIDC sign-in flow. */
  async signIn(): Promise<void> {
    // TODO: implement PKCE flow with Angular Router
  }

  /** Signs out the current user. */
  async signOut(): Promise<void> {
    this._accessToken = undefined;
    this._isAuthenticated = false;
    // TODO: implement session cleanup
  }

  /** Handles the OIDC callback. */
  async handleCallback(_params: Record<string, string>): Promise<void> {
    // TODO: implement code exchange
  }
}
