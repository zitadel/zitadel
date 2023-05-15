import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, of, shareReplay, switchMap, tap, throwError } from 'rxjs';

import { AdminServiceClient } from '../proto/generated/zitadel/AdminServiceClientPb';
import { AuthServiceClient } from '../proto/generated/zitadel/AuthServiceClientPb';
import { ManagementServiceClient } from '../proto/generated/zitadel/ManagementServiceClientPb';
import { ExhaustedService } from './exhausted.service';

interface Environment {
  api: string;
  clientid: string;
  issuer: string;
  customer_portal?: string;
  instance_management_url?: string;
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
  private exhaustedCookieKey = 'zitadel.quota.limiting';
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
  }

  // env returns an `Observable<Environment>` that can be subscribed to whenever needed.
  // It makes the HTTP call exactly once and replays the cached result.
  // If the backend sends the exhaused cookie with the Max-Age property and value 0, the exhaused dialog is shown.
  // If the backend sends the exhaused cookie with the Max-Age property and value above 0, the exhaused dialog is shown.
  get env() {
    return this.environment$;
  }

  // wellKnown returns an `Observable<Environment>` that can be subscribed to whenever needed.
  // It makes the HTTP call exactly once and replays the cached result.
  get wellKnown() {
    return this.wellKnown$;
  }

  private get hasExhaustedCookie() {
    return !!document.cookie
      .split(';')
      .map((c) => c.trim())
      .find(this.isCookie);
  }

  private get cookiesAreEnabled() {
    try {
      document.cookie = 'testcookie';
      const canSetCookie = document.cookie.indexOf('testcookie') != -1;
      document.cookie = 'testcookie; expires=Thu, 01 Jan 1970 00:00:00 UTC';
      return canSetCookie;
    } catch (e) {
      return false;
    }
  }

  // isCookie needs `this` to point to the class instance, so we use an arrow function
  private isCookie = (cookie: string) => {
    return cookie.startsWith(`${this.exhaustedCookieKey}=`);
  };

  createEnvironment() {
    return of(null).pipe(
      // Delete the exhausted cookie before the enviroment is loaded
      tap(() => (document.cookie = `${this.exhaustedCookieKey}=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC"`)),
      switchMap(() => {
        // Wait until the browser deleted the cookie
        return new Promise<void>((resolve, reject) => {
          this.awaitFiveSeconds(() => !this.hasExhaustedCookie, reject);
          resolve();
        });
      }),
      switchMap(() => this.http.get<Environment>(this.environmentJsonPath, { observe: 'response' })),
      catchError((err) => {
        console.error('Getting environment.json failed', err);
        return throwError(() => err);
      }),
      switchMap((resp) => {
        return new Promise<Environment>((resolve, reject) => {
          if (resp.body === null) {
            reject('environment.json has no body');
            return;
          }
          const exhaustedResponseCookie = resp.headers.getAll('set-cookie')?.find(this.isCookie);
          if (!exhaustedResponseCookie || !this.cookiesAreEnabled) {
            resolve(resp.body as Environment);
            return;
          }

          // The `/i` in the end of the RegExp matches without case sentitivity
          const maxAgeRegex = /Max-Age=(\d+)/i;
          const match = exhaustedResponseCookie.match(maxAgeRegex);

          // If there is no Max-Age, we don't know if the browser should have the cookie or not.
          if (!match) {
            resolve(resp.body as Environment);
            return;
          }

          // If it is above 0, we show the exhausted dialog that either refreshes
          // the page or redirects the user to the instance management URL.
          if (parseInt(match[1]) > 0) {
            this.exhaustedSvc.showExhaustedDialog(resp.body.instance_management_url);
            return;
          }

          // If the Max age is 0 or below, the browser must not have the cookie.
          // In this case, we wait for the browser to delete the cookie.
          this.awaitFiveSeconds(() => !this.hasExhaustedCookie, reject);
        });
      }),
      // Cache the first response, then replay it
      shareReplay(1),
    );
  }

  createWellKnown(environment$: Observable<Environment>) {
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

  awaitFiveSeconds(condition: () => boolean, failWithTimeoutMessage: (msg: string) => void) {
    let checks = 0;
    const check = () => {
      if (condition()) {
        return;
      }
      checks++;
      if (checks > 500) {
        failWithTimeoutMessage(`after ${checks} checks, the condition did not return true`);
        return;
      }
      setTimeout(check, 10);
    };
  }
}
