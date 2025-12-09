"use server";

import { addSessionToCookie, updateSessionCookie } from "@/lib/cookies";
import {
  createSessionForUserIdAndIdpIntent,
  createSessionFromChecks,
  getSecuritySettings,
  getSession,
  setSession,
} from "@/lib/zitadel";
import { ConnectError, Duration, timestampMs } from "@zitadel/client";
import { CredentialsCheckError, CredentialsCheckErrorSchema, ErrorDetail } from "@zitadel/proto/zitadel/message_pb";
import { Challenges, RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { headers } from "next/headers";
import { getServiceConfig } from "../service-url";

type CustomCookieData = {
  id: string;
  token: string;
  loginName: string;
  organization?: string;
  creationTs: string;
  expirationTs: string;
  changeTs: string;
  requestId?: string; // if its linked to an OIDC flow
};

const passwordAttemptsHandler = (error: ConnectError) => {
  const details = error.findDetails(CredentialsCheckErrorSchema);

  if (details[0] && "failedAttempts" in details[0]) {
    const failedAttempts = details[0].failedAttempts;
    throw {
      error: `Failed to authenticate: You had ${failedAttempts} password attempts.`,
      failedAttempts: failedAttempts,
    };
  }
  throw error;
};

export async function createSessionAndUpdateCookie(command: {
  checks: Checks;
  requestId: string | undefined;
  lifetime?: Duration;
}): Promise<Session> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let sessionLifetime = command.lifetime;

  if (!sessionLifetime || !sessionLifetime.seconds) {
    console.warn("No session lifetime provided, using default of 24 hours.");

    sessionLifetime = {
      seconds: BigInt(24 * 60 * 60), // 24 hours
      nanos: 0,
    } as Duration; // for usecases where the lifetime is not specified (user discovery)
  }

  const createdSession = await createSessionFromChecks({ serviceConfig, checks: command.checks,
    lifetime: sessionLifetime,
  });

  if (createdSession) {
    return getSession({ serviceConfig, sessionId: createdSession.sessionId,
      sessionToken: createdSession.sessionToken,
    }).then(async (response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: CustomCookieData = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationTs: response.session.creationDate ? `${timestampMs(response.session.creationDate)}` : "",
          expirationTs: response.session.expirationDate ? `${timestampMs(response.session.expirationDate)}` : "",
          changeTs: response.session.changeDate ? `${timestampMs(response.session.changeDate)}` : "",
          loginName: response.session.factors.user.loginName ?? "",
        };

        if (command.requestId) {
          sessionCookie.requestId = command.requestId;
        }

        if (response.session.factors.user.organizationId) {
          sessionCookie.organization = response.session.factors.user.organizationId;
        }

        const securitySettings = await getSecuritySettings({ serviceConfig });
        const iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;

        await addSessionToCookie({ session: sessionCookie, iFrameEnabled });

        return response.session as Session;
      } else {
        throw "could not get session or session does not have loginName";
      }
    });
  } else {
    throw "Could not create session";
  }
}

export async function createSessionForIdpAndUpdateCookie({
  userId,
  idpIntent,
  requestId,
  lifetime,
}: {
  userId: string;
  idpIntent: {
    idpIntentId?: string | undefined;
    idpIntentToken?: string | undefined;
  };
  requestId: string | undefined;
  lifetime?: Duration;
}): Promise<Session> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let sessionLifetime = lifetime;

  if (!sessionLifetime || !sessionLifetime.seconds) {
    console.warn("No IDP session lifetime provided, using default of 24 hours.");

    sessionLifetime = {
      seconds: BigInt(24 * 60 * 60), // 24 hours
      nanos: 0,
    } as Duration;
  }

  const createdSession = await createSessionForUserIdAndIdpIntent({ serviceConfig, userId,
    idpIntent,
    lifetime: sessionLifetime,
  }).catch((error: ErrorDetail | CredentialsCheckError) => {
    console.error("Could not set session", error);
    if ("failedAttempts" in error && error.failedAttempts) {
      throw {
        error: `Failed to authenticate: You had ${error.failedAttempts} password attempts.`,
        failedAttempts: error.failedAttempts,
      };
    }
    throw error;
  });

  if (!createdSession) {
    throw "Could not create session";
  }

  const { session } = await getSession({ serviceConfig, sessionId: createdSession.sessionId,
    sessionToken: createdSession.sessionToken,
  });

  if (!session || !session.factors?.user?.loginName) {
    throw "Could not retrieve session";
  }

  const sessionCookie: CustomCookieData = {
    id: createdSession.sessionId,
    token: createdSession.sessionToken,
    creationTs: session.creationDate ? `${timestampMs(session.creationDate)}` : "",
    expirationTs: session.expirationDate ? `${timestampMs(session.expirationDate)}` : "",
    changeTs: session.changeDate ? `${timestampMs(session.changeDate)}` : "",
    loginName: session.factors.user.loginName ?? "",
    organization: session.factors.user.organizationId ?? "",
  };

  if (requestId) {
    sessionCookie.requestId = requestId;
  }

  if (session.factors.user.organizationId) {
    sessionCookie.organization = session.factors.user.organizationId;
  }

  const securitySettings = await getSecuritySettings({ serviceConfig });
  const iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;

  return addSessionToCookie({ session: sessionCookie, iFrameEnabled }).then(() => {
    return session as Session;
  });
}

export type SessionWithChallenges = Session & {
  challenges: Challenges | undefined;
};

export async function setSessionAndUpdateCookie(command: {
  recentCookie: CustomCookieData;
  checks?: Checks;
  challenges?: RequestChallenges;
  requestId?: string;
  lifetime: Duration;
}) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  return setSession({ serviceConfig, sessionId: command.recentCookie.id,
    sessionToken: command.recentCookie.token,
    challenges: command.challenges,
    checks: command.checks,
    lifetime: command.lifetime,
  })
    .then((updatedSession) => {
      if (updatedSession) {
        const sessionCookie: CustomCookieData = {
          id: command.recentCookie.id,
          token: updatedSession.sessionToken,
          creationTs: command.recentCookie.creationTs,
          expirationTs: command.recentCookie.expirationTs,
          // just overwrite the changeDate with the new one
          changeTs: updatedSession.details?.changeDate ? `${timestampMs(updatedSession.details.changeDate)}` : "",
          loginName: command.recentCookie.loginName,
          organization: command.recentCookie.organization,
        };

        if (command.requestId) {
          sessionCookie.requestId = command.requestId;
        }

        return getSession({ serviceConfig, sessionId: sessionCookie.id,
          sessionToken: sessionCookie.token,
        }).then(async (response) => {
          if (!response?.session || !response.session.factors?.user?.loginName) {
            throw "could not get session or session does not have loginName";
          }

          const { session } = response;
          const newCookie: CustomCookieData = {
            id: sessionCookie.id,
            token: updatedSession.sessionToken,
            creationTs: sessionCookie.creationTs,
            expirationTs: sessionCookie.expirationTs,
            // just overwrite the changeDate with the new one
            changeTs: updatedSession.details?.changeDate ? `${timestampMs(updatedSession.details.changeDate)}` : "",
            loginName: session.factors?.user?.loginName ?? "",
            organization: session.factors?.user?.organizationId ?? "",
          };

          if (sessionCookie.requestId) {
            newCookie.requestId = sessionCookie.requestId;
          }

          const securitySettings = await getSecuritySettings({ serviceConfig });
          const iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;

          return updateSessionCookie({
            id: sessionCookie.id,
            session: newCookie,
            iFrameEnabled,
          }).then(() => {
            return { challenges: updatedSession.challenges, ...session };
          });
        });
      } else {
        throw "Session not be set";
      }
    })
    .catch(passwordAttemptsHandler);
}
