import "server-only";

import { createLogger } from "@/lib/server/logger";
import { hasVerifiedPrimaryFactor, hasVerifiedSecondFactor, SECOND_FACTOR_METHODS } from "@/lib/server/session";
import { checkUserVerification } from "@/lib/server/verify-helper";
import { getLoginSettings, listAuthenticationMethodTypes, ServiceConfig } from "@/lib/server/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { PasskeysType, SecondFactorType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

const logger = createLogger("enrollment-guard");

export const SECOND_FACTOR_REQUIRED_ERROR = "You have to verify your second factor first";
export const PASSKEYS_NOT_ALLOWED_ERROR = "Passkeys are not allowed";
export const SECOND_FACTOR_NOT_ALLOWED_ERROR = "This second factor is not allowed";

/**
 * Central authorization gate for credential-enrollment entry points.
 *
 * Returns `null` when enrollment is allowed, or an error message string explaining why it
 * is not. Enrollment is allowed when EITHER:
 *  - the session proves authentication (a verified primary factor, not expired) AND, if the
 *    user already has a second factor configured, the session has also verified one of those
 *    (or was authenticated with a user-verified passkey). Without the second part, a session
 *    that only passed the first factor could enroll a new authenticator and use it to bypass
 *    the user's existing second factor. Users without a second factor (e.g. the forced
 *    first-time MFA setup) only need the primary factor, OR
 *  - the user has no authentication methods yet AND a prior user-verification check passed
 *    (the legitimate email/init-code onboarding path, mirroring `registerPasskeyLink`).
 *
 * This only decides whether the session may add an authenticator to the account at all. It
 * doesn't know which authenticator is being added and therefore doesn't check the login
 * settings. Setting up a passkey, TOTP, OTP SMS / email or U2F must use
 * `getSetupAuthorizationError`, which adds that check. IdP linking uses this guard directly;
 * which IdPs may be used is checked by the IdP flow itself.
 */
export async function getEnrollmentAuthorizationError({
  serviceConfig,
  session,
  userId,
}: {
  serviceConfig: ServiceConfig;
  session: Partial<Session>;
  userId: string;
}): Promise<string | null> {
  if (hasVerifiedPrimaryFactor(session).valid) {
    // A user-verified passkey is multi-factor on its own.
    if (session.factors?.webAuthN?.verifiedAt && session.factors.webAuthN.userVerified) {
      return null;
    }

    const authmethods = await listAuthenticationMethodTypes({ serviceConfig, userId });
    const hasSecondFactor = authmethods.authMethodTypes.some((method) => SECOND_FACTOR_METHODS.includes(method));

    if (hasSecondFactor && !hasVerifiedSecondFactor(session, authmethods.authMethodTypes)) {
      return SECOND_FACTOR_REQUIRED_ERROR;
    }

    return null;
  }

  // The session is only "identified", not authenticated.
  const authmethods = await listAuthenticationMethodTypes({ serviceConfig, userId });

  // If the user already has any auth method configured, enrollment requires a real
  // authentication (or a valid user-verification check) — reject the identify-only session.
  if (authmethods.authMethodTypes.length !== 0) {
    return "You have to authenticate or have a valid User Verification Check";
  }

  // No auth methods yet: allow only if a prior user-verification check was completed
  // (e.g. the user just verified their email / redeemed an invite in this browser).
  const hasValidUserVerificationCheck = await checkUserVerification(userId);

  logger.info("hasValidUserVerificationCheck", { hasValidUserVerificationCheck });
  if (!hasValidUserVerificationCheck) {
    return "User Verification Check has to be done";
  }

  return null;
}

export type SetupMethod = "passkey" | "totp" | "u2f" | "otp_sms" | "otp_email";

const SECOND_FACTOR_SETTING: Record<Exclude<SetupMethod, "passkey">, SecondFactorType> = {
  totp: SecondFactorType.TOTP,
  u2f: SecondFactorType.U2F,
  otp_sms: SecondFactorType.OTP_SMS,
  otp_email: SecondFactorType.OTP_EMAIL,
};

/**
 * Authorization gate for setting up a passkey, TOTP, OTP SMS / email or U2F during login.
 * Two checks, both must pass:
 *  1. the session may add an authenticator at all (`getEnrollmentAuthorizationError`), and
 *  2. the login settings of the user's organization allow this kind of authenticator
 *     (`passkeysType` for passkeys, `secondFactors` for TOTP, U2F, OTP SMS and OTP email).
 * Check 2 matters because the setup pages can be opened directly by their URL, not only from
 * the choices the login offers.
 *
 * Returns `null` when the setup is allowed, or an error message string.
 */
export async function getSetupAuthorizationError({
  serviceConfig,
  session,
  userId,
  method,
}: {
  serviceConfig: ServiceConfig;
  session: Partial<Session>;
  userId: string;
  method: SetupMethod;
}): Promise<string | null> {
  const enrollmentError = await getEnrollmentAuthorizationError({ serviceConfig, session, userId });
  if (enrollmentError) {
    return enrollmentError;
  }

  const loginSettings = await getLoginSettings({
    serviceConfig,
    organization: session.factors?.user?.organizationId,
  });
  if (!loginSettings) {
    return "Could not load login settings";
  }

  if (method === "passkey") {
    return loginSettings.passkeysType === PasskeysType.ALLOWED ? null : PASSKEYS_NOT_ALLOWED_ERROR;
  }

  if (!loginSettings.secondFactors.includes(SECOND_FACTOR_SETTING[method])) {
    return SECOND_FACTOR_NOT_ALLOWED_ERROR;
  }
  return null;
}

/**
 * Maps an error of `getSetupAuthorizationError` to a key of the `error` translations, so that
 * the setup pages can show a translated message. Any other reason (no or an expired
 * authentication, settings not loadable) asks the user to sign in again.
 */
export function setupErrorTranslationKey(
  error: string,
): "secondFactorRequired" | "authenticatorNotAllowed" | "sessionExpired" {
  switch (error) {
    case SECOND_FACTOR_REQUIRED_ERROR:
      return "secondFactorRequired";
    case PASSKEYS_NOT_ALLOWED_ERROR:
    case SECOND_FACTOR_NOT_ALLOWED_ERROR:
      return "authenticatorNotAllowed";
    default:
      return "sessionExpired";
  }
}
