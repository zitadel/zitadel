import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getUserByID, ServiceConfig } from "./zitadel";

export async function checkMFAFactors(
  serviceConfig: ServiceConfig,
  session: Session,
  loginSettings: LoginSettings | undefined,
  authMethods: AuthenticationMethodType[],
  organization?: string,
  requestId?: string,
) {
  const availableMultiFactors = authMethods?.filter(
    (m: AuthenticationMethodType) =>
      m === AuthenticationMethodType.TOTP ||
      m === AuthenticationMethodType.OTP_SMS ||
      m === AuthenticationMethodType.OTP_EMAIL ||
      m === AuthenticationMethodType.U2F,
  );

  const hasAuthenticatedWithPasskey = session.factors?.webAuthN?.verifiedAt && session.factors?.webAuthN?.userVerified;

  // escape further checks if user has authenticated with passkey
  if (hasAuthenticatedWithPasskey) {
    return;
  }

  // if user has not authenticated with passkey and has only one additional mfa factor, redirect to that
  if (availableMultiFactors?.length == 1) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append("organization", organization ?? (session.factors?.user?.organizationId as string));
    }

    const factor = availableMultiFactors[0];
    // if passkey is other method, but user selected password as alternative, perform a login
    if (factor === AuthenticationMethodType.TOTP) {
      return { redirect: `/otp/time-based?` + params };
    } else if (factor === AuthenticationMethodType.OTP_SMS) {
      return { redirect: `/otp/sms?` + params };
    } else if (factor === AuthenticationMethodType.OTP_EMAIL) {
      return { redirect: `/otp/email?` + params };
    } else if (factor === AuthenticationMethodType.U2F) {
      return { redirect: `/u2f?` + params };
    }
  } else if (availableMultiFactors?.length > 1) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append("organization", organization ?? (session.factors?.user?.organizationId as string));
    }

    return { redirect: `/mfa?` + params };
  } else if (shouldEnforceMFA(session, loginSettings) && !availableMultiFactors.length) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      force: "true", // this defines if the mfa is forced in the settings
      checkAfter: "true", // this defines if the check is directly made after the setup
    });

    if (session.id) {
      params.append("sessionId", session.id);
    }

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append("organization", organization ?? (session.factors?.user?.organizationId as string));
    }

    // TODO: provide a way to setup passkeys on mfa page?
    return { redirect: `/mfa/set?` + params };
    /* eslint-disable no-dupe-else-if -- TODO: this branch is unreachable, conditions overlap with the previous branch */
  } else if (
    loginSettings?.mfaInitSkipLifetime &&
    (loginSettings.mfaInitSkipLifetime.nanos > 0 || loginSettings.mfaInitSkipLifetime.seconds > 0) &&
    !availableMultiFactors.length &&
    session?.factors?.user?.id &&
    shouldEnforceMFA(session, loginSettings)
  ) {
    /* eslint-enable no-dupe-else-if */
    const userResponse = await getUserByID({ serviceConfig, userId: session.factors?.user?.id });

    const humanUser = userResponse?.user?.type.case === "human" ? userResponse?.user.type.value : undefined;

    if (humanUser?.mfaInitSkipped) {
      const mfaInitSkippedTimestamp = timestampDate(humanUser.mfaInitSkipped);

      const mfaInitSkipLifetimeMillis =
        Number(loginSettings.mfaInitSkipLifetime.seconds) * 1000 + loginSettings.mfaInitSkipLifetime.nanos / 1000000;
      const currentTime = Date.now();
      const mfaInitSkippedTime = mfaInitSkippedTimestamp.getTime();
      const timeDifference = currentTime - mfaInitSkippedTime;

      if (!(timeDifference > mfaInitSkipLifetimeMillis)) {
        // if the time difference is smaller than the lifetime, skip the mfa setup
        return;
      }
    }

    // the user has never skipped the mfa init but we have a setting so we redirect

    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      force: "false", // this defines if the mfa is not forced in the settings and can be skipped
      checkAfter: "true", // this defines if the check is directly made after the setup
    });

    if (session.id) {
      params.append("sessionId", session.id);
    }

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append("organization", organization ?? (session.factors?.user?.organizationId as string));
    }

    // TODO: provide a way to setup passkeys on mfa page?
    return { redirect: `/mfa/set?` + params };
  }
}

/**
 * Determines if MFA should be enforced based on the authentication method used and login settings
 * @param session - The current session
 * @param loginSettings - The login settings containing MFA enforcement rules
 * @returns true if MFA should be enforced, false otherwise
 */
export function shouldEnforceMFA(session: Session, loginSettings: LoginSettings | undefined): boolean {
  if (!loginSettings) {
    return false;
  }

  // Check if user authenticated with passkey (passkeys are inherently multi-factor)
  const authenticatedWithPasskey = session.factors?.webAuthN?.verifiedAt && session.factors?.webAuthN?.userVerified;

  // If user authenticated with passkey, MFA is not required regardless of settings
  if (authenticatedWithPasskey) {
    return false;
  }

  // If forceMfa is enabled, MFA is required for ALL authentication methods (except passkeys)
  if (loginSettings.forceMfa) {
    return true;
  }

  // If forceMfaLocalOnly is enabled, MFA is only required for local/password authentication
  if (loginSettings.forceMfaLocalOnly) {
    // Check if user authenticated with password (local authentication)
    const authenticatedWithPassword = !!session.factors?.password?.verifiedAt;

    // Check if user authenticated with IDP (external authentication)
    const authenticatedWithIDP = !!session.factors?.intent?.verifiedAt;

    // If user authenticated with IDP, MFA is not required for forceMfaLocalOnly
    if (authenticatedWithIDP) {
      return false;
    }

    // If user authenticated with password, MFA is required for forceMfaLocalOnly
    if (authenticatedWithPassword) {
      return true;
    }
  }

  return false;
}
