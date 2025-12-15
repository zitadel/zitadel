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
import { getTranslations } from "next-intl/server";
import { completeFlowOrGetUrl } from "../client";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "../cookies";
import { getServiceConfig } from "../service-url";
import { getPublicHost } from "./host";

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
}): Promise<{ redirect: string } | { error: string }> {
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

  const loginSettings = await getLoginSettings({ serviceConfig, organization: session.factors?.user?.organizationId });

  if (requestId && session.id && session.factors?.user) {
    return completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: requestId,
        organization: session.factors.user.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );
  } else if (session.factors?.user) {
    return completeFlowOrGetUrl(
      {
        loginName: session.factors.user.loginName,
        organization: session.factors.user.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  // Fallback error if we couldn't determine where to redirect
  return { error: t("couldNotContinueSession") };
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
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  if (!host) {
    return { error: "Could not get host" };
  }

  if (host && challenges && challenges.webAuthN && !challenges.webAuthN.domain) {
    const [hostname] = host.split(":");

    challenges.webAuthN.domain = hostname;
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization });

  let lifetime = checks?.webAuthN
    ? loginSettings?.multiFactorCheckLifetime // TODO different lifetime for webauthn u2f/passkey
    : checks?.otpEmail || checks?.otpSms
      ? loginSettings?.secondFactorCheckLifetime
      : undefined;

  if (!lifetime || !lifetime.seconds) {
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
    const response = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });
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
  const { serviceConfig } = getServiceConfig(_headers);

  const { sessionId } = options;

  const sessionCookie = await getSessionCookieById({ sessionId });

  const deleteResponse = await deleteSession({
    serviceConfig,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  const securitySettings = await getSecuritySettings({ serviceConfig });
  const iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;

  if (!deleteResponse) {
    throw new Error("Could not delete session");
  }

  return removeSessionFromCookie({ session: sessionCookie, iFrameEnabled });
}
