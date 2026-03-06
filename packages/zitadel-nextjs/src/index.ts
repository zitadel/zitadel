// Auth — OIDC redirect flow
export {
  signIn,
  handleCallback,
  getOIDCSession,
  signOut,
  type OIDCOptions,
  type OIDCCallbackResult,
  type OIDCSession,
} from "./auth/oidc.js";

// Middleware
export { createZitadelMiddleware } from "./middleware.js";
export type { MiddlewareOptions } from "./middleware.js";

// Actions v2 webhook handler
export { createZitadelWebhookHandler } from "./webhook.js";
export type { WebhookHandlerOptions } from "./webhook.js";

// Server actions
export { protectedAction } from "./server-action.js";

// Session (internal — reads OIDC session cookie)
export { getSession } from "./session.js";
export type { SessionData } from "./session.js";

// ZITADEL v2 API client
export { createZitadelApiClient, withApiClient } from "./api.js";
export type { ApiClientOptions, ZitadelApiClient } from "./api.js";
