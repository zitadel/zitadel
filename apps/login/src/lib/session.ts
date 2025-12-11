import { timestampDate } from "@zitadel/client";
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

export async function loadMostRecentSession({ serviceConfig, sessionParams }: LoadMostRecentSessionParams): Promise<Session | undefined> {
  const recent = await getMostRecentCookieWithLoginname({
    loginName: sessionParams.loginName,
    organization: sessionParams.organization,
  });

  return getSession({ serviceConfig, sessionId: recent.id, sessionToken: recent.token }).then((resp: GetSessionResponse) => resp.session);
}

/**
 * mfa is required, session is not valid anymore (e.g. session expired, user logged out, etc.)
 * to check for mfa for automatically selected session -> const response = await listAuthenticationMethodTypes(userId);
 **/
export async function isSessionValid({ serviceConfig, session }: { serviceConfig: ServiceConfig; session: Session }): Promise<boolean> {
  // session can't be checked without user
  if (!session.factors?.user) {
    console.warn("Session has no user");
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

  // Only enforce MFA validation if MFA is required by policy
  if (isMfaRequired) {
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

    if (mfaMethods && mfaMethods.length > 0) {
      // Check if any of the configured MFA methods have been verified
      const totpValid = mfaMethods.includes(AuthenticationMethodType.TOTP) && !!session.factors.totp?.verifiedAt;
      const otpEmailValid =
        mfaMethods.includes(AuthenticationMethodType.OTP_EMAIL) && !!session.factors.otpEmail?.verifiedAt;
      const otpSmsValid = mfaMethods.includes(AuthenticationMethodType.OTP_SMS) && !!session.factors.otpSms?.verifiedAt;
      const u2fValid = mfaMethods.includes(AuthenticationMethodType.U2F) && !!session.factors.webAuthN?.verifiedAt;

      mfaValid = totpValid || otpEmailValid || otpSmsValid || u2fValid;

      if (!mfaValid) {
        console.warn("Session has no valid MFA factor. Configured methods:", mfaMethods, "Session factors:", {
          totp: session.factors.totp?.verifiedAt,
          otpEmail: session.factors.otpEmail?.verifiedAt,
          otpSms: session.factors.otpSms?.verifiedAt,
          webAuthN: session.factors.webAuthN?.verifiedAt,
        });
      }
    } else {
      // No specific MFA methods configured, but MFA is forced - check for any verified MFA factors
      // (excluding IDP which should be handled separately)
      const otpEmail = session.factors.otpEmail?.verifiedAt;
      const otpSms = session.factors.otpSms?.verifiedAt;
      const totp = session.factors.totp?.verifiedAt;
      const webAuthN = session.factors.webAuthN?.verifiedAt;
      // Note: Removed IDP (session.factors.intent?.verifiedAt) as requested

      mfaValid = !!(otpEmail || otpSms || totp || webAuthN);
      if (!mfaValid) {
        console.warn("Session has no valid multifactor", session.factors);
      }
    }
  }

  // If MFA is not required by policy, mfaValid remains true

  const stillValid = session.expirationDate ? timestampDate(session.expirationDate).getTime() > new Date().getTime() : true;

  if (!stillValid) {
    console.warn(
      "Session is expired",
      session.expirationDate ? timestampDate(session.expirationDate).toDateString() : "no expiration date",
    );
    return false;
  }

  const validChecks = !!(validPassword || validPasskey || validIDP);

  if (!validChecks) {
    return false;
  }

  if (!mfaValid) {
    return false;
  }

  // Check email verification if EMAIL_VERIFICATION environment variable is enabled
  if (process.env.EMAIL_VERIFICATION === "true") {
    const userResponse = await getUserByID({ serviceConfig, userId: session.factors.user.id });

    const humanUser = userResponse?.user?.type.case === "human" ? userResponse?.user.type.value : undefined;

    if (humanUser && !humanUser.email?.isVerified) {
      console.warn("Session invalid: Email not verified and EMAIL_VERIFICATION is enabled", session.factors.user.id);
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
}: {
  serviceConfig: ServiceConfig;
  sessions: Session[];
  authRequest?: AuthRequest;
  samlRequest?: SAMLRequest;
}): Promise<Session | undefined> {
  const sessionsWithHint = sessions.filter((s) => {
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
