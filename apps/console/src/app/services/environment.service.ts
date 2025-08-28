import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, map, Observable, of, shareReplay, switchMap, throwError } from 'rxjs';

import { AdminServiceClient } from '../proto/generated/zitadel/AdminServiceClientPb';
import { AuthServiceClient } from '../proto/generated/zitadel/AuthServiceClientPb';
import { ManagementServiceClient } from '../proto/generated/zitadel/ManagementServiceClientPb';
import { ExhaustedService } from './exhausted.service';

export interface Environment {
  api: string;
  clientid: string;
  issuer: string;
  customer_portal?: string;
  instance_management_url?: string;
  posthog_token?: string;
  posthog_url?: string;
  exhausted?: boolean;
}

interface WellKnown {
  authorization_endpoint: string;
  device_authorization_endpoint: string;
  end_session_endpoint: string;
  introspection_endpoint: string;
  token_endpoint: string;
  userinfo_endpoint: string;
  jwks_uri: string;
  revocation_endpoint: string;
}
@Injectable({
  providedIn: 'root',
})
export class EnvironmentService {
  private environmentJsonPath = './assets/environment.json';
  private wellknownPath = '/.well-known/openid-configuration';
  public auth!: AuthServiceClient;
  public mgmt!: ManagementServiceClient;
  public admin!: AdminServiceClient;

  private environment$: Observable<Environment>;
  private wellknown$: Observable<WellKnown>;

  constructor(
    private http: HttpClient,
    private exhaustedSvc: ExhaustedService,
  ) {
    this.environment$ = this.createEnvironment();
    this.wellknown$ = this.createWellKnown(this.environment$);
  }

  // env returns an `Observable<Environment>` that can be subscribed to whenever needed.
  // It makes the HTTP call exactly once and replays the cached result.
  // If the responses exhausted property is true, the exhaused dialog is shown.
  get env() {
    return this.environment$;
  }

  // wellknown returns an `Observable<Environment>` that can be subscribed to whenever needed.
  // It makes the HTTP call exactly once and replays the cached result.
  get wellknown() {
    return this.wellknown$;
  }

  private createEnvironment() {
    return this.http.get<Environment>(this.environmentJsonPath).pipe(
      catchError((err) => {
        console.error('Getting environment.json failed', err);
        return throwError(() => err);
      }),
      switchMap((env) => {
        const env$ = of(env);
        if (env.exhausted) {
          return this.exhaustedSvc.showExhaustedDialog(env$).pipe(map(() => env));
        }
        return env$;
      }),
      // Cache the first response, then replay it
      shareReplay(1),
    );
  }

  private createWellKnown(environment$: Observable<Environment>) {
    return environment$.pipe(
      catchError((err) => {
        console.error('Getting well-known OIDC configuration failed', err);
        return throwError(() => err);
      }),
      switchMap((env) => {
        return this.http.get<WellKnown>(`${env.issuer}${this.wellknownPath}`);
      }),
      // Cache the first response, then replay it
      shareReplay(1),
    );
  }
}
