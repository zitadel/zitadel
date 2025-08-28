import { inject } from '@angular/core';
import { CanActivateFn } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';

import { AuthenticationService } from '../services/authentication.service';

export const authGuard: CanActivateFn = (route) => {
  const auth = inject(AuthenticationService);

  if (!auth.authenticated) {
    if (route.queryParams && route.queryParams['login_hint']) {
      const hint = route.queryParams['login_hint'];
      const configWithPrompt: Partial<AuthConfig> = {
        customQueryParams: {
          login_hint: hint,
        },
      };
      console.log(`authenticate with login_hint: ${hint}`);
      auth.authenticate(configWithPrompt).then();
    } else {
      return auth.authenticate();
    }
  }

  return auth.authenticated;
};
