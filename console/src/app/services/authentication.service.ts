import { Injectable } from '@angular/core';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, Observable } from 'rxjs';

import { GrpcService } from './grpc.service';
import { StatehandlerService } from './statehandler.service';

@Injectable({
    providedIn: 'root',
})
export class AuthenticationService {
    private _authenticated: boolean = false;
    private readonly _authenticationChanged: BehaviorSubject<
        boolean
    > = new BehaviorSubject(this.authenticated);

    constructor(
        private grpcService: GrpcService,
        private config: AuthConfig,
        private oauthService: OAuthService,
        private statehandler: StatehandlerService,
    ) { }

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
        this.config.issuer = config?.issuer || this.grpcService.issuer;
        this.config.clientId = config?.clientId || this.grpcService.clientid;
        this.config.redirectUri = config?.redirectUri || this.grpcService.redirectUri;
        this.config.postLogoutRedirectUri = config?.postLogoutRedirectUri || this.grpcService.postLogoutRedirectUri;
        this.config.customQueryParams = config?.customQueryParams;
        this.oauthService.configure(this.config);

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

