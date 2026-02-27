/**
 * ZITADEL authentication module.
 *
 * Re-exports from the available auth sub-modules:
 *
 * - **`auth/oidc`** — OIDC redirect-based login ("add login to your app")
 * - **`auth/session`** — Session API for custom login UIs *(planned)*
 *
 * @module
 */
export {
  signIn,
  handleCallback,
  getOIDCSession,
  signOut,
  type OIDCOptions,
  type OIDCCallbackResult,
  type OIDCSession,
} from "./oidc.js";
