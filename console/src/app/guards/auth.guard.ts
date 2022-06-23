import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';
import { Observable } from 'rxjs';

import { AuthenticationService } from '../services/authentication.service';
import { GrpcAuthService } from '../services/grpc-auth.service';

@Injectable({
  providedIn: 'root',
})
export class AuthGuard implements CanActivate {
  constructor(private auth: AuthenticationService, private authService: GrpcAuthService) {}

  public canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot,
  ): Observable<boolean> | Promise<boolean> | Promise<any> | boolean {
    if (!this.auth.authenticated) {
      if (route.queryParams && route.queryParams['login_hint']) {
        const hint = route.queryParams['login_hint'];
        const configWithPrompt: Partial<AuthConfig> = {
          customQueryParams: {
            login_hint: hint,
          },
        };
        console.log(`authenticate with login_hint: ${hint}`);
        this.auth.authenticate(configWithPrompt);
      } else {
        return this.auth.authenticate();
      }
    }
    return this.auth.authenticated;
  }
}
