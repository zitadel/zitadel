import { Injectable } from '@angular/core';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, Observable } from 'rxjs';

import { StatehandlerService } from './statehandler.service';

@Injectable({
    providedIn: 'root',
})
export class AuthenticationService {
    private authConfig!: AuthConfig;
    private _authenticated: boolean = false;
    private readonly _authenticationChanged: BehaviorSubject<
        boolean
    > = new BehaviorSubject(this.authenticated);

    constructor(
        private oauthService: OAuthService,
        private statehandler: StatehandlerService,
    ) { }

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

    public async authenticate(
        partialConfig?: Partial<AuthConfig>,
    ): Promise<boolean> {
        if (partialConfig) {
            Object.assign(this.authConfig, partialConfig);
        }
        this.oauthService.configure(this.authConfig);

        this.oauthService.strictDiscoveryDocumentValidation = false;
        await this.oauthService.loadDiscoveryDocumentAndTryLogin();

        this._authenticated = this.oauthService.hasValidAccessToken();

        if (!this.oauthService.hasValidIdToken() || !this.authenticated || partialConfig) {
            const newState = await this.statehandler.createState().toPromise();
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

