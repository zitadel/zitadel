import "server-only";

import { Cookie, getMostRecentCookieWithLoginname, getSessionCookieById } from "@/lib/server/cookies";
import { createLogger } from "@/lib/server/logger";
import { parseAndVerifySessions } from "@/lib/server/session-cookie-signature";
import { shouldEnforceMFA } from "@/lib/server/verify-helper";
import {
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  ServiceConfig,
} from "@/lib/server/zitadel";
import { Timestamp, timestampDate } from "@zitadel/client";
import { AuthRequest } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { SAMLRequest } from "@zitadel/proto/zitadel/saml/v2/authorization_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { cookies, headers } from "next/headers";

const logger = createLogger("session");

/**
 * Why a session that a flow expected could not be resolved from the browser's `sessions` cookie.
 *
 * - `cookie_missing`: the browser sent no `sessions` cookie
 * - `cookie_invalid`: the cookie is unreadable, or the wanted entry has a missing or invalid signature
 * - `session_entry_not_found`: the cookie is fine but holds no entry for the wanted session
 * - `server_session_not_found`: the cookie entry exists, the API no longer knows the session
 * - `session_token_rejected`: the session exists, the API rejects the token of the cookie entry
 * - `session_lookup_failed`: the API call failed for another reason
 */
type SessionResolutionFailure =
  | "cookie_missing"
  | "cookie_invalid"
  | "session_entry_not_found"
  | "server_session_not_found"
  | "session_token_rejected"
  | "session_lookup_failed";

type SessionResolution =
  { session: Session; reason?: undefined } | { session?: undefined; reason: SessionResolutionFailure };

// Stable codes to search the logs for, one per failure reason.
const SESSION_RESOLUTION_LOG_CODES: Record<SessionResolutionFailure, string> = {
  cookie_missing: "LOGIN_SESSION_COOKIE_MISSING",
  cookie_invalid: "LOGIN_SESSION_COOKIE_INVALID",
  session_entry_not_found: "LOGIN_SESSION_ENTRY_NOT_FOUND",
  server_session_not_found: "LOGIN_SESSION_NOT_FOUND",
  session_token_rejected: "LOGIN_SESSION_TOKEN_REJECTED",
  session_lookup_failed: "LOGIN_SESSION_LOOKUP_FAILED",
};

/**
 * Keys in the `error` namespace. An invalid signature and a rejected token share one generic
 * message, the log code tells them apart.
 */
export const SESSION_RESOLUTION_MESSAGE_KEYS: Record<SessionResolutionFailure, string> = {
  cookie_missing: "sessionResolution.cookieMissing",
  cookie_invalid: "sessionResolution.cookieInvalid",
  session_entry_not_found: "sessionResolution.entryNotFound",
  server_session_not_found: "sessionResolution.sessionNotFound",
  session_token_rejected: "sessionResolution.cookieInvalid",
  session_lookup_failed: "sessionResolution.lookupFailed",
};

type SessionLookup = {
  loginName?: string;
  organization?: string;
  sessionId?: string;
};

type SessionResolutionFailureDetails = {
  reason: SessionResolutionFailure;
  context: Record<string, unknown>;
  error?: unknown;
};

function logSessionResolutionFailure(source: string, { reason, context, error }: SessionResolutionFailureDetails) {
  const code = SESSION_RESOLUTION_LOG_CODES[reason];
  const meta = { code, reason, source, ...context };

  if (reason === "session_lookup_failed") {
    logger.error(`${code}: could not resolve the session`, { ...meta, error });
  } else {
    logger.warn(`${code}: could not resolve the session`, meta);
  }
}

function matchesLookup(entry: Partial<Cookie>, { loginName, organization, sessionId }: SessionLookup): boolean {
  if (organization && entry.organization !== organization) {
    return false;
  }
  if (sessionId) {
    return entry.id === sessionId;
  }
  return loginName ? entry.loginName === loginName : true;
}

// The entries of the cookie before any signature check, undefined when the value is not a JSON array.
function readUnverifiedEntries(value: string | undefined): Partial<Cookie>[] | undefined {
  if (!value) {
    return [];
  }
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed.filter((entry) => !!entry && typeof entry === "object") : undefined;
  } catch {
    return undefined;
  }
}

/**
 * Explains why the `sessions` cookie holds no usable entry for the lookup. Only counts and sizes
 * end up in the context: no token, signature or login name.
 */
async function describeSessionCookieFailure(lookup: SessionLookup): Promise<SessionResolutionFailureDetails> {
  const value = (await cookies()).get("sessions")?.value;
  const cookieHeader = (await headers()).get("cookie") ?? "";

  const unverified = readUnverifiedEntries(value);
  const entries = unverified ?? [];
  const verified = parseAndVerifySessions<Cookie>(value);

  const context: Record<string, unknown> = {
    hasLoginName: !!lookup.loginName,
    hasOrganization: !!lookup.organization,
    hasSessionId: !!lookup.sessionId,
    sessionId: lookup.sessionId,
    cookieEntries: entries.length,
    cookieEntriesDropped: entries.length - verified.length,
    // as the browser sent it (URL-encoded), browsers cap a single cookie at ~4096 bytes
    cookieSize: value ? `sessions=${encodeURIComponent(value)}`.length : 0,
    cookieHeaderSize: cookieHeader.length,
    // more than one means cookies with different domain or path attributes shadow each other
    sessionsCookieCount: cookieHeader.match(/(?:^|;\s*)sessions=/g)?.length ?? 0,
  };

  if (lookup.loginName) {
    const loginName = lookup.loginName;
    // an entry exists for the login name but not for the organization, or only with different casing
    context.entriesMatchingLoginName = verified.filter((entry) => entry.loginName === loginName).length;
    context.entriesMatchingLoginNameIgnoringCase = verified.filter(
      (entry) => entry.loginName?.toLowerCase() === loginName.toLowerCase(),
    ).length;
  }

  if (!value) {
    return { reason: "cookie_missing", context };
  }

  // the wanted entry is in the cookie, but it was dropped for its missing or invalid signature
  const wantedEntryDropped =
    entries.some((entry) => matchesLookup(entry, lookup)) && !verified.some((entry) => matchesLookup(entry, lookup));
  if (!unverified || wantedEntryDropped) {
    return { reason: "cookie_invalid", context };
  }

  return { reason: "session_entry_not_found", context };
}

/**
 * Classifies and logs why no entry of the `sessions` cookie matches the lookup.
 * For callers that looked the cookie entry up themselves and need a session to continue.
 */
export async function reportMissingSessionCookie(source: string, lookup: SessionLookup): Promise<SessionResolutionFailure> {
  const failure = await describeSessionCookieFailure(lookup);
  logSessionResolutionFailure(source, failure);
  return failure.reason;
}

async function fetchSession(
  serviceConfig: ServiceConfig,
  cookie: { id: string; token: string },
): Promise<{ session: Session } | SessionResolutionFailureDetails> {
  try {
    const { session } = await getSession({ serviceConfig, sessionId: cookie.id, sessionToken: cookie.token });
    if (session) {
      return { session };
    }
    return { reason: "server_session_not_found", context: { sessionId: cookie.id } };
  } catch (error) {
    const { Code, ConnectError } = await import("@connectrpc/connect");
    const grpcCode = error instanceof ConnectError ? error.code : undefined;
    const context = { sessionId: cookie.id, grpcCode: grpcCode !== undefined ? Code[grpcCode] : undefined };

    // The `sessions` cookie has no maxAge, so it can outlive the server-side session
    // and reference one the session projection no longer holds, e.g. the user logged out, an
    // admin/API terminated the session, or the user/org/instance was removed. In that
    // case getSession throws `not_found`.
    if (grpcCode === Code.NotFound) {
      return { reason: "server_session_not_found", context };
    }
    // The session still exists but the cookie's token no longer matches it (e.g. the token was
    // rotated by a later session update whose cookie never reached the browser): getSession
    // throws `permission_denied` (QUERY-dsfr3).
    if (grpcCode === Code.PermissionDenied) {
      return { reason: "session_token_rejected", context };
    }
    return { reason: "session_lookup_failed", context, error };
  }
}

/**
 * Resolves the session of a cookie entry. Never throws: every failure is logged with its code
 * and returned as a reason.
 */
export async function resolveSession(
  serviceConfig: ServiceConfig,
  cookie: { id: string; token: string },
  source: string,
): Promise<SessionResolution> {
  const result = await fetchSession(serviceConfig, cookie);
  if ("session" in result) {
    return { session: result.session };
  }

  logSessionResolutionFailure(source, result);
  return { reason: result.reason };
}

type LoadMostRecentSessionParams = {
  serviceConfig: ServiceConfig;
  sessionParams: {
    loginName?: string;
    organization?: string;
  };
};

/**
 * Like `loadMostRecentSession`, for pages that cannot continue without the session:
 * it keeps the reason why the session could not be resolved and logs it.
 */
export async function resolveMostRecentSession(
  { serviceConfig, sessionParams }: LoadMostRecentSessionParams,
  source: string,
): Promise<SessionResolution> {
  const recent = await getMostRecentCookieWithLoginname(sessionParams);

  if (!recent) {
    return { reason: await reportMissingSessionCookie(source, sessionParams) };
  }

  return resolveSession(serviceConfig, recent, source);
}

/**
 * Like `loadSessionById`, for pages that cannot continue without the session:
 * it keeps the reason why the session could not be resolved and logs it.
 */
export async function resolveSessionById(
  serviceConfig: ServiceConfig,
  sessionId: string,
  organization: string | undefined,
  source: string,
): Promise<SessionResolution> {
  const recent = await getSessionCookieById({ sessionId, organization });

  if (!recent) {
    return { reason: await reportMissingSessionCookie(source, { sessionId, organization }) };
  }

  return resolveSession(serviceConfig, recent, source);
}

export async function loadMostRecentSession({
  serviceConfig,
  sessionParams,
}: LoadMostRecentSessionParams): Promise<Session | undefined> {
  const recent = await getMostRecentCookieWithLoginname({
    loginName: sessionParams.loginName,
    organization: sessionParams.organization,
  });

  if (!recent) {
    return undefined;
  }

  return loadSession(serviceConfig, recent);
}

export async function loadSessionById(
  serviceConfig: ServiceConfig,
  sessionId: string,
  organization?: string,
): Promise<Session | undefined> {
  const recent = await getSessionCookieById({ sessionId, organization });

  if (!recent) {
    return undefined;
  }

  return loadSession(serviceConfig, recent);
}

/**
 * Returns undefined when the API no longer knows the session or rejects the cookie's token,
 * instead of letting the error crash the caller's render. Every call site handles an undefined
 * session, either with a graceful fallback or by throwing its own explicit "no session" error.
 * Any other failure is thrown.
 */
export async function loadSession(
  serviceConfig: ServiceConfig,
  cookie: { id: string; token: string },
): Promise<Session | undefined> {
  const result = await fetchSession(serviceConfig, cookie);
  if ("session" in result) {
    return result.session;
  }

  if (result.reason === "session_lookup_failed") {
    throw result.error;
  }

  logSessionResolutionFailure("load_session", result);
  return undefined;
}

/**
 * A session proves genuine (primary) authentication when a primary factor — password,
 * passkey/WebAuthn or IDP intent — has been verified and the session has not expired.
 *
 * A WebAuthn factor only counts as a *primary* factor when it was user-verified
 * (`userVerified` true, i.e. passwordless login). A presence-only WebAuthn assertion is a
 * second-factor (U2F) check and must NOT be treated as primary authentication — mirroring
 * how `mfa-helper` distinguishes "authenticated with passkey" from a 2FA WebAuthn check.
 *
 * This is the shared gate behind both the credential-enrollment guard
 * (`getEnrollmentAuthorizationError`) and the MFA-setup page, so an "identify-only"
 * session (only `factors.user` set, produced by submitting a login name) is never enough
 * to attach a new authenticator (GHSA-45f2-5q3r-xgg6).
 */
export function hasVerifiedPrimaryFactor(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey =
    session?.factors?.webAuthN?.verifiedAt && !!session?.factors?.webAuthN?.userVerified
      ? session?.factors?.webAuthN?.verifiedAt
      : undefined;
  const validIDP = session?.factors?.intent?.verifiedAt;

  const stillValid = session.expirationDate ? timestampDate(session.expirationDate) > new Date() : true;

  const verifiedAt = validPassword || validPasskey || validIDP;
  const valid = !!(verifiedAt && stillValid);

  return { valid, verifiedAt };
}

/**
 * mfa is required, session is not valid anymore (e.g. session expired, user logged out, etc.)
 * to check for mfa for automatically selected session -> const response = await listAuthenticationMethodTypes(userId);
 **/
export async function isSessionValid({
  serviceConfig,
  session,
}: {
  serviceConfig: ServiceConfig;
  session: Session;
}): Promise<boolean> {
  // session can't be checked without user
  if (!session.factors?.user) {
    return false;
  }

  let mfaValid = true;

  // Check if user authenticated via different methods
  const validIDP = session?.factors?.intent?.verifiedAt;
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;

  // Get login settings to determine if MFA is actually required by policy
  const loginSettings = await getLoginSettings({ serviceConfig, organization: session.factors?.user?.organizationId });

  // Use the existing shouldEnforceMFA function to determine if MFA is required
  const isMfaRequired = shouldEnforceMFA(session, loginSettings);

  // Always check auth methods to see if the user has MFA factors configured
  const authMethodTypes = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });

  const authMethods = authMethodTypes.authMethodTypes;
  // Filter to only MFA methods (exclude PASSWORD and PASSKEY)
  const mfaMethods = authMethods?.filter(
    (method) =>
      method === AuthenticationMethodType.TOTP ||
      method === AuthenticationMethodType.OTP_EMAIL ||
      method === AuthenticationMethodType.OTP_SMS ||
      method === AuthenticationMethodType.U2F,
  );

  // A user-verified passkey (possession + user verification) is inherently multi-factor.
  // Mirrors the passkey escape in checkMFAFactors, which never prompts for an additional
  // second factor after a passkey login. Without this, users with a passkey plus a
  // configured TOTP/OTP factor get stuck in a redirect loop on /passkey, because the
  // flow never asks for the second factor but this check would deem the session invalid.
  const hasAuthenticatedWithPasskey = !!session.factors.webAuthN?.verifiedAt && !!session.factors.webAuthN?.userVerified;

  if (mfaMethods && mfaMethods.length > 0) {
    // User has MFA methods configured — they must be verified regardless of policy,
    // unless the session was already authenticated with a user-verified passkey
    const totpValid = mfaMethods.includes(AuthenticationMethodType.TOTP) && !!session.factors.totp?.verifiedAt;
    const otpEmailValid = mfaMethods.includes(AuthenticationMethodType.OTP_EMAIL) && !!session.factors.otpEmail?.verifiedAt;
    const otpSmsValid = mfaMethods.includes(AuthenticationMethodType.OTP_SMS) && !!session.factors.otpSms?.verifiedAt;
    const u2fValid = mfaMethods.includes(AuthenticationMethodType.U2F) && !!session.factors.webAuthN?.verifiedAt;

    mfaValid = hasAuthenticatedWithPasskey || totpValid || otpEmailValid || otpSmsValid || u2fValid;
  } else if (isMfaRequired) {
    // No MFA methods configured, but MFA is forced by policy — check for any verified MFA factors
    const otpEmail = session.factors.otpEmail?.verifiedAt;
    const otpSms = session.factors.otpSms?.verifiedAt;
    const totp = session.factors.totp?.verifiedAt;
    const webAuthN = session.factors.webAuthN?.verifiedAt;

    mfaValid = !!(otpEmail || otpSms || totp || webAuthN);
  }

  // If user has no MFA methods and MFA is not required by policy, mfaValid remains true

  const stillValid = session.expirationDate ? timestampDate(session.expirationDate).getTime() > new Date().getTime() : true;

  if (!stillValid) {
    console.warn(
      "[Session] Session is expired",
      session.expirationDate ? timestampDate(session.expirationDate).toDateString() : "no expiration date",
    );
    return false;
  }

  const validChecks = !!(validPassword || validPasskey || validIDP);

  if (!validChecks) {
    return false;
  }

  if (!mfaValid) {
    console.warn("[Session] MFA is required but not valid");
    return false;
  }

  // Check email verification if EMAIL_VERIFICATION environment variable is enabled
  if (process.env.EMAIL_VERIFICATION === "true") {
    const userResponse = await getUserByID({ serviceConfig, userId: session.factors.user.id });

    const humanUser = userResponse?.user?.type.case === "human" ? userResponse?.user.type.value : undefined;

    if (humanUser && !humanUser.email?.isVerified) {
      console.warn("[Session] Email is not verified");
      return false;
    }
  }

  return true;
}

export async function findValidSession({
  serviceConfig,
  sessions,
  authRequest,
  samlRequest,
  organization,
}: {
  serviceConfig: ServiceConfig;
  sessions: Session[];
  authRequest?: AuthRequest;
  samlRequest?: SAMLRequest;
  organization?: string;
}): Promise<Session | undefined> {
  let sessionsWithHint = sessions.filter((s) => {
    if (authRequest && authRequest.hintUserId) {
      return s.factors?.user?.id === authRequest.hintUserId;
    }
    if (authRequest && authRequest.loginHint) {
      return s.factors?.user?.loginName === authRequest.loginHint;
    }
    if (samlRequest) {
      // SAML requests don't contain user hints like OIDC (hintUserId/loginHint)
      // so we return all sessions for further processing
      return true;
    }
    return true;
  });

  if (organization) {
    sessionsWithHint = sessionsWithHint.filter((s) => s.factors?.user?.organizationId === organization);
  }

  if (sessionsWithHint.length === 0) {
    return undefined;
  }

  // sort by change date descending
  sessionsWithHint.sort((a, b) => {
    const dateA = a.changeDate ? timestampDate(a.changeDate).getTime() : 0;
    const dateB = b.changeDate ? timestampDate(b.changeDate).getTime() : 0;
    return dateB - dateA;
  });

  // return the first valid session according to settings
  for (const session of sessionsWithHint) {
    if (await isSessionValid({ serviceConfig, session })) {
      return session;
    }
  }

  return undefined;
}
