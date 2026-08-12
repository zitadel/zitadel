import { isClassifiedError } from "@/lib/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/logger";
import { hasVerifiedPrimaryFactor } from "@/lib/session";
import { checkUserVerification } from "@/lib/verify-helper";
import { listAuthenticationMethodTypes, ServiceConfig } from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

const logger = createLogger("enrollment-guard");

/**
 * Central authorization gate for credential-enrollment entry points.
 *
 * Returns `null` when enrollment is allowed, or an error message string explaining why it
 * is not. Enrollment is allowed when EITHER:
 *  - the session already proves authentication (a verified primary factor, not expired), OR
 *  - the user has no authentication methods yet AND a prior user-verification check passed
 *    (the legitimate email/init-code onboarding path, mirroring `registerPasskeyLink`).
 *
 * This mirrors the gate that `apps/login/src/lib/server/passkeys.ts` already applies to
 * passkey registration, so that U2F, TOTP and OTP enrollment share the exact same rules.
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
    return null;
  }

  // The session is only "identified", not authenticated.
  let authmethods;
  try {
    authmethods = await listAuthenticationMethodTypes({ serviceConfig, userId });
  } catch (error) {
    // Fail closed: without knowing the user's methods (e.g. the user was
    // deleted mid-flow), enrollment must not proceed. Genuine server faults
    // keep throwing so callers surface them.
    if (isClassifiedError(error) && error.isUserError) {
      logger.warn("Could not list authentication methods for enrollment authorization", { grpcCode: error.code });
      return "You have to authenticate or have a valid User Verification Check";
    }
    throw error;
  }

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
