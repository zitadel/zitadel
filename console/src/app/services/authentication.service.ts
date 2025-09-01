import { Injectable } from '@angular/core';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, lastValueFrom, Observable } from 'rxjs';

import { StatehandlerService } from './statehandler/statehandler.service';
import { ToastService } from './toast.service';

@Injectable({
  providedIn: 'root',
})
export class AuthenticationService {
  private authConfig!: AuthConfig;
  private _authenticated = false;
  private readonly _authenticationChanged = new BehaviorSubject(this.authenticated);

  constructor(
    private oauthService: OAuthService,
    private statehandler: StatehandlerService,
    private toast: ToastService,
  ) {}

  public initConfig(data: AuthConfig): void {
    this.authConfig = data;
  }

  public get authenticated(): boolean {
    return this._authenticated;
  }

  public get authenticationChanged(): Observable<boolean> {
    return this._authenticationChanged;
  }

  public getOIDCUser(): Observable<any> {
    return from(this.oauthService.loadUserProfile());
  }

  public getIdToken(): string {
    return this.oauthService.getIdToken();
  }

  public async authenticate(partialConfig?: Partial<AuthConfig>, force: boolean = false): Promise<boolean> {
    if (partialConfig) {
      Object.assign(this.authConfig, partialConfig);
    }
    this.oauthService.configure(this.authConfig);
    this.oauthService.strictDiscoveryDocumentValidation = false;
    try {
      await this.oauthService.loadDiscoveryDocumentAndTryLogin();
    } catch (error) {
      this.toast.showError(error, false, false);
    }

    this._authenticated = this.oauthService.hasValidAccessToken();
    if (!this.oauthService.hasValidIdToken() || !this.authenticated || partialConfig || force) {
      const newState = await lastValueFrom(this.statehandler.createState());
      this.oauthService.initCodeFlow(newState);
    }
    this._authenticationChanged.next(this.authenticated);

    return this.authenticated;
  }

  public signout(): void {
    this.oauthService.logOut();
    this._authenticated = false;
    this._authenticationChanged.next(false);
  }
}
