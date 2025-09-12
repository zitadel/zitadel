import { timestampDate } from "@zitadel/client";
import { AuthRequest } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { SAMLRequest } from "@zitadel/proto/zitadel/saml/v2/authorization_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getMostRecentCookieWithLoginname } from "./cookies";
import { getLoginSettings, getSession, getUserByID, listAuthenticationMethodTypes } from "./zitadel";

type LoadMostRecentSessionParams = {
  serviceUrl: string;
  sessionParams: {
    loginName?: string;
    organization?: string;
  };
};

export async function loadMostRecentSession({
  serviceUrl,
  sessionParams,
}: LoadMostRecentSessionParams): Promise<Session | undefined> {
  const recent = await getMostRecentCookieWithLoginname({
    loginName: sessionParams.loginName,
    organization: sessionParams.organization,
  });

  return getSession({
    serviceUrl,
    sessionId: recent.id,
    sessionToken: recent.token,
  }).then((resp: GetSessionResponse) => resp.session);
}

/**
 * mfa is required, session is not valid anymore (e.g. session expired, user logged out, etc.)
 * to check for mfa for automatically selected session -> const response = await listAuthenticationMethodTypes(userId);
 **/
export async function isSessionValid({ serviceUrl, session }: { serviceUrl: string; session: Session }): Promise<boolean> {
  // session can't be checked without user
  if (!session.factors?.user) {
    console.warn("Session has no user");
    return false;
  }

  let mfaValid = true;

  // Get login settings to determine if MFA is actually required by policy
  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: session.factors?.user?.organizationId,
  });

  const isMfaRequired = loginSettings?.forceMfa || loginSettings?.forceMfaLocalOnly;

  // Only enforce MFA validation if MFA is required by policy
  if (isMfaRequired) {
    // Check if any MFA factors are verified
    const otpEmail = session.factors.otpEmail?.verifiedAt;
    const otpSms = session.factors.otpSms?.verifiedAt;
    const totp = session.factors.totp?.verifiedAt;
    const webAuthN = session.factors.webAuthN?.verifiedAt;
    const idp = session.factors.intent?.verifiedAt; // TODO: forceMFA should not consider this as valid factor

    mfaValid = !!(otpEmail || otpSms || totp || webAuthN || idp);

    if (!mfaValid) {
      console.warn("Session has no valid multifactor", session.factors);
    }
  }
  // If MFA is not required by policy, mfaValid remains true

  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const validIDP = session?.factors?.intent?.verifiedAt;

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
    const userResponse = await getUserByID({
      serviceUrl,
      userId: session.factors.user.id,
    });

    const humanUser = userResponse?.user?.type.case === "human" ? userResponse?.user.type.value : undefined;

    if (humanUser && !humanUser.email?.isVerified) {
      console.warn("Session invalid: Email not verified and EMAIL_VERIFICATION is enabled", session.factors.user.id);
      return false;
    }
  }

  return true;
}

export async function findValidSession({
  serviceUrl,
  sessions,
  authRequest,
  samlRequest,
}: {
  serviceUrl: string;
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
      // TODO: do whatever
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
    if (await isSessionValid({ serviceUrl, session })) {
      return session;
    }
  }

  return undefined;
}
