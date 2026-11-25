"use server";

import { completeFlowOrGetUrl } from "@/lib/actions/flow";
import { hasZitadelErrorId, isRejectedByUser, zitadelErrorId, ZitadelErrorIds } from "@/lib/server/api-error";
import { getOTPSMSResendCooldownSeconds, getOTPSMSRetryAfterSeconds, SessionUpdatedError } from "@/lib/server/challenges";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "@/lib/server/cookies";
import { isClassifiedError } from "@/lib/server/grpc/interceptors/error-classification";
import { getPublicHost, getPublicHostWithProtocol } from "@/lib/server/host";
import { createLogger } from "@/lib/server/logger";
import { shouldEnforceMFA } from "@/lib/server/mfa-helper";
import { reauthenticate } from "@/lib/server/reauthenticate";
import { getServiceConfig } from "@/lib/server/service-url";
import { hasVerifiedPrimaryFactor, isSessionValid, loadSession } from "@/lib/server/session";
import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "@/lib/server/session-cookie";
import {
  deleteSession,
  getLoginSettings,
  getSecuritySettings,
  humanMFAInitSkipped,
  listUsers,
  ServiceConfig,
} from "@/lib/server/zitadel";
import { Code, create, Duration } from "@zitadel/client";
import {
  Challenges,
  RequestChallenges,
  RequestChallenges_OTPEmail_SendCodeSchema,
} from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

const logger = createLogger("session");

/**
 * Server actions are public endpoints, so every argument is caller-controlled. The session an
 * action acts on is therefore always resolved from the caller's own signed `sessions` cookie and
 * loaded from ZITADEL with that cookie's token, never taken from the request body.
 */
async function loadCallerSession({
  serviceConfig,
  sessionId,
  loginName,
  organization,
}: {
  serviceConfig: ServiceConfig;
  sessionId?: string;
  loginName?: string;
  organization?: string;
}): Promise<Session | undefined> {
  const sessionCookie = sessionId
    ? await getSessionCookieById({ sessionId })
    : loginName
      ? await getSessionCookieByLoginName({ loginName, organization })
      : undefined;

  if (!sessionCookie) {
    return undefined;
  }

  return loadSession(serviceConfig, sessionCookie);
}

export async function skipMFAAndContinueWithNextUrl({
  requestId,
  loginName,
  sessionId,
  organization,
}: {
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
}): Promise<{ redirect: string } | { error: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const session = await loadCallerSession({ serviceConfig, sessionId, loginName, organization });
  const userId = session?.factors?.user?.id;

  // only a user who just authenticated may skip the MFA setup prompt, and only for themselves
  if (!session || !userId || !hasVerifiedPrimaryFactor(session).valid) {
    return { error: "Could not skip MFA and continue" };
  }

  const userLoginSettings = await getLoginSettings({ serviceConfig, organization: session.factors?.user?.organizationId });

  if (shouldEnforceMFA(session, userLoginSettings)) {
    return { error: "Could not skip MFA and continue" };
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: organization });

  await humanMFAInitSkipped({ serviceConfig, userId });

  if (requestId && session.id) {
    return completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: requestId,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  } else if (session.factors?.user?.loginName) {
    return completeFlowOrGetUrl(
      {
        loginName: session.factors.user.loginName,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  return { error: "Could not skip MFA and continue" };
}

export type ContinueWithSessionCommand = { sessionId: string; requestId?: string };

export async function continueWithSession({ sessionId, requestId }: ContinueWithSessionCommand) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const t = await getTranslations("error");

  const session = await loadCallerSession({ serviceConfig, sessionId });

  if (!session?.factors?.user) {
    return { error: t("couldNotContinueSession") };
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: session.factors.user.organizationId });

  // Validate session (including MFA) before completing the flow
  const valid = await isSessionValid({ serviceConfig, session });

  if (!valid) {
    logger.warn("continueWithSession: session is not valid (e.g. MFA not completed), redirecting to re-authenticate", {
      sessionId: session.id,
    });

    // Redirect user to re-authenticate (will route to MFA page if password is still valid)
    return (await reauthenticate(session, requestId)) ?? { error: t("couldNotContinueSession") };
  }

  if (requestId && session.id) {
    return completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: requestId,
        organization: session.factors.user.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  return completeFlowOrGetUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );
}

export type UpdateSessionCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  checks?: Checks;
  requestId?: string;
  challenges?: RequestChallenges;
};

/**
 * This action only completes second factors and passkeys. The primary checks (user, password,
 * IDP intent) run through their own actions, which enforce the login policy first, so they are
 * dropped here instead of being forwarded to the session API.
 */
function pickSecondFactorChecks(checks?: Checks): Checks | undefined {
  if (!checks) {
    return undefined;
  }

  return create(ChecksSchema, {
    webAuthN: checks.webAuthN,
    totp: checks.totp,
    otpSms: checks.otpSms,
    otpEmail: checks.otpEmail,
  });
}

/**
 * The WebAuthN relying party and the link in the OTP email are always derived from the host this
 * request was served on (the same host that is forwarded to and validated by ZITADEL), never taken
 * from the caller. Otherwise anyone could make ZITADEL send a genuine OTP email whose link points
 * to a host of their choice.
 */
function bindChallengesToHost(
  challenges: RequestChallenges | undefined,
  hostWithProtocol: string,
  requestId?: string,
): RequestChallenges | undefined {
  if (!challenges) {
    return undefined;
  }

  const bound = { ...challenges };

  if (bound.webAuthN) {
    bound.webAuthN = { ...bound.webAuthN, domain: new URL(hostWithProtocol).hostname };
  }

  if (bound.otpEmail) {
    const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";
    let urlTemplate = `${hostWithProtocol}${basePath}/otp/email?code={{.Code}}&userId={{.UserID}}&sessionId={{.SessionID}}`;

    if (requestId) {
      urlTemplate += `&requestId=${encodeURIComponent(requestId)}`;
    }

    bound.otpEmail = {
      ...bound.otpEmail,
      deliveryType: { case: "sendCode", value: create(RequestChallenges_OTPEmail_SendCodeSchema, { urlTemplate }) },
    };
  }

  return bound;
}

export async function updateOrCreateSession(options: UpdateSessionCommand) {
  const { loginName, sessionId, organization, requestId } = options;
  const checks = pickSecondFactorChecks(options.checks);

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  const t = await getTranslations("verify.errors");

  if (!host) {
    return { error: "Could not get host" }; // Technical error, maybe leave or translate if key exists
  }

  const challenges = bindChallengesToHost(options.challenges, getPublicHostWithProtocol(_headers), requestId);

  let recentSession = sessionId
    ? await getSessionCookieById({ sessionId })
    : loginName
      ? await getSessionCookieByLoginName({ loginName, organization })
      : await getMostRecentSessionCookie();

  if (!recentSession) {
    if (!loginName) {
      return { error: t("couldNotFindSession") };
    }

    const checks = create(ChecksSchema, {
      user: { search: { case: "loginName", value: loginName } },
    });

    // the challenges are requested by the setSession call below, requesting them here as well would send every code twice
    const result = await createSessionAndUpdateCookie({
      checks,
      requestId,
    }).catch((error) => {
      if (isClassifiedError(error) && error.isUserError) {
        logger.warn("Could not create session (client error)", { grpcCode: error.code, httpStatus: error.httpStatus });
      } else {
        logger.error("Could not create session (server error)", { error });
      }
      return undefined;
    });

    if (result && "sessionCookie" in result) {
      recentSession = result.sessionCookie;
    }

    if (!recentSession) {
      return {
        error: t("couldNotFindSession"),
      };
    }
  }

  if (challenges?.otpSms) {
    const retryAfterSeconds = getOTPSMSRetryAfterSeconds(recentSession.otpSmsSentTs);
    if (retryAfterSeconds > 0) {
      // no error message: the UI shows the countdown on the resend button instead
      return { retryAfterSeconds };
    }
  }

  // The lifetime always comes from the login policy of the session's user (the cookie entry is
  // signed and its organization set from the session), never from the caller.
  const loginSettings = await getLoginSettings({ serviceConfig, organization: recentSession.organization ?? organization });

  let lifetime = checks?.webAuthN
    ? loginSettings?.multiFactorCheckLifetime // TODO different lifetime for webauthn u2f/passkey
    : checks?.totp || checks?.otpEmail || checks?.otpSms
      ? loginSettings?.secondFactorCheckLifetime
      : undefined;

  if (!lifetime || !lifetime.seconds) {
    logger.warn("No lifetime provided for session, defaulting to 24 hours");
    lifetime = {
      seconds: BigInt(60 * 60 * 24), // default to 24 hours
      nanos: 0,
    } as Duration;
  }

  const hasCodeCheck = !!(checks?.totp || checks?.otpSms || checks?.otpEmail);

  // A rejected check (most often a mistyped code) is answered with a message for the form: a throw
  // would fail the whole action with a 500 although nothing is wrong on the server. Server faults,
  // and the API rejecting the login's own credentials, are rethrown so they stay visible.
  const rejectedCheckError = (error: unknown): { error: string } => {
    if (!isRejectedByUser(error)) {
      throw error;
    }

    logger.warn("Session check rejected", { grpcCode: error.code, errorId: zitadelErrorId(error) });

    if (hasZitadelErrorId(error, ZitadelErrorIds.userLocked)) {
      return { error: t("accountLocked") };
    }

    if (hasCodeCheck) {
      if (hasZitadelErrorId(error, ZitadelErrorIds.totpReused)) {
        return { error: t("codeReused") };
      }

      if (hasZitadelErrorId(error, ZitadelErrorIds.codeExpired)) {
        return { error: t("codeExpired") };
      }

      if (error.code === Code.InvalidArgument) {
        return { error: t("invalidCode") };
      }
    }

    return { error: t("couldNotUpdateSession") };
  };

  let session;
  try {
    session = await setSessionAndUpdateCookie({
      recentCookie: recentSession,
      checks,
      challenges,
      requestId,
      lifetime,
    });
  } catch (error) {
    // the code was already sent for this session: creating a new session with the same challenge would send
    // a second code, and the first one could not be verified against the new session
    if (error instanceof SessionUpdatedError && (challenges?.otpSms || challenges?.otpEmail)) {
      throw error;
    }

    // A new session only helps when this one is gone: deleted or expired server-side while the
    // cookie survived (NotFound), or its token is no longer accepted (PermissionDenied). Any other
    // failure means the check was rejected, and a new session would run the same check again.
    const sessionGone = isClassifiedError(error) && (error.code === Code.NotFound || error.code === Code.PermissionDenied);

    if (!sessionGone) {
      return rejectedCheckError(error);
    }

    const loginNameForCreation = options.loginName || recentSession?.loginName;
    const orgForCreation = options.organization || recentSession?.organization;

    if (!loginNameForCreation) {
      return { error: t("couldNotFindSession") };
    }

    const users = await listUsers({
      serviceConfig,
      loginName: loginNameForCreation,
      organizationId: orgForCreation,
    });

    const userId = users.details?.totalResult === BigInt(1) ? users.result[0]?.userId : undefined;

    if (!userId) {
      return { error: t("couldNotFindSession") };
    }

    // List the fields instead of spreading `checks`: it is already a message, and create() returns
    // a message unchanged, so a plain `user` object spread into it would never be converted and
    // fails to serialize ("cannot use field zitadel.session.v2.CheckUser.user_id").
    const newChecks = create(ChecksSchema, {
      user: { search: { case: "userId", value: userId } },
      webAuthN: checks?.webAuthN,
      totp: checks?.totp,
      otpSms: checks?.otpSms,
      otpEmail: checks?.otpEmail,
    });

    try {
      const result = await createSessionAndUpdateCookie({
        checks: newChecks,
        requestId,
        lifetime,
        challenges,
      });
      // @ts-ignore
      session = { ...result.session, challenges: result.challenges };
    } catch (createError) {
      return rejectedCheckError(createError);
    }
  }

  if (!session || ("error" in session && session.error)) {
    return { error: t("couldNotUpdateSession") };
  }

  // @ts-ignore
  const challengeResponse: Challenges | undefined = session.challenges;

  return {
    sessionId: session.id,
    factors: session.factors,
    // Only the WebAuthN challenge may reach the browser. `sanitizeChallenges` already stops
    // `returnCode` from being requested, so `otpSms`/`otpEmail` should always be unset here —
    // keeping the projection explicit means no OTP code can leak through this boundary even
    // if that ever regresses (GHSA-3gwm-5wx8-4gm6).
    challenges: challengeResponse?.webAuthN ? { webAuthN: challengeResponse.webAuthN } : undefined,
    // lets the UI count down until the next OTP SMS may be requested
    retryAfterSeconds: challenges?.otpSms ? getOTPSMSResendCooldownSeconds() : undefined,
  };
}

type ClearSessionOptions = {
  sessionId: string;
};

export async function clearSession(options: ClearSessionOptions): Promise<{ error?: string } | void> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const { sessionId } = options;

  const sessionCookie = await getSessionCookieById({ sessionId });

  if (!sessionCookie) {
    return;
  }

  try {
    const deleteResponse = await deleteSession({
      serviceConfig,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    });
    if (!deleteResponse) {
      throw new Error("Could not delete session");
    }
  } catch (error) {
    // The backend verifies the cookie's token before it looks at the session
    // state, so PermissionDenied is what a session that no longer exists answers
    // (nothing to verify against). It is also what a live session with a stale
    // cookie token would answer; we deliberately prune the cookie entry in that
    // case too: the user explicitly asked to remove the account from this
    // browser, and without its token the server-side session cannot be used by
    // anyone and simply expires. Every other failure keeps the entry and is
    // reported so the caller can show it.
    // (A session terminated by an RP-initiated logout keeps its token id, so
    // deleting it again succeeds and never reaches this branch.)
    if (!isClassifiedError(error) || error.code !== Code.PermissionDenied) {
      logger.error("clearSession: could not delete session", { sessionId: sessionCookie.id, error });
      const t = await getTranslations("error");
      return { error: t("couldNotClearSession") };
    }

    logger.warn("clearSession: session rejected the cookie token (gone or stale), pruning cookie entry", {
      sessionId: sessionCookie.id,
    });
  }

  // Only needed for the cookie rewrite, so it must not gate the delete above.
  const securitySettings = await getSecuritySettings({ serviceConfig });
  const iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;

  await removeSessionFromCookie({ session: sessionCookie, iFrameEnabled });
}
