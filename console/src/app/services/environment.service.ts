import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, from, map, Observable, of, shareReplay, switchMap, throwError } from 'rxjs';

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
  exhausted?: boolean;
}

interface WellKnown {
  authorization_endpoint: string;
  end_session_endpoint: string;
  introspection_endpoint: string;
  token_endpoint: string;
  userinfo_endpoint: string;
}
@Injectable({
  providedIn: 'root',
})
export class EnvironmentService {
  private exhaustedCookieKey = 'zitadel.quota.exhausted';
  private environmentJsonPath = './assets/environment.json';
  private wellknownPath = '/.well-known/openid-configuration`';
  public auth!: AuthServiceClient;
  public mgmt!: ManagementServiceClient;
  public admin!: AdminServiceClient;

  private environment$: Observable<Environment>;
  private wellKnown$: Observable<WellKnown>;

  constructor(private http: HttpClient, private exhaustedSvc: ExhaustedService) {
    this.environment$ = this.createEnvironment();
    this.wellKnown$ = this.createWellKnown(this.environment$);
    setInterval(() => {
      console.debug('setInterval', 'document.cookie', document.cookie);
    }, 10);
  }

  // env returns an `Observable<Environment>` that can be subscribed to whenever needed.
  // It makes the HTTP call exactly once and replays the cached result.
  // If the responses exhausted property is true, the exhaused dialog is shown.
  // If it is false, the observable waits until the browser has the cookie set before emitting.
  get env() {
    return this.environment$;
  }

  // wellKnown returns an `Observable<Environment>` that can be subscribed to whenever needed.
  // It makes the HTTP call exactly once and replays the cached result.
  get wellKnown() {
    return this.wellKnown$;
  }

  private createEnvironment() {
    return this.http.get<Environment>(this.environmentJsonPath).pipe(
      catchError((err) => {
        console.error('Getting environment.json failed', err);
        return throwError(() => err);
      }),
      switchMap((env) => {
        if (env.exhausted) {
          return this.exhaustedSvc.showExhaustedDialog(of(env)).pipe(map(() => env));
        }
        if (!navigator.cookieEnabled) {
          return of(env);
        }
        return from(this.awaitFiveSeconds(() => !document.cookie.includes(this.exhaustedCookieKey))).pipe(map(() => env));
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

  private awaitFiveSeconds(condition: () => boolean): Promise<void> {
    return new Promise((resolve, reject) => {
      let checks = 0;
      const check = () => {
        console.debug('awaitFiveSeconds', 'document.cookie', document.cookie);
        if (condition()) {
          resolve();
        } else {
          checks++;
          if (checks > 500) {
            reject(`after ${checks} checks, the condition did not return true`);
          } else {
            setTimeout(check, 10);
          }
        }
      };
      check();
    });
  }
}
