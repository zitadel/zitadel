import type { ZitadelAuthService } from "./auth.service.js";

/**
 * Route guard that checks if the user is authenticated.
 * Placeholder — to be implemented as an Angular CanActivateFn.
 */
export function zitadelAuthGuard(_authService: ZitadelAuthService) {
  return (): boolean => {
    // TODO: implement with Angular Router canActivate
    return _authService.isAuthenticated;
  };
}
