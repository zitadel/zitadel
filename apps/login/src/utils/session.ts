"use server";

import {
  createSessionFromChecks,
  createSessionForUserIdAndIdpIntent,
  getSession,
  server,
  setSession,
} from "@/lib/zitadel";
import {
  SessionCookie,
  addSessionToCookie,
  updateSessionCookie,
} from "./cookies";
import {
  Session,
  Challenges,
  RequestChallenges,
  Checks,
} from "@zitadel/server";

export async function createSessionAndUpdateCookie(
  loginName: string,
  password: string | undefined,
  challenges: RequestChallenges | undefined,
  organization?: string,
  authRequestId?: string,
): Promise<Session> {
  const createdSession = await createSessionFromChecks(
    server,
    password
      ? {
          user: { loginName },
          password: { password },
          // totp: { code: totpCode },
        }
      : { user: { loginName } },
    challenges,
  );

  if (createdSession) {
    return getSession(
      server,
      createdSession.sessionId,
      createdSession.sessionToken,
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: SessionCookie = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationDate: `${response.session.creationDate?.getTime() ?? ""}`,
          expirationDate: `${response.session.expirationDate?.getTime() ?? ""}`,
          changeDate: `${response.session.changeDate?.getTime() ?? ""}`,
          loginName: response.session.factors.user.loginName ?? "",
          organization: response.session.factors.user.organizationId ?? "",
        };

        if (authRequestId) {
          sessionCookie.authRequestId = authRequestId;
        }

        if (organization) {
          sessionCookie.organization = organization;
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

export async function createSessionForUserIdAndUpdateCookie(
  userId: string,
  password: string | undefined,
  challenges: RequestChallenges | undefined,
  authRequestId: string | undefined,
): Promise<Session> {
  const createdSession = await createSessionFromChecks(
    server,
    password
      ? {
          user: { userId },
          password: { password },
          // totp: { code: totpCode },
        }
      : { user: { userId } },
    challenges,
  );

  if (createdSession) {
    return getSession(
      server,
      createdSession.sessionId,
      createdSession.sessionToken,
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: SessionCookie = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationDate: `${response.session.creationDate?.getTime() ?? ""}`,
          expirationDate: `${response.session.expirationDate?.getTime() ?? ""}`,
          changeDate: `${response.session.changeDate?.getTime() ?? ""}`,
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
  organization: string | undefined,
  authRequestId: string | undefined,
): Promise<Session> {
  const createdSession = await createSessionForUserIdAndIdpIntent(
    server,
    userId,
    idpIntent,
  );

  if (createdSession) {
    return getSession(
      server,
      createdSession.sessionId,
      createdSession.sessionToken,
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: SessionCookie = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationDate: `${response.session.creationDate?.getTime() ?? ""}`,
          expirationDate: `${response.session.expirationDate?.getTime() ?? ""}`,
          changeDate: `${response.session.changeDate?.getTime() ?? ""}`,
          loginName: response.session.factors.user.loginName ?? "",
          organization: response.session.factors.user.organizationId ?? "",
        };

        if (authRequestId) {
          sessionCookie.authRequestId = authRequestId;
        }

        if (organization) {
          sessionCookie.organization = organization;
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

export type SessionWithChallenges = Session & {
  challenges: Challenges | undefined;
};

export async function setSessionAndUpdateCookie(
  recentCookie: SessionCookie,
  checks: Checks,
  challenges: RequestChallenges | undefined,
  authRequestId: string | undefined,
): Promise<SessionWithChallenges> {
  return setSession(
    server,
    recentCookie.id,
    recentCookie.token,
    challenges,
    checks,
  ).then((updatedSession) => {
    if (updatedSession) {
      const sessionCookie: SessionCookie = {
        id: recentCookie.id,
        token: updatedSession.sessionToken,
        creationDate: recentCookie.creationDate,
        expirationDate: recentCookie.expirationDate,
        changeDate: `${updatedSession.details?.changeDate?.getTime() ?? ""}`,
        loginName: recentCookie.loginName,
        organization: recentCookie.organization,
      };

      if (authRequestId) {
        sessionCookie.authRequestId = authRequestId;
      }

      return getSession(server, sessionCookie.id, sessionCookie.token).then(
        (response) => {
          if (response?.session && response.session.factors?.user?.loginName) {
            const { session } = response;
            const newCookie: SessionCookie = {
              id: sessionCookie.id,
              token: updatedSession.sessionToken,
              creationDate: sessionCookie.creationDate,
              expirationDate: sessionCookie.expirationDate,
              changeDate: `${session.changeDate?.getTime() ?? ""}`,
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
        },
      );
    } else {
      throw "Session not be set";
    }
  });
}
