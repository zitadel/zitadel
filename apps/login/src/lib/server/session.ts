"use server";

import { setSessionAndUpdateCookie } from "@/lib/server/cookie";
import {
  deleteSession,
  getLoginSettings,
  getSecuritySettings,
  humanMFAInitSkipped,
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
import { getServiceUrlFromHeaders } from "../service-url";

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
}) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: organization,
  });

  await humanMFAInitSkipped({ serviceUrl, userId });

  const url =
    requestId && sessionId
      ? await getNextUrl(
          {
            sessionId: sessionId,
            requestId: requestId,
            organization: organization,
          },
          loginSettings?.defaultRedirectUri,
        )
      : loginName
        ? await getNextUrl(
            {
              loginName: loginName,
              organization: organization,
            },
            loginSettings?.defaultRedirectUri,
          )
        : null;
  if (url) {
    return { redirect: url };
  }
}

export async function continueWithSession({ requestId, ...session }: Session & { requestId?: string }) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: session.factors?.user?.organizationId,
  });

  const url =
    requestId && session.id && session.factors?.user
      ? await getNextUrl(
          {
            sessionId: session.id,
            requestId: requestId,
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
  requestId?: string;
  challenges?: RequestChallenges;
  lifetime?: Duration;
};

export async function updateSession(options: UpdateSessionCommand) {
  let { loginName, sessionId, organization, checks, requestId, challenges } = options;
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
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = _headers.get("host");

  if (!host) {
    return { error: "Could not get host" };
  }

  if (host && challenges && challenges.webAuthN && !challenges.webAuthN.domain) {
    const [hostname] = host.split(":");

    challenges.webAuthN.domain = hostname;
  }

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization,
  });

  let lifetime = checks?.webAuthN
    ? loginSettings?.multiFactorCheckLifetime // TODO different lifetime for webauthn u2f/passkey
    : checks?.otpEmail || checks?.otpSms
      ? loginSettings?.secondFactorCheckLifetime
      : undefined;

  if (!lifetime) {
    console.warn("No lifetime provided for session, defaulting to 24 hours");
    lifetime = {
      seconds: BigInt(60 * 60 * 24), // default to 24 hours
      nanos: 0,
    } as Duration;
  }

  const session = await setSessionAndUpdateCookie({
    recentCookie: recentSession,
    checks,
    challenges,
    requestId,
    lifetime,
  });

  if (!session) {
    return { error: "Could not update session" };
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (checks && checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({
      serviceUrl,
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
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const { sessionId } = options;

  const sessionCookie = await getSessionCookieById({ sessionId });

  const deleteResponse = await deleteSession({
    serviceUrl,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  const securitySettings = await getSecuritySettings({ serviceUrl });
  const iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;

  if (!deleteResponse) {
    throw new Error("Could not delete session");
  }

  return removeSessionFromCookie({ session: sessionCookie, iFrameEnabled });
}
