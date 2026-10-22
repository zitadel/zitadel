import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { QueryClient } from '@tanstack/angular-query-experimental';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, firstValueFrom } from 'rxjs';
import { AccountStorage } from './storage';
import { safeAccountPath } from './routes';

export interface AccountEnvironment {
  api: string;
  issuer: string;
  clientid: string;
  loginOrgId: string;
}
interface Transaction {
  id: string;
  path: string;
  created: number;
}

/** Account-only configuration around the Console's OIDC library. No login implementation. */
@Injectable({ providedIn: 'root' })
export class AccountAuthService {
  readonly authenticationChanged = new BehaviorSubject(false);
  environment!: AccountEnvironment;
  error = false;
  private ready?: Promise<void>;
  private resumePath = safeAccountPath(location.pathname);
  constructor(
    private oauth: OAuthService,
    private http: HttpClient,
    private storage: AccountStorage,
    private router: Router,
    private dialogs: MatDialog,
    private queries: QueryClient,
  ) {}
  // GrpcService supplies Console defaults; the Account always applies its own strict contract below.
  initConfig(_config: AuthConfig): void {}
  get authenticated(): boolean {
    return this.valid();
  }
  getIdToken(): string {
    return this.oauth.getIdToken();
  }
  initialize(): Promise<void> {
    return (this.ready ??= this.start());
  }
  private async start(): Promise<void> {
    try {
      const env = await firstValueFrom(this.http.get<AccountEnvironment>('/assets/environment.json'));
      for (const value of [env.api, env.issuer]) {
        const url = new URL(value);
        if (url.protocol !== 'https:' || url.origin !== value) throw new Error('configuration');
      }
      if (!env.clientid || !/^\d+$/.test(env.loginOrgId)) throw new Error('configuration');
      this.environment = env;
      this.storage.setItem('organizationId', env.loginOrgId);
      this.oauth.configure({
        issuer: env.issuer,
        clientId: env.clientid,
        responseType: 'code',
        oidc: true,
        scope: `openid profile email urn:zitadel:iam:user:resourceowner urn:zitadel:iam:org:id:${env.loginOrgId} urn:zitadel:iam:org:project:id:zitadel:aud`,
        redirectUri: `${location.origin}/auth/callback`,
        postLogoutRedirectUri: `${location.origin}/signedout`,
        requireHttps: true,
        strictDiscoveryDocumentValidation: true,
        clearHashAfterLogin: true,
        disablePKCE: false,
        showDebugInformation: false,
        timeoutFactor: 1,
      });
      if (location.pathname === '/signedout') {
        this.oauth.logOut(true);
        this.storage.clear();
        return;
      }
      await this.oauth.loadDiscoveryDocument();
      if (location.pathname === '/auth/callback') {
        const saved = this.storage.getItem('transaction');
        this.storage.removeItem('transaction'); // consume even on error; never replay a failed callback
        const tx: Transaction | undefined = saved ? JSON.parse(saved) : undefined;
        if (
          !tx ||
          typeof tx.id !== 'string' ||
          !Number.isFinite(tx.created) ||
          Date.now() - tx.created > 600000 ||
          Date.now() < tx.created
        )
          throw new Error('transaction');
        this.resumePath = safeAccountPath(tx.path);
        await this.oauth.tryLoginCodeFlow(); // library verifies state, nonce, issuer, audience, expiry, signature handler
        if (this.oauth.state !== tx.id || !this.valid()) throw new Error('context');
        this.storage.removeItem('nonce');
        this.storage.removeItem('PKCE_verifier');
        history.replaceState(null, '', safeAccountPath(tx.path));
      }
      if (this.oauth.hasValidIdToken() && !this.valid()) throw new Error('context');
      // The native password form calls history.back() on success. Seed an internal
      // predecessor for a direct entry/callback so it never returns to the IdP.
      if (this.valid() && location.pathname === '/seguranca/senha') {
        history.replaceState(null, '', '/seguranca');
        history.pushState(null, '', '/seguranca/senha');
      }
      this.authenticationChanged.next(this.valid());
      this.oauth.events.subscribe((event) => {
        if (event.type === 'token_expires' || event.type === 'session_terminated') {
          this.resumePath = safeAccountPath(location.pathname);
          this.error = true;
          this.oauth.logOut(true);
          this.storage.clear();
          this.authenticationChanged.next(false);
          this.dialogs.closeAll();
          this.queries.clear();
          void this.router.navigateByUrl('/auth/error');
        }
      });
    } catch {
      this.oauth.logOut(true);
      this.storage.clear();
      this.error = true;
      history.replaceState(null, '', '/auth/error');
    }
  }
  private valid(): boolean {
    const c = this.oauth.getIdentityClaims() as Record<string, unknown> | null;
    const aud = c?.['aud'];
    return (
      !!this.environment &&
      this.oauth.hasValidIdToken() &&
      this.oauth.hasValidAccessToken() &&
      c?.['iss'] === this.environment.issuer &&
      typeof c?.['sub'] === 'string' &&
      !!c['sub'] &&
      (aud === this.environment.clientid || (Array.isArray(aud) && aud.includes(this.environment.clientid))) &&
      (!c['azp'] || c['azp'] === this.environment.clientid) &&
      c['urn:zitadel:iam:user:resourceowner:id'] === this.environment.loginOrgId
    );
  }
  async authenticate(_partial?: Partial<AuthConfig>, force = false): Promise<boolean> {
    await this.initialize();
    if (this.valid() && !force) return true;
    if (!this.error) this.login(location.pathname, force);
    return false;
  }
  login(path = '/', force = false): void {
    if (!this.environment) return;
    const existing = this.storage.getItem('transaction');
    if (existing) {
      this.error = true;
      location.assign('/auth/error');
      return;
    }
    this.resumePath = safeAccountPath(path);
    const tx: Transaction = { id: crypto.randomUUID(), path: this.resumePath, created: Date.now() };
    this.storage.setItem('transaction', JSON.stringify(tx));
    this.oauth.initCodeFlow(tx.id, force ? { prompt: 'login', max_age: '0' } : {});
  }
  retry(): void {
    this.storage.removeItem('transaction');
    // A new bootstrap retries discovery/configuration and creates a fresh transaction.
    location.assign(this.resumePath);
  }
  signout(): void {
    this.authenticationChanged.next(false);
    this.dialogs.closeAll();
    this.queries.clear();
    this.oauth.logOut();
    this.storage.clear();
  }
}
