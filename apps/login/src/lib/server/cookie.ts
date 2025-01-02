"use server";

import { addSessionToCookie, updateSessionCookie } from "@/lib/cookies";
import {
  createSessionForUserIdAndIdpIntent,
  createSessionFromChecks,
  getSession,
  setSession,
} from "@/lib/zitadel";
import { Duration, timestampMs } from "@zitadel/client";
import {
  Challenges,
  RequestChallenges,
} from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";

type CustomCookieData = {
  id: string;
  token: string;
  loginName: string;
  organization?: string;
  creationTs: string;
  expirationTs: string;
  changeTs: string;
  authRequestId?: string; // if its linked to an OIDC flow
};

export async function createSessionAndUpdateCookie(
  checks: Checks,
  challenges: RequestChallenges | undefined,
  authRequestId: string | undefined,
  lifetime?: Duration,
): Promise<Session> {
  const createdSession = await createSessionFromChecks(checks, challenges);

  if (createdSession) {
    return getSession({
      sessionId: createdSession.sessionId,
      sessionToken: createdSession.sessionToken,
    }).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: CustomCookieData = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationTs: response.session.creationDate
            ? `${timestampMs(response.session.creationDate)}`
            : "",
          expirationTs: response.session.expirationDate
            ? `${timestampMs(response.session.expirationDate)}`
            : "",
          changeTs: response.session.changeDate
            ? `${timestampMs(response.session.changeDate)}`
            : "",
          loginName: response.session.factors.user.loginName ?? "",
        };

        if (authRequestId) {
          sessionCookie.authRequestId = authRequestId;
        }

        if (response.session.factors.user.organizationId) {
          sessionCookie.organization =
            response.session.factors.user.organizationId;
        }

        return addSessionToCookie(sessionCookie).then(() => {
          return response.session as Session;
        });
      } else {
        throw "could not get session or session does not have loginName";
      }
    });
  } else {
    throw "Could not create session";
  }
}

export async function createSessionForIdpAndUpdateCookie(
  userId: string,
  idpIntent: {
    idpIntentId?: string | undefined;
    idpIntentToken?: string | undefined;
  },
  authRequestId: string | undefined,
  lifetime?: Duration,
): Promise<Session> {
  const createdSession = await createSessionForUserIdAndIdpIntent(
    userId,
    idpIntent,
    lifetime,
  );

  if (!createdSession) {
    throw "Could not create session";
  }

  const { session } = await getSession({
    sessionId: createdSession.sessionId,
    sessionToken: createdSession.sessionToken,
  });

  if (!session || !session.factors?.user?.loginName) {
    throw "Could not retrieve session";
  }

  const sessionCookie: CustomCookieData = {
    id: createdSession.sessionId,
    token: createdSession.sessionToken,
    creationTs: session.creationDate
      ? `${timestampMs(session.creationDate)}`
      : "",
    expirationTs: session.expirationDate
      ? `${timestampMs(session.expirationDate)}`
      : "",
    changeTs: session.changeDate ? `${timestampMs(session.changeDate)}` : "",
    loginName: session.factors.user.loginName ?? "",
    organization: session.factors.user.organizationId ?? "",
  };

  if (authRequestId) {
    sessionCookie.authRequestId = authRequestId;
  }

  if (session.factors.user.organizationId) {
    sessionCookie.organization = session.factors.user.organizationId;
  }

  return addSessionToCookie(sessionCookie).then(() => {
    return session as Session;
  });
}

export type SessionWithChallenges = Session & {
  challenges: Challenges | undefined;
};

export async function setSessionAndUpdateCookie(
  recentCookie: CustomCookieData,
  checks?: Checks,
  challenges?: RequestChallenges,
  authRequestId?: string,
  lifetime?: Duration,
) {
  return setSession(
    recentCookie.id,
    recentCookie.token,
    challenges,
    checks,
    lifetime,
  ).then((updatedSession) => {
    if (updatedSession) {
      const sessionCookie: CustomCookieData = {
        id: recentCookie.id,
        token: updatedSession.sessionToken,
        creationTs: recentCookie.creationTs,
        expirationTs: recentCookie.expirationTs,
        // just overwrite the changeDate with the new one
        changeTs: updatedSession.details?.changeDate
          ? `${timestampMs(updatedSession.details.changeDate)}`
          : "",
        loginName: recentCookie.loginName,
        organization: recentCookie.organization,
      };

      if (authRequestId) {
        sessionCookie.authRequestId = authRequestId;
      }

      return getSession({
        sessionId: sessionCookie.id,
        sessionToken: sessionCookie.token,
      }).then((response) => {
        if (response?.session && response.session.factors?.user?.loginName) {
          const { session } = response;
          const newCookie: CustomCookieData = {
            id: sessionCookie.id,
            token: updatedSession.sessionToken,
            creationTs: sessionCookie.creationTs,
            expirationTs: sessionCookie.expirationTs,
            // just overwrite the changeDate with the new one
            changeTs: updatedSession.details?.changeDate
              ? `${timestampMs(updatedSession.details.changeDate)}`
              : "",
            loginName: session.factors?.user?.loginName ?? "",
            organization: session.factors?.user?.organizationId ?? "",
          };

          if (sessionCookie.authRequestId) {
            newCookie.authRequestId = sessionCookie.authRequestId;
          }

          return updateSessionCookie(sessionCookie.id, newCookie).then(() => {
            return { challenges: updatedSession.challenges, ...session };
          });
        } else {
          throw "could not get session or session does not have loginName";
        }
      });
    } else {
      throw "Session not be set";
    }
  });
}
