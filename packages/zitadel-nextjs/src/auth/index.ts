/**
 * ZITADEL authentication module.
 *
 * Re-exports from the available auth sub-modules:
 *
 * - **`auth/oidc`** — OIDC redirect-based login ("add login to your app")
 * - **`auth/session`** — Session API helpers for custom login UIs
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

export {
  createSession,
  setSession,
  getSession,
  deleteSession,
  createCallback,
  type SessionAuthOptions,
  type SessionLifetime,
  type CreateSessionOptions,
  type SetSessionOptions,
  type GetSessionOptions,
  type DeleteSessionOptions,
  type CreateCallbackOptions,
} from "./session.js";
