"use server";

import { isClassifiedError } from "@/lib/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/logger";
import { recordAuthAttempt, recordAuthFailure, recordAuthSuccess } from "@/lib/metrics";
import { setSessionAndUpdateCookie } from "@/lib/server/cookie";
import { getLoginSettings } from "@/lib/zitadel";
import { Code, create, Duration } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { completeFlowOrGetUrl } from "../client";
import { getSessionCookieById, getSessionCookieByLoginName } from "../cookies";
import { getServiceConfig } from "../service-url";

const logger = createLogger("recovery-code");

export type SendRecoveryCodeCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  requestId?: string;
  code: string;
};

/**
 * Maps an error thrown while checking a recovery code to a translation key and a metrics reason.
 *
 * The backend translates error messages before building the connect status, so the resulting
 * message is `"<translated text> (<error id>)"`. Matching on the raw `Errors.User.…` key is
 * therefore unreliable. We map primarily on the gRPC code and only use the stable error ids
 * (and the raw keys, for untranslated deployments) as hints.
 */
function mapRecoveryCodeError(error: unknown): { key: string; reason: string } {
  // defensive: passwordAttemptsHandler rethrows a plain object for credential check details
  if (error && typeof error === "object" && "failedAttempts" in error) {
    return { key: "verify.errors.invalidCode", reason: "invalid_code" };
  }

  if (isClassifiedError(error)) {
    const hint = `${error.rawMessage ?? ""} ${error.message ?? ""}`;

    if (error.code === Code.InvalidArgument) {
      return { key: "verify.errors.invalidCode", reason: "invalid_code" };
    }

    if (error.code === Code.FailedPrecondition) {
      if (/Errors\.User\.Locked|COMMAND-2w6oa|COMMAND-ASV12/.test(hint)) {
        return { key: "verify.errors.userLocked", reason: "user_locked" };
      }
      if (/RecoveryCodes\.NotReady|COMMAND-84rgg/.test(hint)) {
        return { key: "verify.errors.notReady", reason: "not_ready" };
      }
    }
  }

  return { key: "verify.errors.couldNotVerifyCode", reason: "session_update_failed" };
}

/**
 * Verifies a recovery code as a second factor on an existing session.
 *
 * Unlike `updateOrCreateSession`, this never falls back to creating a new session: a failed
 * check already counts towards the lockout policy, so retrying it on a fresh session would
 * burn a second attempt.
 */
export async function sendRecoveryCode(
  command: SendRecoveryCodeCommand,
): Promise<{ error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const t = await getTranslations("recoveryCode");

  recordAuthAttempt("recovery_code", command.organization);

  const sessionCookie = command.sessionId
    ? await getSessionCookieById({ sessionId: command.sessionId, organization: command.organization })
    : await getSessionCookieByLoginName({ loginName: command.loginName, organization: command.organization });

  if (!sessionCookie) {
    recordAuthFailure("recovery_code", "session_not_found", command.organization);
    return { error: t("verify.errors.couldNotFindSession") };
  }

  const loginSettings = await getLoginSettings({
    serviceConfig,
    organization: command.organization ?? sessionCookie.organization,
  });

  let lifetime = loginSettings?.secondFactorCheckLifetime;

  if (!lifetime || !lifetime.seconds) {
    logger.warn("No second factor lifetime provided, defaulting to 24 hours");
    lifetime = {
      seconds: BigInt(60 * 60 * 24), // default to 24 hours
      nanos: 0,
    } as Duration;
  }

  const checks = create(ChecksSchema, {
    recoveryCode: { code: command.code.trim() },
  });

  let session;
  try {
    session = await setSessionAndUpdateCookie({
      recentCookie: sessionCookie,
      checks,
      requestId: command.requestId,
      lifetime,
    });
  } catch (error) {
    const { key, reason } = mapRecoveryCodeError(error);

    if (isClassifiedError(error) && error.isUserError) {
      logger.warn("Could not verify recovery code (client error)", {
        grpcCode: error.code,
        httpStatus: error.httpStatus,
      });
    } else {
      logger.error("Could not verify recovery code", { error });
    }

    recordAuthFailure("recovery_code", reason, command.organization);
    return { error: t(key) };
  }

  if (!session?.factors?.user?.loginName) {
    recordAuthFailure("recovery_code", "session_invalid", command.organization);
    return { error: t("verify.errors.couldNotVerifyCode") };
  }

  const organization = command.organization ?? session.factors.user.organizationId;

  const result = await completeFlowOrGetUrl(
    command.requestId && session.id
      ? {
          sessionId: session.id,
          requestId: command.requestId,
          organization,
        }
      : {
          loginName: session.factors.user.loginName,
          organization,
        },
    loginSettings?.defaultRedirectUri,
  );

  if (result && "error" in result) {
    recordAuthFailure("recovery_code", "flow_error", command.organization);
    return result;
  }

  if (result) {
    recordAuthSuccess("recovery_code", command.organization);
    return result;
  }

  recordAuthFailure("recovery_code", "navigation_failed", command.organization);
  return { error: t("verify.errors.couldNotDetermineRedirect") };
}
