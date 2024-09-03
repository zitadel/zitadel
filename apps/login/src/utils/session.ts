"use server";

import {
  createSessionFromChecks,
  createSessionForUserIdAndIdpIntent,
  getSession,
  setSession,
} from "@/lib/zitadel";
import { addSessionToCookie, updateSessionCookie } from "@zitadel/next";
import {
  Challenges,
  RequestChallenges,
} from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import {
  Checks,
  ChecksSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { toDate } from "@zitadel/client";
import { create } from "@zitadel/client";

type CustomCookieData = {
  id: string;
  token: string;
  loginName: string;
  organization?: string;
  creationDate: string;
  expirationDate: string;
  changeDate: string;
  authRequestId?: string; // if its linked to an OIDC flow
};

export async function createSessionAndUpdateCookie(
  loginName: string,
  password: string | undefined,
  challenges: RequestChallenges | undefined,
  organization?: string,
  authRequestId?: string,
) {
  const checks = create(
    ChecksSchema,
    password
      ? {
          user: { search: { case: "loginName", value: loginName } },
          password: { password },
        }
      : { user: { search: { case: "loginName", value: loginName } } },
  );

  const createdSession = await createSessionFromChecks(checks, challenges);

  if (createdSession) {
    return getSession(
      createdSession.sessionId,
      createdSession.sessionToken,
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: CustomCookieData = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationDate: `${toDate(response.session.creationDate)?.getTime() ?? ""}`,
          expirationDate: `${toDate(response.session.expirationDate)?.getTime() ?? ""}`,
          changeDate: `${toDate(response.session.changeDate)?.getTime() ?? ""}`,
          loginName: response.session.factors.user.loginName ?? "",
          organization: response.session.factors.user.organizationId ?? "",
        };

        if (authRequestId) {
          sessionCookie.authRequestId = authRequestId;
        }

        if (organization) {
          sessionCookie.organization = organization;
        }

        return addSessionToCookie<CustomCookieData>(sessionCookie).then(() => {
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
  const checks = create(
    ChecksSchema,
    password
      ? {
          user: { search: { case: "userId", value: userId } },
          password: { password },
          // totp: { code: totpCode },
        }
      : { user: { search: { case: "userId", value: userId } } },
  );
  const createdSession = await createSessionFromChecks(checks, challenges);

  if (createdSession) {
    return getSession(
      createdSession.sessionId,
      createdSession.sessionToken,
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: CustomCookieData = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationDate: `${toDate(response.session.creationDate)?.getTime() ?? ""}`,
          expirationDate: `${toDate(response.session.expirationDate)?.getTime() ?? ""}`,
          changeDate: `${toDate(response.session.changeDate)?.getTime() ?? ""}`,
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
    userId,
    idpIntent,
  );

  if (createdSession) {
    return getSession(
      createdSession.sessionId,
      createdSession.sessionToken,
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: CustomCookieData = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          creationDate: `${toDate(response.session.creationDate)?.getTime() ?? ""}`,
          expirationDate: `${toDate(response.session.expirationDate)?.getTime() ?? ""}`,
          changeDate: `${toDate(response.session.changeDate)?.getTime() ?? ""}`,
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
  recentCookie: CustomCookieData,
  checks?: Checks,
  challenges?: RequestChallenges,
  authRequestId?: string,
) {
  return setSession(
    recentCookie.id,
    recentCookie.token,
    challenges,
    checks,
  ).then((updatedSession) => {
    if (updatedSession) {
      const sessionCookie: CustomCookieData = {
        id: recentCookie.id,
        token: updatedSession.sessionToken,
        creationDate: recentCookie.creationDate,
        expirationDate: recentCookie.expirationDate,
        changeDate: `${toDate(updatedSession.details?.changeDate)?.getTime() ?? ""}`,
        loginName: recentCookie.loginName,
        organization: recentCookie.organization,
      };

      if (authRequestId) {
        sessionCookie.authRequestId = authRequestId;
      }

      return getSession(sessionCookie.id, sessionCookie.token).then(
        (response) => {
          if (response?.session && response.session.factors?.user?.loginName) {
            const { session } = response;
            const newCookie: CustomCookieData = {
              id: sessionCookie.id,
              token: updatedSession.sessionToken,
              creationDate: sessionCookie.creationDate,
              expirationDate: sessionCookie.expirationDate,
              changeDate: `${toDate(session.changeDate)?.getTime() ?? ""}`,
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
