"use server";

import { isClassifiedError } from "@/lib/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/logger";
import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "@/lib/server/cookie";
import {
  deleteSession,
  getLoginSettings,
  getSecuritySettings,
  humanMFAInitSkipped,
  listAuthenticationMethodTypes,
  listUsers,
} from "@/lib/zitadel";
import { Code, create, Duration, timestampDate } from "@zitadel/client";
import { Challenges, RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { completeFlowOrGetUrl } from "../client";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "../cookies";
import { getServiceConfig } from "../service-url";
import {
  hasVerifiedPrimaryFactor,
  isSessionValid,
  reportMissingSessionCookie,
  resolveSession,
  SESSION_RESOLUTION_MESSAGE_KEYS,
} from "../session";
import { getPublicHost } from "./host";
import { sendLoginname } from "./loginname";

const logger = createLogger("session");

export async function skipMFAAndContinueWithNextUrl({
  userId,
  requestId,
  loginName,
  sessionId,
  organization,
}: {
  userId: string;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
}): Promise<{ redirect: string } | { error: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const loginSettings = await getLoginSettings({ serviceConfig, organization: organization });

  await humanMFAInitSkipped({ serviceConfig, userId });

  if (requestId && sessionId) {
    return completeFlowOrGetUrl(
      {
        sessionId: sessionId,
        requestId: requestId,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  } else if (loginName) {
    return completeFlowOrGetUrl(
      {
        loginName: loginName,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  return { error: "Could not skip MFA and continue" };
}

export type ContinueWithSessionCommand = Session & { requestId?: string };

export async function continueWithSession({ requestId, ...session }: ContinueWithSessionCommand) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const t = await getTranslations("error");

  if (!session.factors?.user) {
    return { error: t("couldNotContinueSession") };
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: session.factors.user.organizationId });

  // Validate session (including MFA) before completing the flow
  const valid = await isSessionValid({ serviceConfig, session: session as Session });

  if (!valid) {
    logger.warn("continueWithSession: session is not valid (e.g. MFA not completed), redirecting to re-authenticate", {
      sessionId: session.id,
    });

    // Redirect user to re-authenticate (will route to MFA page if password is still valid)
    const res = await sendLoginname({
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
      requestId: requestId,
    });

    if (res && "redirect" in res && res.redirect) {
      return { redirect: res.redirect };
    }

    if (res && "samlData" in res && res.samlData) {
      return { samlData: res.samlData };
    }

    return { error: t("couldNotContinueSession") };
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
  lifetime?: Duration;
};

// A second factor can only complete a session that already holds a verified first factor.
function isSecondFactorStep(checks?: Checks, challenges?: RequestChallenges): boolean {
  return !!(checks?.totp || checks?.otpSms || checks?.otpEmail || challenges?.otpSms || challenges?.otpEmail);
}

export async function updateOrCreateSession(options: UpdateSessionCommand) {
  let { loginName, sessionId, organization, checks, requestId, challenges, lifetime } = options;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  const t = await getTranslations("verify.errors");
  const tError = await getTranslations("error");

  if (!host) {
    return { error: "Could not get host" }; // Technical error, maybe leave or translate if key exists
  }

  if (challenges && challenges.webAuthN && !challenges.webAuthN.domain) {
    const [hostname] = host.split(":");

    challenges.webAuthN.domain = hostname;
  }

  let recentSession = sessionId
    ? await getSessionCookieById({ sessionId })
    : loginName
      ? await getSessionCookieByLoginName({ loginName, organization })
      : await getMostRecentSessionCookie();

  const secondFactorStep = isSecondFactorStep(checks, challenges);

  if (!recentSession) {
    // A new session only holds the user check. That starts a passkey flow, but it can never complete
    // a second factor step: the verified first factor of the session that got lost would be missing.
    if (!loginName || secondFactorStep) {
      const reason = await reportMissingSessionCookie("update_session", { loginName, organization, sessionId });
      // the reason stays in the logs, the caller only learns that retrying cannot succeed
      return { error: tError(SESSION_RESOLUTION_MESSAGE_KEYS[reason]), sessionLost: true };
    }

    const checks = create(ChecksSchema, {
      user: { search: { case: "loginName", value: loginName } },
    });

    const result = await createSessionAndUpdateCookie({
      checks,
      challenges,
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

  const loginSettings = await getLoginSettings({ serviceConfig, organization });

  if (!lifetime) {
    lifetime = checks?.webAuthN
      ? loginSettings?.multiFactorCheckLifetime // TODO different lifetime for webauthn u2f/passkey
      : checks?.otpEmail || checks?.otpSms
        ? loginSettings?.secondFactorCheckLifetime
        : undefined;
  }

  if (!lifetime || !lifetime.seconds) {
    logger.warn("No lifetime provided for session, defaulting to 24 hours");
    lifetime = {
      seconds: BigInt(60 * 60 * 24), // default to 24 hours
      nanos: 0,
    } as Duration;
  }

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
    // Creating a new session below replaces the cookie entry and drops the factors of the existing session.
    const existing = await resolveSession(serviceConfig, recentSession, "update_session");
    const expired = existing.session?.expirationDate ? timestampDate(existing.session.expirationDate) < new Date() : false;

    // The session is still usable, so the update itself was rejected (e.g. a wrong code). The caller
    // reports that, and a session with a verified first factor is never replaced by one without.
    if (existing.session && !expired && (secondFactorStep || hasVerifiedPrimaryFactor(existing.session).verifiedAt)) {
      throw error;
    }

    // The session is gone or expired. A second factor needs its verified first factor, which a new
    // session cannot have, so only the steps that start a login continue with a new session.
    if (secondFactorStep) {
      const reason = existing.reason ?? "server_session_not_found";
      logger.warn("LOGIN_SESSION_UPDATE_FAILED: could not update the session and did not replace it with a new one", {
        code: "LOGIN_SESSION_UPDATE_FAILED",
        reason,
        sessionId: recentSession.id,
        expired,
        secondFactorStep,
        grpcCode: isClassifiedError(error) ? error.code : undefined,
      });
      return { error: tError(SESSION_RESOLUTION_MESSAGE_KEYS[reason]), sessionLost: true };
    }

    const loginNameForCreation = options.loginName || recentSession?.loginName;
    const orgForCreation = options.organization || recentSession?.organization;

    if (!loginNameForCreation) {
      throw error;
    }

    const users = await listUsers({
      serviceConfig,
      loginName: loginNameForCreation,
      organizationId: orgForCreation,
    });

    if (users.details?.totalResult === BigInt(1) && users.result[0].userId) {
      const user = users.result[0];
      const newChecks = create(ChecksSchema, {
        ...(checks || {}),
        user: { search: { case: "userId", value: user.userId } } as any,
      });

      const result = await createSessionAndUpdateCookie({
        checks: newChecks,
        requestId,
        lifetime,
        challenges,
      });
      // @ts-ignore
      session = { ...result.session, challenges: result.challenges };
    } else {
      throw error;
    }
  }

  if (!session || ("error" in session && session.error)) {
    return { error: t("couldNotUpdateSession") };
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (checks && checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
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
    authMethods,
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
