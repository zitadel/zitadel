import type { ZitadelAuthService } from "./auth.service.js";

/**
 * HTTP interceptor that attaches the ZITADEL access token to outgoing requests.
 * Placeholder — to be implemented as an Angular HttpInterceptorFn.
 */
export function zitadelAuthInterceptor(_authService: ZitadelAuthService) {
  return (req: unknown, next: (req: unknown) => unknown) => {
    // TODO: implement with Angular HttpInterceptorFn
    return next(req);
  };
}
