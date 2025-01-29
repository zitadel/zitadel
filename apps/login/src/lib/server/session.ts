"use server";

import { setSessionAndUpdateCookie } from "@/lib/server/cookie";
import {
  deleteSession,
  getLoginSettings,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { Duration } from "@zitadel/client";
import { RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { headers } from "next/headers";
import { getNextUrl } from "../client";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "../cookies";
import { getServiceUrlFromHeaders } from "../service";

export async function continueWithSession({
  authRequestId,
  ...session
}: Session & { authRequestId?: string }) {
  const _headers = await headers();
  const { serviceUrl, serviceRegion } = getServiceUrlFromHeaders(_headers);

  const loginSettings = await getLoginSettings({
    serviceUrl,
    serviceRegion,
    organization: session.factors?.user?.organizationId,
  });

  const url =
    authRequestId && session.id && session.factors?.user
      ? await getNextUrl(
          {
            sessionId: session.id,
            authRequestId: authRequestId,
            organization: session.factors.user.organizationId,
          },
          loginSettings?.defaultRedirectUri,
        )
      : session.factors?.user
        ? await getNextUrl(
            {
              loginName: session.factors.user.loginName,
              organization: session.factors.user.organizationId,
            },
            loginSettings?.defaultRedirectUri,
          )
        : null;
  if (url) {
    return { redirect: url };
  }
}

export type UpdateSessionCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  checks?: Checks;
  authRequestId?: string;
  challenges?: RequestChallenges;
  lifetime?: Duration;
};

export async function updateSession(options: UpdateSessionCommand) {
  let {
    loginName,
    sessionId,
    organization,
    checks,
    authRequestId,
    challenges,
  } = options;
  const recentSession = sessionId
    ? await getSessionCookieById({ sessionId })
    : loginName
      ? await getSessionCookieByLoginName({ loginName, organization })
      : await getMostRecentSessionCookie();

  if (!recentSession) {
    return {
      error: "Could not find session",
    };
  }

  const _headers = await headers();
  const { serviceUrl, serviceRegion } = getServiceUrlFromHeaders(_headers);
  const host = _headers.get("host");

  if (!host) {
    return { error: "Could not get host" };
  }

  if (
    host &&
    challenges &&
    challenges.webAuthN &&
    !challenges.webAuthN.domain
  ) {
    const [hostname, port] = host.split(":");

    challenges.webAuthN.domain = hostname;
  }

  const loginSettings = await getLoginSettings({ serviceUrl, organization });

  const lifetime = checks?.webAuthN
    ? loginSettings?.multiFactorCheckLifetime // TODO different lifetime for webauthn u2f/passkey
    : checks?.otpEmail || checks?.otpSms
      ? loginSettings?.secondFactorCheckLifetime
      : undefined;

  const session = await setSessionAndUpdateCookie(
    recentSession,
    checks,
    challenges,
    authRequestId,
    lifetime,
  );

  if (!session) {
    return { error: "Could not update session" };
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (checks && checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({
      serviceUrl,
      serviceRegion,
      userId: session.factors.user.id,
    });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  return {
    sessionId: session.id,
    factors: session.factors,
    challenges: session.challenges,
    authMethods,
  };
}

type ClearSessionOptions = {
  sessionId: string;
};

export async function clearSession(options: ClearSessionOptions) {
  const _headers = await headers();
  const { serviceUrl, serviceRegion } = getServiceUrlFromHeaders(_headers);

  const { sessionId } = options;

  const session = await getSessionCookieById({ sessionId });

  const deletedSession = await deleteSession({
    serviceUrl,
    serviceRegion,
    sessionId: session.id,
    sessionToken: session.token,
  });

  if (deletedSession) {
    return removeSessionFromCookie(session);
  }
}

type CleanupSessionCommand = {
  sessionId: string;
};

export async function cleanupSession({ sessionId }: CleanupSessionCommand) {
  const _headers = await headers();
  const { serviceUrl, serviceRegion } = getServiceUrlFromHeaders(_headers);

  const sessionCookie = await getSessionCookieById({ sessionId });

  const deleteResponse = await deleteSession({
    serviceUrl,
    serviceRegion,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!deleteResponse) {
    throw new Error("Could not delete session");
  }

  return removeSessionFromCookie(sessionCookie);
}
