import { PlatformLocation } from '@angular/common';
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
        private platformLocation: PlatformLocation,
        private oauthService: OAuthService,
        private statehandler: StatehandlerService,
    ) { }

    public authInit = (
        envDeps: Promise<any>, // (() => Function),
    ): () => Promise<any> => {
        return (): Promise<any> => {
            return envDeps.then(data => {
                this.authConfig = {
                    scope: 'openid profile email', // offline_access
                    responseType: 'code',
                    oidc: true,
                    clientId: data.clientid,
                    redirectUri: window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'auth/callback',
                    postLogoutRedirectUri: window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'signedout',
                };
            });
        };
    };

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
        config?: Partial<AuthConfig>,
        setState: boolean = true,
        force: boolean = false,
    ): Promise<boolean> {
        if (config) {
            this.authConfig = config;
        }
        this.oauthService.configure(this.authConfig);

        this.oauthService.strictDiscoveryDocumentValidation = false;
        await this.oauthService.loadDiscoveryDocumentAndTryLogin();

        this._authenticated = this.oauthService.hasValidAccessToken();

        if (!this.oauthService.hasValidIdToken() || !this.authenticated || config || force) {
            const newState = setState ? await this.statehandler.createState().toPromise() : undefined;
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

