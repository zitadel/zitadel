import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { LoginSettings, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

/**
 * Reason why an existing session cannot be reused to complete the current
 * authentication request. Used to render an actionable hint in the account
 * picker instead of letting the flow dead-end (see issue #11805).
 */
export type SessionReuseBlockReason = "orgMismatch" | "localAuthNotAllowed" | "externalIdpNotAllowed" | "passkeysNotAllowed";

export type SessionReuseResult = { reusable: true } | { reusable: false; reason: SessionReuseBlockReason };

/**
 * Determines the primary (first-factor) authentication method of a session.
 *
 * Mirrors the backend `nextSteps` first-factor detection: a session that was
 * established through an external IdP (`intent`) is treated as an IdP login even
 * if it also carries a password factor; otherwise password, otherwise passkey.
 */
function getPrimaryFactor(session: Session): "idp" | "password" | "passkey" | undefined {
  if (session.factors?.intent?.verifiedAt) {
    return "idp";
  }
  if (session.factors?.password?.verifiedAt) {
    return "password";
  }
  if (session.factors?.webAuthN?.verifiedAt) {
    return "passkey";
  }
  return undefined;
}

/**
 * Replicates the subset of the backend login-policy / auth-method gating
 * (`internal/auth/repository/eventsourcing/eventstore/auth_request.go` →
 * `nextSteps`) that the account picker can evaluate client-side, so an
 * incompatible session is not silently offered as a reusable account.
 *
 * The two gates mirrored here are:
 *  1. Organization match — `LinkSessionToAuthRequest` rejects a session whose
 *     user resource-owner differs from the auth request's pinned organization
 *     (`Errors.User.NotAllowedOrg`). Only enforced when a target organization
 *     is known.
 *  2. Auth-method vs. login policy — the session's primary factor must be
 *     permitted by the target org's login settings:
 *       - password  → `allowLocalAuthentication`
 *       - external IdP → `allowExternalIdp`
 *       - passkey   → `passkeysType === ALLOWED`
 *
 * It intentionally errs towards "reusable": a block is only returned when there
 * is a definitive policy reason. When the target organization or its login
 * settings are unknown, the session is treated as reusable and the server-side
 * graceful re-authentication fallback in `loginWithOIDCAndSession` /
 * `loginWithSAMLAndSession` remains the safety net.
 *
 * NOTE (follow-up): the backend also intersects the session's IdP with the
 * org's *allowed* IdP set (`checkForAllowedIDPs`). Replicating that requires the
 * active identity-provider list (`getActiveIdentityProviders`) and is not yet
 * covered here — a session using an IdP that the target org no longer allows
 * will currently fall through to the server-side fallback.
 */
export function checkSessionReuse({
  session,
  targetOrganization,
  loginSettings,
}: {
  session: Session;
  targetOrganization?: string;
  loginSettings?: LoginSettings;
}): SessionReuseResult {
  const sessionOrg = session.factors?.user?.organizationId;

  if (targetOrganization && sessionOrg && sessionOrg !== targetOrganization) {
    return { reusable: false, reason: "orgMismatch" };
  }

  if (!loginSettings) {
    return { reusable: true };
  }

  const primary = getPrimaryFactor(session);

  if (primary === "password" && !loginSettings.allowLocalAuthentication) {
    return { reusable: false, reason: "localAuthNotAllowed" };
  }

  if (primary === "idp" && !loginSettings.allowExternalIdp) {
    return { reusable: false, reason: "externalIdpNotAllowed" };
  }

  if (primary === "passkey" && loginSettings.passkeysType !== PasskeysType.ALLOWED) {
    return { reusable: false, reason: "passkeysNotAllowed" };
  }

  return { reusable: true };
}
