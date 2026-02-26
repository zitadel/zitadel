import type { ZitadelConfig } from "./config.js";
import { ZitadelAuthService } from "./auth.service.js";

/**
 * Provides ZITADEL services for Angular dependency injection.
 * Placeholder — to be implemented with Angular's provide* pattern.
 */
export function provideZitadel(config: ZitadelConfig) {
  return {
    authService: new ZitadelAuthService(config),
  };
}
