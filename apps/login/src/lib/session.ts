import { Timestamp, timestampDate } from "@zitadel/client";
import { AuthRequest } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { SAMLRequest } from "@zitadel/proto/zitadel/saml/v2/authorization_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getMostRecentCookieWithLoginname } from "./cookies";
import { shouldEnforceMFA } from "./verify-helper";
import { getLoginSettings, getSession, getUserByID, listAuthenticationMethodTypes, ServiceConfig } from "./zitadel";

type LoadMostRecentSessionParams = {
  serviceConfig: ServiceConfig;
  sessionParams: {
    loginName?: string;
    organization?: string;
  };
};

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

  return getSession({ serviceConfig, sessionId: recent.id, sessionToken: recent.token })
    .then((resp: GetSessionResponse) => resp.session)
    .catch(async (error) => {
      const { Code, ConnectError } = await import("@connectrpc/connect");
      const isNotFound = error instanceof ConnectError && error.code === Code.NotFound;

      // The `sessions` cookie has no maxAge, so it can outlive the server-side session
      // and reference one the session projection no longer holds — e.g. the user logged out, an
      // admin/API terminated the session, or the user/org/instance was removed. In that
      // case getSession throws `not_found`. Treat it like the no-cookie case above and
      // return undefined instead of letting the error crash the caller's render. Every
      // call site already handles an undefined session — either with a graceful fallback
      // or by throwing its own explicit "no session" error.
      if (isNotFound) {
        console.warn("[Session] Could not load most recent session", error);
        return undefined;
      }

      throw error;
    });
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
